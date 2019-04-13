package ouralerts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateMessageXML implicitely tests unmarshallMessageXML() and
// alert.validate().
// TODO: Improve this test
func TestValidateMessageXML(t *testing.T) {
	assert := assert.New(t)
	var err error

	// CAP 1.2 specification examples
	err = ValidateMessageXML([]byte(testHomelandSecurityAdvisorySystemAlert))
	assert.Nil(err)
	err = ValidateMessageXML([]byte(testSevereThunderstormWarning))
	assert.Nil(err)
	err = ValidateMessageXML([]byte(testEarthquakeReportUpdateMessage))
	assert.Nil(err)
	err = ValidateMessageXML([]byte(testAmberAlertMultilingualMessage))
	assert.Nil(err)

	// Actual NWS examples are invalid due to empty polygon
	err = ValidateMessageXML([]byte(testNWSHydrologicOutlook))
	assert.NotNil(err)
	err = ValidateMessageXML([]byte(testNWSWinterWeatherAdvisory))
	assert.NotNil(err)
	err = ValidateMessageXML([]byte(testNWSWinterStormWarning))
	assert.NotNil(err)
	err = ValidateMessageXML([]byte(testNWSAirQualityAlert))
	assert.NotNil(err)
}

func TestIsValidAddressesString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidAddressesString(testAddressesStringValid))
	assert.False(isValidAddressesString(testAddressesStringEmpty))
}

func TestIsValidIncidentsString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidIncidentsString(testIncidentsStringValid))
	assert.False(isValidIncidentsString(testIncidentsStringEmpty))
}

func TestIsValidSpaceDelimitedQuotedStrings(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidSpaceDelimitedQuotedStrings(testIsValidSpaceDelimitedQuotedStringValid))
	assert.False(isValidSpaceDelimitedQuotedStrings(testIsValidSpaceDelimitedQuotedStringOddNumberOfQuotes))
	assert.False(isValidSpaceDelimitedQuotedStrings(testIsValidSpaceDelimitedQuotedStringEmpty))
}

func TestIsValidTimeString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidTimeString(testTimeStringValid))
	assert.False(isValidTimeString(testTimeStringBadZone))
}

func TestIsValidURLString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidURLString(testURLStringFullValid))
	assert.True(isValidURLString(testURLStringRelativeValid))
	assert.False(isValidURLString(testURLStringInvalid))
}

func TestIsValidReferencesString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidReferencesString(testReferencesStringValid))
	assert.False(isValidReferencesString(testReferencesStringMissingPart))
	assert.False(isValidReferencesString(testReferencesStringBadTime))
	assert.False(isValidReferencesString(testReferencesStringEmpty))
}

func TestIsValidPolygonString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidPolygonString(testPolygonStringValid))
	assert.False(isValidPolygonString(testPolygonStringShort))
	assert.False(isValidPolygonString(testPolygonStringOpen))
	assert.False(isValidPolygonString(testPolygonStringBadPoint))
	assert.False(isValidPolygonString(testPolygonStringEmpty))
}

func TestIsValidCircleString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidCircleString(testCircleStringValid))
	assert.False(isValidCircleString(testCircleStringBadPoint))
	assert.False(isValidCircleString(testCircleStringNoPoint))
	assert.False(isValidCircleString(testCircleStringNoRadius))
}
