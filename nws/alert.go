// Copyright 2019 Michael Camilleri <mike@mikecamilleri.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nws

import (
	"net/http"
	"time"
)

var (
	// AlertStatuses ...
	AlertStatuses = map[string]string{
		"Actual":   "Actionable by all targeted recipients",
		"Exercise": "Actionable only by designated exercise participants; exercise identifier SHOULD appear in <note>",
		"System":   "For messages that support alert network internal functions",
		"Test":     "Technical testing only, all recipients disregard",
		"Draft":    "A preliminary template or draft, not actionable in its current form",
	}

	// AlertMessageTypes ...
	AlertMessageTypes = map[string]string{
		"Alert":  "Initial information requiring attention by targeted recipients",
		"Update": "Updates and supercedes the earlier message(s) identified in <references>",
		"Cancel": "Cancels the earlier message(s) identified in <references>",
		"Ack":    "Acknowledges receipt and acceptance of the message(s) identified in <references>",
		"Error":  "Indicates rejection of the message(s) identified in <references>; explanation SHOULD appear in <note>",
	}

	//AlertCategories ...
	AlertCategories = map[string]string{
		"Geo":       "Geophysical (inc. landslide)",
		"Met":       "Meteorological (inc. flood)",
		"Safety":    "General emergency and public safety",
		"Security":  "Law enforcement, military, homeland and local/private security",
		"Rescue":    "Rescue and recovery",
		"Fire":      "Fire suppression and rescue",
		"Health":    "Medical and public health",
		"Env":       "Pollution and other environmental",
		"Transport": "Public and private transportation",
		"Infra":     "Utility, telecommunication, other non-transport infrastructure",
		"CBRNE":     "Chemical, Biological, Radiological, Nuclear or High-Yield Explosive threat or attack",
		"Other":     "Other events",
	}

	// AlertSeverities ...
	AlertSeverities = map[string]string{
		"Extreme":  "Extraordinary threat to life or property",
		"Severe":   "Significant threat to life or property",
		"Moderate": "Possible threat to life or property",
		"Minor":    "Minimal to no known threat to life or property",
		"Unknown":  "Severity unknown",
	}

	// AlertCertainties ...
	AlertCertainties = map[string]string{
		"Observed": "Determined to have occurred or to be ongoing",
		"Likely":   "Likely (p > ~50%)",
		"Possible": "Possible but not likely (p <= ~50%)",
		"Unlikely": "Not expected to occur (p ~ 0)",
		"Unknown":  "Certainty unknown",
	}

	// AlertUrgencies ...
	AlertUrgencies = map[string]string{
		"Immediate": "Responsive action SHOULD be taken immediately",
		"Expected":  "Responsive action SHOULD be taken soon (within next hour)",
		"Future":    "Responsive action SHOULD be taken in the near future",
		"Past":      "Responsive action is no longer required",
		"Unknown":   "Urgency not known",
	}

	// AlertResponses ...
	AlertResponses = map[string]string{
		"Shelter":  "Take shelter in place or per <instruction>",
		"Evacuate": "Relocate as instructed in the <instruction>",
		"Prepare":  "Make preparations per the <instruction>",
		"Execute":  "Execute a pre-planned activity identified in <instruction>",
		"Avoid":    "Avoid the subject event as per the <instruction>",
		"Monitor":  "Attend to information sources as described in <instruction>",
		"Assess":   "Evaluate the information in this message",
		"AllClear": "The subject event no longer poses a threat or concern and any follow on action is described in <instruction>",
		"None":     "No action recommended",
	}
)

// Alert ...
type Alert struct {
	ID string

	TimeRetrieved time.Time // when the client retrieved this alert
	TimeSent      time.Time // when this alert was sent
	TimeEffective time.Time // when the information in this messgae becomes effective
	TimeExpires   time.Time // when the information in this messgae expires
	TimeOnset     time.Time // when the beginning of the hazard is expected
	TimeEnds      time.Time // not in CAP spec, likely when the end of the hazard is expected

	SenderID        string // appears to usually be an email address
	SenderName      string
	Status          string // must be a key in AlertStatuses
	MessageType     string // must be a key in AlertMessageTypes
	Category        string // must ge a key in AlertCategories
	Severity        string // must be a key in AlertSeverities
	Certainty       string // must be a key in AlertCertainties
	Urgency         string // must be a key in Alert Urgencies
	Event           string
	AreaDescription string
	Headline        string
	Description     string
	Instruction     string
	Response        string // must be a key in AlerResponses
}

// getActiveAlertsForPoint ...
func getActiveAlertsForPoint(httpClinet *http.Client, httpUserAgentString string, point Point) ([]Alert, error) {
	return nil, nil
}

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - note different location (needed a place with an active alert)
// - looks like JSONized CAP
// - does `/alerts/active` endpoint include cacellations? Use `/alerts` instead?
//   - just assume cancelled if no longer present?

