/*
	Package ouralerts implements the ability to parse and validate OASIS Common
	Alerting Protocol Alert Messages
*/
package ouralerts

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat           = "2006-01-02T15:04:05-07:00"
	restrictedCharacters = " ,<&"
)

var (
	XMLNamespaces = map[string]string{
		"urn:oasis:names:tc:emergency:cap:1.2": "1.2",
		"urn:oasis:names:tc:emergency:cap:1.1": "1.1",
		"urn:oasis:names:tc:emergency:cap:1.0": "1.0",
	}
	AlertStatuses = map[string]string{
		"Actual":   "Actionable by all targeted recipients",
		"Exercise": "Actionable only by designated exercise participants; exercise identifier SHOULD appear in <note>",
		"System":   "For messages that support alert network internal functions",
		"Test":     "Technical testing only, all recipients disregard",
		"Draft":    "A preliminary template or draft, not actionable in its current form",
	}
	AlertMsgTypes = map[string]string{
		"Alert":  "Initial information requiring attention by targeted recipients",
		"Update": "Updates and supercedes the earlier message(s) identified in <references>",
		"Cancel": "Cancels the earlier message(s) identified in <references>",
		"Ack":    "Acknowledges receipt and acceptance of the message(s) identified in <references>",
		"Error":  "Indicates rejection of the message(s) identified in <references>; explanation SHOULD appear in <note>",
	}
	AlertScopes = map[string]string{
		"Public":     "For general dissemination to unrestricted audiences",
		"Restricted": "For dissemination only to users with a known operational requirement (see <restriction>)",
		"Private":    "For dissemination only to specified addresses (see <addresses>)",
	}
	AlertInfoCategories = map[string]string{
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
	AlertInfoResponseTypes = map[string]string{
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
	AlertInfoUrgencies = map[string]string{
		"Immediate": "Responsive action SHOULD be taken immediately",
		"Expected":  "Responsive action SHOULD be taken soon (within next hour)",
		"Future":    "Responsive action SHOULD be taken in the near future",
		"Past":      "Responsive action is no longer required",
		"Unknown":   "Urgency not known",
	}
	AlertInfoSeverities = map[string]string{
		"Extreme":  "Extraordinary threat to life or property",
		"Severe":   "Significant threat to life or property",
		"Moderate": "Possible threat to life or property",
		"Minor":    "Minimal to no known threat to life or property",
		"Unknown":  "Severity unknown",
	}
	AlertInfoCertainties = map[string]string{
		"Observed": "Determined to have occurred or to be ongoing",
		"Likely":   "Likely (p > ~50%)",
		"Possible": "Possible but not likely (p <= ~50%)",
		"Unlikely": "Not expected to occur (p ~ 0)",
		"Unknown":  "Certainty unknown",
	}
)

// Alert represents a parsed and validated CAP alert message
type Alert struct {
	// TODO: add namespace support to distiguish CAP versions
	// CAPVersion  string
	Identifier  string
	Sender      string
	Sent        time.Time
	Status      string
	MsgType     string
	Source      string
	Scope       string
	Restriction string
	Addresses   []string
	Codes       []string
	Note        string
	References  []Reference
	Incidents   []string
	Infos       []Info
}

// Info
type Info struct {
	Language      string
	Categories    []string
	Event         string
	ResponseTypes []string
	Urgency       string
	Severity      string
	Certainty     string
	Audience      string
	EventCodes    []NamedValue
	Effective     time.Time
	Onset         time.Time
	Expires       time.Time
	SenderName    string
	Headline      string
	Description   string
	Instruction   string
	Web           *url.URL
	Contact       string
	Parameters    []NamedValue
	Resources     []Resource
	Areas         []Area
}

// Resource
type Resource struct {
	ResourceDesc string
	MIMEType     string
	Size         int // approximate size in bytes
	URI          *url.URL
	DerefURI     string // base-64 encoded binary
	Digest       string // SHA-1 hash
}

// Area
type Area struct {
	AreaDesc string
	// Polygon is already a reference type, but we want a pointer here for
	// consistency with []*Cricle
	Polygons []Polygon
	Circles  []Circle
	Geocodes []NamedValue
	// because golang has no built in support for decimals, these values
	// are being left as `string` so the caller can handle as necessary.
	Altitude string // feet above mean sea level
	Ceiling  string // feet above mean sea level
}

// NamedValue
type NamedValue struct {
	ValueName string
	Value     string
}

// Reference holds a reference to another alert
type Reference struct {
	Sender     string
	Identifier string
	Sent       time.Time
}

// parseReferencesString parses a references string
func parseReferencesString(referencesString string) ([]Reference, error) {
	if len(referencesString) == 0 {
		return nil, errors.New("referencesString is empty")
	}
	refStrings := strings.Fields(referencesString)
	var refs []Reference
	for _, rs := range refStrings {
		parts := strings.Split(rs, ",")
		if len(parts) != 3 {
			return nil, errors.New("reference must contain three parts")
		}
		t, err := time.Parse(timeFormat, parts[2])
		if err != nil {
			return nil, errors.New("invalid time string")
		}
		refs = append(refs, Reference{Sender: parts[0], Identifier: parts[1], Sent: t})
	}
	return refs, nil
}

// isValidReferencesString tests whether a references string is valid
func isValidReferencesString(referenceString string) bool {
	if _, err := parseReferencesString(referenceString); err != nil {
		return false
	}
	return true
}

// Point defines a WGS 84 coordinate point on the earth
type Point struct {
	// because golang has no built in support for decimals, these values
	// are being left as `string` so the caller can handle as necessary.
	Latitude  string
	Longitude string
}

// Polygon defines a polygonal area
type Polygon []Point

// parsePolygonString parses a polygon string
func parsePolygonString(polygonString string) (Polygon, error) {
	// TODO: test validity of numbers
	pointStrings := strings.Fields(polygonString)
	if len(pointStrings) < 4 {
		return nil, errors.New("a polygon must contain al least four points")
	}
	var poly Polygon
	for _, ps := range pointStrings {
		vals := strings.Split(ps, ",")
		if len(vals) != 2 {
			return nil, errors.New("point must contain exactly two values")
		}
		poly = append(poly, Point{Latitude: vals[0], Longitude: vals[1]})
	}
	if poly[0] != poly[len(poly)-1] {
		return nil, errors.New("first and last points must be equal")
	}
	return poly, nil
}

// isValidPolygonString tests whether a polygon string is valid
func isValidPolygonString(polygonString string) bool {
	if _, err := parsePolygonString(polygonString); err != nil {
		return false
	}
	return true
}

// Circle defines a circular area
type Circle struct {
	// because golang has no built in support for decimals, these values
	// are being left as `string` so the caller can handle as necessary.
	// The coordinate system used is WGS 84.
	Point  Point
	Radius string // kilometers
}

// parseCircleString
func parseCircleString(circleString string) (Circle, error) {
	// TODO: test validity of numbers
	var circle Circle
	pointRadiusStrings := strings.Fields(circleString)
	if len(pointRadiusStrings) != 2 {
		return circle, errors.New("a circle must contain a central point and a radius")
	}
	psVals := strings.Split(pointRadiusStrings[0], ",")
	if len(psVals) != 2 {
		return circle, errors.New("central point must contain exactly two values")
	}
	circle.Point = Point{Latitude: psVals[0], Longitude: psVals[1]}
	circle.Radius = pointRadiusStrings[1]
	return circle, nil
}

// isValidCircleString tests whether a circle string is valid
func isValidCircleString(circleString string) bool {
	if _, err := parseCircleString(circleString); err != nil {
		return false
	}
	return true
}

// parseAddressesString parses an addresses string
func parseAddressesString(addressesString string) []string {
	return splitSpaceDelimitedQuotedStrings(addressesString)
}

// isValidAddressesString tests whether an addresses string is valid
func isValidAddressesString(addressesString string) bool {
	if len(parseAddressesString(addressesString)) == 0 {
		return false
	}
	return true
}

// parseIncidentsString parases an incidents string
func parseIncidentsString(indidentsString string) []string {
	return splitSpaceDelimitedQuotedStrings(indidentsString)
}

// isValidIncidentsString tests whether an incidentas string is valid
func isValidIncidentsString(incidentsString string) bool {
	if len(parseIncidentsString(incidentsString)) == 0 {
		return false
	}
	return true
}

// splitSpaceDelimitedQuotedStrings splits space delimited quoted strings into
// a slice of strings
func splitSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings string) []string {
	// we use strings.SplitAfter to retain multiple whitespace in quoted
	// sections
	var fields []string
	if len(spaceDelimitedQuotedStrings) == 0 {
		return nil
	}
	words := strings.SplitAfter(spaceDelimitedQuotedStrings, ` `)
	var currField string
	for _, word := range words {
		if strings.HasPrefix(word, `"`) {
			// first word of quoted section
			trimmed := strings.TrimPrefix(word, `"`)
			currField = trimmed
		} else if len(currField) == 0 {
			// this block handles words not in a quoted section
			fields = append(fields, strings.TrimSuffix(word, ` `))
		} else if strings.HasSuffix(word, `" `) {
			// last word of quoted section
			trimmed := strings.TrimSuffix(word, `" `)
			currField += trimmed
			fields = append(fields, currField)
			currField = "" // triggers start of new field on next iteration
		} else if strings.HasSuffix(word, `"`) {
			// last word of quoted section and string
			trimmed := strings.TrimSuffix(word, `"`)
			currField += trimmed
			fields = append(fields, currField)
			currField = "" // triggers start of new field on next iteration
		} else {
			// intermediate word of quoted section
			currField += word
		}
	}
	return fields
}

// isValidTimeString tests whether a time string is valid
func isValidTimeString(timeString string) bool {
	if _, err := time.Parse(timeFormat, timeString); err != nil {
		return false
	}
	return true
}

// isValidURLString tests whether a URL string is valid
func isValidURLString(urlString string) bool {
	if _, err := url.Parse(urlString); err != nil {
		return false
	}
	return true
}

// alert is an unexported struct used internally to unmarshal a CAP alert
// message XML into
type alert struct {
	// TODO: need to add namespace support to distiguish CAP versions
	Identifier  string   `xml:"identifier"`
	Sender      string   `xml:"sender"`
	Sent        string   `xml:"sent"`
	Status      string   `xml:"status"`
	MsgType     string   `xml:"msgType"`
	Source      string   `xml:"source"`
	Scope       string   `xml:"scope"`
	Restriction string   `xml:"restriction"`
	Addresses   string   `xml:"addresses"`
	Codes       []string `xml:"code"`
	Note        string   `xml:"note"`
	References  string   `xml:"references"`
	Incidents   string   `xml:"incidents"`
	Infos       []struct {
		Language      string   `xml:"language"`
		Categories    []string `xml:"category"`
		Event         string   `xml:"event"`
		ResponseTypes []string `xml:"responseType"`
		Urgency       string   `xml:"urgency"`
		Severity      string   `xml:"severity"`
		Certainty     string   `xml:"certainty"`
		Audience      string   `xml:"audience"`
		EventCodes    []struct {
			ValueName string `xml:"valueName"`
			Value     string `xml:"value"`
		} `xml:"eventCode"`
		Effective   string `xml:"effective"`
		Onset       string `xml:"onset"`
		Expires     string `xml:"expires"`
		SenderName  string `xml:"senderName"`
		Headline    string `xml:"headline"`
		Description string `xml:"description"`
		Instruction string `xml:"instruction"`
		Web         string `xml:"web"`
		Contact     string `xml:"contact"`
		Parameters  []struct {
			ValueName string `xml:"valueName"`
			Value     string `xml:"value"`
		} `xml:"parameter"`
		Resources []struct {
			ResourceDesc string `xml:"resourceDesc"`
			MIMEType     string `xml:"mimeType"`
			Size         string `xml:"size"`
			URI          string `xml:"uri"`
			DerefURI     string `xml:"derefUri"`
			Digest       string `xml:"digest"`
		} `xml:"resource"`
		Areas []struct {
			AreaDesc string   `xml:"areaDesc"`
			Polygons []string `xml:"polygon"`
			Circles  []string `xml:"circle"`
			Geocodes []struct {
				ValueName string `xml:"valueName"`
				Value     string `xml:"value"`
			} `xml:"geocode"`
			Altitude string `xml:"altitude"`
			Ceiling  string `xml:"ceiling"`
		} `xml:"area"`
	} `xml:"info"`
}

// validate validates that the content of an alert struct conforms to the CAP
// 1.2 specification
// TODO: implement validation for CAP 1.1 and 1.0
func (a *alert) validate() error {
	var errStrs []string
	var missingElements []string

	if len(a.Identifier) == 0 {
		missingElements = append(missingElements, "alert.identifier")
	} else if strings.ContainsAny(a.Identifier, restrictedCharacters) {
		errStrs = append(errStrs, "alert.identifier contains contains one or more restricted characters")
	}

	if len(a.Sender) == 0 {
		missingElements = append(missingElements, "alert.sender")
	} else if strings.ContainsAny(a.Sender, restrictedCharacters) {
		errStrs = append(errStrs, "alert.sender contains contains one or more restricted characters")
	}

	if len(a.Sent) == 0 {
		missingElements = append(missingElements, "alert.sent")
	} else if !isValidTimeString(a.Sent) {
		errStrs = append(errStrs, "invalid alert.sent time")
	}

	if len(a.Status) == 0 {
		missingElements = append(missingElements, "alert.status")
	} else if _, ok := AlertStatuses[a.Status]; !ok {
		errStrs = append(errStrs, "invalid alert.status")
	}

	if len(a.MsgType) == 0 {
		missingElements = append(missingElements, "alert.msgType")
	} else if _, ok := AlertMsgTypes[a.MsgType]; !ok {
		errStrs = append(errStrs, "invalid alert.msgType")
	}

	if len(a.Scope) == 0 {
		missingElements = append(missingElements, "alert.scope")
	} else if _, ok := AlertScopes[a.Scope]; !ok {
		errStrs = append(errStrs, "invalid alert.scope")
	} else if a.Scope == "Restricted" && len(a.Restriction) == 0 {
		errStrs = append(errStrs, "if alert.scope is Restricted must have alert.restriction")
	} else if a.Scope != "Restricted" && len(a.Restriction) > 0 {
		errStrs = append(errStrs, "if alert.scope is not Restricted must not have alert.restriction")
	} else if a.Scope == "Private" && len(a.Addresses) == 0 {
		errStrs = append(errStrs, "if alert.scope is Private must have alert.addresses")
	}

	if len(a.Addresses) > 0 {
		if !isValidAddressesString(a.Addresses) {
			errStrs = append(errStrs, "invalid alert.addresses")
		}
	}

	if len(a.References) > 0 {
		if !isValidReferencesString(a.References) {
			errStrs = append(errStrs, "invalid alert.info[%d].area[%d].circle[%d]")
		}
	}

	if len(a.Incidents) > 0 {
		if !isValidIncidentsString(a.Incidents) {
			errStrs = append(errStrs, "invalid alert.incidents")
		}
	}

	for i, info := range a.Infos {
		if len(info.Categories) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].category", i))
		} else {
			for j, cat := range info.Categories {
				if _, ok := AlertInfoCategories[cat]; !ok {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info.category[%d]", j))
				}
			}
		}

		if len(info.Event) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].event", i))
		}

		for i, respType := range info.ResponseTypes {
			if _, ok := AlertInfoResponseTypes[respType]; !ok {
				errStrs = append(errStrs, fmt.Sprintf("invalid alert.info.responseType[%d]", i))
			}
		}

		if len(info.Urgency) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].urgency", i))
		} else if _, ok := AlertInfoUrgencies[info.Urgency]; !ok {
			errStrs = append(errStrs, "invalid alert.info.urgency")
		}

		if len(info.Severity) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].severity", i))
		} else if _, ok := AlertInfoSeverities[info.Severity]; !ok {
			errStrs = append(errStrs, "invalid alert.info.severity")
		}

		if len(info.Certainty) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].certainty", i))
		} else if _, ok := AlertInfoCertainties[info.Certainty]; !ok {
			errStrs = append(errStrs, "invalid alert.info.certainty")
		}

		if len(info.Effective) > 0 && !isValidTimeString(info.Effective) {
			errStrs = append(errStrs, "invalid alert.info.effective time")
		}

		if len(info.Onset) > 0 && !isValidTimeString(info.Onset) {
			errStrs = append(errStrs, "invalid alert.info.onset time")
		}

		if len(info.Expires) > 0 && !isValidTimeString(info.Expires) {
			errStrs = append(errStrs, "invalid alert.info.expires time")
		}

		if len(info.Web) > 0 && !isValidURLString(info.Web) {
			errStrs = append(errStrs, "invalid alert.info.web URL")
		}

		for j, resource := range info.Resources {
			if len(resource.ResourceDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].category", i, j))
			}

			if len(resource.MIMEType) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].mimeType", i, j))
			}

			if len(resource.Size) > 0 {
				if size, err := strconv.Atoi(resource.Size); err != nil || size < 0 {
					errStrs = append(errStrs, "invalid alert.info.resource.size")
				}
			}

			if len(resource.URI) > 0 && !isValidURLString(resource.URI) {
				errStrs = append(errStrs, "invalid alert.info.resource.uri URL")
			}

			if !(len(resource.URI) > 0 || len(resource.DerefURI) > 0) {
				errStrs = append(errStrs, "invalid alert.info.resource.uri URL")
			}
		}

		for j, area := range info.Areas {
			if len(area.AreaDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].area[%d] must have uri or derefUri", i, j))
			}
			for _, p := range area.Polygons {
				if !isValidPolygonString(p) {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].area[%d].polygon[%d]", i, j, p))
				}
			}
			for _, c := range area.Circles {
				if !isValidCircleString(c) {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].area[%d].circle[%d]", i, j, c))
				}
			}
		}
	}

	if len(missingElements) > 0 {
		errStrs = append(errStrs, fmt.Sprintf("missing elements ", strings.Join(missingElements, ", ")))
	}
	if len(errStrs) > 0 {
		return errors.New(strings.Join(errStrs, "; "))
	}

	return nil
}

