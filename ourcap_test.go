package ouralerts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// CAP Alert Message Examples from specification
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2-os.html
	HomelandSecurityAdvisorySystemAlert = `<?xml version = "1.0" encoding = "UTF-8"?>
<alert xmlns = "urn:oasis:names:tc:emergency:cap:1.2">
  <identifier>43b080713727</identifier> 
  <sender>hsas@dhs.gov</sender> 
  <sent>2003-04-02T14:39:01-05:00</sent>
  <status>Actual</status> 
  <msgType>Alert</msgType>
  <scope>Public</scope>  
  <info>
    <category>Security</category>   
    <event>Homeland Security Advisory System Update</event>   
    <urgency>Immediate</urgency>   
    <severity>Severe</severity>   
    <certainty>Likely</certainty>   
    <senderName>U.S. Government, Department of Homeland Security</senderName>
    <headline>Homeland Security Sets Code ORANGE</headline>
    <description>The Department of Homeland Security has elevated the Homeland Security Advisory System threat level to ORANGE / High in response to intelligence which may indicate a heightened threat of terrorism.</description>
    <instruction> A High Condition is declared when there is a high risk of terrorist attacks. In addition to the Protective Measures taken in the previous Threat Conditions, Federal departments and agencies should consider agency-specific Protective Measures in accordance with their existing plans.</instruction> 
    <web>http://www.dhs.gov/dhspublic/display?theme=29</web>
    <parameter>
      <valueName>HSAS</valueName>
      <value>ORANGE</value>
    </parameter>   
    <resource>
      <resourceDesc>Image file (GIF)</resourceDesc>
      <mimeType>image/gif</mimeType>   
      <uri>http://www.dhs.gov/dhspublic/getAdvisoryImage</uri>
    </resource>   
    <area>       
      <areaDesc>U.S. nationwide and interests worldwide</areaDesc>   
    </area>
  </info>
</alert>
`
	SevereThunderstormWarning = `<?xml version = "1.0" encoding = "UTF-8"?>
<alert xmlns = "urn:oasis:names:tc:emergency:cap:1.2">
  <identifier>KSTO1055887203</identifier> 
  <sender>KSTO@NWS.NOAA.GOV</sender> 
  <sent>2003-06-17T14:57:00-07:00</sent>
  <status>Actual</status> 
  <msgType>Alert</msgType>
  <scope>Public</scope> 
  <info>
    <category>Met</category>   
    <event>SEVERE THUNDERSTORM</event>
    <responseType>Shelter</responseType> 
    <urgency>Immediate</urgency>   
    <severity>Severe</severity>   
    <certainty>Observed</certainty>
    <eventCode>
      <valueName>SAME</valueName>
      <value>SVR</value>
    </eventCode>
    <expires>2003-06-17T16:00:00-07:00</expires>  
    <senderName>NATIONAL WEATHER SERVICE SACRAMENTO CA</senderName>
    <headline>SEVERE THUNDERSTORM WARNING</headline>
    <description> AT 254 PM PDT...NATIONAL WEATHER SERVICE DOPPLER RADAR INDICATED A SEVERE THUNDERSTORM OVER SOUTH CENTRAL ALPINE COUNTY...OR ABOUT 18 MILES SOUTHEAST OF KIRKWOOD...MOVING SOUTHWEST AT 5 MPH. HAIL...INTENSE RAIN AND STRONG DAMAGING WINDS ARE LIKELY WITH THIS STORM.</description>
    <instruction>TAKE COVER IN A SUBSTANTIAL SHELTER UNTIL THE STORM PASSES.</instruction>
    <contact>BARUFFALDI/JUSKIE</contact>
    <area>       
      <areaDesc>EXTREME NORTH CENTRAL TUOLUMNE COUNTY IN CALIFORNIA, EXTREME NORTHEASTERN CALAVERAS COUNTY IN CALIFORNIA, SOUTHWESTERN ALPINE COUNTY IN CALIFORNIA</areaDesc>
      <polygon>38.47,-120.14 38.34,-119.95 38.52,-119.74 38.62,-119.89 38.47,-120.14</polygon>
      <geocode>
        <valueName>SAME</valueName>
        <value>006109</value>
      </geocode>
      <geocode>
        <valueName>SAME</valueName>
        <value>006009</value>
      </geocode>
      <geocode>
        <valueName>SAME</valueName>
        <value>006003</value>
      </geocode>
    </area>
  </info>
</alert>
`
	EarthquakeReportUpdateMessage = `<?xml version = "1.0" encoding = "UTF-8"?>
<alert xmlns = "urn:oasis:names:tc:emergency:cap:1.2">
  <identifier>TRI13970876.2</identifier> 
  <sender>trinet@caltech.edu</sender> 
  <sent>2003-06-11T20:56:00-07:00</sent>
  <status>Actual</status> 
  <msgType>Update</msgType>
  <scope>Public</scope>
  <references>trinet@caltech.edu,TRI13970876.1,2003-06-11T20:30:00-07:00</references>
  <info>
    <category>Geo</category>
    <event>Earthquake</event>   
    <urgency>Past</urgency>   
    <severity>Minor</severity>   
    <certainty>Observed</certainty>
    <senderName>Southern California Seismic Network (TriNet) operated by Caltech and USGS</senderName>
    <headline>EQ 3.4 Imperial County CA</headline>
    <description>A minor earthquake measuring 3.4 on the Richter scale occurred near Brawley, California at 8:30 PM Pacific Daylight Time on Wednesday, June 11, 2003. (This event has now been reviewed by a seismologist)</description>
    <web>http://www.trinet.org/scsn/scsn.html</web>
    <parameter>
      <valueName>EventID</valueName>
      <value>13970876</value>
    </parameter>
    <parameter>
      <valueName>Version</valueName>
      <value>1</value>
    </parameter>
    <parameter>
      <valueName>Magnitude</valueName>
      <value>3.4 Ml</value>
    </parameter>
    <parameter>
      <valueName>Depth</valueName>
      <value>11.8 mi.</value>
    </parameter>
    <parameter>
      <valueName>Quality</valueName>
      <value>Excellent</value>
    </parameter>
    <area>       
      <areaDesc>1 mi. WSW of Brawley, CA; 11 mi. N of El Centro, CA; 30 mi. E of OCOTILLO (quarry); 1 mi. N of the Imperial Fault</areaDesc>
      <circle>32.9525,-115.5527 0</circle>  
    </area>
  </info>
</alert>
`
	AmberAlertMultilingualMessage = `<?xml version = "1.0" encoding = "UTF-8"?>
<alert xmlns = "urn:oasis:names:tc:emergency:cap:1.2">
   <identifier>KAR0-0306112239-SW</identifier> 
   <sender>KARO@CLETS.DOJ.CA.GOV</sender>
   <sent>2003-06-11T22:39:00-07:00</sent>
   <status>Actual</status> 
   <msgType>Alert</msgType>
   <source>SW</source>
   <scope>Public</scope>
   <info>
     <language>en-US</language>
     <category>Rescue</category>   
     <event>Child Abduction</event>   
     <urgency>Immediate</urgency>   
     <severity>Severe</severity>   
     <certainty>Likely</certainty>
     <eventCode>
        <valueName>SAME</valueName>
        <value>CAE</value>
     </eventCode>
     <senderName>Los Angeles Police Dept - LAPD</senderName>
     <headline>Amber Alert in Los Angeles County</headline>
     <description>DATE/TIME: 06/11/03, 1915 HRS.  VICTIM(S): KHAYRI DOE JR. M/B BLK/BRO 3'0", 40 LBS. LIGHT COMPLEXION.  DOB 06/24/01. WEARING RED SHORTS, WHITE T-SHIRT, W/BLUE COLLAR.  LOCATION: 5721 DOE ST., LOS ANGELES, CA.  SUSPECT(S): KHAYRI DOE SR. DOB 04/18/71 M/B, BLK HAIR, BRO EYE. VEHICLE: 81' BUICK 2-DR, BLUE (4XXX000).</description>
     <contact>DET. SMITH, 77TH DIV, LOS ANGELES POLICE DEPT-LAPD AT 213 485-2389</contact>
     <area>
        <areaDesc>Los Angeles County</areaDesc>
        <geocode>
           <valueName>SAME</valueName>
           <value>006037</value>
        </geocode>
     </area>
   </info>
   <info>
     <language>es-US</language>
     <category>Rescue</category>   
     <event>Abducción de Niño</event>
     <urgency>Immediate</urgency>   
     <severity>Severe</severity>   
     <certainty>Likely</certainty>
     <eventCode>
        <valueName>SAME</valueName>
        <value>CAE</value>
     </eventCode>
     <senderName>Departamento de Policía de Los Ángeles - LAPD</senderName>
     <headline>Alerta Amber en el condado de Los Ángeles</headline>
     <description>DATE/TIME: 06/11/03, 1915 HORAS. VÍCTIMAS: KHAYRI DOE JR. M/B BLK/BRO 3'0", 40 LIBRAS. TEZ LIGERA. DOB 06/24/01. CORTOCIRCUITOS ROJOS QUE USAN, CAMISETA BLANCA, COLLAR DE W/BLUE. LOCALIZACIÓN: 5721 DOE ST., LOS ÁNGELES. SOSPECHOSO: KHAYRI DOE ST. DOB 04/18/71 M/B, PELO DEL NEGRO, OJO DE BRO. VEHÍCULO: 81' BUICK 2-DR, AZUL (4XXX000)</description>
     <contact>DET. SMITH, 77TH DIV, LOS ANGELES POLICE DEPT-LAPD AT 213 485-2389</contact>
     <area>
        <areaDesc>condado de Los Ángeles</areaDesc>
        <geocode>
           <valueName>SAME</valueName>
           <value>006037</value>
        </geocode>
     </area>
   </info>
</alert>`

	// constants for unit tests
	ReferencesStringValid       = `user@example.com,XX1122333,2017-01-01T10:43:00-08:00 user2@example.com,2XX1122333,2017-01-01T10:43:00-08:00`
	ReferencesStringMissingPart = `user@example.com,2016-01-01T10:43:00-08:00`
	ReferencesStringBadTime     = `user@example.com,XX1122333,2016-01-01T10:43:00`
	ReferencesStringEmpty       = ``

	PolygonStringValid    = `38.47,-120.14 38.52,-119.74 38.62,-119.89 38.47,-120.14`
	PolygonStringShort    = `38.47,-120.14 38.62,-119.89 38.47,-120.14`
	PolygonStringOpen     = `38.47,-120.14 38.34,-119.95 38.52,-119.74 38.62,-119.89`
	PolygonStringBadPoint = `38.47,-120.14 38.52 38.62,-119.89 38.47,-120.14`
	PolygonStringEmpty    = ``

	CircleStringValid    = `32.9525,-115.5527 1`
	CircleStringBadPoint = `-115.5527 1`
	CircleStringNoPoint  = `1`
	CircleStringNoRadius = `32.9525,-115.5527`

	AddressesStringValid = `one@example.com two@example.com`
	AddressesStringEmpty = ``

	IncidentsStringValid = `XXXX1 XXXX2`
	IncidentsStringEmpty = ``

	SpaceDelimitedQuotedStringValid  = `"hello world" live "goodbye world"`
	SpaceDelimitedQuotedStringValid2 = `one two "three ... (3)" four`
	SpaceDelimitedQuotedStringValid3 = `one`
	SpaceDelimitedQuotedStringEmpty  = ``

	TimeStringValid   = `2017-01-01T10:43:00-08:00`
	TimeStringBadZone = `2017-01-01T10:43:00Z`

	URLStringFullValid     = `http://mikcamilleri.com/`
	URLStringRelativeValid = `hello`
	URLStringInvalid       = `http://example.com\`
)

func TestParseReferencesString(t *testing.T) {
	assert := assert.New(t)
	var refs []Reference
	var err error

	refs, err = parseReferencesString(ReferencesStringValid)
	assert.Nil(err)
	assert.Len(refs, 2)
	assert.Equal("user@example.com", refs[0].Sender)
	assert.Equal("XX1122333", refs[0].Identifier)
	tm, _ := time.Parse("2006-01-02T15:04:05-07:00", "2017-01-01T10:43:00-08:00")
	assert.Equal(tm, refs[0].Sent)

	refs, err = parseReferencesString(ReferencesStringMissingPart)
	assert.NotNil(err)
	assert.Nil(refs)

	refs, err = parseReferencesString(ReferencesStringBadTime)
	assert.NotNil(err)
	assert.Nil(refs)

	refs, err = parseReferencesString(ReferencesStringEmpty)
	assert.NotNil(err)
	assert.Nil(refs)
}

func TestIsValidReferencesString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidReferencesString(ReferencesStringValid))
	assert.False(isValidReferencesString(ReferencesStringMissingPart))
	assert.False(isValidReferencesString(ReferencesStringBadTime))
	assert.False(isValidReferencesString(ReferencesStringEmpty))
}

func TestParsePolygonString(t *testing.T) {
	assert := assert.New(t)
	var poly Polygon
	var err error

	poly, err = parsePolygonString(PolygonStringValid)
	assert.Nil(err)
	assert.Len(poly, 4)
	assert.Equal(Point{Latitude: "38.47", Longitude: "-120.14"}, poly[0])

	poly, err = parsePolygonString(PolygonStringShort)
	assert.NotNil(err)
	assert.Len(poly, 0)

	poly, err = parsePolygonString(PolygonStringOpen)
	assert.NotNil(err)
	assert.Len(poly, 0)

	poly, err = parsePolygonString(PolygonStringBadPoint)
	assert.NotNil(err)
	assert.Len(poly, 0)

	poly, err = parsePolygonString(PolygonStringEmpty)
	assert.NotNil(err)
	assert.Len(poly, 0)
}

func TestIsValidPolygonString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidPolygonString(PolygonStringValid))
	assert.False(isValidPolygonString(PolygonStringShort))
	assert.False(isValidPolygonString(PolygonStringOpen))
	assert.False(isValidPolygonString(PolygonStringBadPoint))
	assert.False(isValidPolygonString(PolygonStringEmpty))
}

func TestParseCircleString(t *testing.T) {
	assert := assert.New(t)
	var circle Circle
	var err error

	circle, err = parseCircleString(CircleStringValid)
	assert.Nil(err)
	assert.Equal(Circle{Point: Point{Latitude: "32.9525", Longitude: "-115.5527"}, Radius: "1"}, circle)

	circle, err = parseCircleString(CircleStringBadPoint)
	assert.NotNil(err)
	assert.Equal(Circle{}, circle)

	circle, err = parseCircleString(CircleStringNoPoint)
	assert.NotNil(err)
	assert.Equal(Circle{}, circle)

	circle, err = parseCircleString(CircleStringNoRadius)
	assert.NotNil(err)
	assert.Equal(Circle{}, circle)
}

func TestIsValidCircleString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidCircleString(CircleStringValid))
	assert.False(isValidCircleString(CircleStringBadPoint))
	assert.False(isValidCircleString(CircleStringNoPoint))
	assert.False(isValidCircleString(CircleStringNoRadius))
}

func TestParseAddressesString(t *testing.T) {
	assert := assert.New(t)
	var addrs []string

	addrs = parseAddressesString(AddressesStringValid)
	assert.Equal([]string{"one@example.com", "two@example.com"}, addrs)

	addrs = parseAddressesString(AddressesStringEmpty)
	assert.Len(addrs, 0)
}

func TestIsValidAddressesString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidAddressesString(AddressesStringValid))
	assert.False(isValidAddressesString(AddressesStringEmpty))
}

func TestParseIncidentsString(t *testing.T) {
	assert := assert.New(t)
	var incidents []string

	incidents = parseIncidentsString(IncidentsStringValid)
	assert.Equal([]string{"XXXX1", "XXXX2"}, incidents)

	incidents = parseIncidentsString(IncidentsStringEmpty)
	assert.Len(incidents, 0)
}

func TestIsValidIncidentsString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidIncidentsString(IncidentsStringValid))
	assert.False(isValidIncidentsString(IncidentsStringEmpty))
}

func TestSplitSpaceDelimitedQuotedStrings(t *testing.T) {
	assert := assert.New(t)
	var strs []string

	strs = splitSpaceDelimitedQuotedStrings(SpaceDelimitedQuotedStringValid)
	assert.Equal([]string{"hello world", "live", "goodbye world"}, strs)

	strs = splitSpaceDelimitedQuotedStrings(SpaceDelimitedQuotedStringValid2)
	assert.Equal([]string{"one", "two", "three ... (3)", "four"}, strs)

	strs = splitSpaceDelimitedQuotedStrings(SpaceDelimitedQuotedStringValid3)
	assert.Equal([]string{"one"}, strs)

	strs = splitSpaceDelimitedQuotedStrings(SpaceDelimitedQuotedStringEmpty)
	assert.Len(strs, 0)
}

func TestIsValidTimeString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidTimeString(TimeStringValid))
	assert.False(isValidTimeString(TimeStringBadZone))
}

func TestIsValidURLString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidURLString(URLStringFullValid))
	assert.True(isValidURLString(URLStringRelativeValid))
	assert.False(isValidURLString(URLStringInvalid))
}

func TestProcessAlertMessageXML(t *testing.T) {
	assert := assert.New(t)
	var err error

	_, err = ProcessAlertMsgXML([]byte(HomelandSecurityAdvisorySystemAlert))
	assert.Nil(err)

	_, err = ProcessAlertMsgXML([]byte(SevereThunderstormWarning))
	assert.Nil(err)

	_, err = ProcessAlertMsgXML([]byte(EarthquakeReportUpdateMessage))
	assert.Nil(err)

	_, err = ProcessAlertMsgXML([]byte(AmberAlertMultilingualMessage))
	assert.Nil(err)
}
