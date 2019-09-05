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

// A Gridpoint represents a single NWS gridpoint
type Gridpoint struct {
	WFO   string // weather forecast office
	GridX int
	GridY int
	City  string
	State string
}

// getGridpointForPoint retrieves from the NWS API the gridpoint that contains a
// particular point.
func getGridpointForPoint(httpClient *http.Client, httpUserAgentString string, apiURLString string, point Point) (*Gridpoint, error) {
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		apiURLString,
		fmt.Sprintf(getGridpointForPointEndpointURLStringFmt, point.Lat, point.Lon),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return newGridpointFromPointRespBody(respBody)
}

// newGridpointFromPointResponseBody returns a Gridpoint pointer, given a
// response body from the NWS API.
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
	if err := json.Unmarshal(respBody, &gpRaw); err != nil {
		return nil, err
	}

	// validate and build returned value
	var err error
	var gp Gridpoint

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
