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

const getStationsForGridpointEndpointURLStringFmt = "gridpoints/%s/%d,%d/stations" // wfo, lat, lon

// A Station represents a single weather station.
type Station struct {
	ID    string // callsign
	Name  string
	Point Point
}

// getStationsForGridpoint retrieves from the NWS API a list of stations that
// are proximal to a particular gridpoint.
func getStationsForGridpoint(httpClient *http.Client, httpUserAgentString string, gridpoint Gridpoint) ([]Station, error) {
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		fmt.Sprintf(
			getStationsForGridpointEndpointURLStringFmt,
			gridpoint.WFO,
			gridpoint.GridX,
			gridpoint.GridY,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return newStationsFromStationsRespBody(respBody)
}

// newStationsFromStationsRespBody returns a slice of stations, given a response
// body from the NWS API.
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
	if err := json.Unmarshal(respBody, &stnsRaw); err != nil {
		return nil, err
	}

	// validate and build returned slice
	var stns []Station

	for _, sRaw := range stnsRaw.Features {
		if sRaw.Properties.StationIdentifier == "" {
			continue // skip if no callsign
		}
		s := Station{
			ID:   strings.ToUpper(sRaw.Properties.StationIdentifier),
			Name: sRaw.Properties.Name,
		}
		if len(sRaw.Geometry.Coordinates) == 2 {
			s.Point.Lat, _ = strconv.ParseFloat(sRaw.Geometry.Coordinates[1], 64)
			s.Point.Lon, _ = strconv.ParseFloat(sRaw.Geometry.Coordinates[0], 64)
		}
		stns = append(stns, s)
	}

	return stns, nil
}
