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

// Package nws implements a client for interacting with the United States
// National Weather Service API Web Service. The client implements a subset of
// available endpoints. This package is location centric. Each client is
// structured around a single point on earth and is able to retrieve data from
// the National Weather Service relating to that point.
package nws

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultAPIURLString   = "https://api.weather.gov/"
	defaultThrottleString = "5m"
)

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

	httpClient          *http.Client
	httpUserAgentString string
	apiURLString        string
	point               Point
	gridpoint           Gridpoint
	stations            []Station
	defaultStationID    string
	alerts              []Alert
	semidailyForecast   Forecast
	hourlyForecast      Forecast
	observations        map[string]ObsTime // key is a station ID

	alertsLastRetrived             time.Time
	semidailyForecastLastRetrieved time.Time
	hourlyForecastLastRetrieved    time.Time
}

// ObsTime holds an observation and the time that it was last retrieved
type ObsTime struct {
	observation              Observation
	observationLastRetrieved time.Time
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

	if err = c.setAPIURLString(defaultAPIURLString); err != nil {
		return nil, err
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

// SetAPIURLString sets the URL of the NWS API Web Service.
//
// The url must begin with `http` (`https` is inherently acceptable) and end
// with a slash (`/`).
func (c *Client) SetAPIURLString(urlString string) error {
	return c.setAPIURLString(urlString)
}

// Point returns the Point for this Client.
func (c *Client) Point() Point {
	return c.point
}

// Gridpoint returns the Gridpoint for this Client.
func (c *Client) Gridpoint() Gridpoint {
	return c.gridpoint
}

// Stations returns the list of weather stations for this client.
//
// These appear to be ordered based on proximity to the Point used to retrieve
// them, but this isn't documented.
func (c *Client) Stations() []Station {
	return c.stations
}

// DefaultStationID returns the ID of the default weather station for this
// Client
func (c *Client) DefaultStationID() string {
	return c.defaultStationID
}

// SetDefaultStationID changes the default station ID.
func (c *Client) SetDefaultStationID(id string) error {
	return c.setDefaultStationID(id)
}

// Alerts returns a slice of alerts containing the currently active alerts as of
// the last time they were retrieved.
func (c *Client) Alerts(id string) []Alert {
	return c.alerts
}

// SemidailyForecast returns the last retrieved semi-daily forecast.
//
// The NWS tends to refer to the semi-daily forecast as simply "forecast."
func (c *Client) SemidailyForecast() Forecast {
	return c.semidailyForecast
}

// HourlyForecast returns the last retrieved hourly forcast.
func (c *Client) HourlyForecast() Forecast {
	return c.hourlyForecast
}

// LatestObservationForDefaultStation returns the last retrieved observation
// for the default station.
func (c *Client) LatestObservationForDefaultStation() Observation {
	// return empty observation if station does not exist in obeservations map
	return c.observations[c.defaultStationID].observation
}

// LatestObservationForStation returns the last retrieved observation for a
// station.
func (c *Client) LatestObservationForStation(id string) Observation {
	// return empty observation if station does not exist in obeservations map
	return c.observations[id].observation
}

// UpdateAlerts updates the active alerts for this Client.
func (c *Client) UpdateAlerts() error {
	alerts, err := getActiveAlertsForPoint(c.httpClient, c.httpUserAgentString, c.apiURLString, c.point)
	if err != nil {
		return err
	}
	c.alerts = alerts
	c.alertsLastRetrived = time.Now()
	return nil
}

// UpdateSemidailyForecast updates the semi-daily forecast for this Client.
func (c *Client) UpdateSemidailyForecast() error {
	f, err := getSemidailyForecastForGridpoint(c.httpClient, c.httpUserAgentString, c.apiURLString, c.gridpoint)
	if err != nil {
		return err
	}
	c.semidailyForecast = *f
	c.semidailyForecastLastRetrieved = f.TimeRetrieved
	return nil
}

// UpdateHourlyForecast updates the hourly forecast for this Client.
func (c *Client) UpdateHourlyForecast() error {
	f, err := getHourlyForecastForGridpoint(c.httpClient, c.httpUserAgentString, c.apiURLString, c.gridpoint)
	if err != nil {
		return err
	}
	c.hourlyForecast = *f
	c.hourlyForecastLastRetrieved = f.TimeRetrieved
	return nil
}

// UpdateLatestObservationForDefaultStation updates the latest observation for
// the default station.
func (c *Client) UpdateLatestObservationForDefaultStation() error {
	o, err := getLatestObservationForStation(c.httpClient, c.httpUserAgentString, c.apiURLString, c.defaultStationID)
	if err != nil {
		return err
	}
	c.observations[c.defaultStationID] = ObsTime{
		observation:              *o,
		observationLastRetrieved: o.TimeRetrieved,
	}
	return nil
}

// UpdateLatestOservationForStation updates the latest observation for
// a station.
func (c *Client) UpdateLatestOservationForStation(id string) error {
	o, err := getLatestObservationForStation(c.httpClient, c.httpUserAgentString, c.apiURLString, id)
	if err != nil {
		return err
	}
	c.observations[id] = ObsTime{
		observation:              *o,
		observationLastRetrieved: o.TimeRetrieved,
	}
	return nil
}

// AlertsLastRetrieved returns the time that alerts waere last successfuly
// retrieved.
func (c *Client) AlertsLastRetrieved(id string) time.Time {
	return c.alertsLastRetrived
}

// SemidailyForecastLastRetrieved returns the time that the semi-daily forecast
// was last successfuly retrieved.
func (c *Client) SemidailyForecastLastRetrieved() time.Time {
	return c.semidailyForecastLastRetrieved
}

// HourlyForecastLastRetrieved returns the time that hourly forecast was last
// successfuly retrieved.
func (c *Client) HourlyForecastLastRetrieved() time.Time {
	return c.hourlyForecastLastRetrieved
}

// LatestObservationForDefaultStationLastRetrieved returns the time that the
// latesst observation for the default station was last successfuly retrieved.
func (c *Client) LatestObservationForDefaultStationLastRetrieved() time.Time {
	// return zero time if station does not exist in obeservations map
	return c.observations[c.defaultStationID].observationLastRetrieved
}

// LatestObservationForStationLastRetrieved returns the time that the latest
// observations for the specified station was last successfuly retrieved.
func (c *Client) LatestObservationForStationLastRetrieved(id string) time.Time {
	// return zero time if station does not exist in obeservations map
	return c.observations[id].observationLastRetrieved
}

// setAPIURLString sets the URL of the NWS API Web Service.
//
// The url must begin with `http` (`https` is inherently acceptable) and end
// with a slash (`/`).
func (c *Client) setAPIURLString(urlString string) error {
	if !strings.HasPrefix(urlString, "http") {
		return fmt.Errorf("urlString must begin with `http`: %s", urlString)
	}
	if !strings.HasSuffix(urlString, "/") {
		return fmt.Errorf("urlString must end with a slash (`/`): %s", urlString)
	}
	c.apiURLString = urlString
	return nil
}

// setGridpointFromPoint set the Client's gridpoint from its point.
func (c *Client) setGridpointFromPoint() error {
	gp, err := getGridpointForPoint(c.httpClient, c.httpUserAgentString, c.apiURLString, c.point)
	if err != nil {
		return err
	}
	c.gridpoint = *gp
	return nil
}

// setStationsFromGridpont sets the Client's stations from its gridpoint.
func (c *Client) setStationsFromGridpont() error {
	stns, err := getStationsForGridpoint(c.httpClient, c.httpUserAgentString, c.apiURLString, c.gridpoint)
	if err != nil {
		return err
	}
	c.stations = stns
	return nil
}

// setDefaultStationID sets the Client's default station to the first station in
// its stations slice.
func (c *Client) setDefaultStationID(id string) error {
	if len(c.stations) < 1 {
		return errors.New("client has no stations")
	}
	c.defaultStationID = c.stations[0].ID
	return nil
}

// doAPIRequest both makes a GET request to the specified endpoint and handles
// non-200 responses. get will only return an *http.Rsponse with a 200 status
// code.
func doAPIRequest(httpClient *http.Client, httpUserAgentString string, apiURLString string, endpoint string, query url.Values) ([]byte, error) {
	// build the request
	req, err := http.NewRequest("GET", apiURLString+endpoint, nil)
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

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check status code, return error if not 200
	// TODO: handle errors like server side timeouts, this is difficult because
	// the API is so sparsely documented.
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", resp.Status, respBody)
	}

	return respBody, nil
}
