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
