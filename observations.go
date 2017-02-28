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

// Observation holds a single meteorological observation. Numeric fields are
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

// observation is a private struct use to export a single meteorological
// observation into. The fields and XML mappings are based on
// http://www.nws.noaa.gov/view/current_observation.xsd. This package is built
// with an eye towards home automation applications and usefullness for most
// people. Some available fields are not included becuase they are outside the
// scope of this package.
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

// convert converts an unexported observation struct to an exported
// observation struct
func (o *observation) convert() (*Observation, error) {
	ret := Observation{}
	ret.Location = o.Location
	ret.StationId = o.StationId
	ret.Latitude = o.Latitude
	ret.Longitude = o.Longitude
	ret.Elevation = o.Elevation
	// the name of the XML field implies that the time will be provided in
	// RFC822 format, but it is actually RFC1123Z time as defined by the Go
	// time package. 😞
	time, err := time.Parse(observationTimeFmt, o.ObservationTimeRFC822)
	if err != nil {
		return nil, err
	}
	ret.Time = time
	ret.Weather = o.Weather
	ret.TempF = o.TempF
	ret.TempC = o.TempC
	ret.RelativeHumidity = o.RelativeHumidity
	ret.WindDir = o.WindDir
	ret.WindDegrees = o.WindDegrees
	ret.WindMph = o.WindMph
	ret.WindKt = o.WindKt
	ret.WindGustMph = o.WindGustMph
	ret.WindGustKt = o.WindGustKt
	ret.PressureMb = o.PressureMb
	ret.PressureIn = o.PressureIn
	ret.DewpointF = o.DewpointF
	ret.DewpointC = o.DewpointC
	ret.HeatIndexF = o.HeatIndexF
	ret.HeatIndexC = o.HeatIndexC
	ret.WindchillF = o.WindchillF
	ret.WindchillC = o.WindchillC
	ret.VisibilityMi = o.VisibilityMi
	return &ret, nil
}

// getCurrentObservationXML gets the most recent observation from a stationId
// and returns a byte array containing XML
func getCurrentObservationXML(httpClient *http.Client, stationId string) ([]byte, error) {
	resp, err := httpClient.Get(fmt.Sprintf(currentObservationsURLFmt, stationId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response had status: %s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

// processObservationXML accepts an observation in XML and returns a pointer to
// an Observation struct
func processObservationXML(observationXML []byte) (*Observation, error) {
	obs := &observation{}
	// creating our own decoder is required since the character set isn't UTF-8
	reader := bytes.NewReader(observationXML)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(obs); err != nil {
		return nil, err
	}
	return obs.convert()
}

// getCurrentObservation gets the most recent observation from a stationId and
// returns a pointer to an Observation struct
func getCurrentObservation(httpClient *http.Client, stationId string) (*Observation, error) {
	obsXML, err := getCurrentObservationXML(httpClient, stationId)
	if err != nil {
		return nil, err
	}
	return processObservationXML(obsXML)
}
