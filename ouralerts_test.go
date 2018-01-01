package ouralerts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// CAP alert message examples from specification
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2-os.html
	testHomelandSecurityAdvisorySystemAlert = `<?xml version = "1.0" encoding = "UTF-8"?>
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
	testSevereThunderstormWarning = `<?xml version = "1.0" encoding = "UTF-8"?>
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
	testEarthquakeReportUpdateMessage = `<?xml version = "1.0" encoding = "UTF-8"?>
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
	testAmberAlertMultilingualMessage = `<?xml version = "1.0" encoding = "UTF-8"?>
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

	// actual CAP alert messages from the National Weather Service. The comment
	// near the top of each message is incorrect. These are the actual alert
	// messages and not part of an Atom feed.
	testNWSHydrologicOutlook = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>
<?xml-stylesheet href='https://alerts.weather.gov/cap/capatomproduct.xsl' type='text/xsl'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<alert xmlns = 'urn:oasis:names:tc:emergency:cap:1.1'>

<!-- http-date = Fri, 03 Mar 2017 06:00:00 GMT -->
<identifier>NOAA-NWS-ALERTS-AR125843BFB820.HydrologicOutlook.125843FE95E0AR.JANESFJAN.1758a41ffe72dccbcd11038a67ab6cd2</identifier>
<sender>w-nws.webmaster@noaa.gov</sender>
<sent>2017-03-03T00:00:00-06:00</sent>
<status>Actual</status>
<msgType>Alert</msgType>
<scope>Public</scope>
<note>Alert for Ashley; Chicot (Arkansas) Issued by the National Weather Service</note>
<info>
<category>Met</category>
<event>Hydrologic Outlook</event>
<urgency>Future</urgency>
<severity>Unknown</severity>
<certainty>Possible</certainty>
<eventCode>
<valueName>SAME</valueName>
<value></value>
</eventCode>
<effective>2017-03-03T00:00:00-06:00</effective>
<expires>2017-03-07T12:00:00-06:00</expires>
<senderName>NWS Jackson (Central Mississippi)</senderName>
<headline>Hydrologic Outlook issued March 03 at 12:00AM CST until March 07 at 12:00PM CST by NWS Jackson</headline>
<description>...SPRING FLOOD POTENTIAL OUTLOOK...
...FLOOD RISK IS BELOW AVERAGE ACROSS THE TOMBIGBEE RIVER SYSTEM...
...FLOOD RISE IS AVERAGE ACROSS THE REMAINDER OF THE FORECAST AREA...
This outlook considers rainfall which has already fallen,
snowpack,soil moisture, streamflow, and the 90 day rainfall and
temperature outlook. The primary factor in the development of
significant river flooding across the WFO Jackson Forecast Area is
the occurrence of excessive rainfall in a relatively short period of
time.
SYNOPSIS...
Over the past several months, below normal precipitation
has occurred over the lower Missouri and middle Mississippi Valleys.
Warmer temperatures have kept areas in the lower Mississippi River
Valley from receiving significant snow this season. Snow depths of 2
to 4 inches are confined to portions of Minnesota, Wisconsin, and
north Iowa. Snow water equivalents are generally 0.5 inches or less.
The remainder of the area is snow free. Soil moisture conditions are
generally below normal over the lower Missouri and middle
Mississippi Valleys and near normal over lower Ohio Valley.
Across the ARKLAMISS Region, temperatures have been well above
normal for the first two months of the year. Vegetation is beginning
to grow earlier than normal. Higher evapotranspiration rates than
normal are already occurring. Rainfall since the first of the year
is running at or below normal with only a few isolated areas having
above normal rainfall. Soil moisture is at or below normal across the
region.
MISSISSIPPI RIVER FROM ARKANSAS CITY TO NATCHEZ...
The flood season has been uneventful on the Ohio and Mississippi
Rivers. Streamflows have been near to below normal and no flooding
has occurred this season. The current forecast shows no flooding
over the next couple of weeks but higher flows will occur later
into March.
See the chart below for specific locations showing percent of normal
streamflows:
3/1
Mississippi River             Thebes IL     111%
Ohio River                     Cairo IL      67%
Mississippi River            Memphis TN      36%
Mississippi River      Arkansas City AR      49%
Mississippi River          Vicksburg MS      66%
Mississippi River            Natchez MS      70%
Mississippi River  Red River Landing LA      75%
Mississippi River        Baton Rouge LA      74%
Mississippi River        New Orleans LA      74%
Based on existing soil moisture, streamflow conditions, and normal
spring rainfall patterns; an Average Flood Potential is expected
along the lower Mississippi and lower Ohio Rivers. The magnitude of
future crests will depend on the amount and extent of any upstream
accumulation of snow cover and resultant snowmelt; coupled with the
frequency, intensity, and extent of spring rains.
OUACHITA/BLACK BASINS OF SOUTHEAST ARKANSAS AND NORTHEAST
LOUISIANA...
Streamflows are running near and below seasonal averages. Soil
moisture content is near normal and no flooding is occurring or
expected at this time.
Observed daily streamflows as a percent of mean are given below:
3/1
Bayou Bartholomew             Portland AR     52%
Bayou Bartholomew                Jones LA     27%
Tensas                          Tendal LA     69%
Bayou Macon                     Eudora AR     28%
Ouachita River                  Monroe LA     82%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns; an Average Flood Potential is expected over the
Ouachita and Black River Basins.
BIG BLACK AND HOMOCHITTO RIVER BASINS...
Soil moisture and streamflows have been seasonal to below seasonal
averages. No flooding is occurring or expected over the next several
days.
Observed daily streamflows as a percent of mean are given below:
3/1
Big Black River                   West MS     64%
Big Black River                 Bovina MS     24%
Homochitto River               Rosetta MS     15%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns; an Average Flood Potential is expected over the
Big Black River Basin.
YAZOO BASIN...
Streamflows are running below seasonal averages. Soil moisture
content is near seasonal averages and no flooding is expected during
the next several days.
Observed Daily Streamflows as a percent of normal:
3/1
Tallahatchie                     Money MS      27%
Big Sunflower                Sunflower MS      30%
Percent of available flood control storage is given below.
3/1
Arkabutla Res. MS     86%
Sardis Res. MS     85%
Enid Res. MS     86%
Grenada Res. MS     91%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns; an Average Flood Potential is expected over the
Yazoo River Basin.
PEARL RIVER BASIN...
Heavy rainfall near the end of January produced some minor flooding
along tributaries in the Upper Pearl and along the mainstem of the
lower Pearl River. For the last month, soil moisture content and
streamflows have been normal to below seasonal averages. No flooding
is occurring or expected over the next several days.
Observed daily streamflows as a percent of mean are given below:
3/1
Pearl River                   Carthage MS      26%
Pearl River                    Jackson MS      38%
Pearl River                 Monticello MS      34%
Pearl River                   Columbia MS      18%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns; an Average Flood Potential is expected over the
Pearl River Basin.
PASCAGOULA RIVER BASIN INCLUDING THE LEAF AND CHICKASAWHAY SUB-
BASINS...
Seasonal rainfall in January produced minor to moderate rises. No
flooding is occurring or expected over the next several days.
Soil moisture content is near seasonal levels while streamflows are
running below normal.
Observed daily streamflows as a percent of mean are given below:
3/1
Leaf River                    Collins MS       50%
Leaf River                Hattiesburg MS       38%
Tallahala Creek                Laurel MS       32%
Chickasawhay River         Enterprise MS       39%
Black Creek                  Brooklyn MS       49%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns; an Average Flood Potential is expected over the
Pascagoula River Basin.
TOMBIGBEE RIVER IN MISSISSIPPI...
Soil moisture content and streamflows have been normal seasonal
averages. No flooding is occurring or expected over the next several
days.
Observed daily streamflows as a percent of mean are given below:
3/1
Tombigbee River                Bigbee MS       77%
Buttahatchee River           Aberdeen MS       62%
Luxapallila Creek            Columbus MS       21%
Noxubee River                   Macon MS       18%
Based on existing soil moisture, streamflows, and normal spring
rainfall patterns;an below average flood potential is expected
across the Tombigbee Basin.
EXTENDED TEMPERATURE AND PRECIPITATION OUTLOOK...
The 8-14 Day Outlook issued by the NWS Climate Prediction Center
indicates chances of above normal temperatures. Southeast of the
Natchez Trace, there are equal chances of above normal, normal, and
below normal rainfall. Northwest of the Trace, there is chance for
above normal rainfall. Across the Lower Mississippi Valley, most of
the area can expect above normal rainfall.
The 30 Day Outlook indicates chances of above normal temperatures
across the area. Equal chances of above/below normal precipitation
is indicated over the lower Mississippi Valley.
The 90 Day Outlook issued by the NWS Climate Prediction Center
indicates chances of above normal temperatures over the lower
Mississippi Valley. Equal chances of above/below normal
precipitation is indicated over the lower Mississippi Valley.
This will be the last scheduled spring flood outlook for 2017.</description>
<instruction></instruction>
<parameter>
<valueName>WMOHEADER</valueName>
<value></value>
</parameter>
<parameter>
<valueName>UGC</valueName>
<value>ARZ074-075-LAZ007>009-015-016-023>026-MSZ018-019-025>066-072>074</value>
</parameter>
<parameter>
<valueName>VTEC</valueName>
<value></value>
</parameter>
<parameter>
<valueName>TIME...MOT...LOC</valueName>
<value></value>
</parameter>
<area>
<areaDesc>Ashley; Chicot</areaDesc>
<polygon></polygon>
<geocode>
<valueName>FIPS6</valueName>
<value>005003</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>005017</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>ARZ074</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>ARZ075</value>
</geocode>
</area>
</info>
</alert>`

	testNWSWinterWeatherAdvisory = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>
<?xml-stylesheet href='https://alerts.weather.gov/cap/capatomproduct.xsl' type='text/xsl'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<alert xmlns = 'urn:oasis:names:tc:emergency:cap:1.1'>

<!-- http-date = Fri, 03 Mar 2017 02:17:00 GMT -->
<identifier>NOAA-NWS-ALERTS-AK125843C0F744.WinterWeatherAdvisory.125843C18CE0AK.AJKWSWAJK.68c3463fbe08420b2fe3cd80b30f55f6</identifier>
<sender>w-nws.webmaster@noaa.gov</sender>
<sent>2017-03-03T05:17:00-09:00</sent>
<status>Actual</status>
<msgType>Alert</msgType>
<scope>Public</scope>
<note>Alert for Dixon Entrance to Cape Decision Coastal Area (Alaska) Issued by the National Weather Service</note>
<info>
<category>Met</category>
<event>Winter Weather Advisory</event>
<urgency>Expected</urgency>
<severity>Minor</severity>
<certainty>Likely</certainty>
<eventCode>
<valueName>SAME</valueName>
<value></value>
</eventCode>
<effective>2017-03-03T05:17:00-09:00</effective>
<expires>2017-03-03T09:00:00-09:00</expires>
<senderName>NWS Juneau (Juneau and surrounding areas)</senderName>
<headline>Winter Weather Advisory issued March 03 at 5:17AM AKST until March 03 at 9:00AM AKST by NWS Juneau</headline>
<description>...WINTER WEATHER ADVISORY REMAINS IN EFFECT UNTIL 9 AM AKST THIS
MORNING...
* SNOW...Additional 1 to 3 inches through mid morning Friday.
Snowfall rates in heavier snow showers could exceed 2 inches
per hour.
* TIMING...The snow showers should begin to diminish by late
Friday morning.
* IMPACTS...Visibility will be reduced below a half mile during
heavier snow showers. Travel may be hazardous.</description>
<instruction>An advisory means that a potentially hazardous event is already
occurring or imminent.
This statement will be updated by 9 AM AKST Friday or sooner if
conditions warrant.</instruction>
<parameter>
<valueName>WMOHEADER</valueName>
<value></value>
</parameter>
<parameter>
<valueName>UGC</valueName>
<value>AKZ027</value>
</parameter>
<parameter>
<valueName>VTEC</valueName>
<value>/O.CON.PAJK.WW.Y.0015.000000T0000Z-170303T1800Z/</value>
</parameter>
<parameter>
<valueName>TIME...MOT...LOC</valueName>
<value></value>
</parameter>
<area>
<areaDesc>Dixon Entrance to Cape Decision Coastal Area</areaDesc>
<polygon></polygon>
<geocode>
<valueName>FIPS6</valueName>
<value>002198</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>AKZ027</value>
</geocode>
</area>
</info>
</alert>`

	testNWSWinterStormWarning = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>
