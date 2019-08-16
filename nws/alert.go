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
)

// Alert ...
type Alert struct {
}

// getActiveAlertsForPoint ...
func getActiveAlertsForPoint(httpClinet *http.Client, point Point) ([]Alert, error) {
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
