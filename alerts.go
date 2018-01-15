package ourwx

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mikecamilleri/ouralerts"
)

const (
	alertByZoneURLFmt = "https://alerts.weather.gov/cap/wwaatmget.php?x=%s&y=0"
)

func getAlerts(httpClient *http.Client, zone string) ([]*ouralerts.Alert, error) {
	// get the feed
	resp, err := httpClient.Get(fmt.Sprintf(alertByZoneURLFmt, zone))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response had status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	f := &feed{}
	if err := xml.Unmarshal(body, f); err != nil {
		return nil, err
	}

	// get the alerts
	// TODO: store the alerts and don't retrieve them multiple times.
	alerts := []*ouralerts.Alert{}
	for _, e := range f.Entries {
		// this is a quick and dirty way to detect entries indicating that there
		// are no active alerts.
		if strings.HasPrefix(e.Link, fmt.Sprintf("https://alerts.weather.gov/cap/wwaatmget.php?x=%s", zone)) {
			break
		}
		a, err := getAlert(httpClient, e.Link)
		if err != nil {
			// TODO: implement logging and log the error
			continue
		}
		alerts = append(alerts, a)
	}

	return alerts, nil
}

// getAlert gets a single alert given a URL.
func getAlert(httpClient *http.Client, url string) (*ouralerts.Alert, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response had status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// valiate alert message. How to convey invalidity and still return message?

	return ouralerts.ProcessMessageXML(body)
}

// feed is a private struct used for unmarshalling the Atom feed containing
// alert entries
type feed struct {
	Entries []struct {
		// because links are unique and don't change we can also use them an an
		// ID; in fact, this is what the NWS does in the feed.
		Link string `xml:"id"`
	} `xml:"entry"`
}
