package nws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	currentObservationsURLFmt = "http://w1.weather.gov/xml/current_obs/%s.xml"
)

// Observation holds a single meteorological observation. The fields and XML
// mappings are based on http://www.nws.noaa.gov/view/current_observation.xsd.
// See http://w1.weather.gov/xml/current_obs/ for general information. Numeric
// fields are type `string` so that empty strings may represent missing data in
// the struct. Some fields in the XSD are not unmarshalled because they are
// unecessary for the purposes of this package.
type Observation struct {
	SuggestedPickup       string `xml:"suggested_pickup"`
	SuggestedPickupPeriod string `xml:"suggested_pickup_period"`
	Location              string `xml:"location"`
	StationId             string `xml:"station_id"`
	Latitude              string `xml:"longitude"`
	Longitude             string `xml:"latitude"`
	Elevation             string `xml:"elevation"`
	observationTimeRFC822 string `xml:"observation_time_rfc822"`
	Time                  time.Time
	Weather               string `xml:"weather"`
	TempF                 string `xml:"temp_f"`
	TempC                 string `xml:"temp_c"`
	WaterTempF            string `xml:"water_temp_f"`
	WaterTempC            string `xml:"water_temp_f"`
	RelativeHumidity      string `xml:"relative_humidity"`
	WindDir               string `xml:"wind_dir"`
	WindDegrees           string `xml:"wind_degrees"`
	WindMph               string `xml:"wind_mph"`
	WindGustMph           string `xml:"wind_gust_mph"`
	WindKt                string `xml:"wind_kt"`
	WindGustKt            string `xml:"wind_gust_kt"`
	PressureMb            string `xml:"pressure_mb"`
	PressureIn            string `xml:"pressure_in"`
	PressureTendencyMb    string `xml:"pressure_tendency_mb"`
	PressureTendencyIn    string `xml:"pressure_tendency_in"`
	DewpointF             string `xml:"dewpoint_f"`
	DewpointC             string `xml:"dewpoint_c"`
	HeatIndexF            string `xml:"heat_index_f"`
	HeatIndexC            string `xml:"heat_index_c"`
	WindchillF            string `xml:"windchill_f"`
	WindchillC            string `xml:"windchill_c"`
	VisibilityMi          string `xml:"visibility_mi"`
	WaveHeightM           string `xml:"wave_height_m"`
	WaveHeightFt          string `xml:"wave_height_ft"`
	DominantPeriodSec     string `xml:"dominat_period_sec"`
	AveragePeriodSec      string `xml:"average_period_sec"`
	MeanWaveDir           string `xml:"mean_wave_dir"`
	MeanWaveDegrees       string `xml:"mean_wave_degrees"`
	TideFt                string `xml:"tide_ft"`
	Steepness             string `xml:"steepness"`
	WaterColumnHeight     string `xml:"water_column_height"`
	SurfHeightFt          string `xml:"surf_height_ft"`
	SwellDir              string `xml:"swell_dir"`
	SwellDegrees          string `xml:"swell_degrees"`
	SwellPeriod           string `xml:"swell_period"`
}

// getCurrentObservationXML gets the most recent observation from a stationId
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

// processObservationXML accepts an observation in XML and returns an
// Observation sruct
func processObservationXML(observationXML []byte) (*Observation, error) {
	obs := &Observation{}

	if err := xml.Unmarshal(observationXML, obs); err != nil {
		return nil, err
	}

	var err error
	obs.Time, err = time.Parse(time.RFC822, obs.observationTimeRFC822)
	if err != nil {
		return nil, err
	}

	return obs, nil
}

// getCurrentObservation gets the most recent observation from a stationId and
// returns it
func getCurrentObservation(httpClient *http.Client, stationId string) (*Observation, error) {
	obsXML, err := getCurrentObservationXML(httpClient, stationId)
	if err != nil {
		return nil, err
	}
	return processObservationXML(obsXML)
}