<?xml-stylesheet href='https://alerts.weather.gov/cap/capatomproduct.xsl' type='text/xsl'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<alert xmlns = 'urn:oasis:names:tc:emergency:cap:1.1'>

<!-- http-date = Fri, 03 Mar 2017 12:30:00 GMT -->
<identifier>NOAA-NWS-ALERTS-WA125843C0AE38.WinterStormWarning.125843CF4880WA.SEWWSWSEW.64b531dc35ffdbfe7d1d2aebe5604882</identifier>
<sender>w-nws.webmaster@noaa.gov</sender>
<sent>2017-03-03T04:30:00-08:00</sent>
<status>Actual</status>
<msgType>Alert</msgType>
<scope>Public</scope>
<note>Alert for West Slopes North Central Cascades and Passes (Washington) Issued by the National Weather Service</note>
<info>
<category>Met</category>
<event>Winter Storm Warning</event>
<urgency>Expected</urgency>
<severity>Moderate</severity>
<certainty>Likely</certainty>
<eventCode>
<valueName>SAME</valueName>
<value>WSW</value>
</eventCode>
<effective>2017-03-03T04:30:00-08:00</effective>
<expires>2017-03-04T00:00:00-08:00</expires>
<senderName>NWS Seattle (Northwest Washington)</senderName>
<headline>Winter Storm Warning issued March 03 at 4:30AM PST until March 04 at 12:00AM PST by NWS Seattle</headline>
<description>...WINTER STORM WARNING REMAINS IN EFFECT UNTIL MIDNIGHT PST
TONIGHT...
* SNOW ACCUMULATIONS...An additional 6 to 14 inches is expected
through this evening...including Stevens and Snoqualmie Passes.
* SOME AFFECTED LOCATIONS...Interstate 90 including Snoqualmie
Pass, and U S Highway 2 including Stevens Pass
* TIMING...Snowfall intensity will generally peak late this
afternoon and this evening. Snoqualmie Pass should see snow
change to rain for a good bit of today, before changing back
over to snow this evening.
* SNOW LEVEL...2500 to 3000 feet through this evening, rising to
3500 feet early Friday. Quickly falling below 1500 feet on this
evening.
* MAIN IMPACT...Travel will become difficult on snow covered
roadways. Expect delays if traveling through the Cascades.</description>
<instruction>A Winter Storm Warning for heavy snow means severe winter weather
conditions are expected or occurring. Significant amounts of
snow are forecast that will make travel dangerous. Only travel in
an emergency. If you must travel, keep an extra flashlight, food,
and water in your vehicle in case of an emergency. Call 5-1-1 for
the latest road conditions in the mountain passes.</instruction>
<parameter>
<valueName>WMOHEADER</valueName>
<value></value>
</parameter>
<parameter>
<valueName>UGC</valueName>
<value>WAZ568</value>
</parameter>
<parameter>
<valueName>VTEC</valueName>
<value>/O.CON.KSEW.WS.W.0007.000000T0000Z-170304T0800Z/</value>
</parameter>
<parameter>
<valueName>TIME...MOT...LOC</valueName>
<value></value>
</parameter>
<area>
<areaDesc>West Slopes North Central Cascades and Passes</areaDesc>
<polygon></polygon>
<geocode>
<valueName>FIPS6</valueName>
<value>053033</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>053053</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>053061</value>
</geocode>
<geocode>
<valueName>FIPS6</valueName>
<value>053067</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>WAZ568</value>
</geocode>
</area>
</info>
</alert>`

	testNWSAirQualityAlert = `<?xml version = '1.0' encoding = 'UTF-8' standalone = 'yes'?>
