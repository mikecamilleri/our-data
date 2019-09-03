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
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURLString         = "https://api.weather.gov/"
	defaultThrottleString = "5m"
)

// A ValueUnit represents a value and its unit (e.g. 32 inches).
type ValueUnit struct {
	Value float64
	Unit  string
}

// A Client is used to interact with the NWS API for a specific location on
// Earth.
type Client struct {
	// AlertsThrottle represeents the minimum time that must elapse between
	// updating the active alerts.
	AlertsThrottle time.Duration

	// SemidailyForecastThrottle represeents the minimum time that must elapse
	// between updating the semi-daily forecast.
	SemidailyForecastThrottle time.Duration

	// HourlyForecastThrottle represeents the minimum time that must elapse
	// between updating the hourly forecast.
	HourlyForecastThrottle time.Duration

	// ObservationsThrottle represents the minimum time that must elapse between
	// updating the latest observation for any station.
	ObservationsThrottle time.Duration

	// // TODO: the channels below will be used by Auto* functions.
	// AutoAlertsChan            chan []Alert
	// AutoSemidailyForecastChan chan Forecast
	// AutoHourlyForecastChan    chan Forecast
	// AutoObservationChans      map[string]chan Observation

	httpClient                     *http.Client
	httpUserAgentString            string
	point                          Point
	gridpoint                      Gridpoint
	stations                       []Station
	defaultStationID               string
	alerts                         []Alert
	alertsLastRetrived             time.Time
	semidailyForecast              Forecast
	semidailyForecastLastRetrieved time.Time
	hourlyForecast                 Forecast
	hourlyForecastLastRetrieved    time.Time
	observations                   map[string]struct {
		observation              Observation
		observarionLastRetrieved time.Time
	} // the observations key is a station ID
}

// NewClientFromCoordinates creates a new client given a WGS 84 (EPSG:4326)
// latitude and longitide.
//
// httpUserAgentString can be set to anything. The NWS API uses User-Agent as a
// quasi-auth type thing and for security logging. There is no default becuase
// it should be unique to your application.
//   "A User Agent is required to identify your application. This string can be
//   anything, and the more unique to your application the less likely it will
//   be affected by a security event. If you include contact information
//   (website or email), we can contact you if your string is associated to a
//   security event. This will be replaced with an API key in the future."
//   -- https://www.weather.gov/documentation/services-web-api
func NewClientFromCoordinates(httpClient *http.Client, httpUserAgentString string, lat float64, lon float64) (*Client, error) {
	var err error

	c := &Client{
		httpClient:          &http.Client{},
		httpUserAgentString: httpUserAgentString,

		// point Lat and Lon are rounded to four decimal places because the API
		// requires that requests be made with at most four decimal places. The
		// API will 301 redirect, but using four in the first place eliminates
		// those extra requests.
		point: Point{
			Lat: math.Round(lat*10000) / 10000,
			Lon: math.Round(lon*10000) / 10000,
		},
	}

	if err = c.setGridpointFromPoint(); err != nil {
		return nil, err
	}

	if err = c.setStationsFromGridpont(); err != nil {
		return nil, err
	}

	if err = c.setDefaultStationID(c.stations[0].ID); err != nil {
		return nil, err
	}

	defaultThrottle, err := time.ParseDuration(defaultThrottleString)
	if err != nil {
		return nil, err
	}
	c.ObservationsThrottle = defaultThrottle
	c.AlertsThrottle = defaultThrottle
	c.SemidailyForecastThrottle = defaultThrottle
	c.HourlyForecastThrottle = defaultThrottle

	return c, nil
}

// TODO:
//
// *LastRetrived functions would help the caller know when to request an update
//
// Also create Auto* functions that automatically update and send new data
// on a channel
//     - add channels to client
//

// Point returns the Point for this Client.
func (c *Client) Point() (Point, error) {
	return c.point, nil
}

// Gridpoint returns the Gridpoint for this Client.
func (c *Client) Gridpoint() (Gridpoint, error) {
	return c.gridpoint, nil
}

