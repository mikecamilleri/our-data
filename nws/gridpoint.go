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
)

const getGridpointForPointEndpointURLStringFmt = "points/%f,%f" // lat, lon

// Gridpoint ...
type Gridpoint struct {
	WFO   string // weather forecast office
	GridX int
	GridY int
	City  string
	State string
}

// getGridpointForPoint ...
func getGridpointForPoint(httpClient *http.Client, httpUserAgentString string, point Point) (*Gridpoint, error) {
	respBody, err := doAPIRequest(httpClient, httpUserAgentString, fmt.Sprintf(getGridpointForPointEndpointURLStringFmt, point.Lat, point.Lon), nil)
	if err != nil {
		return nil, err
	}
	return newGridpointFromPointRespBody(respBody)
}

// newGridpointFromPointResponseBody ...
func newGridpointFromPointRespBody(respBody []byte) (*Gridpoint, error) {
	// unmarshal the body into a temporary struct
	gpRaw := struct {
		Properties struct {
			CWA              string
			GridX            string
			GridY            string
			RelativeLocation struct {
				Properties struct {
					City  string
					State string
				}
			}
		}
	}{}
	if err := json.Unmarshal(respBody, gpRaw); err != nil {
		return nil, err
	}

	// validate and build returned value
	var err error
	gp := Gridpoint{}

	// must have WFO, gridX, and gridY
	if len(gpRaw.Properties.CWA) != 3 {
		return nil, fmt.Errorf("WFO/CWA must be three characters: \"%s\" is %d characters", gpRaw.Properties.CWA, len(gpRaw.Properties.CWA))
	}
	gp.WFO = strings.ToUpper(gpRaw.Properties.CWA)
	if gp.GridX, err = strconv.Atoi(gpRaw.Properties.GridX); err != nil {
		return nil, fmt.Errorf("GridX must be an integer: \"%s\"", gpRaw.Properties.GridX)
	}
	if gp.GridY, err = strconv.Atoi(gpRaw.Properties.GridY); err != nil {
		return nil, fmt.Errorf("GridY must be an integer: \"%s\"", gpRaw.Properties.GridY)
	}

	gp.City = gpRaw.Properties.RelativeLocation.Properties.City
	gp.State = gpRaw.Properties.RelativeLocation.Properties.State

	return &gp, nil
}

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - redirects are 301

// mike@Darwin-D go % curl -X GET "https://api.weather.gov/points/45.45805556%2C-122.66361111"
// {
//     "correlationId": "5d87b573-4c63-4b4c-85be-ae3e5308e8f4",
//     "title": "Adjusting Precision Of Point Coordinate",
//     "type": "https://api.weather.gov/problems/AdjustPointPrecision",
//     "status": 301,
//     "detail": "The precision of latitude/longitude points is limited to 4 decimal digits for efficiency. The location attribute contains your request mapped to the nearest supported point. If your client supports it, you will be redirected.",
//     "instance": "https://api.weather.gov/requests/5d87b573-4c63-4b4c-85be-ae3e5308e8f4"
// }
// mike@Darwin-D go % curl -X GET "https://api.weather.gov/points/45.4580,-122.6636"
// {
//     "correlationId": "c1303b55-6f14-4d8d-9942-9c105c6c4fa8",
//     "title": "Adjusting Trailing Zeros Of Point Coordinate",
//     "type": "https://api.weather.gov/problems/AdjustTrailingZeroRedundancy",
//     "status": 301,
//     "detail": "The coordinates cannot have trailing zeros in the decimal digit. The location attribute contains your request with the redundancy removed. If your client supports it, you will be redirected.",
//     "instance": "https://api.weather.gov/requests/c1303b55-6f14-4d8d-9942-9c105c6c4fa8"
// }
// mike@Darwin-D go % curl -X GET "https://api.weather.gov/points/45.458,-122.6636"
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
//     "id": "https://api.weather.gov/points/45.458,-122.6636",
//     "type": "Feature",
//     "geometry": {
//         "type": "Point",
//         "coordinates": [
//             -122.6636,
//             45.457999999999998
//         ]
//     },
//     "properties": {
//         "@id": "https://api.weather.gov/points/45.458,-122.6636",
//         "@type": "wx:Point",
//         "cwa": "PQR",
//         "forecastOffice": "https://api.weather.gov/offices/PQR",
//         "gridX": 112,
//         "gridY": 100,
//         "forecast": "https://api.weather.gov/gridpoints/PQR/112,100/forecast",
//         "forecastHourly": "https://api.weather.gov/gridpoints/PQR/112,100/forecast/hourly",
//         "forecastGridData": "https://api.weather.gov/gridpoints/PQR/112,100",
//         "observationStations": "https://api.weather.gov/gridpoints/PQR/112,100/stations",
//         "relativeLocation": {
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.620915,
//                     45.444124000000002
//                 ]
//             },
//             "properties": {
//                 "city": "Milwaukie",
//                 "state": "OR",
//                 "distance": {
//                     "value": 3669.782693174343,
//                     "unitCode": "unit:m"
//                 },
//                 "bearing": {
//                     "value": 294,
//                     "unitCode": "unit:degrees_true"
//                 }
//             }
//         },
//         "forecastZone": "https://api.weather.gov/zones/forecast/ORZ006",
//         "county": "https://api.weather.gov/zones/county/ORC051",
//         "fireWeatherZone": "https://api.weather.gov/zones/fire/ORZ604",
//         "timeZone": "America/Los_Angeles",
//         "radarStation": "KRTX"
//     }
// }
// mike@Darwin-D go %
