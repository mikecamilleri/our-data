/*
	Package ouralerts implements the ability to parse and validate OASIS Common
	Alerting Protocol alert messages
*/

package ouralerts

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	timeFormat           = "2006-01-02T15:04:05-07:00"
	restrictedCharacters = " ,<&"
)

var (
	XMLNamespaces = map[string]string{
		"1.2": "urn:oasis:names:tc:emergency:cap:1.2",
		"1.1": "urn:oasis:names:tc:emergency:cap:1.1",
		"1.0": "urn:oasis:names:tc:emergency:cap:1.0",
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
	EventCodes    url.Values
	Effective     time.Time
	Onset         time.Time
	Expires       time.Time
	SenderName    string
	Headline      string
	Description   string
	Instruction   string
	Web           *url.URL
	Contact       string
	Parameters    url.Values
	Resources     []Resource
	Areas         []Area
}

// Resource
type Resource struct {
	ResourceDesc string
	MIMEType     string
	Size         int64 // approximate size in bytes
	URI          *url.URL
	DerefURI     string // base-64 encoded binary
	Digest       string // SHA-1 hash
}

// Area
type Area struct {
	AreaDesc string
	Polygons []Polygon
	Circles  []Circle
	Geocodes url.Values
	Altitude float64 // feet above mean sea level
	Ceiling  float64 // feet above mean sea level
}

// Reference holds a reference to another alert
type Reference struct {
	Sender     string
	Identifier string
	Sent       time.Time
}

// Polygon defines a polygonal area
type Polygon []Point

// Circle defines a circular area
type Circle struct {
	Point  Point
	Radius float64 // kilometers
}

// Point defines a WGS 84 coordinate point on the earth
type Point struct {
	Latitude  float64
	Longitude float64
}

// ValidateMessageXML validates an XML CAP alert message
func ValidateMessageXML(messageXML []byte) error {
	a := &alert{}
	if err := xml.Unmarshal(messageXML, a); err != nil {
		return fmt.Errorf("error unmarshalling alert message XML: %s", err)
	}
	if err := a.validate(); err != nil {
		return fmt.Errorf("error(s) validating alert message: %s", err)
	}
	return nil
}

// ProcessMessageXML takes an XML CAP alert message and returns an Alert struct.
// An effort is made to process invalid messages. If validity is a concern, the
// XML message should be validated separately with ValidateMessageXML.
func ProcessMessageXML(messageXML []byte) (*Alert, error) {
	// creating our own decoder is required since the character set may not be UTF-8
	a := &alert{}
	reader := bytes.NewReader(messageXML)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	if err := decoder.Decode(a); err != nil {
		return nil, fmt.Errorf("error unmarshalling alert message XML: %s", err)
	}

	converted, err := a.convert()
	if err != nil {
		return nil, fmt.Errorf("error(s) converting message to exported Alert struct: %s", err)
	}
	return converted, nil
}

// alert is used internally for unmarshalling a CAP alert message
type alert struct {
	XMLNS       string   `xml:"xmlns,attr"`
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
func (a *alert) validate() error {
	var errStrs []string
	var missingElements []string

	if a.XMLNS != XMLNamespaces["1.2"] {
		errStrs = append(errStrs, fmt.Sprintf("XML namepace is %s, this validater is designed for %s", a.XMLNS, XMLNamespaces["1.2"]))
	}

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

	if len(a.Addresses) > 0 && !isValidAddressesString(a.Addresses) {
		errStrs = append(errStrs, "invalid alert.addresses")
	}

	if len(a.References) > 0 && !isValidReferencesString(a.References) {
		errStrs = append(errStrs, "invalid alert.references")
	}

	if len(a.Incidents) > 0 && !isValidIncidentsString(a.Incidents) {
		errStrs = append(errStrs, "invalid alert.incidents")
	}

	for i, info := range a.Infos {
		if len(info.Categories) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].category", i))
		} else {
			for j, cat := range info.Categories {
				if _, ok := AlertInfoCategories[cat]; !ok {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].category[%d]", i, j))
				}
			}
		}

		if len(info.Event) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].event", i))
		}

		for j, respType := range info.ResponseTypes {
			if _, ok := AlertInfoResponseTypes[respType]; !ok {
				errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].responseType[%d]", i, j))
			}
		}

		if len(info.Urgency) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].urgency", i))
		} else if _, ok := AlertInfoUrgencies[info.Urgency]; !ok {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].urgency", i))
		}

		if len(info.Severity) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].severity", i))
		} else if _, ok := AlertInfoSeverities[info.Severity]; !ok {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].severity", i))
		}

		if len(info.Certainty) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].certainty", i))
		} else if _, ok := AlertInfoCertainties[info.Certainty]; !ok {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].certainty", i))
		}

		if len(info.Effective) > 0 && !isValidTimeString(info.Effective) {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].effective time", i))
		}

		if len(info.Onset) > 0 && !isValidTimeString(info.Onset) {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].onset time", i))
		}

		if len(info.Expires) > 0 && !isValidTimeString(info.Expires) {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].expires time", i))
		}

		if len(info.Web) > 0 && !isValidURLString(info.Web) {
			errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].web URL", i))
		}

		for j, resource := range info.Resources {
			if len(resource.ResourceDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].category", i, j))
			}

			if len(resource.MIMEType) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].mimeType", i, j))
			}

			if len(resource.Size) > 0 {
				if _, err := strconv.ParseInt(resource.Size, 10, 64); err != nil {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].resource[%d].size", i, j))
				}
			}

			if len(resource.URI) > 0 && !isValidURLString(resource.URI) {
				errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].resource[%d].uri URL", i, j))
			}
		}

		for j, area := range info.Areas {
			if len(area.AreaDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].area[%d].areaDesc", i, j))
			}
			for k, p := range area.Polygons {
				if !isValidPolygonString(p) {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].area[%d].polygon[%d]", i, j, k))
				}
			}
			for k, c := range area.Circles {
				if !isValidCircleString(c) {
					errStrs = append(errStrs, fmt.Sprintf("invalid alert.info[%d].area[%d].circle[%d]", i, j, k))
				}
			}
		}
	}

	if len(missingElements) > 0 {
		errStrs = append(errStrs, fmt.Sprintf("missing elements: %s", strings.Join(missingElements, ", ")))
	}

	// TODO: consider whether returning errStrs wold be more useful
	// TODO: consider joining with newlines
	if len(errStrs) > 0 {
		return errors.New(strings.Join(errStrs, "; "))
	}

	return nil
}

