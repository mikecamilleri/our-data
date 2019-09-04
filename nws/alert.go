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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const getActiveAlertsForPointEndpointURLStringFmt = "alerts/active"

var (
	// AlertStatuses are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
	AlertStatuses = map[string]string{
		"Actual":   "Actionable by all targeted recipients",
		"Exercise": "Actionable only by designated exercise participants; exercise identifier SHOULD appear in <note>",
		"System":   "For messages that support alert network internal functions",
		"Test":     "Technical testing only, all recipients disregard",
		"Draft":    "A preliminary template or draft, not actionable in its current form",
	}

	// AlertMessageTypes are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
	AlertMessageTypes = map[string]string{
		"Alert":  "Initial information requiring attention by targeted recipients",
		"Update": "Updates and supercedes the earlier message(s) identified in <references>",
		"Cancel": "Cancels the earlier message(s) identified in <references>",
		"Ack":    "Acknowledges receipt and acceptance of the message(s) identified in <references>",
		"Error":  "Indicates rejection of the message(s) identified in <references>; explanation SHOULD appear in <note>",
	}

	//AlertCategories are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
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

	// AlertSeverities are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
	AlertSeverities = map[string]string{
		"Extreme":  "Extraordinary threat to life or property",
		"Severe":   "Significant threat to life or property",
		"Moderate": "Possible threat to life or property",
		"Minor":    "Minimal to no known threat to life or property",
		"Unknown":  "Severity unknown",
	}

	// AlertCertainties are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
	AlertCertainties = map[string]string{
		"Observed": "Determined to have occurred or to be ongoing",
		"Likely":   "Likely (p > ~50%)",
		"Possible": "Possible but not likely (p <= ~50%)",
		"Unlikely": "Not expected to occur (p ~ 0)",
		"Unknown":  "Certainty unknown",
	}

	// AlertUrgencies are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
	AlertUrgencies = map[string]string{
		"Immediate": "Responsive action SHOULD be taken immediately",
		"Expected":  "Responsive action SHOULD be taken soon (within next hour)",
		"Future":    "Responsive action SHOULD be taken in the near future",
		"Past":      "Responsive action is no longer required",
		"Unknown":   "Urgency not known",
	}

	// AlertResponses are defined in
	// http://docs.oasis-open.org/emergency/cap/v1.2/CAP-v1.2.html
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

// An Alert represents a single alert returned from the NWS API.
type Alert struct {
	ID string

	TimeRetrieved time.Time // when the client retrieved this alert
	TimeSent      time.Time // when this alert was sent
	TimeEffective time.Time // when the information in this messgae becomes effective
	TimeExpires   time.Time // when the information in this messgae expires
	TimeOnset     time.Time // when the beginning of the hazard is expected
	TimeEnds      time.Time // not in CAP spec, likely when the end of the hazard is expected

	SenderID   string // appears to usually be an email address
	SenderName string

	Status      string   // must be a key in AlertStatuses
	MessageType string   // must be a key in AlertMessageTypes
	References  []string // IDs of alerts that this alert affects based on MessageType

	Category        string // must be a key in AlertCategories
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

// getActiveAlertsForPoint retrieves from the NWS API active alerts for a given
// point.
func getActiveAlertsForPoint(httpClient *http.Client, httpUserAgentString string, point Point) ([]Alert, error) {
	// It may be more efficient to use "zone" or "area", but it isn't clear from
	// the limited documentation whish is most appropriate. "Point" seems like it
	// has the best chance of returning appropriate/relevent alerts.
	var query url.Values
	query.Add("point", fmt.Sprintf("%s,%s", point.Lat, point.Lon))
	respBody, err := doAPIRequest(
		httpClient,
		httpUserAgentString,
		fmt.Sprintf(getActiveAlertsForPointEndpointURLStringFmt),
		query,
	)
	if err != nil {
		return nil, err
	}
	return newAlertsFromAlertsRespBody(respBody)
}

// newAlertsFromAlertsRespBody returns a slice of Alerts, given a response body
// from the NWS API.
func newAlertsFromAlertsRespBody(respBody []byte) ([]Alert, error) {
	// unmarshal the body into a temporary struct
	alertsRaw := struct {
		Features []struct {
			Properties struct {
				ID         string
				AreaDesc   string
				References []struct {
					Identifier string
				}
				Sent        string
				Effective   string
				Onset       string
				Expires     string
				Ends        string
				Status      string
				MessageType string
				Category    string
				Severity    string
				Certainty   string
				Urgency     string
				Event       string
				Sender      string
				SenderName  string
				Headline    string
				Description string
				Instruction string
				Response    string
			}
		}
	}{}
	if err := json.Unmarshal(respBody, alertsRaw); err != nil {
		return nil, err
	}

	// validate and build returned slice
	var alerts []Alert

	for _, aRaw := range alertsRaw.Features {
		var ok bool
		var a Alert

		if aRaw.Properties.ID == "" {
			continue // skip if no ID
		}
		a.ID = aRaw.Properties.ID

		// generally, ignore bad data
		// the idea here is to get as complete an alert as possible
		a.TimeRetrieved = time.Now()
		a.TimeSent, _ = time.Parse(time.RFC3339, aRaw.Properties.Sent)
		a.TimeEffective, _ = time.Parse(time.RFC3339, aRaw.Properties.Effective)
		a.TimeExpires, _ = time.Parse(time.RFC3339, aRaw.Properties.Expires)
		a.TimeOnset, _ = time.Parse(time.RFC3339, aRaw.Properties.Onset)
		a.TimeEnds, _ = time.Parse(time.RFC3339, aRaw.Properties.Ends)

		a.SenderID = aRaw.Properties.Sender
		a.SenderName = aRaw.Properties.SenderName

		a.Status = aRaw.Properties.Status
		a.MessageType = aRaw.Properties.MessageType
		for _, ref := range aRaw.Properties.References {
			if ref.Identifier != "" {
				a.References = append(a.References, ref.Identifier)
			}
		}

		if _, ok = AlertCategories[aRaw.Properties.Category]; ok {
			a.Category = aRaw.Properties.Category
		}
		if _, ok = AlertSeverities[aRaw.Properties.Severity]; ok {
			a.Severity = aRaw.Properties.Severity
		}
		if _, ok = AlertCertainties[aRaw.Properties.Certainty]; ok {
			a.Certainty = aRaw.Properties.Certainty
		}
		if _, ok = AlertUrgencies[aRaw.Properties.Urgency]; ok {
			a.Urgency = aRaw.Properties.Urgency
		}
		a.Event = aRaw.Properties.Event
		a.AreaDescription = aRaw.Properties.AreaDesc
		a.Headline = aRaw.Properties.Headline
		a.Description = aRaw.Properties.Description
		a.Instruction = aRaw.Properties.Instruction
		if _, ok = AlertResponses[aRaw.Properties.Response]; ok {
			a.Response = aRaw.Properties.Response
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}

////////////////////////////////////////////////////////////////////////////////
// EXAMPLE request and responses below.
// - note different location (needed a place with an active alert)
// - looks like JSONized CAP
// - does `/alerts/active` endpoint include cacellations? Use `/alerts` instead?
//   - just assume cancelled if no longer present?

// mike@Darwin-D ~ % curl -i -X GET "https://api.weather.gov/alerts/active?point=45.458,-122.6636"
// HTTP/2 200
// server: nginx/1.12.2
// content-type: application/geo+json
// last-modified: Wed, 28 Aug 2019 17:36:22 GMT
// access-control-allow-origin: *
// x-server-id: vm-bldr-nids-apiapp4.ncep.noaa.gov
// x-correlation-id: 3e62c5d1-f0e1-4253-be7c-f3a392d44be3
// x-request-id: 3e62c5d1-f0e1-4253-be7c-f3a392d44be3
// cache-control: public, max-age=30, s-maxage=30
// expires: Wed, 28 Aug 2019 17:38:39 GMT
// date: Wed, 28 Aug 2019 17:38:09 GMT
// content-length: 6553
// vary: Accept,Feature-Flags
// strict-transport-security: max-age=31536000 ; includeSubDomains ; preload

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
//             "id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3789007-3246513",
//             "type": "Feature",
//             "geometry": null,
//             "properties": {
//                 "@id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3789007-3246513",
//                 "@type": "wx:Alert",
//                 "id": "NWS-IDP-PROD-3789007-3246513",
//                 "areaDesc": "Cascade Foothills in Lane County; Northern Oregon Cascade Foothills; Greater Portland Metro Area; Greater Vancouver Area; Lower Columbia; Lower Columbia and I - 5 Corridor in Cowlitz County; Central Willamette Valley; Western Columbia River Gorge; South Washington Cascade Foothills; Western Columbia River Gorge",
//                 "geocode": {
//                     "UGC": [
//                         "ORZ012",
//                         "ORZ010",
//                         "ORZ006",
//                         "WAZ039",
//                         "ORZ005",
//                         "WAZ022",
//                         "ORZ007",
//                         "ORZ015",
//                         "WAZ040",
//                         "WAZ045"
//                     ],
//                     "SAME": [
//                         "041039",
//                         "041005",
//                         "041043",
//                         "041047",
//                         "041051",
//                         "041009",
//                         "041067",
//                         "053011",
//                         "053015",
//                         "053069",
//                         "041053",
//                         "041071",
//                         "041027",
//                         "053059"
//                     ]
//                 },
//                 "affectedZones": [
//                     "https://api.weather.gov/zones/forecast/ORZ012",
//                     "https://api.weather.gov/zones/forecast/ORZ010",
//                     "https://api.weather.gov/zones/forecast/ORZ006",
//                     "https://api.weather.gov/zones/forecast/WAZ039",
//                     "https://api.weather.gov/zones/forecast/ORZ005",
//                     "https://api.weather.gov/zones/forecast/WAZ022",
//                     "https://api.weather.gov/zones/forecast/ORZ007",
//                     "https://api.weather.gov/zones/forecast/ORZ015",
//                     "https://api.weather.gov/zones/forecast/WAZ040",
//                     "https://api.weather.gov/zones/forecast/WAZ045"
//                 ],
//                 "references": [
//                     {
//                         "@id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3788552-3246201",
//                         "identifier": "NWS-IDP-PROD-3788552-3246201",
//                         "sender": "w-nws.webmaster@noaa.gov",
//                         "sent": "2019-08-27T20:12:00-07:00"
//                     },
//                     {
//                         "@id": "https://api.weather.gov/alerts/NWS-IDP-PROD-3788552-3246200",
//                         "identifier": "NWS-IDP-PROD-3788552-3246200",
//                         "sender": "w-nws.webmaster@noaa.gov",
//                         "sent": "2019-08-27T20:12:00-07:00"
//                     }
//                 ],
//                 "sent": "2019-08-28T04:28:00-07:00",
//                 "effective": "2019-08-28T04:28:00-07:00",
//                 "onset": "2019-08-28T09:00:00-07:00",
//                 "expires": "2019-08-28T20:00:00-07:00",
//                 "ends": "2019-08-28T20:00:00-07:00",
//                 "status": "Actual",
//                 "messageType": "Update",
//                 "category": "Met",
//                 "severity": "Moderate",
//                 "certainty": "Likely",
//                 "urgency": "Expected",
//                 "event": "Heat Advisory",
//                 "sender": "w-nws.webmaster@noaa.gov",
//                 "senderName": "NWS Portland OR",
//                 "headline": "Heat Advisory issued August 28 at 4:28AM PDT until August 28 at 8:00PM PDT by NWS Portland OR",
//                 "description": "* HIGH TEMPERATURES...92 to 102 degrees today.\n\n* TIMING...Hottest time of the day will be between 2 and 7 PM.\nSome cooling may occur a little earlier than 7 PM for areas near\ngaps in the Coast Range.\n\n* IMPACTS...Hot temperatures will increase the chance for heat\nrelated illnesses, especially for those who are sensitive to\nheat. People most vulnerable include those who spend a lot of\ntime outdoors, those without air conditioning, those without\nadequate hydration, young children, and the elderly.",
//                 "instruction": "A Heat Advisory means that a period of hot temperatures is\nexpected. Hot temperatures will create a situation in which heat\nrelated illnesses are possible. Drink plenty of fluids, stay in\nan air-conditioned room, stay out of the sunshine, and check up\non relatives and neighbors.\n\nTake extra precautions, if you work or spend time outside. When\npossible, reschedule strenuous activities to early morning or\nevening. Know the signs and symptoms of heat exhaustion and heat\nstroke. Wear light weight and loose fitting clothing when\npossible and drink plenty of water.\n\nTo reduce risk during outdoor work, the Occupational Safety and\nHealth Administration recommends scheduling frequent rest breaks\nin shaded or air conditioned environments. Anyone overcome by\nheat should be moved to a cool and shaded location. Heat stroke\nis an emergency, call 9 1 1.",
//                 "response": "Execute",
//                 "parameters": {
//                     "NWSheadline": [
//                         "HEAT ADVISORY REMAINS IN EFFECT FROM 9 AM THIS MORNING TO 8 PM PDT THIS EVENING"
//                     ],
//                     "VTEC": [
//                         "/O.CON.KPQR.HT.Y.0003.190828T1600Z-190829T0300Z/"
//                     ],
//                     "PIL": [
//                         "PQRNPWPQR"
//                     ],
//                     "BLOCKCHANNEL": [
//                         "CMAS",
//                         "EAS",
//                         "NWEM"
//                     ],
//                     "eventEndingTime": [
//                         "2019-08-28T20:00:00-07:00"
//                     ]
//                 }
//             }
//         }
//     ],
//     "title": "current watches, warnings, and advisories for 45.458 N, 122.6636 W",
//     "updated": "2019-08-28T17:36:22+00:00"
// }
