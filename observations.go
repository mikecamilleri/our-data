package nws

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	observationTimeFmt        = time.RFC1123Z
	currentObservationsURLFmt = "http://w1.weather.gov/xml/current_obs/%s.xml"
)

// An Observation holds a single meteorological observation. Numeric fields are
// type `string` so that empty strings may represent missing data in the struct
// and so that decimal values remain exact.
type Observation struct {
	Location         string `xml:"location"`
	StationId        string `xml:"station_id"`
	Latitude         string `xml:"latitude"`
	Longitude        string `xml:"longitude"`
	Elevation        string `xml:"elevation"`
	Time             time.Time
	Weather          string `xml:"weather"`
	TempF            string `xml:"temp_f"`
	TempC            string `xml:"temp_c"`
	RelativeHumidity string `xml:"relative_humidity"`
	WindDir          string `xml:"wind_dir"`
	WindDegrees      string `xml:"wind_degrees"`
	WindMph          string `xml:"wind_mph"`
	WindKt           string `xml:"wind_kt"`
	WindGustMph      string `xml:"wind_gust_mph"`
	WindGustKt       string `xml:"wind_gust_kt"`
	PressureMb       string `xml:"pressure_mb"`
	PressureIn       string `xml:"pressure_in"`
	DewpointF        string `xml:"dewpoint_f"`
	DewpointC        string `xml:"dewpoint_c"`
	HeatIndexF       string `xml:"heat_index_f"`
	HeatIndexC       string `xml:"heat_index_c"`
	WindchillF       string `xml:"windchill_f"`
	WindchillC       string `xml:"windchill_c"`
	VisibilityMi     string `xml:"visibility_mi"`
}

// observation is a private struct used to parse a single meteorological
// observation into. The fields and XML mappings are based on
// http://www.nws.noaa.gov/view/current_observation.xsd. Some available fields
// are not included becuase they are outside the scope of this package.
type observation struct {
	// SuggestedPickup       string `xml:"suggested_pickup"`
	// SuggestedPickupPeriod string `xml:"suggested_pickup_period"`
	Location              string `xml:"location"`
	StationId             string `xml:"station_id"`
	Latitude              string `xml:"latitude"`
	Longitude             string `xml:"longitude"`
	Elevation             string `xml:"elevation"`
	ObservationTimeRFC822 string `xml:"observation_time_rfc822"`
	Weather               string `xml:"weather"`
	TempF                 string `xml:"temp_f"`
	TempC                 string `xml:"temp_c"`
	RelativeHumidity      string `xml:"relative_humidity"`
	WindDir               string `xml:"wind_dir"`
	WindDegrees           string `xml:"wind_degrees"`
	WindMph               string `xml:"wind_mph"`
	WindKt                string `xml:"wind_kt"`
	WindGustMph           string `xml:"wind_gust_mph"`
	WindGustKt            string `xml:"wind_gust_kt"`
	PressureMb            string `xml:"pressure_mb"`
	PressureIn            string `xml:"pressure_in"`
	DewpointF             string `xml:"dewpoint_f"`
	DewpointC             string `xml:"dewpoint_c"`
	HeatIndexF            string `xml:"heat_index_f"`
	HeatIndexC            string `xml:"heat_index_c"`
	WindchillF            string `xml:"windchill_f"`
	WindchillC            string `xml:"windchill_c"`
	VisibilityMi          string `xml:"visibility_mi"`
}

// newObservationFromXML builds a new Observation from raw XML and returns it.
func newObservationFromXML(xmlBytes []byte) (*Observation, error) {
	// decode the byte array into an observation struct
	// creating our own decoder is required since the character set isn't UTF-8
	raw := &observation{}
	reader := bytes.NewReader(xmlBytes)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	// build the Observation struct
	ret := &Observation{}
	ret.Location = raw.Location
	ret.StationId = raw.StationId
	ret.Latitude = raw.Latitude
	ret.Longitude = raw.Longitude
	ret.Elevation = raw.Elevation
	// the name of the XML field implies that the time will be provided in
	// RFC822 format, but it is actually RFC1123Z time as defined by the Go
	// time package. 😞
	time, err := time.Parse(observationTimeFmt, raw.ObservationTimeRFC822)
	if err != nil {
		return nil, err
	}
	ret.Time = time
	ret.Weather = raw.Weather
	ret.TempF = raw.TempF
	ret.TempC = raw.TempC
	ret.RelativeHumidity = raw.RelativeHumidity
	ret.WindDir = raw.WindDir
	ret.WindDegrees = raw.WindDegrees
	ret.WindMph = raw.WindMph
	ret.WindKt = raw.WindKt
	ret.WindGustMph = raw.WindGustMph
	ret.WindGustKt = raw.WindGustKt
	ret.PressureMb = raw.PressureMb
	ret.PressureIn = raw.PressureIn
	ret.DewpointF = raw.DewpointF
	ret.DewpointC = raw.DewpointC
	ret.HeatIndexF = raw.HeatIndexF
	ret.HeatIndexC = raw.HeatIndexC
	ret.WindchillF = raw.WindchillF
	ret.WindchillC = raw.WindchillC
	ret.VisibilityMi = raw.VisibilityMi

	return ret, nil
}

// getCurrentObservation gets the most recent observation from a station and
// returns a pointer to an Observation struct
func getCurrentObservation(httpClient *http.Client, station string) (*Observation, error) {
	// make the HTTP request and read resp.Body
	resp, err := httpClient.Get(fmt.Sprintf(currentObservationsURLFmt, station))
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

	return newObservationFromXML(body)
}