<?xml-stylesheet href='https://alerts.weather.gov/cap/capatomproduct.xsl' type='text/xsl'?>

<!--
This atom/xml feed is an index to active advisories, watches and warnings 
issued by the National Weather Service.  This index file is not the complete 
Common Alerting Protocol (CAP) alert message.  To obtain the complete CAP 
alert, please follow the links for each entry in this index.  Also note the 
CAP message uses a style sheet to convey the information in a human readable 
format.  Please view the source of the CAP message to see the complete data 
set.  Not all information in the CAP message is contained in this index of 
active alerts.
-->

<alert xmlns = 'urn:oasis:names:tc:emergency:cap:1.1'>

<!-- http-date = Thu, 02 Mar 2017 07:31:00 GMT -->
<identifier>NOAA-NWS-ALERTS-WY125843B27DCC.AirQualityAlert.125843CE3CECWY.RIWAQARIW.81a02965446ba12d76592f0a3ba9d089</identifier>
<sender>w-nws.webmaster@noaa.gov</sender>
<sent>2017-03-02T12:31:00-07:00</sent>
<status>Actual</status>
<msgType>Alert</msgType>
<scope>Public</scope>
<note>Alert for Sublette (Wyoming) Issued by the National Weather Service</note>
<info>
<category>Met</category>
<event>Air Quality Alert</event>
<urgency>Unknown</urgency>
<severity>Unknown</severity>
<certainty>Unknown</certainty>
<eventCode>
<valueName>SAME</valueName>
<value></value>
</eventCode>
<effective>2017-03-02T12:31:00-07:00</effective>
<expires>2017-03-03T18:15:00-07:00</expires>
<senderName>NWS Riverton (Western Wyoming)</senderName>
<headline>Air Quality Alert issued March 02 at 12:31PM MST  by NWS Riverton</headline>
<description>...Ozone Action Day in Effect for Friday, March 3, 2017 for the
Upper Green River Basin Ozone Nonattainment Area...
The following message is transmitted at the request of the
Air Quality Division of the Wyoming Department of Environmental
Quality.
The Air Quality Division has issued an Ozone Action Day for
Friday, March 3, 2017 for the Upper Green River Basin Ozone
Nonattainment Area. Ozone is an air pollutant that can cause
respiratory distress especially to children, the elderly, and
people with existing respiratory conditions such as asthma. People
in these sensitive groups should limit strenuous or extended
outdoor activities, especially in the afternoon and evening.
For more information, please visit the websites for the Wyoming
Department of Environmental Quality at deq.wyoming.gov and the
Wyoming Department of Health at www.health.wyo.gov.</description>
<instruction></instruction>
<parameter>
<valueName>WMOHEADER</valueName>
<value></value>
</parameter>
<parameter>
<valueName>UGC</valueName>
<value>WYC035</value>
</parameter>
<parameter>
<valueName>VTEC</valueName>
<value></value>
</parameter>
<parameter>
<valueName>TIME...MOT...LOC</valueName>
<value></value>
</parameter>
<area>
<areaDesc>Sublette</areaDesc>
<polygon></polygon>
<geocode>
<valueName>FIPS6</valueName>
<value>056035</value>
</geocode>
<geocode>
<valueName>UGC</valueName>
<value>WYC035</value>
</geocode>
</area>
</info>
</alert>`

	// constants for unit tests
	testReferencesStringValid       = `user@example.com,XX1122333,2017-01-01T10:43:00-08:00 user2@example.com,2XX1122333,2017-01-01T10:43:00-08:00`
	testReferencesStringMissingPart = `user@example.com,2016-01-01T10:43:00-08:00`
	testReferencesStringBadTime     = `user@example.com,XX1122333,2016-01-01T10:43:00`
	testReferencesStringEmpty       = ``

	testPolygonStringValid    = `38.47,-120.14 38.52,-119.74 38.62,-119.89 38.47,-120.14`
	testPolygonStringShort    = `38.47,-120.14 38.62,-119.89 38.47,-120.14`
	testPolygonStringOpen     = `38.47,-120.14 38.34,-119.95 38.52,-119.74 38.62,-119.89`
	testPolygonStringBadPoint = `38.47,-120.14 38.52 38.62,-119.89 38.47,-120.14`
	testPolygonStringEmpty    = ``

	testCircleStringValid    = `32.9525,-115.5527 1`
	testCircleStringBadPoint = `-115.5527 1`
	testCircleStringNoPoint  = `1`
	testCircleStringNoRadius = `32.9525,-115.5527`

	testAddressesStringValid = `one@example.com two@example.com`
	testAddressesStringEmpty = ``

	testIncidentsStringValid = `XXXX1 XXXX2`
	testIncidentsStringEmpty = ``

	testSpaceDelimitedQuotedStringValid  = `"hello world" live "goodbye world"`
	testSpaceDelimitedQuotedStringValid2 = `one two "three ... (3)" four`
	testSpaceDelimitedQuotedStringValid3 = `one`
	testSpaceDelimitedQuotedStringEmpty  = ``

	testTimeStringValid   = `2017-01-01T10:43:00-08:00`
	testTimeStringBadZone = `2017-01-01T10:43:00Z`

	testURLStringFullValid     = `http://mikcamilleri.com/`
	testURLStringRelativeValid = `hello`
	testURLStringInvalid       = `http://example.com\`
)

// TestValidateMessageXML implicitely tests alert.validate().
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

// TestProcessMessageXML tests that messages are processed as expected. This
// implicitely tests alert.convert()
// TODO: Improve this test
func TestProcessMessageXML(t *testing.T) {
	assert := assert.New(t)
	var err error

	// CAP 1.2 specification examples are all valid
	_, err = ProcessMessageXML([]byte(testHomelandSecurityAdvisorySystemAlert))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testSevereThunderstormWarning))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testEarthquakeReportUpdateMessage))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testAmberAlertMultilingualMessage))
	assert.Nil(err)

	// Actual NWS examples are invalid due to empty polygon
	_, err = ProcessMessageXML([]byte(testNWSHydrologicOutlook))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testNWSWinterWeatherAdvisory))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testNWSWinterStormWarning))
	assert.Nil(err)
	_, err = ProcessMessageXML([]byte(testNWSAirQualityAlert))
	assert.Nil(err)
}

func TestRremoveEmptyStringsFromSlice(t *testing.T) {

}

func TestParseAddressesString(t *testing.T) {
	assert := assert.New(t)
	var addrs []string
	var err error

	addrs, err = parseAddressesString(testAddressesStringValid)
	assert.Equal([]string{"one@example.com", "two@example.com"}, addrs)
	assert.Nil(err)

	addrs, err = parseAddressesString(testAddressesStringEmpty)
	assert.Nil(addrs)
	assert.NotNil(err)
}

func TestIsValidAddressesString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidAddressesString(testAddressesStringValid))
	assert.False(isValidAddressesString(testAddressesStringEmpty))
}

func TestParseIncidentsString(t *testing.T) {
	assert := assert.New(t)
	var incidents []string
	var err error

	incidents, err = parseIncidentsString(testIncidentsStringValid)
	assert.Equal([]string{"XXXX1", "XXXX2"}, incidents)
	assert.Nil(err)

	incidents, err = parseIncidentsString(testIncidentsStringEmpty)
	assert.Nil(incidents)
	assert.NotNil(err)
}

func TestIsValidIncidentsString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidIncidentsString(testIncidentsStringValid))
	assert.False(isValidIncidentsString(testIncidentsStringEmpty))
}

func TestSplitSpaceDelimitedQuotedStrings(t *testing.T) {
	assert := assert.New(t)
	var strs []string
	var err error

	strs, err = splitSpaceDelimitedQuotedStrings(testSpaceDelimitedQuotedStringValid)
	assert.Equal([]string{"hello world", "live", "goodbye world"}, strs)
	assert.Nil(err)

	strs, err = splitSpaceDelimitedQuotedStrings(testSpaceDelimitedQuotedStringValid2)
	assert.Equal([]string{"one", "two", "three ... (3)", "four"}, strs)
	assert.Nil(err)

	strs, err = splitSpaceDelimitedQuotedStrings(testSpaceDelimitedQuotedStringValid3)
	assert.Equal([]string{"one"}, strs)
	assert.Nil(err)

	strs, err = splitSpaceDelimitedQuotedStrings(testSpaceDelimitedQuotedStringEmpty)
	assert.Nil(strs)
	assert.NotNil(err)
}

func TestIsValidSpaceDelimitedQuotedStrings(t *testing.T) {

}

func TestParseTimeString(t *testing.T) {

}

func TestIsValidTimeString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidTimeString(testTimeStringValid))
	assert.False(isValidTimeString(testTimeStringBadZone))
}

func TestParseURLString(t *testing.T) {

}

func TestIsValidURLString(t *testing.T) {
	assert := assert.New(t)
	assert.True(isValidURLString(testURLStringFullValid))
	assert.True(isValidURLString(testURLStringRelativeValid))
	assert.False(isValidURLString(testURLStringInvalid))
}

func TestParseReferencesString(t *testing.T) {
	assert := assert.New(t)
	var refs []Reference
	var err error

	refs, err = parseReferencesString(testReferencesStringValid)
	assert.Len(refs, 2)
	assert.Equal("user@example.com", refs[0].Sender)
	assert.Equal("XX1122333", refs[0].Identifier)
	tm, _ := time.Parse("2006-01-02T15:04:05-07:00", "2017-01-01T10:43:00-08:00")
	assert.Equal(tm, refs[0].Sent)
	assert.Nil(err)

	refs, err = parseReferencesString(testReferencesStringMissingPart)
	assert.Nil(refs)
	assert.NotNil(err)

	refs, err = parseReferencesString(testReferencesStringBadTime)
	assert.Nil(refs)
	assert.NotNil(err)

	refs, err = parseReferencesString(testReferencesStringEmpty)
	assert.Nil(refs)
	assert.NotNil(err)
}

func TestIsValidReferencesString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidReferencesString(testReferencesStringValid))
	assert.False(isValidReferencesString(testReferencesStringMissingPart))
	assert.False(isValidReferencesString(testReferencesStringBadTime))
	assert.False(isValidReferencesString(testReferencesStringEmpty))
}

func TestParseSingleReferencesString(t *testing.T) {

}

func TestParsePolygonString(t *testing.T) {
	assert := assert.New(t)
	var poly Polygon
	var err error

	poly, err = parsePolygonString(testPolygonStringValid)
	assert.Len(poly, 4)
	assert.Equal(Point{Latitude: 38.47, Longitude: -120.14}, poly[0])
	assert.Nil(err)

	poly, err = parsePolygonString(testPolygonStringShort)
	assert.Len(poly, 0)
	assert.NotNil(err)

	poly, err = parsePolygonString(testPolygonStringOpen)
	assert.Len(poly, 0)
	assert.NotNil(err)

	poly, err = parsePolygonString(testPolygonStringBadPoint)
	assert.Len(poly, 0)
	assert.NotNil(err)

	poly, err = parsePolygonString(testPolygonStringEmpty)
	assert.Len(poly, 0)
	assert.NotNil(err)
}

func TestIsValidPolygonString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidPolygonString(testPolygonStringValid))
	assert.False(isValidPolygonString(testPolygonStringShort))
	assert.False(isValidPolygonString(testPolygonStringOpen))
	assert.False(isValidPolygonString(testPolygonStringBadPoint))
	assert.False(isValidPolygonString(testPolygonStringEmpty))
}

func TestParseCircleString(t *testing.T) {
	assert := assert.New(t)
	var circle Circle
	var err error

	circle, err = parseCircleString(testCircleStringValid)
	assert.Equal(Circle{Point: Point{Latitude: 32.9525, Longitude: -115.5527}, Radius: 1}, circle)
	assert.Nil(err)

	circle, err = parseCircleString(testCircleStringBadPoint)
	assert.Equal(Circle{}, circle)
	assert.NotNil(err)

	circle, err = parseCircleString(testCircleStringNoPoint)
	assert.Equal(Circle{}, circle)
	assert.NotNil(err)

	circle, err = parseCircleString(testCircleStringNoRadius)
	assert.Equal(Circle{}, circle)
	assert.NotNil(err)
}

func TestIsValidCircleString(t *testing.T) {
	assert := assert.New(t)

	assert.True(isValidCircleString(testCircleStringValid))
	assert.False(isValidCircleString(testCircleStringBadPoint))
	assert.False(isValidCircleString(testCircleStringNoPoint))
	assert.False(isValidCircleString(testCircleStringNoRadius))
}
