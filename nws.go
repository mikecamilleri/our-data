/*
Package nws implements a client to interact with several United Sates
National Weather Service (NWS) APIs.

This package is built with an eye towards home automation and similar
applications of interest to the general public. This is reflected in the
available APIs and subset of data retrieved from each.
*/
package nws

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	htmlForecastURLFmt = "http://forecast.weather.gov/MapClick.php?lat=%s&lon=%s"
	// alertByZoneURLFmt = "https://alerts.weather.gov/cap/wwaatmget.php?x=%s&y=0"
)

// Client is a client that is used to retrieve data from several NWS APIs. Each
// client is built for a specific location on the earth.
type Client struct {
	httpClient *http.Client
	// coordinates are to be WGS 84 values
	// (https://en.wikipedia.org/wiki/World_Geodetic_System)
	latitude  string
	longitude string
	// zone is used to get alerts (https://alerts.weather.gov)
	zone string
	// station is used to get current observations
	// (http://w1.weather.gov/xml/current_obs/)
	station string
}

// NewClient returns a new Client given latitude, longitude, zone, and station.
func NewClient(latitude, longitude, zone, station string) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{},
		latitude:   latitude,
		longitude:  longitude,
		zone:       zone,
		station:    station,
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// NewClientFromCoordinates returns a new Client given a latitude and longitude.
func NewClientFromCoordinates(latitude, longitude string) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{},
		latitude:   latitude,
		longitude:  longitude,
	}
	if err := c.setZoneAndStationFromCoordinates(); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// CurrentObservation returns the current conditions at the Client's location.
func (c *Client) CurrentObservation() (*Observation, error) {
	// TODO: store the most recent observation in the Client and avoid making
	// API requests too ofton.
	return getCurrentObservation(c.httpClient, c.station)
}

func (c *Client) setZoneAndStationFromCoordinates() error {
	// make the HTTP request and read resp.Body
	resp, err := c.httpClient.Get(fmt.Sprintf(htmlForecastURLFmt, c.latitude, c.longitude))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("http response had status: %s", resp.Status)
	}
	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// regex for the station
	stationRegexp, err := regexp.Compile(` \(([[:upper:]]{4})\)</h2>`)
	if err != nil {
		return err
	}
	stationMatch := stationRegexp.FindSubmatch(htmlBytes)
	if len(stationMatch) < 2 {
		return errors.New("error getting station from HTML forecast")
	}
	stationBytes := stationMatch[1]
	if len(stationBytes) == 0 {
		return errors.New("error getting station from HTML forecast")
	}

	// regex for the zone
	zoneRegexp, err := regexp.Compile(`<p class="myforecast-location"><a href="MapClick.php\?zoneid=([[:upper:]]{3}[[:digit:]]{3})">`)
	if err != nil {
		return err
	}
	zoneMatch := zoneRegexp.FindSubmatch(htmlBytes)
	if len(zoneMatch) < 2 {
		return errors.New("error getting zone from HTML forecast")
	}
	zoneBytes := zoneMatch[1]
	if len(zoneBytes) == 0 {
		return errors.New("error getting zone from HTML forecast")
	}

	c.station = string(stationBytes)
	c.zone = string(zoneBytes)
	return nil
}

func (c *Client) validate() error {
	// verify that c.station is valid
	if _, err := c.CurrentObservation(); err != nil {
		return fmt.Errorf("error validating client: %v", err)
	}

	// verify the c.latitude and c.longitude are valid
	// TODO: use NDFD to do this when implemented

	// verify that c.zone is valid
	// TODO: use CAP alerts to do this when implemented

	return nil
}
