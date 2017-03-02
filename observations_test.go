package nws

import (
	"net/http"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/mikecamilleri/nws/mock"
	"github.com/stretchr/testify/require"
)

var (
	parsedTestObservationTime, _ = time.Parse(time.RFC1123Z, "Mon, 27 Feb 2017 08:53:00 -0800")
	parsedTestObservation        = Observation{
		// SuggestedPickup:       "15 minutes after the hour",
		// SuggestedPickupPeriod: "60",
		Location:         "Portland, Portland International Airport, OR",
		StationId:        "KPDX",
		Latitude:         "45.59578",
		Longitude:        "-122.60917",
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

// TestGetCurrentObservation also inherently tests newObservationFromXML
func TestGetCurrentObservation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &http.Client{}
	obs, err := getCurrentObservation(c, "KPDX")

	require.Nil(t, err)
	require.Equal(t, *obs, parsedTestObservation)
}
