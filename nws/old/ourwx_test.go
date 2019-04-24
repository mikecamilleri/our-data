package ourwx

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/mikecamilleri/ourwx/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testLat          = "45.53"
	testLon          = "-122.67"
	testZone         = "ORZ006"
	testZoneNoAlerts = "ORZ006NoAlerts"
	testStation      = "KPDX"
)

func TestNewClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c, err := NewClient(testLat, testLon, testZone, testStation)

	require.Nil(t, err)
	assert.Equal(t, testLat, c.latitude)
	assert.Equal(t, testLon, c.longitude)
	assert.Equal(t, testZone, c.zone)
	assert.Equal(t, testStation, c.station)
}

// TestNewClientFromCoordinates also inherently tests
// setZoneAndStationFromCoordinates
func TestNewClientFromCoordinates(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c, err := NewClientFromCoordinates(testLat, testLon)

	require.Nil(t, err)
	assert.Equal(t, testLat, c.latitude)
	assert.Equal(t, testLon, c.longitude)
	assert.Equal(t, testZone, c.zone)
	assert.Equal(t, testStation, c.station)
}

func TestCurrentObservation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &Client{
		httpClient: &http.Client{},
		station:    testStation,
	}
	o, err := c.CurrentObservation()

	require.Nil(t, err)
	assert.Equal(t, &parsedTestObservation, o)
}

func TestCurrentAlerts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &Client{
		httpClient: &http.Client{},
		zone:       testZone,
	}
	a, err := c.CurrentAlerts()

	assert.Nil(t, err)
	assert.Len(t, a, 1)
}

func TestValidate(t *testing.T) {
	// TODO once method is more complete
}