// convert converts an unexported `alert` to an exported `Alert`. An effort is
// made to convert invalid messages and fields.
func (a *alert) convert() (*Alert, error) {
	var ret Alert // the Alert to be returned

	ret.Identifier = a.Identifier
	ret.Sender = a.Sender
	ret.Sent, _ = parseTimeString(a.Sent)
	ret.Status = a.Status
	ret.MsgType = a.MsgType
	ret.Source = a.Source
	ret.Scope = a.Scope
	ret.Restriction = a.Restriction
	ret.Addresses, _ = parseAddressesString(a.Addresses)
	ret.Codes = removeEmptyStringsFromSlice(a.Codes)
	ret.Note = a.Note
	ret.References, _ = parseReferencesString(a.References)
	ret.Incidents, _ = parseIncidentsString(a.Incidents)

	for _, aInfo := range a.Infos {
		var retInfo Info

		// per the spec en-US is assumed if no Language is defined
		if len(aInfo.Language) == 0 {
			retInfo.Language = "en-US"
		} else {
			retInfo.Language = aInfo.Language
		}

		retInfo.Categories = removeEmptyStringsFromSlice(aInfo.Categories)
		retInfo.Event = aInfo.Event
		retInfo.ResponseTypes = removeEmptyStringsFromSlice(aInfo.ResponseTypes)
		retInfo.Urgency = aInfo.Urgency
		retInfo.Severity = aInfo.Severity
		retInfo.Certainty = aInfo.Certainty
		retInfo.Audience = aInfo.Audience

		for _, ec := range aInfo.EventCodes {
			if len(ec.ValueName) > 0 {
				if retInfo.EventCodes == nil {
					retInfo.EventCodes = make(url.Values)
				}
				retInfo.EventCodes.Add(ec.ValueName, ec.Value)
			}
		}

		retInfo.Effective, _ = parseTimeString(aInfo.Effective)
		retInfo.Onset, _ = parseTimeString(aInfo.Onset)
		retInfo.Expires, _ = parseTimeString(aInfo.Expires)
		retInfo.SenderName = aInfo.SenderName
		retInfo.Headline = aInfo.Headline
		retInfo.Description = aInfo.Description
		retInfo.Instruction = aInfo.Instruction
		retInfo.Web, _ = parseURLString(aInfo.Web)
		retInfo.Contact = aInfo.Contact

		for _, p := range aInfo.Parameters {
			if len(p.ValueName) > 0 {
				if retInfo.Parameters == nil {
					retInfo.Parameters = make(url.Values)
				}
				retInfo.Parameters.Add(p.ValueName, p.Value)
			}
		}

		for _, aiResource := range aInfo.Resources {
			var retResource Resource

			retResource.ResourceDesc = aiResource.ResourceDesc
			retResource.MIMEType = aiResource.MIMEType
			retResource.Size, _ = strconv.ParseInt(aiResource.Size, 10, 64)
			retResource.URI, _ = parseURLString(aiResource.URI)
			retResource.DerefURI = aiResource.DerefURI
			retResource.Digest = aiResource.Digest

			retInfo.Resources = append(retInfo.Resources, retResource)
		}

		for _, aiArea := range aInfo.Areas {
			var retArea Area

			retArea.AreaDesc = aiArea.AreaDesc
			for _, p := range aiArea.Polygons {
				if parsed, err := parsePolygonString(p); err == nil {
					retArea.Polygons = append(retArea.Polygons, parsed)
				}
			}
			for _, c := range aiArea.Circles {
				if parsed, err := parseCircleString(c); err == nil {
					retArea.Circles = append(retArea.Circles, parsed)
				}
			}

			for _, g := range aiArea.Geocodes {
				if len(g.ValueName) > 0 {
					if retArea.Geocodes == nil {
						retArea.Geocodes = make(url.Values)
					}
					retArea.Geocodes.Add(g.ValueName, g.Value)
				}
			}

			retArea.Altitude, _ = strconv.ParseFloat(aiArea.Altitude, 64)
			retArea.Ceiling, _ = strconv.ParseFloat(aiArea.Ceiling, 64)

			retInfo.Areas = append(retInfo.Areas, retArea)
		}

		ret.Infos = append(ret.Infos, retInfo)
	}

	return &ret, nil
}

