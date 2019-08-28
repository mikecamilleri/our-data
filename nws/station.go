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

const getStationsForGridpointEndpointURLStringFmt = "gridpoints/%s/%f,%f/stations" // lat, lon

// Station ...
type Station struct {
	ID    string // callsign
	Name  string
	Point Point
}

// getStationsForGridpoint ...
func getStationsForGridpoint(httpClient *http.Client, httpUserAgentString string, gridpoint Gridpoint) ([]Station, error) {
	respBody, err := doAPIRequest(httpClient, httpUserAgentString, fmt.Sprintf(getStationsForGridpointEndpointURLStringFmt, gridpoint.WFO, gridpoint.GridX, gridpoint.GridY), nil)
	if err != nil {
		return nil, err
	}
	return newStationsFromStationsRespBody(respBody)
}

// newStationsFromStationsRespBody ...
func newStationsFromStationsRespBody(respBody []byte) ([]Station, error) {
	// unmarshal the body into a temporary struct
	stnsRaw := struct {
		Features []struct {
			Geometry struct {
				Coordinates []string // lon, lat (annoying)
			}
			Properties struct {
				StationIdentifier string // callsign
				Name              string
			}
		}
	}{}
	if err := json.Unmarshal(respBody, stnsRaw); err != nil {
		return nil, err
	}

	// validate and build returned slice
	var stns []Station
	for _, sRaw := range stnsRaw.Features {
		if sRaw.Properties.StationIdentifier == "" {
			continue // skip if no callsign
		}
		s := Station{
			Name: sRaw.Properties.Name,
			ID:   strings.ToUpper(sRaw.Properties.StationIdentifier),
		}
		if len(sRaw.Geometry.Coordinates) == 2 {
			s.Point.Lon, _ = strconv.ParseFloat(sRaw.Geometry.Coordinates[1], 64)
			s.Point.Lat, _ = strconv.ParseFloat(sRaw.Geometry.Coordinates[0], 64)
		}
	}

	return stns, nil
}

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - very long, lots of stations
// - how are these sorted? - They appear to be sorted by distance :/