// convert converts an unexported `alert` to an exported `Alert`
func (a *alert) convert() (*Alert, error) {
	var ret Alert // the Alert to be returned
	var err error

	ret.Identifier = a.Identifier
	ret.Sender = a.Sender
	if len(a.Sent) > 0 {
		if ret.Sent, err = time.Parse(timeFormat, a.Sent); err != nil {
			return nil, err
		}
	}
	ret.Status = a.Status
	ret.MsgType = a.MsgType
	ret.Source = a.Source
	ret.Scope = a.Scope
	ret.Restriction = a.Restriction
	if len(a.Addresses) > 0 {
		ret.Addresses = parseAddressesString(a.Addresses)
	}
	ret.Codes = a.Codes
	ret.Note = a.Note
	if len(a.References) > 0 {
		if ret.References, err = parseReferencesString(a.References); err != nil {
			return nil, err
		}
	}
	if len(a.Incidents) > 0 {
		ret.Incidents = parseIncidentsString(a.Incidents)
	}

	for _, aInfo := range a.Infos {
		var retInfo Info

		// en-US is assumed if no Language is defined
		if len(aInfo.Language) == 0 {
			retInfo.Language = "en-US"
		} else {
			retInfo.Language = aInfo.Language
		}
		retInfo.Categories = aInfo.Categories
		retInfo.Event = aInfo.Event
		retInfo.ResponseTypes = aInfo.ResponseTypes
		retInfo.Urgency = aInfo.Urgency
		retInfo.Severity = aInfo.Severity
		retInfo.Certainty = aInfo.Certainty
		retInfo.Audience = aInfo.Audience
		for _, ec := range aInfo.EventCodes {
			retInfo.EventCodes = append(retInfo.EventCodes, NamedValue{ValueName: ec.ValueName, Value: ec.Value})
		}
		if len(aInfo.Effective) > 0 {
			if retInfo.Effective, err = time.Parse(timeFormat, aInfo.Effective); err != nil {
				return nil, err
			}
		}
		if len(aInfo.Onset) > 0 {
			if retInfo.Onset, err = time.Parse(timeFormat, aInfo.Onset); err != nil {
				return nil, err
			}
		}
		if len(aInfo.Expires) > 0 {
			if retInfo.Expires, err = time.Parse(timeFormat, aInfo.Expires); err != nil {
				return nil, err
			}
		}
		retInfo.SenderName = aInfo.SenderName
		retInfo.Headline = aInfo.Headline
		retInfo.Description = aInfo.Description
		retInfo.Instruction = aInfo.Instruction
		if len(aInfo.Web) > 0 {
			if retInfo.Web, err = url.Parse(aInfo.Web); err != nil {
				return nil, err
			}
		}
		retInfo.Contact = aInfo.Contact
		for _, p := range aInfo.Parameters {
			retInfo.Parameters = append(retInfo.Parameters, NamedValue{ValueName: p.ValueName, Value: p.Value})
		}

		for _, aiResource := range aInfo.Resources {
			var retResource Resource

			retResource.ResourceDesc = aiResource.ResourceDesc
			retResource.MIMEType = aiResource.MIMEType
			if len(aiResource.Size) > 0 {
				if retResource.Size, err = strconv.Atoi(aiResource.Size); err != nil {
					return nil, err
				}
			}
			if len(aiResource.URI) > 0 {
				if retResource.URI, err = url.Parse(aiResource.URI); err != nil {
					return nil, err
				}
			}
			retResource.DerefURI = aiResource.Digest
			retResource.Digest = aiResource.Digest

			retInfo.Resources = append(retInfo.Resources, retResource)
		}

		for _, aiArea := range aInfo.Areas {
			var retArea Area

			retArea.AreaDesc = aiArea.AreaDesc
			for _, p := range aiArea.Polygons {
				if parsed, err := parsePolygonString(p); err != nil {
					return nil, err
				} else {
					retArea.Polygons = append(retArea.Polygons, parsed)
				}
			}
			for _, c := range aiArea.Circles {
				if parsed, err := parseCircleString(c); err != nil {
					return nil, err
				} else {
					retArea.Circles = append(retArea.Circles, parsed)
				}
			}
			for _, g := range aiArea.Geocodes {
				retArea.Geocodes = append(retArea.Geocodes, NamedValue{ValueName: g.ValueName, Value: g.Value})
			}
			retArea.Altitude = aiArea.Altitude
			retArea.Ceiling = aiArea.Ceiling

			retInfo.Areas = append(retInfo.Areas, retArea)
		}

		ret.Infos = append(ret.Infos, retInfo)
	}

	return &ret, nil
}

// ProcessAlertMsgXML takes an XML CAP alert message and returns an Alert struct
func ProcessAlertMsgXML(alertMsgXML []byte) (*Alert, error) {
	raw := &alert{}
	if err := xml.Unmarshal(alertMsgXML, raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling alert message XML: %s", err)
	}
	if err := raw.validate(); err != nil {
		return nil, fmt.Errorf("error(s) validating alert: %s", err)
	}
	processed, err := raw.convert()
	if err != nil {
		return nil, fmt.Errorf("error(s) converting alert: %s", err)
	}

	return processed, nil
}
