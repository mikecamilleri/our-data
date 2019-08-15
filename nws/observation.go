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
	"net/http"
)

// Observation ...
type Observation struct {
}

// getLatestObservationForStation ...
func getLatestObservationForStation(httpClinet *http.Client, station string) (Observation, error) {
	return Observation{}, nil
}

// newObservationFromJSON ?

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
