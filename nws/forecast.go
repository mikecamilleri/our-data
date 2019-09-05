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

const (
	getSemidailyForecastForGridpointEndpointURLStringFmt = "gridpoints/%s/%d,%d/forecast"        // wfo, lat, lon
	getHourlyForecastForGridpointEndpointURLStringFmt    = "gridpoints/%s/%d,%d/forecast/hourly" // wfo, lat, lon
)

// A Forecast represents a forecast for a specific place on Earth returned from
// the NWS API.
//
// Forecasts contain a variable number of Periods, each representing an
// arbitrary length of time.
type Forecast struct {
	// Gridpoint Gridpoint

	TimeRetrieved time.Time
	TimeForecast  time.Time

	Periods []Period
}

// A Period represents the forecast for a particular range of time at a
// a particular place on Earth.
type Period struct {
	Number int
	Name   string

	TimeStart time.Time
	TimeEnd   time.Time

	IsDaytime        bool
	Temperature      ValueUnit
	TemperatureTrend string
	WindSpeedMin     ValueUnit
	WindSpeedMax     ValueUnit
	WindDirection    string
	ForecastShort    string
	ForecastDetailed string
}

// getSemidailyForceastForGridpoint retrieves from the NWS API the latest
// semni-daily forecast for a particular gridpoint.
//
// The NWS tends to refer to semni-daily forecasts simply as "forecast."
func getSemidailyForecastForGridpoint(httpClient *http.Client, httpUserAgentString string, apiURLString string, gridpoint Gridpoint) (*Forecast, error) {
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		apiURLString,
		fmt.Sprintf(
			getSemidailyForecastForGridpointEndpointURLStringFmt,
			gridpoint.WFO,
			gridpoint.GridX,
			gridpoint.GridY,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return newForecastFromForecastRespBody(respBody)
}

// getHourlyForecastForGridpoint retrieves from the NWS API the latest
// hourly forecast for a particular gridpoint.
func getHourlyForecastForGridpoint(httpClient *http.Client, httpUserAgentString string, apiURLString string, gridpoint Gridpoint) (*Forecast, error) {
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		apiURLString,
		fmt.Sprintf(getHourlyForecastForGridpointEndpointURLStringFmt, gridpoint.WFO, gridpoint.GridX, gridpoint.GridY),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return newForecastFromForecastRespBody(respBody)
}

// newForecastFromForecastRespBody returns a Forecast pointer, given a response
// body from the NWS API.
func newForecastFromForecastRespBody(respBody []byte) (*Forecast, error) {
	// unmarshal the body into a temporary struct
	fRaw := struct {
		Properties struct {
			UpdateTime string
			Periods    []struct {
				Number           string
				Name             string
				StartTime        string
				EndTime          string
				IsDaytime        bool
				Temperature      string
				TemperatureUnit  string
				TemperatureTrend string
				WindSpeed        string // "2 to 7 mph" or "5 mph"
				WindDirection    string
				ShortForecast    string
				DetailedForecast string
			}
		}
	}{}
	if err := json.Unmarshal(respBody, &fRaw); err != nil {
		return nil, err
	}

	// validate and build returned slice
	var err error
	var f Forecast

	// must have valid times
	f.TimeRetrieved = time.Now()
	f.TimeForecast, err = time.Parse(time.RFC3339, fRaw.Properties.UpdateTime)
	if err != nil {
		return nil, err
	}

	// iterate through periods
	for _, pRaw := range fRaw.Properties.Periods {
		p := Period{}

		p.Number, err = strconv.Atoi(pRaw.Number)
		if err != nil {
			continue // skip if no number
		}
		p.TimeStart, err = time.Parse(time.RFC3339, pRaw.StartTime)
		if err != nil {
			continue // skip if bad start time
		}
		p.TimeEnd, err = time.Parse(time.RFC3339, pRaw.EndTime)
		if err != nil {
			continue // skip if bad end time
		}

		// ignore any missing or invalid fields
		p.Name = pRaw.Name
		p.IsDaytime = pRaw.IsDaytime

		tv, err := strconv.ParseFloat(pRaw.Temperature, 64)
		if err == nil && (pRaw.TemperatureUnit == "F" || pRaw.TemperatureUnit == "C") {
			p.Temperature.Value = tv
			p.Temperature.Unit = pRaw.TemperatureUnit
		}

		p.TemperatureTrend = pRaw.TemperatureTrend

		wsTokens := strings.Split(pRaw.WindSpeed, " ")
		if len(wsTokens) == 4 {
			p.WindSpeedMin.Value, err = strconv.ParseFloat(wsTokens[0], 64)
			if err == nil && wsTokens[3] == "mph" {
				p.WindSpeedMin.Unit = wsTokens[3]
			}
			p.WindSpeedMax.Value, err = strconv.ParseFloat(wsTokens[2], 64)
			if err == nil && wsTokens[3] == "mph" {
				p.WindSpeedMax.Unit = wsTokens[3]
			}
		}
		if len(wsTokens) == 2 {
			p.WindSpeedMin.Value, err = strconv.ParseFloat(wsTokens[0], 64)
			if err == nil && wsTokens[1] == "mph" {
				p.WindSpeedMin.Unit = wsTokens[1]
			}
			p.WindSpeedMax = p.WindSpeedMin
		}

		p.WindDirection = pRaw.WindDirection
		p.ForecastShort = pRaw.ShortForecast
		p.ForecastDetailed = pRaw.DetailedForecast

		f.Periods = append(f.Periods, p)
	}

	return &f, nil
}
