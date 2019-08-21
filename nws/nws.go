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
	httpClient    *http.Client
	baseURLString string

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
func NewClientFromCoordinates(httpClient *http.Client, lat float64, lon float64) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{},

		// point is rounded to four decimal places because the API requires
		// that requests be made with at most four decimal places. The API will
		// 301 redirect, but using four in the first place eliminates the need
		// for those.
		point: Point{
			Lat: math.Round(lat*10000) / 10000,
			Lon: math.Round(lon*10000) / 10000,
		},
	}

	if err := c.setBaseURLString(baseURLString); err != nil {
		return nil, err
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

// BaseURLString ...
func (c *Client) BaseURLString() string {
	return c.baseURLString
}

// SetBaseURLString sets the base URL used by the client. The base URL is set
// to the default during Client construction. This is in place to facilitate
// testing and so that the user can change the base URL if the NWS does.
func (c *Client) SetBaseURLString(url string) error {
	return c.setBaseURLString(url)
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
	f, err := getSemidailyForcastsForGridpoint(c.httpClient, c.gridpoint)
	return *f, err
}

// HourlyForecast ...
func (c *Client) HourlyForecast() (Forecast, error) {
	// update LastRetrieved
	// set value in c
	f, err := getHourlyForcastsForGridpoint(c.httpClient, c.gridpoint)
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

// setBaseURLString ...
func (c *Client) setBaseURLString(url string) error {
	c.baseURLString = url
	return nil
}

// setGridpointFromPoint ...
func (c *Client) setGridpointFromPoint() error {
	gp, err := getGridpointForPoint(c.httpClient, c.point)
	if err != nil {
		return err
	}
	c.gridpoint = *gp
	return nil
}

// setStationsFromGridpont ...
func (c *Client) setStationsFromGridpont() error {
	stns, err := getStationsForGridpoint(c.httpClient, c.gridpoint)
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
func (c *Client) get(path string, query url.Values) (*http.Response, error) {

	return nil, nil
}
