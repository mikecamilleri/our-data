package nws

import (
	"net/http"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/stretchr/testify/require"
)

const (
	testObservation = `<?xml version="1.0" encoding="ISO-8859-1"?> 
<?xml-stylesheet href="latest_ob.xsl" type="text/xsl"?>
<current_observation version="1.0"
	 xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	 xsi:noNamespaceSchemaLocation="http://www.weather.gov/view/current_observation.xsd">
	<credit>NOAA's National Weather Service</credit>
	<credit_URL>http://weather.gov/</credit_URL>
	<image>
		<url>http://weather.gov/images/xml_logo.gif</url>
		<title>NOAA's National Weather Service</title>
		<link>http://weather.gov</link>
	</image>
	<suggested_pickup>15 minutes after the hour</suggested_pickup>
	<suggested_pickup_period>60</suggested_pickup_period>
	<location>Fakeland, Fake International Airport, FK</location>
	<station_id>KFAK</station_id>
	<latitude>43.450101</latitude>
	<longitude>-87.222019</longitude>
	<elevation>0</elevation>
	<observation_time>Last Updated on Feb 27 2017, 8:53 am PST</observation_time>
    <observation_time_rfc822>Mon, 27 Feb 2017 08:53:00 -0800</observation_time_rfc822>
	<weather>Overcast</weather>
	<temperature_string>38.0 F (3.3 C)</temperature_string>
	<temp_f>38.0</temp_f>
	<temp_c>3.3</temp_c>
	<relative_humidity>86</relative_humidity>
	<wind_string>Southwest at 6.9 MPH (6 KT)</wind_string>
	<wind_dir>Southwest</wind_dir>
	<wind_degrees>230</wind_degrees>
	<wind_mph>6.9</wind_mph>
	<wind_kt>6</wind_kt>
	<wind_gust_mph>200</wind_gust_mph>
	<wind_gust_kt>173.8</wind_gust_kt>
	<pressure_string>1009.9 mb</pressure_string>
	<pressure_mb>1009.9</pressure_mb>
	<pressure_in>29.82</pressure_in>
	<dewpoint_string>34.0 F (1.1 C)</dewpoint_string>
	<dewpoint_f>34.0</dewpoint_f>
	<dewpoint_c>1.1</dewpoint_c>
	<windchill_string>33 F (1 C)</windchill_string>
    <windchill_f>33</windchill_f>
    <windchill_c>1</windchill_c>
    <heat_index_string>38.0 F (3.3 C)</heat_index_string>
    <heat_index_f>38.0</heat_index_f>
    <heat_index_c>3.3</heat_index_c>
	<visibility_mi>10.00</visibility_mi>
 	<icon_url_base>http://forecast.weather.gov/images/wtf/small/</icon_url_base>
	<two_day_history_url>http://www.weather.gov/data/obhistory/KPDX.html</two_day_history_url>
	<icon_url_name>ovc.png</icon_url_name>
	<ob_url>http://www.weather.gov/data/METAR/KPDX.1.txt</ob_url>
	<disclaimer_url>http://weather.gov/disclaimer.html</disclaimer_url>
	<copyright_url>http://weather.gov/disclaimer.html</copyright_url>
	<privacy_policy_url>http://weather.gov/notice.html</privacy_policy_url>
</current_observation>`
)

var (
	parsedTestObservationTime, _ = time.Parse(time.RFC1123Z, "Mon, 27 Feb 2017 08:53:00 -0800")
	parsedTestObservation        = Observation{
		// SuggestedPickup:       "15 minutes after the hour",
		// SuggestedPickupPeriod: "60",
		Location:         "Fakeland, Fake International Airport, FK",
		StationId:        "KFAK",
		Latitude:         "43.450101",
		Longitude:        "-87.222019",
		Elevation:        "0",
		Time:             parsedTestObservationTime,
		Weather:          "Overcast",
		TempF:            "38.0",
		TempC:            "3.3",
		RelativeHumidity: "86",
		WindDir:          "Southwest",
		WindDegrees:      "230",
		WindMph:          "6.9",
		WindKt:           "6",
		WindGustMph:      "200",
		WindGustKt:       "173.8",
		PressureMb:       "1009.9",
		PressureIn:       "29.82",
		DewpointF:        "34.0",
		DewpointC:        "1.1",
		HeatIndexF:       "38.0",
		HeatIndexC:       "3.3",
		WindchillF:       "33",
		WindchillC:       "1",
		VisibilityMi:     "10.00",
	}
)

func TestGetCurrentObservationXML(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET",
		"http://w1.weather.gov/xml/current_obs/KFAK.xml",
		httpmock.NewStringResponder(200, testObservation))

	c := &http.Client{}
	obs, err := getCurrentObservationXML(c, "KFAK")
	require.Nil(t, err)
	require.Equal(t, obs, []byte(testObservation))
}

func TestProcessObservationXML(t *testing.T) {
	obs, err := processObservationXML([]byte(testObservation))
	require.Nil(t, err)
	require.Equal(t, *obs, parsedTestObservation)
}

func TestGetCurrentObservation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET",
		"http://w1.weather.gov/xml/current_obs/KFAK.xml",
		httpmock.NewStringResponder(200, testObservation))

	c := &http.Client{}
	obs, err := getCurrentObservation(c, "KFAK")
	require.Nil(t, err)
	require.Equal(t, *obs, parsedTestObservation)
}
