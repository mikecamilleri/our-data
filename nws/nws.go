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

// Package nws ...
package nws

import (
	"math"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURLString = "https://api.weather.gov/"
)

// ValueUnit ...
type ValueUnit struct {
	Value float64
	Unit  string
}

// Client ...
type Client struct {
	httpClient *http.Client
	// The NWS API uses User-Agent as a quasi-auth type thing and for security
	// logging. It needs to be set in each request. There is no default becuase
	// it should be unique to your application.
	//   "A User Agent is required to identify your application. This string can
	//   be anything, and the more unique to your application the less likely it
	//   will be affected by a security event. If you include contact
	//   information (website or email), we can contact you if your string is
	//   associated to a security event. This will be replaced with an API key
	//   in the future."
	//   -- https://www.weather.gov/documentation/services-web-api
	httpUserAgentString string

	point            Point
	gridpoint        Gridpoint
	stations         []Station
	defaultStationID string

	alerts                         []Alert
	alertsLastRetrived             time.Time
	semidailyForecast              Forecast
	semidailyForecastLastRetrieved time.Time
	hourlyForecast                 Forecast
	hourlyForecastLastRetrieved    time.Time
	observations                   map[string]Observation // key is stationID
}

// NewClientFromCoordinates ...
func NewClientFromCoordinates(httpClient *http.Client, httpUserAgentString string, lat float64, lon float64) (*Client, error) {
	c := &Client{
		httpClient:          &http.Client{},
		httpUserAgentString: httpUserAgentString,

		// point is rounded to four decimal places because the API requires
		// that requests be made with at most four decimal places. The API will
		// 301 redirect, but using four in the first place eliminates the need
		// for those.
		point: Point{
			Lat: math.Round(lat*10000) / 10000,
			Lon: math.Round(lon*10000) / 10000,
		},
	}

	if err := c.setGridpointFromPoint(); err != nil {
		return nil, err
	}

	if err := c.setStationsFromGridpont(); err != nil {
		return nil, err
	}

	if err := c.setDefaultStationID(c.stations[0].ID); err != nil {
		return nil, err
	}

	return c, nil
}

// Point ...
func (c *Client) Point() (Point, error) {
	return c.point, nil
}

// Gridpoint ...
func (c *Client) Gridpoint() (Gridpoint, error) {
	return c.gridpoint, nil
}

// Stations ...
func (c *Client) Stations() ([]Station, error) {
	return c.stations, nil
}

// DefaultStationID ...
func (c *Client) DefaultStationID() (string, error) {
	return c.defaultStationID, nil
}

// SetDefaultStationID is set to the first station in the list of stations
// returned when the Client is constructed.
func (c *Client) SetDefaultStationID(id string) error {
	return c.setDefaultStationID(id)
}

// Alerts ...
func (c *Client) Alerts(id string) ([]Alert, error) {
	// update LastRetrieved if there is no error from getActiveAlertsForPoint()
	// in case of error, still return alerts and log the error?
	return c.alerts, nil
}

// SemidailyForecast ...
func (c *Client) SemidailyForecast() (Forecast, error) {
	// update LastRetrieved
	// set value in c
	f, err := getSemidailyForcastsForGridpoint(c.httpClient, c.httpUserAgentString, c.gridpoint)
	return *f, err
}

// HourlyForecast ...
func (c *Client) HourlyForecast() (Forecast, error) {
	// update LastRetrieved
	// set value in c
	f, err := getHourlyForcastsForGridpoint(c.httpClient, c.httpUserAgentString, c.gridpoint)
	return *f, err
}

// LatestObservationForDefaultStation ...
func (c *Client) LatestObservationForDefaultStation() (Observation, error) {
	// TODO: Figure out how to recrod last retrieved time for each.
	// Perhaps these should be within c.stations[i].observation.
	// That would likeley be best.
	// The interface of these private attributes doesn't matter as much.
	return c.observations[c.defaultStationID], nil
}

// LatestObservationForStation ...
func (c *Client) LatestObservationForStation(id string) (Observation, error) {
	return c.observations[id], nil
}

// setGridpointFromPoint ...
func (c *Client) setGridpointFromPoint() error {
	gp, err := getGridpointForPoint(c.httpClient, c.httpUserAgentString, c.point)
	if err != nil {
		return err
	}
	c.gridpoint = *gp
	return nil
}

// setStationsFromGridpont ...
func (c *Client) setStationsFromGridpont() error {
	stns, err := getStationsForGridpoint(c.httpClient, c.httpUserAgentString, c.gridpoint)
	if err != nil {
		return err
	}
	c.stations = stns
	return nil
}

// setDefaultStationID ...
func (c *Client) setDefaultStationID(id string) error {
	c.defaultStationID = c.stations[0].ID
	return nil
}

// get ...
// It is the caller's responsibility to read and close the http.Response.Body.
func get(httpClient *http.Client, httpUserAgentString string, endpoint string, query url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", baseURLString+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", httpUserAgentString)
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	return httpClient.Do(req)
}
