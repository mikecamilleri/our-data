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
	"unit:unit:m_s-1":     "m/s",
	"unit:Pa":             "Pa",
	"unit:m":              "m",
	"unit:percent":        "percent",
}

// Observation ...
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

// getLatestObservationForStation ...
func getLatestObservationForStation(httpClient *http.Client, httpUserAgentString string, stationID string) (*Observation, error) {
	respBody, err := doAPIRequest(httpClient, httpUserAgentString, fmt.Sprintf(getLatestObeservationForStationEndpointURLStringFmt, stationID), nil)
	if err != nil {
		return nil, err
	}
	return newObservationFromStationObservationRespBody(respBody)
}

// newObservationFromStationObservationRespBody ...
func newObservationFromStationObservationRespBody(respBody []byte) (*Observation, error) {
	// TODO: Eventually it probably makes sense to just parse the METAR. This
	// endpoint seems to be converting everything to SI units which doesn't
	// make sense given the source (METAR) and typical use of these data.

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
	if err := json.Unmarshal(respBody, oRaw); err != nil {
		return nil, err
	}

	// validate and build returned value
	var u string
	var uok bool
	var v float64
	var err error
	o := Observation{}

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

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - not all fields are populated in this example.
// - https://api.weather.gov/ontology returns 404.
// - should get more examples.

// mike@Darwin-D nws % curl -i -X GET "https://api.weather.gov/stations/KPDX/observations/latest"
// HTTP/2 200
// server: nginx/1.12.2
// content-type: application/geo+json
// last-modified: Wed, 14 Aug 2019 15:53:00 GMT
// access-control-allow-origin: *
// x-server-id: vm-bldr-nids-apiapp5.ncep.noaa.gov
// x-correlation-id: 77e17b0e-ad70-4bad-8cd2-f51e61df7ba0
// x-request-id: 77e17b0e-ad70-4bad-8cd2-f51e61df7ba0
// cache-control: public, max-age=299, s-maxage=300
// expires: Wed, 14 Aug 2019 16:51:30 GMT
// date: Wed, 14 Aug 2019 16:46:31 GMT
// content-length: 4720
// vary: Accept,Feature-Flags
// strict-transport-security: max-age=31536000 ; includeSubDomains ; preload

// {
//     "@context": [
//         "https://raw.githubusercontent.com/geojson/geojson-ld/master/contexts/geojson-base.jsonld",
//         {
//             "wx": "https://api.weather.gov/ontology#",
//             "s": "https://schema.org/",
//             "geo": "http://www.opengis.net/ont/geosparql#",
//             "unit": "http://codes.wmo.int/common/unit/",
//             "@vocab": "https://api.weather.gov/ontology#",
//             "geometry": {
//                 "@id": "s:GeoCoordinates",
//                 "@type": "geo:wktLiteral"
//             },
//             "city": "s:addressLocality",
//             "state": "s:addressRegion",
//             "distance": {
//                 "@id": "s:Distance",
//                 "@type": "s:QuantitativeValue"
//             },
//             "bearing": {
//                 "@type": "s:QuantitativeValue"
//             },
//             "value": {
//                 "@id": "s:value"
//             },
//             "unitCode": {
//                 "@id": "s:unitCode",
//                 "@type": "@id"
//             },
//             "forecastOffice": {
//                 "@type": "@id"
//             },
//             "forecastGridData": {
//                 "@type": "@id"
//             },
//             "publicZone": {
//                 "@type": "@id"
//             },
//             "county": {
//                 "@type": "@id"
//             }
//         }
//     ],
//     "id": "https://api.weather.gov/stations/KPDX/observations/2019-08-14T15:53:00+00:00",
//     "type": "Feature",
//     "geometry": {
//         "type": "Point",
//         "coordinates": [
//             -122.59999999999999,
//             45.600000000000001
//         ]
//     },
//     "properties": {
//         "@id": "https://api.weather.gov/stations/KPDX/observations/2019-08-14T15:53:00+00:00",
//         "@type": "wx:ObservationStation",
//         "elevation": {
//             "value": 12,
//             "unitCode": "unit:m"
//         },
//         "station": "https://api.weather.gov/stations/KPDX",
//         "timestamp": "2019-08-14T15:53:00+00:00",
//         "rawMessage": "KPDX 141553Z 30005KT 10SM SCT250 21/13 A3018 RMK AO2 SLP217 T02060128",
//         "textDescription": "Partly Cloudy",
//         "icon": "https://api.weather.gov/icons/land/day/sct?size=medium",
//         "presentWeather": [],
//         "temperature": {
//             "value": 20.600000000000023,
//             "unitCode": "unit:degC",
//             "qualityControl": "qc:V"
//         },
//         "dewpoint": {
//             "value": 12.800000000000011,
//             "unitCode": "unit:degC",
//             "qualityControl": "qc:V"
//         },
//         "windDirection": {
//             "value": 300,
//             "unitCode": "unit:degree_(angle)",
//             "qualityControl": "qc:V"
//         },
//         "windSpeed": {
//             "value": 2.6000000000000001,
//             "unitCode": "unit:m_s-1",
//             "qualityControl": "qc:V"
//         },
//         "windGust": {
//             "value": null,
//             "unitCode": "unit:m_s-1",
//             "qualityControl": "qc:Z"
//         },
//         "barometricPressure": {
//             "value": 102200,
//             "unitCode": "unit:Pa",
//             "qualityControl": "qc:V"
//         },
//         "seaLevelPressure": {
//             "value": 102170,
//             "unitCode": "unit:Pa",
//             "qualityControl": "qc:V"
//         },
//         "visibility": {
//             "value": 16090,
//             "unitCode": "unit:m",
//             "qualityControl": "qc:C"
//         },
//         "maxTemperatureLast24Hours": {
//             "value": null,
//             "unitCode": "unit:degC",
//             "qualityControl": null
//         },
//         "minTemperatureLast24Hours": {
//             "value": null,
//             "unitCode": "unit:degC",
//             "qualityControl": null
//         },
//         "precipitationLastHour": {
//             "value": null,
//             "unitCode": "unit:m",
//             "qualityControl": "qc:Z"
//         },
//         "precipitationLast3Hours": {
//             "value": null,
//             "unitCode": "unit:m",
//             "qualityControl": "qc:Z"
//         },
//         "precipitationLast6Hours": {
//             "value": null,
//             "unitCode": "unit:m",
//             "qualityControl": "qc:Z"
//         },
//         "relativeHumidity": {
//             "value": 60.935094465891581,
//             "unitCode": "unit:percent",
//             "qualityControl": "qc:C"
//         },
//         "windChill": {
//             "value": null,
//             "unitCode": "unit:degC",
//             "qualityControl": "qc:V"
//         },
//         "heatIndex": {
//             "value": null,
//             "unitCode": "unit:degC",
//             "qualityControl": "qc:V"
//         },
//         "cloudLayers": [
//             {
//                 "base": {
//                     "value": 7620,
//                     "unitCode": "unit:m"
//                 },
//                 "amount": "SCT"
//             }
//         ]
//     }
// }
// mike@Darwin-D nws %
