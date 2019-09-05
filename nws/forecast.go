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

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - semidaily and hourly

// mike@Darwin-D nws % curl -i -X GET "https://api.weather.gov/gridpoints/PQR/112,100/forecast"
// HTTP/2 200
// server: nginx/1.12.2
// content-type: application/geo+json
// access-control-allow-origin: *
// x-server-id: vm-bldr-nids-apiapp11.ncep.noaa.gov
// x-correlation-id: fadf7cd7-79d6-4c5a-b244-b74356019d03
// x-request-id: fadf7cd7-79d6-4c5a-b244-b74356019d03
// cache-control: public, max-age=897, s-maxage=120
// expires: Wed, 14 Aug 2019 17:47:37 GMT
// date: Wed, 14 Aug 2019 17:32:40 GMT
// content-length: 11332
// vary: Accept,Feature-Flags
// strict-transport-security: max-age=31536000 ; includeSubDomains ; preload

// {
//     "@context": [
//         "https://raw.githubusercontent.com/geojson/geojson-ld/master/contexts/geojson-base.jsonld",
//         {
//             "wx": "https://api.weather.gov/ontology#",
//             "geo": "http://www.opengis.net/ont/geosparql#",
//             "unit": "http://codes.wmo.int/common/unit/",
//             "@vocab": "https://api.weather.gov/ontology#"
//         }
//     ],
//     "type": "Feature",
//     "geometry": {
//         "type": "GeometryCollection",
//         "geometries": [
//             {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.65165829999999,
//                     45.463786800000001
//                 ]
//             },
//             {
//                 "type": "Polygon",
//                 "coordinates": [
//                     [
//                         [
//                             -122.6695926,
//                             45.4720613
//                         ],
//                         [
//                             -122.6634372,
//                             45.451199799999998
//                         ],
//                         [
//                             -122.63372720000001,
//                             45.455510099999998
//                         ],
//                         [
//                             -122.63987630000001,
//                             45.476371799999995
//                         ],
//                         [
//                             -122.6695926,
//                             45.4720613
//                         ]
//                     ]
//                 ]
//             }
//         ]
//     },
//     "properties": {
//         "updated": "2019-08-14T17:01:39+00:00",
//         "units": "us",
//         "forecastGenerator": "BaselineForecastGenerator",
//         "generatedAt": "2019-08-14T17:32:40+00:00",
//         "updateTime": "2019-08-14T17:01:39+00:00",
//         "validTimes": "2019-08-14T11:00:00+00:00/P8DT1H",
//         "elevation": {
//             "value": 60.960000000000001,
//             "unitCode": "unit:m"
//         },
//         "periods": [
//             {
//                 "number": 1,
//                 "name": "Today",
//                 "startTime": "2019-08-14T10:00:00-07:00",
//                 "endTime": "2019-08-14T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 86,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 to 9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=medium",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": "Sunny, with a high near 86. North northwest wind 5 to 9 mph."
//             },
//             {
//                 "number": 2,
//                 "name": "Tonight",
//                 "startTime": "2019-08-14T18:00:00-07:00",
//                 "endTime": "2019-08-15T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 to 9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=medium",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": "Mostly clear, with a low around 60. North northwest wind 3 to 9 mph."
//             },
//             {
//                 "number": 3,
//                 "name": "Thursday",
//                 "startTime": "2019-08-15T06:00:00-07:00",
//                 "endTime": "2019-08-15T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 83,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 to 8 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=medium",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": "Mostly sunny, with a high near 83. Northwest wind 3 to 8 mph."
//             },
//             {
//                 "number": 4,
//                 "name": "Thursday Night",
//                 "startTime": "2019-08-15T18:00:00-07:00",
//                 "endTime": "2019-08-16T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 to 8 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct/rain?size=medium",
//                 "shortForecast": "Partly Cloudy then Patchy Drizzle",
//                 "detailedForecast": "Patchy drizzle after 5am. Partly cloudy, with a low around 60. West northwest wind 3 to 8 mph."
//             },
//             {
//                 "number": 5,
//                 "name": "Friday",
//                 "startTime": "2019-08-16T06:00:00-07:00",
//                 "endTime": "2019-08-16T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 to 7 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/rain/bkn?size=medium",
//                 "shortForecast": "Patchy Drizzle then Partly Sunny",
//                 "detailedForecast": "Patchy drizzle before 11am. Partly sunny, with a high near 78. West southwest wind 3 to 7 mph."
//             },
//             {
//                 "number": 6,
//                 "name": "Friday Night",
//                 "startTime": "2019-08-16T18:00:00-07:00",
//                 "endTime": "2019-08-17T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 to 7 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=medium",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": "Mostly cloudy, with a low around 60."
//             },
//             {
//                 "number": 7,
//                 "name": "Saturday",
//                 "startTime": "2019-08-17T06:00:00-07:00",
//                 "endTime": "2019-08-17T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=medium",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": "Mostly cloudy, with a high near 77."
//             },
//             {
//                 "number": 8,
//                 "name": "Saturday Night",
//                 "startTime": "2019-08-17T18:00:00-07:00",
//                 "endTime": "2019-08-18T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=medium",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": "Partly cloudy, with a low around 61."
//             },
//             {
//                 "number": 9,
//                 "name": "Sunday",
//                 "startTime": "2019-08-18T06:00:00-07:00",
//                 "endTime": "2019-08-18T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=medium",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": "Partly sunny, with a high near 78."
//             },
//             {
//                 "number": 10,
//                 "name": "Sunday Night",
//                 "startTime": "2019-08-18T18:00:00-07:00",
//                 "endTime": "2019-08-19T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=medium",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": "Partly cloudy, with a low around 60."
//             },
//             {
//                 "number": 11,
//                 "name": "Monday",
//                 "startTime": "2019-08-19T06:00:00-07:00",
//                 "endTime": "2019-08-19T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 81,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=medium",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": "Mostly sunny, with a high near 81."
//             },
//             {
//                 "number": 12,
//                 "name": "Monday Night",
//                 "startTime": "2019-08-19T18:00:00-07:00",
//                 "endTime": "2019-08-20T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=medium",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": "Partly cloudy, with a low around 60."
//             },
//             {
//                 "number": 13,
//                 "name": "Tuesday",
//                 "startTime": "2019-08-20T06:00:00-07:00",
//                 "endTime": "2019-08-20T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=medium",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": "Partly sunny, with a high near 78."
//             },
//             {
//                 "number": 14,
//                 "name": "Tuesday Night",
//                 "startTime": "2019-08-20T18:00:00-07:00",
//                 "endTime": "2019-08-21T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=medium",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": "Mostly cloudy, with a low around 60."
//             }
//         ]
//     }
// }%
// mike@Darwin-D nws %
// mike@Darwin-D nws %
// mike@Darwin-D nws % curl -i -X GET "https://api.weather.gov/gridpoints/PQR/112,100/forecast/hourly"
// HTTP/2 200
// server: nginx/1.12.2
// content-type: application/geo+json
// access-control-allow-origin: *
// x-server-id: vm-bldr-nids-apiapp3.ncep.noaa.gov
// x-correlation-id: 8c6d135b-9a13-4639-b358-0d0016f1a43a
// x-request-id: 8c6d135b-9a13-4639-b358-0d0016f1a43a
// cache-control: public, max-age=894, s-maxage=120
// expires: Wed, 14 Aug 2019 17:47:41 GMT
// date: Wed, 14 Aug 2019 17:32:47 GMT
// vary: Accept,Feature-Flags
// strict-transport-security: max-age=31536000 ; includeSubDomains ; preload