// removeEmptyStringsFromSlice returns a slice of strings that is the imput
// slice with the empty values removed, or nil if empty
func removeEmptyStringsFromSlice(sliceOfStrings []string) []string {
	var out []string
	for _, s := range sliceOfStrings {
		if len(s) > 0 {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// parseAddressesString parses an addresses string
func parseAddressesString(addressesString string) ([]string, error) {
	if !isValidAddressesString(addressesString) {
		return nil, errors.New("error parsing addresses string")
	}
	return splitSpaceDelimitedQuotedStrings(addressesString)
}

// isValidAddressesString tests whether an addresses string is valid
func isValidAddressesString(addressesString string) bool {
	return isValidSpaceDelimitedQuotedStrings(addressesString)
}

// parseIncidentsString parases an incidents string
func parseIncidentsString(incidentsString string) ([]string, error) {
	if !isValidIncidentsString(incidentsString) {
		return nil, errors.New("error parsing incidents string")
	}
	return splitSpaceDelimitedQuotedStrings(incidentsString)
}

// isValidIncidentsString tests whether an incidents string is valid
func isValidIncidentsString(incidentsString string) bool {
	return isValidSpaceDelimitedQuotedStrings(incidentsString)
}

// splitSpaceDelimitedQuotedStrings splits space delimited quoted strings into
// a slice of strings
func splitSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings string) ([]string, error) {
	if !isValidSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings) {
		return nil, errors.New("error splitting space delimited quoted string")
	}
	var fields []string
	// we use strings.SplitAfter to retain multiple whitespace in quoted
	// sections
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
	return fields, nil
}

// isValidSpaceDelimitedQuotedStrings tests whether a space delimited quoted
// string is valid
func isValidSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings string) bool {
	if len(spaceDelimitedQuotedStrings) == 0 {
		return false
	}
	if strings.Count(spaceDelimitedQuotedStrings, `"`)%2 != 0 {
		return false
	}
	return true
}

// parseTimeString parses a time string
func parseTimeString(timeString string) (time.Time, error) {
	if !isValidTimeString(timeString) {
		return time.Time{}, errors.New("error parsing time string")
	}
	t, err := time.Parse(timeFormat, timeString)
	if err != nil {
		return time.Time{}, errors.New("error parsing time string")
	}
	return t, nil
}

// isValidTimeString tests whether a time string is valid
func isValidTimeString(timeString string) bool {
	if len(timeString) == 0 {
		return false
	}
	if _, err := time.Parse(timeFormat, timeString); err != nil {
		return false
	}
	return true
}

// parseURLString parses a URL string
func parseURLString(urlString string) (*url.URL, error) {
	if !isValidURLString(urlString) {
		return nil, errors.New("error parsing URL string")
	}
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, errors.New("error parsing URL string")
	}
	return url, nil
}

// isValidURLString tests whether a URL string is valid
func isValidURLString(urlString string) bool {
	if len(urlString) == 0 {
		return false
	}
	if _, err := url.Parse(urlString); err != nil {
		return false
	}
	return true
}

// parseReferencesString parses a references string
func parseReferencesString(referencesString string) ([]Reference, error) {
	if !isValidReferencesString(referencesString) {
		return nil, errors.New("error parsing references string")
	}
	refStrings := strings.Fields(referencesString)
	var refs []Reference
	for _, rs := range refStrings {
		r, err := parseSingleReferenceString(rs)
		if err != nil {
			continue
		}
		refs = append(refs, r)
	}
	if len(refs) == 0 {
		return nil, errors.New("error parsing references string")
	}
	return refs, nil
}

