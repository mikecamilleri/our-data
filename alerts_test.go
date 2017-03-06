package ourwx

import (
	"net/http"
	"testing"

	"github.com/mikecamilleri/ourwx/mock"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

// TestGetAlerts also inherently tests getAlert
func TestGetAlerts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &http.Client{}
	alerts, err := getAlerts(c, testZone)

	assert.Nil(t, err)
	assert.Len(t, alerts, 1)
}

// TestGetAlertsNoAlerts is to test our handling of the strange way NWS builds
// a feed that contains no alerts.
func TestGetAlertsNoAlerts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mock.RegisterResponders()

	c := &http.Client{}
	alerts, err := getAlerts(c, testZoneNoAlerts)

	assert.Nil(t, err)
	assert.Len(t, alerts, 0)
}