// Stations returns the list of weather stations for this client.
//
// These appear to be ordered based on proximity to the Point used to retrieve
// them, but this isn't documented.
func (c *Client) Stations() ([]Station, error) {
	return c.stations, nil
}

// DefaultStationID returns the ID of the default weather station for this
// Client
func (c *Client) DefaultStationID() (string, error) {
	return c.defaultStationID, nil
}

// SetDefaultStationID changes the default station ID.
func (c *Client) SetDefaultStationID(id string) error {
	return c.setDefaultStationID(id)
}

// Alerts returns a slice of alerts containing the currently active alerts as of
// the last time they were retrieved.
func (c *Client) Alerts(id string) ([]Alert, error) {
	// update LastRetrieved if there is no error from getActiveAlertsForPoint()
	// in case of error, still return alerts and log the error?
	return c.alerts, nil
}

// SemidailyForecast returns the last retrieved semi-daily forecast.
//
// The NWS tends to refer to the semi-daily forecast as simply "forecast."
func (c *Client) SemidailyForecast() (Forecast, error) {
	// update LastRetrieved
	// set value in c
	f, err := getSemidailyForecastForGridpoint(c.httpClient, c.httpUserAgentString, c.gridpoint)
	return *f, err
}

// HourlyForecast returns the last retrieved hourly forcast.
func (c *Client) HourlyForecast() (Forecast, error) {
	// update LastRetrieved
	// set value in c
	f, err := getHourlyForecastForGridpoint(c.httpClient, c.httpUserAgentString, c.gridpoint)
	return *f, err
}

// LatestObservationForDefaultStation returns the last retrieved observation
// for the default station.
func (c *Client) LatestObservationForDefaultStation() (Observation, error) {
	// TODO: Figure out how to recrod last retrieved time for each.
	// Perhaps these should be within c.stations[i].observation.
	// That would likeley be best.
	// The interface of these private attributes doesn't matter as much.
	return c.observations[c.defaultStationID].observation, nil
}

// LatestObservationForStation returns the last retrieved observation for a
// station.
func (c *Client) LatestObservationForStation(id string) (Observation, error) {
	return c.observations[id].observation, nil
}

// UpdateAlerts updates the active alerts for this Client.
func (c *Client) UpdateAlerts() error {
	return nil
}

// UpdateSemidailyForecast updates the semi-daily forecast for this Client.
func (c *Client) UpdateSemidailyForecast() error {
	return nil
}

// UpdateHourlyForecast updates the hourly forecast for this Client.
func (c *Client) UpdateHourlyForecast() error {
	return nil
}

// UpdateLatestObservationForDefaultStation updates the latest observation for
// the default station.
func (c *Client) UpdateLatestObservationForDefaultStation() error {
	return nil
}

// UpdateLatestOservationForStation updates the latest observation for
// a station..
func (c *Client) UpdateLatestOservationForStation() error {
	return nil
}

// // AutoAlerts will automatically update and emit updated slices of Alerts on
// // the Client's AutoAlertsChan
// //
// // AutoAlerts should be run as a goroutine. AutoAlerts will instatiate the
// // Client's AutoAlertsChan upon execution and will return when it detects that
// // the channel is closed.
// func (c *Client) AutoAlerts() error {
// 	return nil
// }

// // StopAutoAlerts will stop an active AutoAlerts goroutine by closing its
// // channel.
// func (c *Client) StopAutoAlerts() error {
// 	return nil
// }

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

// doAPIRequest both makes a GET request to the specified endpoint and handles
// non-200 responses. get will only return an *http.Rsponse with a 200 status
// code.
func doAPIRequest(httpClient *http.Client, httpUserAgentString string, endpoint string, query url.Values) ([]byte, error) {
	// build the request
	req, err := http.NewRequest("GET", baseURLString+endpoint, nil)
	if err != nil {
		return nil, err
	}
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	req.Header.Set("User-Agent", httpUserAgentString)

	// make the request, return error if error
	// TODO: handle errors like client side timeouts
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// see below for why this is done here instead of after checking status
	// code. Setting this up for TODOs below.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check status code, return error if not 200
	// TODO: handle errors like server side timeouts
	// TODO: do something with the response body if error
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return respBody, nil
}