// isValidReferencesString tests whether a references string is valid
func isValidReferencesString(referencesString string) bool {
	if len(referencesString) == 0 {
		return false
	}
	refStrings := strings.Fields(referencesString)
	for _, rs := range refStrings {
		if _, err := parseSingleReferenceString(rs); err != nil {
			return false
		}
	}
	return true
}

// parseSingleReferenceString parses a single reference string
func parseSingleReferenceString(singleReferenceString string) (Reference, error) {
	parts := strings.Split(singleReferenceString, ",")
	if len(parts) != 3 {
		return Reference{}, errors.New("reference must contain three parts")
	}
	t, err := parseTimeString(parts[2])
	if err != nil {
		return Reference{}, errors.New("invalid time string")
	}
	return Reference{Sender: parts[0], Identifier: parts[1], Sent: t}, nil
}

// parsePolygonString parses a polygon string
func parsePolygonString(polygonString string) (Polygon, error) {
	if !isValidPolygonString(polygonString) {
		return Polygon{}, errors.New("error parsing polygon string")
	}
	var polygon Polygon
	pointStrings := strings.Fields(polygonString)
	for _, ps := range pointStrings {
		vals := strings.Split(ps, ",")
		var lat, lon float64
		var err error
		if lat, err = strconv.ParseFloat(vals[0], 64); err != nil {
			return Polygon{}, errors.New("error parsing polygon string")
		}
		if lon, err = strconv.ParseFloat(vals[1], 64); err != nil {
			return Polygon{}, errors.New("error parsing polygon string")
		}
		polygon = append(polygon, Point{lat, lon})
	}
	return polygon, nil
}

// isValidPolygonString tests whether a polygon string is valid
func isValidPolygonString(polygonString string) bool {
	if len(polygonString) == 0 {
		return false
	}
	pointStrings := strings.Fields(polygonString)
	// a polygon must contain at least four points
	if len(pointStrings) < 4 {
		return false
	}
	var polygon Polygon
	for _, ps := range pointStrings {
		vals := strings.Split(ps, ",")
		// a point must contain exactly two values
		if len(vals) != 2 {
			return false
		}
		var lat, lon float64
		var err error
		if lat, err = strconv.ParseFloat(vals[0], 64); err != nil {
			return false
		}
		if lon, err = strconv.ParseFloat(vals[1], 64); err != nil {
			return false
		}
		polygon = append(polygon, Point{lat, lon})
	}
	// first and last points must be equal
	if polygon[0] != polygon[len(polygon)-1] {
		return false
	}
	return true
}

// parseCircleString parses a circle string
func parseCircleString(circleString string) (Circle, error) {
	if !isValidCircleString(circleString) {
		return Circle{}, errors.New("error parsing circle string")
	}
	var circle Circle
	var lat, lon, rad float64
	var err error
	pointRadiusStrings := strings.Fields(circleString)
	if len(pointRadiusStrings) != 2 {
		return Circle{}, errors.New("error parsing circle string")
	}
	pVals := strings.Split(pointRadiusStrings[0], ",")
	if len(pVals) != 2 {
		return Circle{}, errors.New("error parsing circle string")
	}
	if lat, err = strconv.ParseFloat(pVals[0], 64); err != nil {
		return Circle{}, errors.New("error parsing circle string")
	}
	if lon, err = strconv.ParseFloat(pVals[1], 64); err != nil {
		return Circle{}, errors.New("error parsing circle string")
	}
	if rad, err = strconv.ParseFloat(pointRadiusStrings[1], 64); err != nil {
		return Circle{}, errors.New("error parsing circle string")
	}
	circle.Point = Point{Latitude: lat, Longitude: lon}
	circle.Radius = rad
	return circle, nil
}

// isValidCircleString tests whether a circle string is valid
func isValidCircleString(circleString string) bool {
	if len(circleString) == 0 {
		return false
	}
	pointRadiusStrings := strings.Fields(circleString)
	if len(pointRadiusStrings) != 2 {
		// a circle must contain a central point and a radius
		return false
	}
	pVals := strings.Split(pointRadiusStrings[0], ",")
	if len(pVals) != 2 {
		// central point must contain exactly two values
		return false
	}
	if _, err := strconv.ParseFloat(pVals[0], 64); err != nil {
		return false
	}
	if _, err := strconv.ParseFloat(pVals[1], 64); err != nil {
		return false
	}
	if _, err := strconv.ParseFloat(pointRadiusStrings[1], 64); err != nil {
		return false
	}
	return true
}
