package nws

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/mikecamilleri/nws/mock"
	"github.com/stretchr/testify/assert"
)

func TestSetZoneAndStationFromCoordinates(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &Client{
		httpClient: &http.Client{},
		latitude:   "45.53",
		longitude:  "-122.67",
	}
	err := c.setZoneAndStationFromCoordinates()

	assert.Nil(t, err)
	assert.Equal(t, "KPDX", c.station)
	assert.Equal(t, "ORZ006", c.zone)
}