// mike@Darwin-D nws % curl -X GET "https://api.weather.gov/alerts/active?point=32.2111,-81.4178"
// {
//     "@context": [
//         "https://raw.githubusercontent.com/geojson/geojson-ld/master/contexts/geojson-base.jsonld",
//         {
//             "wx": "https://api.weather.gov/ontology#",
//             "@vocab": "https://api.weather.gov/ontology#"
//         }
//     ],
//     "type": "FeatureCollection",
//     "features": [
//         {
//             "id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3762936-3231541",
//             "type": "Feature",
//             "geometry": null,
//             "properties": {
//                 "@id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3762936-3231541",
//                 "@type": "wx:Alert",
//                 "id": "NWS-IDP-PROD-3762936-3231541",
//                 "areaDesc": "Inland Liberty; Evans; Long; Coastal Liberty; Coastal Bryan; Tattnall; Inland Bryan; Effingham; Bulloch; Coastal McIntosh; Inland Chatham; Candler; Inland McIntosh; Coastal Chatham",
//                 "geocode": {
//                     "UGC": [
//                         "GAZ138",
//                         "GAZ115",
//                         "GAZ137",
//                         "GAZ139",
//                         "GAZ117",
//                         "GAZ114",
//                         "GAZ116",
//                         "GAZ101",
//                         "GAZ100",
//                         "GAZ141",
//                         "GAZ118",
//                         "GAZ099",
//                         "GAZ140",
//                         "GAZ119"
//                     ],
//                     "SAME": [
//                         "013179",
//                         "013109",
//                         "013183",
//                         "013029",
//                         "013267",
//                         "013103",
//                         "013031",
//                         "013191",
//                         "013051",
//                         "013043"
//                     ]
//                 },
//                 "affectedZones": [
//                     "https://api.weather.gov/zones/forecast/GAZ138",
//                     "https://api.weather.gov/zones/forecast/GAZ115",
//                     "https://api.weather.gov/zones/forecast/GAZ137",
//                     "https://api.weather.gov/zones/forecast/GAZ139",
//                     "https://api.weather.gov/zones/forecast/GAZ117",
//                     "https://api.weather.gov/zones/forecast/GAZ114",
//                     "https://api.weather.gov/zones/forecast/GAZ116",
//                     "https://api.weather.gov/zones/forecast/GAZ101",
//                     "https://api.weather.gov/zones/forecast/GAZ100",
//                     "https://api.weather.gov/zones/forecast/GAZ141",
//                     "https://api.weather.gov/zones/forecast/GAZ118",
//                     "https://api.weather.gov/zones/forecast/GAZ099",
//                     "https://api.weather.gov/zones/forecast/GAZ140",
//                     "https://api.weather.gov/zones/forecast/GAZ119"
//                 ],
//                 "references": [
//                     {
//                         "@id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3762718-3231437",
//                         "identifier": "NWS-IDP-PROD-3762718-3231437",
//                         "sender": "w-nws.webmaster@noaa.gov",
//                         "sent": "2019-08-14T03:34:00-04:00"
//                     }
//                 ],
//                 "sent": "2019-08-14T07:20:00-04:00",
//                 "effective": "2019-08-14T07:20:00-04:00",
//                 "onset": "2019-08-14T11:00:00-04:00",
//                 "expires": "2019-08-14T16:00:00-04:00",
//                 "ends": "2019-08-14T18:00:00-04:00",
//                 "status": "Actual",
//                 "messageType": "Update",
//                 "category": "Met",
//                 "severity": "Moderate",
//                 "certainty": "Likely",
//                 "urgency": "Expected",
//                 "event": "Heat Advisory",
//                 "sender": "w-nws.webmaster@noaa.gov",
//                 "senderName": "NWS Charleston SC",
//                 "headline": "Heat Advisory issued August 14 at 7:20AM EDT until August 14 at 6:00PM EDT by NWS Charleston SC",
//                 "description": "* HEAT INDEX VALUES...Around 110 due to temperatures in the upper\n90s, and dewpoints in the upper 70s.\n\n* TIMING...11 AM to 6 PM\n\n* IMPACTS...Dangerously high temperatures and humidity could\nquickly cause heat stress or heat stroke if precautions are\nnot taken.",
//                 "instruction": "If you must be outdoors, drink plenty of fluids, wear light\nweight clothing and stay out of direct sunshine. In addition,\nknow the signs of heat illnesses and be sure to check on those\nwho are most vulnerable to the heat such as young children and\nthe elderly. Never leave children or pets in a vehicle.\n\nTo reduce risk during outdoor work, the Occupational Safety and\nHealth Administration recommends scheduling frequent rest breaks\nin shaded or air conditioned environments. Anyone overcome by\nheat should be moved to a cool and shaded location. Heat stroke\nis an emergency - call 9 1 1.",
//                 "response": "Execute",
//                 "parameters": {
//                     "NWSheadline": [
//                         "HEAT ADVISORY REMAINS IN EFFECT FROM 11 AM THIS MORNING TO 6 PM EDT THIS EVENING"
//                     ],
//                     "VTEC": [
//                         "/O.CON.KCHS.HT.Y.0006.190814T1500Z-190814T2200Z/"
//                     ],
//                     "PIL": [
//                         "CHSNPWCHS"
//                     ],
//                     "BLOCKCHANNEL": [
//                         "CMAS",
//                         "EAS",
//                         "NWEM"
//                     ],
//                     "eventEndingTime": [
//                         "2019-08-14T18:00:00-04:00"
//                     ]
//                 }
//             }
//         }
//     ],
//     "title": "current watches, warnings, and advisories for 32.2111 N, 81.4178 W",
//     "updated": "2019-08-14T17:51:02+00:00"
// }%
// mike@Darwin-D nws %
