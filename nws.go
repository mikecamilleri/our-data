/*
	Package nws implements a client to interact with several National Weather
	Service APIs
*/
package nws

import "net/http"

const (
	alertByStateURLFmt        = "https://alerts.weather.gov/cap/%s.php?x=0"
	alertByZoneOrCountyURLFmt = "https://alerts.weather.gov/cap/wwaatmget.php?x=%s&y=0"
)

type Client struct {
	// latitude  string
	// longitude string
	// stateCode       string // CAP - two letter state abbreviation (or)
	// zone            string // CAP - (ORZ006)
	// county          string // CAP - (ORC051)

	// lastAlertTime time.Time

	stationId  string // Current Conditions - (KPDX)
	httpClient *http.Client
}

func NewClient(stationId string) *Client {
	c := &Client{
		stationId:  stationId,
		httpClient: &http.Client{},
	}

	return c
}

func (c *Client) CurrentObservation() (*Observation, error) {
	return getCurrentObservation(c.httpClient, c.stationId)
}

// func (c *Client) Alerts() ([]cap.Alert, error) {
// 	return []cap.Alert{}, nil
// }
