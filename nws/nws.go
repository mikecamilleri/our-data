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
	"time"
)

// Client ...
type Client struct {
	httpClient *http.Client
	point      Point
	gridpoint  Gridpoint

	stations         []Station
	defaultStationID string

	alerts             []Alert
	alertsLastRetrived time.Time

	semidailyForecast              Forecast
	semidailyForecastLastRetrieved time.Time
	hourlyForecast                 Forecast
	hourlyForecastLastRetrieved    time.Time

	// observations maps Observation to station ID (callsigns)
	observations map[string]Observation
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

	if err := c.setGridpointFromPoint(); err != nil {
		return nil, err
	}
	// setStationsFromGridpont() also sets defaultStationID
	if err := c.setStationsFromGridpont(); err != nil {
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
	return "", nil
}

// SetDefaultStationID ...
func (c *Client) SetDefaultStationID() error {
	return nil
}

// SemidailyForecast ...
func (c *Client) SemidailyForecast() (Forecast, error) {
	// update LastRetrieved
	return c.semidailyForecast, nil
}

// Alerts ...
func (c *Client) Alerts(id string) ([]Alert, error) {
	// update LastRetrieved
	return c.alerts, nil
}

// HourlyForecast ...
func (c *Client) HourlyForecast() (Forecast, error) {
	// update LastRetrieved
	return c.hourlyForecast, nil
}

// ObservationForDefaultStation ...
func (c *Client) ObservationForDefaultStation() (Observation, error) {
	return c.observations[c.defaultStationID], nil
}

// ObservationForStation ...
func (c *Client) ObservationForStation(id string) (Observation, error) {
	return c.observations[id], nil
}

func (c *Client) setGridpointFromPoint() error {
	gp, err := getGridpointForPoint(c.httpClient, c.point)
	if err != nil {
		return err
	}
	c.gridpoint = gp
	return nil
}

func (c *Client) setStationsFromGridpont() error {
	stns, err := getStationsForGridpoint(c.httpClient, c.gridpoint)
	if err != nil {
		return err
	}
	// setDefaultStationID() here to stns[0].ID
	c.stations = stns
	return nil
}

func (c *Client) setDefaultStationID() error {
	return nil
}
