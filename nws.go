/*
	Package nws implements a client to interact with several National Weather
	Service APIs
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
	htmlForecastURLFmt        = "http://forecast.weather.gov/MapClick.php?lat=%s&lon=%s"
	alertByStateURLFmt        = "https://alerts.weather.gov/cap/%s.php?x=0"
	alertByZoneOrCountyURLFmt = "https://alerts.weather.gov/cap/wwaatmget.php?x=%s&y=0"
)

type Client struct {
	httpClient *http.Client
	// latitude and longitude are used for forecasts (not yet implemented) and
	// go get the local zone and station if not specified
	latitude  string
	longitude string
	zone      string // CAP - (ORZ006)
	station   string // Current Conditions - (KPDX)
}

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

func (c *Client) CurrentObservation() (*Observation, error) {
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
