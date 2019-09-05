// Copyright 2019 Michael Camilleri <mike@mikecamilleri.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const getLatestObeservationForStationEndpointURLStringFmt = "stations/%s/observations/latest" // id

var observationUnitCodes = map[string]string{
	"unit:degC":           "C",
	"unit:degree_(angle)": "degrees true",
	"unit:m_s-1":          "m/s",
	"unit:Pa":             "Pa",
	"unit:m":              "m",
	"unit:percent":        "percent",
}

// A Observation represents the weather at a particular a particular station
// at a particular point in time returned from the NWS API.
type Observation struct {
	StationID string

	TimeRetrieved time.Time
	TimeObserved  time.Time

	Temperature               ValueUnit
	Dewpoint                  ValueUnit
	WindDirection             ValueUnit
	WindSpeed                 ValueUnit
	WindGust                  ValueUnit
	BarometricPressure        ValueUnit
	SeaLevelPressure          ValueUnit
	Visibility                ValueUnit
	TemperatureLast24HoursMin ValueUnit
	TemperatureLast24HoursMax ValueUnit
	PrecipitationLastHour     ValueUnit
	PrecipitationLast3Hours   ValueUnit
	PrecipitationLast6Hours   ValueUnit
	RelativeHumidity          ValueUnit
	WindChill                 ValueUnit
	HeatIndex                 ValueUnit
	// CloudLayers

	METAR string // raw METAR string
}

// getLatestObservationForStation retrieves from the NWS API the latest
// observation from a particular station.
func getLatestObservationForStation(httpClient *http.Client, httpUserAgentString string, stationID string) (*Observation, error) {
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		fmt.Sprintf(getLatestObeservationForStationEndpointURLStringFmt, stationID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return newObservationFromStationObservationRespBody(respBody)
}

// newObservationFromStationObservationRespBody returns an Obsevation pointer,
// given a response body from the NWS API.
func newObservationFromStationObservationRespBody(respBody []byte) (*Observation, error) {
	// TODO: Eventually it probably makes sense to just parse the METAR. This
	// endpoint seems to be converting everything to SI units which doesn't
	// make sense given the source (METAR) and typical use of these data.

	// TODO: Currently the WMO uit codes are converted to easier to read unit
	// names. Eventually these should be standardized among packages in this
	// Git repo. These are also inconsistant with the forecast data from NWS.

	// unmarshal the body into a temporary struct
	oRaw := struct {
		Properties struct {
			Station     string // URL
			Timestamp   string // time observed
			RawMessage  string // raw METAR
			Temperature struct {
				Value    string
				UnitCode string
			}
			Dewpoint struct {
				Value    string
				UnitCode string
			}
			WindDirection struct {
				Value    string
				UnitCode string
			}
			WindSpeed struct {
				Value    string
				UnitCode string
			}
			WindGust struct {
				Value    string
				UnitCode string
			}
			BarometricPressure struct {
				Value    string
				UnitCode string
			}
			SeaLevelPressure struct {
				Value    string
				UnitCode string
			}
			Visibility struct {
				Value    string
				UnitCode string
			}
			MaxTemperatureLast24Hours struct {
				Value    string
				UnitCode string
			}
			MinTemperatureLast24Hours struct {
				Value    string
				UnitCode string
			}
			PrecipitationLastHour struct {
				Value    string
				UnitCode string
			}
			PrecipitationLast3Hours struct {
				Value    string
				UnitCode string
			}
			PrecipitationLast6Hours struct {
				Value    string
				UnitCode string
			}
			RelativeHumidity struct {
				Value    string
				UnitCode string
			}
			WindChill struct {
				Value    string
				UnitCode string
			}
			HeatIndex struct {
				Value    string
				UnitCode string
			}
		}
	}{}
	if err := json.Unmarshal(respBody, &oRaw); err != nil {
		return nil, err
	}

	// validate and build returned value
	var u string
	var uok bool
	var v float64
	var err error
	var o Observation

	// must have valid station ID and times
	o.StationID = strings.TrimPrefix(oRaw.Properties.Station, "https://api.weather.gov/stations/")
	if o.StationID == "" {
		return nil, fmt.Errorf("station string invalid: \"%s\"", oRaw.Properties.Station)
	}
	o.TimeRetrieved = time.Now()
	o.TimeObserved, err = time.Parse(time.RFC3339, oRaw.Properties.Timestamp)
	if err != nil {
		return nil, err
	}

	// ignore any properties that are null, malformed, or have unrecognized units
	v, err = strconv.ParseFloat(oRaw.Properties.Temperature.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.Temperature.UnitCode]
	if uok && err == nil {
		o.Temperature.Value = v
		o.Temperature.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.Dewpoint.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.Dewpoint.UnitCode]
	if uok && err == nil {
		o.Dewpoint.Value = v
		o.Dewpoint.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.WindDirection.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.WindDirection.UnitCode]
	if uok && err == nil {
		o.WindDirection.Value = v
		o.WindDirection.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.WindSpeed.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.WindSpeed.UnitCode]
	if uok && err == nil {
		o.WindSpeed.Value = v
		o.WindSpeed.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.WindGust.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.WindGust.UnitCode]
	if uok && err == nil {
		o.WindGust.Value = v
		o.WindGust.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.BarometricPressure.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.BarometricPressure.UnitCode]
	if uok && err == nil {
		o.BarometricPressure.Value = v
		o.BarometricPressure.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.SeaLevelPressure.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.SeaLevelPressure.UnitCode]
	if uok && err == nil {
		o.SeaLevelPressure.Value = v
		o.SeaLevelPressure.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.Visibility.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.Visibility.UnitCode]
	if uok && err == nil {
		o.Visibility.Value = v
		o.Visibility.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.MinTemperatureLast24Hours.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.MinTemperatureLast24Hours.UnitCode]
	if uok && err == nil {
		o.TemperatureLast24HoursMin.Value = v
		o.TemperatureLast24HoursMin.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.MaxTemperatureLast24Hours.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.MaxTemperatureLast24Hours.UnitCode]
	if uok && err == nil {
		o.TemperatureLast24HoursMax.Value = v
		o.TemperatureLast24HoursMax.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.PrecipitationLastHour.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.PrecipitationLastHour.UnitCode]
	if uok && err == nil {
		o.PrecipitationLastHour.Value = v
		o.PrecipitationLastHour.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.PrecipitationLast3Hours.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.PrecipitationLast3Hours.UnitCode]
	if uok && err == nil {
		o.PrecipitationLast3Hours.Value = v
		o.PrecipitationLast3Hours.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.PrecipitationLast6Hours.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.PrecipitationLast6Hours.UnitCode]
	if uok && err == nil {
		o.PrecipitationLast6Hours.Value = v
		o.PrecipitationLast6Hours.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.RelativeHumidity.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.RelativeHumidity.UnitCode]
	if uok && err == nil {
		o.RelativeHumidity.Value = v
		o.RelativeHumidity.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.WindChill.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.WindChill.UnitCode]
	if uok && err == nil {
		o.WindChill.Value = v
		o.WindChill.Unit = u
	}
	v, err = strconv.ParseFloat(oRaw.Properties.HeatIndex.Value, 64)
	u, uok = observationUnitCodes[oRaw.Properties.HeatIndex.UnitCode]
	if uok && err == nil {
		o.HeatIndex.Value = v
		o.HeatIndex.Unit = u
	}

	o.METAR = oRaw.Properties.RawMessage

	return &o, nil
}