// {
//     "@context": [
//         "https://raw.githubusercontent.com/geojson/geojson-ld/master/contexts/geojson-base.jsonld",
//         {
//             "wx": "https://api.weather.gov/ontology#",
//             "geo": "http://www.opengis.net/ont/geosparql#",
//             "unit": "http://codes.wmo.int/common/unit/",
//             "@vocab": "https://api.weather.gov/ontology#"
//         }
//     ],
//     "type": "Feature",
//     "geometry": {
//         "type": "GeometryCollection",
//         "geometries": [
//             {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.65165829999999,
//                     45.463786800000001
//                 ]
//             },
//             {
//                 "type": "Polygon",
//                 "coordinates": [
//                     [
//                         [
//                             -122.6695926,
//                             45.4720613
//                         ],
//                         [
//                             -122.6634372,
//                             45.451199799999998
//                         ],
//                         [
//                             -122.63372720000001,
//                             45.455510099999998
//                         ],
//                         [
//                             -122.63987630000001,
//                             45.476371799999995
//                         ],
//                         [
//                             -122.6695926,
//                             45.4720613
//                         ]
//                     ]
//                 ]
//             }
//         ]
//     },
//     "properties": {
//         "updated": "2019-08-14T17:01:39+00:00",
//         "units": "us",
//         "forecastGenerator": "HourlyForecastGenerator",
//         "generatedAt": "2019-08-14T17:32:47+00:00",
//         "updateTime": "2019-08-14T17:01:39+00:00",
//         "validTimes": "2019-08-14T11:00:00+00:00/P8DT1H",
//         "elevation": {
//             "value": 60.960000000000001,
//             "unitCode": "unit:m"
//         },
//         "periods": [
//             {
//                 "number": 1,
//                 "name": "",
//                 "startTime": "2019-08-14T10:00:00-07:00",
//                 "endTime": "2019-08-14T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 2,
//                 "name": "",
//                 "startTime": "2019-08-14T11:00:00-07:00",
//                 "endTime": "2019-08-14T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 74,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/skc?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 3,
//                 "name": "",
//                 "startTime": "2019-08-14T12:00:00-07:00",
//                 "endTime": "2019-08-14T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 4,
//                 "name": "",
//                 "startTime": "2019-08-14T13:00:00-07:00",
//                 "endTime": "2019-08-14T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 80,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 5,
//                 "name": "",
//                 "startTime": "2019-08-14T14:00:00-07:00",
//                 "endTime": "2019-08-14T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 83,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 6,
//                 "name": "",
//                 "startTime": "2019-08-14T15:00:00-07:00",
//                 "endTime": "2019-08-14T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 85,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 7,
//                 "name": "",
//                 "startTime": "2019-08-14T16:00:00-07:00",
//                 "endTime": "2019-08-14T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 86,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 8,
//                 "name": "",
//                 "startTime": "2019-08-14T17:00:00-07:00",
//                 "endTime": "2019-08-14T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 86,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/skc?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 9,
//                 "name": "",
//                 "startTime": "2019-08-14T18:00:00-07:00",
//                 "endTime": "2019-08-14T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 85,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 10,
//                 "name": "",
//                 "startTime": "2019-08-14T19:00:00-07:00",
//                 "endTime": "2019-08-14T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 83,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 11,
//                 "name": "",
//                 "startTime": "2019-08-14T20:00:00-07:00",
//                 "endTime": "2019-08-14T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 80,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 12,
//                 "name": "",
//                 "startTime": "2019-08-14T21:00:00-07:00",
//                 "endTime": "2019-08-14T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 13,
//                 "name": "",
//                 "startTime": "2019-08-14T22:00:00-07:00",
//                 "endTime": "2019-08-14T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "9 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 14,
//                 "name": "",
//                 "startTime": "2019-08-14T23:00:00-07:00",
//                 "endTime": "2019-08-15T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 15,
//                 "name": "",
//                 "startTime": "2019-08-15T00:00:00-07:00",
//                 "endTime": "2019-08-15T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 16,
//                 "name": "",
//                 "startTime": "2019-08-15T01:00:00-07:00",
//                 "endTime": "2019-08-15T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 66,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 17,
//                 "name": "",
//                 "startTime": "2019-08-15T02:00:00-07:00",
//                 "endTime": "2019-08-15T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 64,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 18,
//                 "name": "",
//                 "startTime": "2019-08-15T03:00:00-07:00",
//                 "endTime": "2019-08-15T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 19,
//                 "name": "",
//                 "startTime": "2019-08-15T04:00:00-07:00",
//                 "endTime": "2019-08-15T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "6 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 20,
//                 "name": "",
//                 "startTime": "2019-08-15T05:00:00-07:00",
//                 "endTime": "2019-08-15T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 21,
//                 "name": "",
//                 "startTime": "2019-08-15T06:00:00-07:00",
//                 "endTime": "2019-08-15T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 22,
//                 "name": "",
//                 "startTime": "2019-08-15T07:00:00-07:00",
//                 "endTime": "2019-08-15T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 23,
//                 "name": "",
//                 "startTime": "2019-08-15T08:00:00-07:00",
//                 "endTime": "2019-08-15T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 24,
//                 "name": "",
//                 "startTime": "2019-08-15T09:00:00-07:00",
//                 "endTime": "2019-08-15T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 25,
//                 "name": "",
//                 "startTime": "2019-08-15T10:00:00-07:00",
//                 "endTime": "2019-08-15T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 66,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 26,
//                 "name": "",
//                 "startTime": "2019-08-15T11:00:00-07:00",
//                 "endTime": "2019-08-15T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 27,
//                 "name": "",
//                 "startTime": "2019-08-15T12:00:00-07:00",
//                 "endTime": "2019-08-15T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 28,
//                 "name": "",
//                 "startTime": "2019-08-15T13:00:00-07:00",
//                 "endTime": "2019-08-15T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 29,
//                 "name": "",
//                 "startTime": "2019-08-15T14:00:00-07:00",
//                 "endTime": "2019-08-15T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 79,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 30,
//                 "name": "",
//                 "startTime": "2019-08-15T15:00:00-07:00",
//                 "endTime": "2019-08-15T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 81,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 31,
//                 "name": "",
//                 "startTime": "2019-08-15T16:00:00-07:00",
//                 "endTime": "2019-08-15T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 82,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 32,
//                 "name": "",
//                 "startTime": "2019-08-15T17:00:00-07:00",
//                 "endTime": "2019-08-15T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 83,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 33,
//                 "name": "",
//                 "startTime": "2019-08-15T18:00:00-07:00",
//                 "endTime": "2019-08-15T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 81,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 34,
//                 "name": "",
//                 "startTime": "2019-08-15T19:00:00-07:00",
//                 "endTime": "2019-08-15T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 79,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 35,
//                 "name": "",
//                 "startTime": "2019-08-15T20:00:00-07:00",
//                 "endTime": "2019-08-15T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 36,
//                 "name": "",
//                 "startTime": "2019-08-15T21:00:00-07:00",
//                 "endTime": "2019-08-15T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 37,
//                 "name": "",
//                 "startTime": "2019-08-15T22:00:00-07:00",
//                 "endTime": "2019-08-15T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 71,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "8 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 38,
//                 "name": "",
//                 "startTime": "2019-08-15T23:00:00-07:00",
//                 "endTime": "2019-08-16T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 68,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 39,
//                 "name": "",
//                 "startTime": "2019-08-16T00:00:00-07:00",
//                 "endTime": "2019-08-16T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 66,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 40,
//                 "name": "",
//                 "startTime": "2019-08-16T01:00:00-07:00",
//                 "endTime": "2019-08-16T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 41,
//                 "name": "",
//                 "startTime": "2019-08-16T02:00:00-07:00",
//                 "endTime": "2019-08-16T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 42,
//                 "name": "",
//                 "startTime": "2019-08-16T03:00:00-07:00",
//                 "endTime": "2019-08-16T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 43,
//                 "name": "",
//                 "startTime": "2019-08-16T04:00:00-07:00",
//                 "endTime": "2019-08-16T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 44,
//                 "name": "",
//                 "startTime": "2019-08-16T05:00:00-07:00",
//                 "endTime": "2019-08-16T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 45,
//                 "name": "",
//                 "startTime": "2019-08-16T06:00:00-07:00",
//                 "endTime": "2019-08-16T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 46,
//                 "name": "",
//                 "startTime": "2019-08-16T07:00:00-07:00",
//                 "endTime": "2019-08-16T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 47,
//                 "name": "",
//                 "startTime": "2019-08-16T08:00:00-07:00",
//                 "endTime": "2019-08-16T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 48,
//                 "name": "",
//                 "startTime": "2019-08-16T09:00:00-07:00",
//                 "endTime": "2019-08-16T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 49,
//                 "name": "",
//                 "startTime": "2019-08-16T10:00:00-07:00",
//                 "endTime": "2019-08-16T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/rain?size=small",
//                 "shortForecast": "Patchy Drizzle",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 50,
//                 "name": "",
//                 "startTime": "2019-08-16T11:00:00-07:00",
//                 "endTime": "2019-08-16T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 68,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 51,
//                 "name": "",
//                 "startTime": "2019-08-16T12:00:00-07:00",
//                 "endTime": "2019-08-16T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 71,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 52,
//                 "name": "",
//                 "startTime": "2019-08-16T13:00:00-07:00",
//                 "endTime": "2019-08-16T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 53,
//                 "name": "",
//                 "startTime": "2019-08-16T14:00:00-07:00",
//                 "endTime": "2019-08-16T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 54,
//                 "name": "",
//                 "startTime": "2019-08-16T15:00:00-07:00",
//                 "endTime": "2019-08-16T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 55,
//                 "name": "",
//                 "startTime": "2019-08-16T16:00:00-07:00",
//                 "endTime": "2019-08-16T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "WSW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 56,
//                 "name": "",
//                 "startTime": "2019-08-16T17:00:00-07:00",
//                 "endTime": "2019-08-16T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 57,
//                 "name": "",
//                 "startTime": "2019-08-16T18:00:00-07:00",
//                 "endTime": "2019-08-16T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 58,
//                 "name": "",
//                 "startTime": "2019-08-16T19:00:00-07:00",
//                 "endTime": "2019-08-16T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 59,
//                 "name": "",
//                 "startTime": "2019-08-16T20:00:00-07:00",
//                 "endTime": "2019-08-16T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 74,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 60,
//                 "name": "",
//                 "startTime": "2019-08-16T21:00:00-07:00",
//                 "endTime": "2019-08-16T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 71,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 61,
//                 "name": "",
//                 "startTime": "2019-08-16T22:00:00-07:00",
//                 "endTime": "2019-08-16T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 69,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "7 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 62,
//                 "name": "",
//                 "startTime": "2019-08-16T23:00:00-07:00",
//                 "endTime": "2019-08-17T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 63,
//                 "name": "",
//                 "startTime": "2019-08-17T00:00:00-07:00",
//                 "endTime": "2019-08-17T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 64,
//                 "name": "",
//                 "startTime": "2019-08-17T01:00:00-07:00",
//                 "endTime": "2019-08-17T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 64,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 65,
//                 "name": "",
//                 "startTime": "2019-08-17T02:00:00-07:00",
//                 "endTime": "2019-08-17T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 66,
//                 "name": "",
//                 "startTime": "2019-08-17T03:00:00-07:00",
//                 "endTime": "2019-08-17T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 67,
//                 "name": "",
//                 "startTime": "2019-08-17T04:00:00-07:00",
//                 "endTime": "2019-08-17T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 68,
//                 "name": "",
//                 "startTime": "2019-08-17T05:00:00-07:00",
//                 "endTime": "2019-08-17T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 69,
//                 "name": "",
//                 "startTime": "2019-08-17T06:00:00-07:00",
//                 "endTime": "2019-08-17T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 70,
//                 "name": "",
//                 "startTime": "2019-08-17T07:00:00-07:00",
//                 "endTime": "2019-08-17T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 71,
//                 "name": "",
//                 "startTime": "2019-08-17T08:00:00-07:00",
//                 "endTime": "2019-08-17T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 72,
//                 "name": "",
//                 "startTime": "2019-08-17T09:00:00-07:00",
//                 "endTime": "2019-08-17T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 73,
//                 "name": "",
//                 "startTime": "2019-08-17T10:00:00-07:00",
//                 "endTime": "2019-08-17T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 74,
//                 "name": "",
//                 "startTime": "2019-08-17T11:00:00-07:00",
//                 "endTime": "2019-08-17T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/ovc?size=small",
//                 "shortForecast": "Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 75,
//                 "name": "",
//                 "startTime": "2019-08-17T12:00:00-07:00",
//                 "endTime": "2019-08-17T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 69,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 76,
//                 "name": "",
//                 "startTime": "2019-08-17T13:00:00-07:00",
//                 "endTime": "2019-08-17T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 72,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 77,
//                 "name": "",
//                 "startTime": "2019-08-17T14:00:00-07:00",
//                 "endTime": "2019-08-17T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 74,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 78,
//                 "name": "",
//                 "startTime": "2019-08-17T15:00:00-07:00",
//                 "endTime": "2019-08-17T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 79,
//                 "name": "",
//                 "startTime": "2019-08-17T16:00:00-07:00",
//                 "endTime": "2019-08-17T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 80,
//                 "name": "",
//                 "startTime": "2019-08-17T17:00:00-07:00",
//                 "endTime": "2019-08-17T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 81,
//                 "name": "",
//                 "startTime": "2019-08-17T18:00:00-07:00",
//                 "endTime": "2019-08-17T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 82,
//                 "name": "",
//                 "startTime": "2019-08-17T19:00:00-07:00",
//                 "endTime": "2019-08-17T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 83,
//                 "name": "",
//                 "startTime": "2019-08-17T20:00:00-07:00",
//                 "endTime": "2019-08-17T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 84,
//                 "name": "",
//                 "startTime": "2019-08-17T21:00:00-07:00",
//                 "endTime": "2019-08-17T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 71,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 85,
//                 "name": "",
//                 "startTime": "2019-08-17T22:00:00-07:00",
//                 "endTime": "2019-08-17T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 69,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 86,
//                 "name": "",
//                 "startTime": "2019-08-17T23:00:00-07:00",
//                 "endTime": "2019-08-18T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 87,
//                 "name": "",
//                 "startTime": "2019-08-18T00:00:00-07:00",
//                 "endTime": "2019-08-18T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 88,
//                 "name": "",
//                 "startTime": "2019-08-18T01:00:00-07:00",
//                 "endTime": "2019-08-18T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 64,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 89,
//                 "name": "",
//                 "startTime": "2019-08-18T02:00:00-07:00",
//                 "endTime": "2019-08-18T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 90,
//                 "name": "",
//                 "startTime": "2019-08-18T03:00:00-07:00",
//                 "endTime": "2019-08-18T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 91,
//                 "name": "",
//                 "startTime": "2019-08-18T04:00:00-07:00",
//                 "endTime": "2019-08-18T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 92,
//                 "name": "",
//                 "startTime": "2019-08-18T05:00:00-07:00",
//                 "endTime": "2019-08-18T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 93,
//                 "name": "",
//                 "startTime": "2019-08-18T06:00:00-07:00",
//                 "endTime": "2019-08-18T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 94,
//                 "name": "",
//                 "startTime": "2019-08-18T07:00:00-07:00",
//                 "endTime": "2019-08-18T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 95,
//                 "name": "",
//                 "startTime": "2019-08-18T08:00:00-07:00",
//                 "endTime": "2019-08-18T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 96,
//                 "name": "",
//                 "startTime": "2019-08-18T09:00:00-07:00",
//                 "endTime": "2019-08-18T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 97,
//                 "name": "",
//                 "startTime": "2019-08-18T10:00:00-07:00",
//                 "endTime": "2019-08-18T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NNW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 98,
//                 "name": "",
//                 "startTime": "2019-08-18T11:00:00-07:00",
//                 "endTime": "2019-08-18T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 99,
//                 "name": "",
//                 "startTime": "2019-08-18T12:00:00-07:00",
//                 "endTime": "2019-08-18T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 72,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 100,
//                 "name": "",
//                 "startTime": "2019-08-18T13:00:00-07:00",
//                 "endTime": "2019-08-18T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 101,
//                 "name": "",
//                 "startTime": "2019-08-18T14:00:00-07:00",
//                 "endTime": "2019-08-18T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 102,
//                 "name": "",
//                 "startTime": "2019-08-18T15:00:00-07:00",
//                 "endTime": "2019-08-18T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 103,
//                 "name": "",
//                 "startTime": "2019-08-18T16:00:00-07:00",
//                 "endTime": "2019-08-18T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 104,
//                 "name": "",
//                 "startTime": "2019-08-18T17:00:00-07:00",
//                 "endTime": "2019-08-18T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 105,
//                 "name": "",
//                 "startTime": "2019-08-18T18:00:00-07:00",
//                 "endTime": "2019-08-18T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 106,
//                 "name": "",
//                 "startTime": "2019-08-18T19:00:00-07:00",
//                 "endTime": "2019-08-18T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 107,
//                 "name": "",
//                 "startTime": "2019-08-18T20:00:00-07:00",
//                 "endTime": "2019-08-18T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 108,
//                 "name": "",
//                 "startTime": "2019-08-18T21:00:00-07:00",
//                 "endTime": "2019-08-18T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 109,
//                 "name": "",
//                 "startTime": "2019-08-18T22:00:00-07:00",
//                 "endTime": "2019-08-18T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 67,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "5 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 110,
//                 "name": "",
//                 "startTime": "2019-08-18T23:00:00-07:00",
//                 "endTime": "2019-08-19T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 111,
//                 "name": "",
//                 "startTime": "2019-08-19T00:00:00-07:00",
//                 "endTime": "2019-08-19T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 64,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 112,
//                 "name": "",
//                 "startTime": "2019-08-19T01:00:00-07:00",
//                 "endTime": "2019-08-19T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 113,
//                 "name": "",
//                 "startTime": "2019-08-19T02:00:00-07:00",
//                 "endTime": "2019-08-19T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 114,
//                 "name": "",
//                 "startTime": "2019-08-19T03:00:00-07:00",
//                 "endTime": "2019-08-19T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 115,
//                 "name": "",
//                 "startTime": "2019-08-19T04:00:00-07:00",
//                 "endTime": "2019-08-19T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 116,
//                 "name": "",
//                 "startTime": "2019-08-19T05:00:00-07:00",
//                 "endTime": "2019-08-19T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 117,
//                 "name": "",
//                 "startTime": "2019-08-19T06:00:00-07:00",
//                 "endTime": "2019-08-19T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 118,
//                 "name": "",
//                 "startTime": "2019-08-19T07:00:00-07:00",
//                 "endTime": "2019-08-19T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 119,
//                 "name": "",
//                 "startTime": "2019-08-19T08:00:00-07:00",
//                 "endTime": "2019-08-19T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 120,
//                 "name": "",
//                 "startTime": "2019-08-19T09:00:00-07:00",
//                 "endTime": "2019-08-19T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 121,
//                 "name": "",
//                 "startTime": "2019-08-19T10:00:00-07:00",
//                 "endTime": "2019-08-19T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 68,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 122,
//                 "name": "",
//                 "startTime": "2019-08-19T11:00:00-07:00",
//                 "endTime": "2019-08-19T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 123,
//                 "name": "",
//                 "startTime": "2019-08-19T12:00:00-07:00",
//                 "endTime": "2019-08-19T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 72,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 124,
//                 "name": "",
//                 "startTime": "2019-08-19T13:00:00-07:00",
//                 "endTime": "2019-08-19T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 125,
//                 "name": "",
//                 "startTime": "2019-08-19T14:00:00-07:00",
//                 "endTime": "2019-08-19T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 126,
//                 "name": "",
//                 "startTime": "2019-08-19T15:00:00-07:00",
//                 "endTime": "2019-08-19T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 79,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 127,
//                 "name": "",
//                 "startTime": "2019-08-19T16:00:00-07:00",
//                 "endTime": "2019-08-19T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 81,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "WNW",
//                 "icon": "https://api.weather.gov/icons/land/day/sct?size=small",
//                 "shortForecast": "Mostly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 128,
//                 "name": "",
//                 "startTime": "2019-08-19T17:00:00-07:00",
//                 "endTime": "2019-08-19T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 81,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/day/few?size=small",
//                 "shortForecast": "Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 129,
//                 "name": "",
//                 "startTime": "2019-08-19T18:00:00-07:00",
//                 "endTime": "2019-08-19T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 80,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 130,
//                 "name": "",
//                 "startTime": "2019-08-19T19:00:00-07:00",
//                 "endTime": "2019-08-19T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 131,
//                 "name": "",
//                 "startTime": "2019-08-19T20:00:00-07:00",
//                 "endTime": "2019-08-19T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 132,
//                 "name": "",
//                 "startTime": "2019-08-19T21:00:00-07:00",
//                 "endTime": "2019-08-19T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 133,
//                 "name": "",
//                 "startTime": "2019-08-19T22:00:00-07:00",
//                 "endTime": "2019-08-19T23:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "3 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/few?size=small",
//                 "shortForecast": "Mostly Clear",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 134,
//                 "name": "",
//                 "startTime": "2019-08-19T23:00:00-07:00",
//                 "endTime": "2019-08-20T00:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 68,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 135,
//                 "name": "",
//                 "startTime": "2019-08-20T00:00:00-07:00",
//                 "endTime": "2019-08-20T01:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 66,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 136,
//                 "name": "",
//                 "startTime": "2019-08-20T01:00:00-07:00",
//                 "endTime": "2019-08-20T02:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 64,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 137,
//                 "name": "",
//                 "startTime": "2019-08-20T02:00:00-07:00",
//                 "endTime": "2019-08-20T03:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 138,
//                 "name": "",
//                 "startTime": "2019-08-20T03:00:00-07:00",
//                 "endTime": "2019-08-20T04:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 61,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 139,
//                 "name": "",
//                 "startTime": "2019-08-20T04:00:00-07:00",
//                 "endTime": "2019-08-20T05:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "NW",
//                 "icon": "https://api.weather.gov/icons/land/night/sct?size=small",
//                 "shortForecast": "Partly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 140,
//                 "name": "",
//                 "startTime": "2019-08-20T05:00:00-07:00",
//                 "endTime": "2019-08-20T06:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 141,
//                 "name": "",
//                 "startTime": "2019-08-20T06:00:00-07:00",
//                 "endTime": "2019-08-20T07:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 60,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 142,
//                 "name": "",
//                 "startTime": "2019-08-20T07:00:00-07:00",
//                 "endTime": "2019-08-20T08:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 62,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 143,
//                 "name": "",
//                 "startTime": "2019-08-20T08:00:00-07:00",
//                 "endTime": "2019-08-20T09:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 63,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 144,
//                 "name": "",
//                 "startTime": "2019-08-20T09:00:00-07:00",
//                 "endTime": "2019-08-20T10:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 65,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 145,
//                 "name": "",
//                 "startTime": "2019-08-20T10:00:00-07:00",
//                 "endTime": "2019-08-20T11:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 68,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 146,
//                 "name": "",
//                 "startTime": "2019-08-20T11:00:00-07:00",
//                 "endTime": "2019-08-20T12:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 70,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 147,
//                 "name": "",
//                 "startTime": "2019-08-20T12:00:00-07:00",
//                 "endTime": "2019-08-20T13:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 73,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 148,
//                 "name": "",
//                 "startTime": "2019-08-20T13:00:00-07:00",
//                 "endTime": "2019-08-20T14:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 75,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 149,
//                 "name": "",
//                 "startTime": "2019-08-20T14:00:00-07:00",
//                 "endTime": "2019-08-20T15:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 150,
//                 "name": "",
//                 "startTime": "2019-08-20T15:00:00-07:00",
//                 "endTime": "2019-08-20T16:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 151,
//                 "name": "",
//                 "startTime": "2019-08-20T16:00:00-07:00",
//                 "endTime": "2019-08-20T17:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "2 mph",
//                 "windDirection": "N",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 152,
//                 "name": "",
//                 "startTime": "2019-08-20T17:00:00-07:00",
//                 "endTime": "2019-08-20T18:00:00-07:00",
//                 "isDaytime": true,
//                 "temperature": 78,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/day/bkn?size=small",
//                 "shortForecast": "Partly Sunny",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 153,
//                 "name": "",
//                 "startTime": "2019-08-20T18:00:00-07:00",
//                 "endTime": "2019-08-20T19:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 77,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 154,
//                 "name": "",
//                 "startTime": "2019-08-20T19:00:00-07:00",
//                 "endTime": "2019-08-20T20:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 76,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 155,
//                 "name": "",
//                 "startTime": "2019-08-20T20:00:00-07:00",
//                 "endTime": "2019-08-20T21:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 74,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             },
//             {
//                 "number": 156,
//                 "name": "",
//                 "startTime": "2019-08-20T21:00:00-07:00",
//                 "endTime": "2019-08-20T22:00:00-07:00",
//                 "isDaytime": false,
//                 "temperature": 72,
//                 "temperatureUnit": "F",
//                 "temperatureTrend": null,
//                 "windSpeed": "1 mph",
//                 "windDirection": "W",
//                 "icon": "https://api.weather.gov/icons/land/night/bkn?size=small",
//                 "shortForecast": "Mostly Cloudy",
//                 "detailedForecast": ""
//             }
//         ]
//     }
// }
// mike@Darwin-D nws %