// mike@Darwin-D nws % curl -i -X GET "https://api.weather.gov/gridpoints/PQR/112,100/stations"
// HTTP/2 200
// server: nginx/1.12.2
// content-type: application/geo+json
// access-control-allow-origin: *
// x-server-id: vm-bldr-nids-apiapp11.ncep.noaa.gov
// x-correlation-id: c1a3bd75-1bfa-4163-bc33-a68b6f80fcc4
// x-request-id: c1a3bd75-1bfa-4163-bc33-a68b6f80fcc4
// cache-control: public, max-age=16458, s-maxage=120
// expires: Wed, 14 Aug 2019 21:39:05 GMT
// date: Wed, 14 Aug 2019 17:04:47 GMT
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
//             },
//             "observationStations": {
//                 "@container": "@list",
//                 "@type": "@id"
//             }
//         }
//     ],
//     "type": "FeatureCollection",
//     "features": [
//         {
//             "id": "https://api.weather.gov/stations/KPDX",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.60917000000001,
//                     45.595779999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KPDX",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 6.0960000000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KPDX",
//                 "name": "Portland, Portland International Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KVUO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.65419,
//                     45.621029999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KVUO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 6.0960000000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KVUO",
//                 "name": "Pearson Airfield",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KTTD",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.40889,
//                     45.551110000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KTTD",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 10.972800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KTTD",
//                 "name": "Portland, Portland-Troutdale Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ARAO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.75028,
//                     45.281939999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ARAO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 45.719999999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ARAO",
//                 "name": "AURORA",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KHIO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.95444000000001,
//                     45.54806
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KHIO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 61.874400000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KHIO",
//                 "name": "Portland, Portland-Hillsboro Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KUAO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.76555999999999,
//                     45.248890000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KUAO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 59.1312,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KUAO",
//                 "name": "Aurora, Aurora State Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/EGKO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.33111100000001,
//                     45.368056000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/EGKO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 223.11360000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "EGKO3",
//                 "name": "EAGLE CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/RRWO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.28,
//                     45.539999999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/RRWO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 9.1440000000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "RRWO3",
//                 "name": "Rooster Rock",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT95",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.8707,
//                     45.695399999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT95",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 17.0688,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT95",
//                 "name": "US30 at Rocky Point (US 30 MP 16.5)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/FOGO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.08360999999999,
//                     45.553060000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/FOGO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 54.864000000000004,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "FOGO",
//                 "name": "FOREST GROVE",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TBELL",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.203,
//                     45.569000000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TBELL",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 235.0008,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TBELL",
//                 "name": "Cape Horn",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TR951",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.34527799999999,
//                     45.723056
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TR951",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 350.52000000000004,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TR951",
//                 "name": "LARCH MT.",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KSPB",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.86221999999999,
//                     45.769170000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KSPB",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 15.849600000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KSPB",
//                 "name": "Scappoose, Scappoose Industrial Airpark",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TPARA",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.7099,
//                     45.871499999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TPARA",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 14.9352,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TPARA",
//                 "name": "I-5 @ Paradise Point",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KMMV",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.13222,
//                     45.196109999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KMMV",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 47.8536,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KMMV",
//                 "name": "McMinnville, McMinnville Municipal Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT15",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.0348,
//                     45.376330000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT15",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 325.83120000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT15",
//                 "name": "Brightwood Weigh Station (US 26 MP 36.5)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/WPKO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.195278,
//                     45.109444000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/WPKO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1325.8800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "WPKO3",
//                 "name": "WANDERER'S PEAK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/LGFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.895278,
//                     45.496667000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/LGFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 853.44000000000005,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "LGFO3",
//                 "name": "LOG CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/BNDW",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.93111,
//                     45.647779999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/BNDW",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 24.384,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "BNDW",
//                 "name": "Bonneville Dam",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT13",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.30139,
//                     45.759599999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT13",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 237.4392,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT13",
//                 "name": "Timber Junction (US 26 MP 37.7)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/HSFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.400806,
//                     44.940806000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/HSFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1036.9296000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "HSFO3",
//                 "name": "HORSE CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/CYFW1",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.202778,
//                     45.929443900000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/CYFW1",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 762,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "CYFW1",
//                 "name": "CANYON CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT78",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.88484,
//                     45.669310000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT78",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 63.703200000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT78",
//                 "name": "Cascade Locks (I-84 MP 44.55)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/SFKO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.483611,
//                     45.595278
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/SFKO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 687.93360000000007,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "SFKO3",
//                 "name": "SOUTH FORK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KSLE",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.995,
//                     44.907780000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KSLE",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 64.00800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KSLE",
//                 "name": "Salem, McNary Field",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/GVT50",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.78274999999999,
//                     45.28857
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/GVT50",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1527.048,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "GVT50",
//                 "name": "Ski Bowl Summit",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT96",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.46193,
//                     45.7971
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT96",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 416.05200000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT96",
//                 "name": "US26 at Sunset Rest Area (US 26 MP 28.6)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT75",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.744,
//                     45.301839999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT75",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1216.152,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT75",
//                 "name": "Government Camp (US-26 MP 54.2)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/DCCW1",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.9875,
//                     45.943610999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/DCCW1",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 822.96000000000004,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "DCCW1",
//                 "name": "DRY CRK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/BDFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.53577799999999,
//                     45.217306000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/BDFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 609.60000000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "BDFO3",
//                 "name": "RYE MOUNTAIN",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT93",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.96925,
//                     46.095939999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT93",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 15.24,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT93",
//                 "name": "US30 at Lewis and Clark Bridge (US 30 MP 48.6",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TIM70",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.71174999999999,
//                     45.345370000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TIM70",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 2130.5520000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TIM70",
//                 "name": "Timberline Magic Mile",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TIM59",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.71133,
//                     45.329970000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TIM59",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1792.2240000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TIM59",
//                 "name": "Timberline Lodge",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/RXFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.921111,
//                     45.027500000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/RXFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 990.60000000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "RXFO3",
//                 "name": "RED BOX",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KKLS",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.90000000000001,
//                     46.116669999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KKLS",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 6.0960000000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KKLS",
//                 "name": "Kelso, Kelso-Longview Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MHM73",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.68163,
//                     45.349269999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MHM73",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 2225.04,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MHM73",
//                 "name": "Mt Hood Meadows Cascade Express",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MHM66",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.67227,
//                     45.34357
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MHM66",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1993.3920000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MHM66",
//                 "name": "Mt Hood Meadows Blue",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MHM54",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.66603000000001,
//                     45.332630000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MHM54",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1639.8240000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MHM54",
//                 "name": "Mt Hood Meadows Base",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MLLO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.27166699999999,
//                     46.022500000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MLLO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 314.24880000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MLLO3",
//                 "name": "MILLER",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/DEFO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.64055999999999,
//                     45.586390000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/DEFO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 350.52000000000004,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "DEFO",
//                 "name": "DEE FLAT",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/PARO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.61667,
//                     45.544440000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/PARO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 438.91200000000003,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "PARO",
//                 "name": "PARKDALE",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/SYNO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.69225,
//                     44.717556000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/SYNO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 226.16160000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "SYNO3",
//                 "name": "JORDAN",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MMRO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.59694399999999,
//                     45.579444000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MMRO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 775.41120000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MMRO3",
//                 "name": "MIDDLE MTN",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/RKHO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.469444,
//                     44.924999999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/RKHO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 547.72559999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "RKHO3",
//                 "name": "ROCKHOUSE 1",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TMKO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.80249999999999,
//                     45.457222000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TMKO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 3.3528000000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TMKO3",
//                 "name": "TILLAMOOK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KTMK",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.81444,
//                     45.418059900000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KTMK",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 10.972800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KTMK",
//                 "name": "Tillamook, Tillamook Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/PNGO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.50917,
//                     45.65222
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/PNGO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 188.976,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "PNGO",
//                 "name": "PINEGROVE",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/HOXO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.51806000000001,
//                     45.684440000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/HOXO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 155.44800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "HOXO",
//                 "name": "HOOD RIVER",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TR950",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.891944,
//                     46.271388999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TR950",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 64.92240000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TR950",
//                 "name": "CASTLE ROCK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/CDFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.771944,
//                     45.211666999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/CDFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 676.65600000000006,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "CDFO3",
//                 "name": "CEDAR",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MTRO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.55627800000001,
//                     46.011471999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MTRO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 620.26800000000003,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MTRO3",
//                 "name": "TIDEWATER",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT57",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.71536,
//                     45.072069900000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT57",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 188.976,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT57",
//                 "name": "Murphy Hill (OR 18 MP 15.4)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/BOFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.003056,
//                     44.721944000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/BOFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1088.136,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "BOFO3",
//                 "name": "BOULDER CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/SPMW1",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.92661,
//                     46.179499999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/SPMW1",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1036.3200000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "SPMW1",
//                 "name": "SPENCER MEADOW",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/MSH33",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.26503,
//                     46.303330000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/MSH33",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 993.64800000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "MSH33",
//                 "name": "Mt St Helens - Coldwater",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/YEFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.427778,
//                     44.592222
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/YEFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 938.78399999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "YEFO3",
//                 "name": "YELLOWSTONE MTN.",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT46",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.43891000000001,
//                     46.170909999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT46",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 197.8152,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT46",
//                 "name": "Bradley Wayside (US 30 MP 74.9)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/CRVO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.19,
//                     44.634169999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/CRVO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 70.103999999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "CRVO",
//                 "name": "CORVALLIS",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ABNW1",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.083333,
//                     46.342778000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ABNW1",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 609.60000000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ABNW1",
//                 "name": "ABERNATHY MTN.",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/AT297",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.96733,
//                     45.196829999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/AT297",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 8.5343999999999998,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "AT297",
//                 "name": "PACCTY-2 Pacific City",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/BKRW1",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.538611,
//                     46.056666999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/BKRW1",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 819.91200000000003,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "BKRW1",
//                 "name": "BUCK CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/AP682",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -124.0061,
//                     45.010809999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/AP682",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 27.127200000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "AP682",
//                 "name": "W7KKE-3 Road's End",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KCVO",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.28333000000001,
//                     44.5
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KCVO",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 74.980800000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KCVO",
//                 "name": "Corvallis, Corvallis Municipal Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KAST",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.88249999999999,
//                     46.156939999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KAST",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 3.048,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KAST",
//                 "name": "Astoria, Astoria Regional Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT50",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.14151,
//                     44.395479999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT50",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1292.3520000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT50",
//                 "name": "Tombstone Summit (US 20 MP 64)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/FNWO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.325278,
//                     44.418332900000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/FNWO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 93.878399999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "FNWO3",
//                 "name": "CORVALLIS",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/BRUO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.849417,
//                     44.284278
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/BRUO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 649.22400000000005,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "BRUO3",
//                 "name": "BRUSH CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/PEFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -121.995,
//                     44.236666999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/PEFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1051.5599999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "PEFO3",
//                 "name": "PEBBLE",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KONP",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -124.05806,
//                     44.580280000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KONP",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 49.072800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KONP",
//                 "name": "Newport, Newport Municipal Airport",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/VCFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.46386099999999,
//                     44.252389000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/VCFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 477.012,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "VCFO3",
//                 "name": "VILLAGE CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/TCFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.57725000000001,
//                     44.111139000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/TCFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 691.28640000000007,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "TCFO3",
//                 "name": "TROUT CREEK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/KEUG",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.21444,
//                     44.133330000000001
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/KEUG",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 110.94720000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "KEUG",
//                 "name": "Eugene, Mahlon Sweet Field",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/CNFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.886667,
//                     44.348889
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/CNFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 591.00720000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "CNFO3",
//                 "name": "CANNIBAL MOUNTAIN",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/AS531",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -124.09211999999999,
//                     44.343029999999999
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/AS531",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 21.945600000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "AS531",
//                 "name": "YACHTS Yachats",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/LOFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.37955599999999,
//                     43.906388999999997
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/LOFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 589.78800000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "LOFO3",
//                 "name": "HIGH POINT",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/CCRO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.805556,
//                     43.725278000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/CCRO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 933.9072000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "CCRO3",
//                 "name": "GREEN MOUNTAIN",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/ODT49",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.14475,
//                     43.75224
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/ODT49",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 219.45600000000002,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "ODT49",
//                 "name": "Wards Butte / Cottage Grove (I-5 MP 170)",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/GPFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -123.890278,
//                     43.928055999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/GPFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 548.63999999999999,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "GPFO3",
//                 "name": "GOODWIN PEAK",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/SGFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.629167,
//                     43.663611000000003
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/SGFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1319.1744000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "SGFO3",
//                 "name": "SUGARLOAF",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/FEFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.30202800000001,
//                     43.680528000000002
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/FEFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1028.0904,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "FEFO3",
//                 "name": "FIELDS",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/DUNO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -124.119722,
//                     43.957777999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/DUNO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 36.576000000000001,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "DUNO3",
//                 "name": "DUNES",
//                 "timeZone": "America/Los_Angeles"
//             }
//         },
//         {
//             "id": "https://api.weather.gov/stations/EMFO3",
//             "type": "Feature",
//             "geometry": {
//                 "type": "Point",
//                 "coordinates": [
//                     -122.230278,
//                     43.483055999999998
//                 ]
//             },
//             "properties": {
//                 "@id": "https://api.weather.gov/stations/EMFO3",
//                 "@type": "wx:ObservationStation",
//                 "elevation": {
//                     "value": 1170.432,
//                     "unitCode": "unit:m"
//                 },
//                 "stationIdentifier": "EMFO3",
//                 "name": "EMIGRANT",
//                 "timeZone": "America/Los_Angeles"
//             }
//         }
//     ],
//     "observationStations": [
//         "https://api.weather.gov/stations/KPDX",
//         "https://api.weather.gov/stations/KVUO",
//         "https://api.weather.gov/stations/KTTD",
//         "https://api.weather.gov/stations/ARAO",
//         "https://api.weather.gov/stations/KHIO",
//         "https://api.weather.gov/stations/KUAO",
//         "https://api.weather.gov/stations/EGKO3",
//         "https://api.weather.gov/stations/RRWO3",
//         "https://api.weather.gov/stations/ODT95",
//         "https://api.weather.gov/stations/FOGO",
//         "https://api.weather.gov/stations/TBELL",
//         "https://api.weather.gov/stations/TR951",
//         "https://api.weather.gov/stations/KSPB",
//         "https://api.weather.gov/stations/TPARA",
//         "https://api.weather.gov/stations/KMMV",
//         "https://api.weather.gov/stations/ODT15",
//         "https://api.weather.gov/stations/WPKO3",
//         "https://api.weather.gov/stations/LGFO3",
//         "https://api.weather.gov/stations/BNDW",
//         "https://api.weather.gov/stations/ODT13",
//         "https://api.weather.gov/stations/HSFO3",
//         "https://api.weather.gov/stations/CYFW1",
//         "https://api.weather.gov/stations/ODT78",
//         "https://api.weather.gov/stations/SFKO3",
//         "https://api.weather.gov/stations/KSLE",
//         "https://api.weather.gov/stations/GVT50",
//         "https://api.weather.gov/stations/ODT96",
//         "https://api.weather.gov/stations/ODT75",
//         "https://api.weather.gov/stations/DCCW1",
//         "https://api.weather.gov/stations/BDFO3",
//         "https://api.weather.gov/stations/ODT93",
//         "https://api.weather.gov/stations/TIM70",
//         "https://api.weather.gov/stations/TIM59",
//         "https://api.weather.gov/stations/RXFO3",
//         "https://api.weather.gov/stations/KKLS",
//         "https://api.weather.gov/stations/MHM73",
//         "https://api.weather.gov/stations/MHM66",
//         "https://api.weather.gov/stations/MHM54",
//         "https://api.weather.gov/stations/MLLO3",
//         "https://api.weather.gov/stations/DEFO",
//         "https://api.weather.gov/stations/PARO",
//         "https://api.weather.gov/stations/SYNO3",
//         "https://api.weather.gov/stations/MMRO3",
//         "https://api.weather.gov/stations/RKHO3",
//         "https://api.weather.gov/stations/TMKO3",
//         "https://api.weather.gov/stations/KTMK",
//         "https://api.weather.gov/stations/PNGO",
//         "https://api.weather.gov/stations/HOXO",
//         "https://api.weather.gov/stations/TR950",
//         "https://api.weather.gov/stations/CDFO3",
//         "https://api.weather.gov/stations/MTRO3",
//         "https://api.weather.gov/stations/ODT57",
//         "https://api.weather.gov/stations/BOFO3",
//         "https://api.weather.gov/stations/SPMW1",
//         "https://api.weather.gov/stations/MSH33",
//         "https://api.weather.gov/stations/YEFO3",
//         "https://api.weather.gov/stations/ODT46",
//         "https://api.weather.gov/stations/CRVO",
//         "https://api.weather.gov/stations/ABNW1",
//         "https://api.weather.gov/stations/AT297",
//         "https://api.weather.gov/stations/BKRW1",
//         "https://api.weather.gov/stations/AP682",
//         "https://api.weather.gov/stations/KCVO",
//         "https://api.weather.gov/stations/KAST",
//         "https://api.weather.gov/stations/ODT50",
//         "https://api.weather.gov/stations/FNWO3",
//         "https://api.weather.gov/stations/BRUO3",
//         "https://api.weather.gov/stations/PEFO3",
//         "https://api.weather.gov/stations/KONP",
//         "https://api.weather.gov/stations/VCFO3",
//         "https://api.weather.gov/stations/TCFO3",
//         "https://api.weather.gov/stations/KEUG",
//         "https://api.weather.gov/stations/CNFO3",
//         "https://api.weather.gov/stations/AS531",
//         "https://api.weather.gov/stations/LOFO3",
//         "https://api.weather.gov/stations/CCRO3",
//         "https://api.weather.gov/stations/ODT49",
//         "https://api.weather.gov/stations/GPFO3",
//         "https://api.weather.gov/stations/SGFO3",
//         "https://api.weather.gov/stations/FEFO3",
//         "https://api.weather.gov/stations/DUNO3",
//         "https://api.weather.gov/stations/EMFO3"
//     ]
// }%                                                                                                                                mike@Darwin-D nws %
