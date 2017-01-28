package cap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
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

// alert is an unexported struct used to unmarshal a CAP alert message into
type alert struct {
	// TODO: need to add namespace support to distiguish CAP versions
	Identifier  string   `xml:"indentifier"`
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
	Info        []struct {
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

// Alert represents a parsed and validated CAP alert
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
	Info        []struct {
		Language      string
		Categories    []string
		Event         string
		ResponseTypes []string
		Urgency       string
		Severity      string
		Certainty     string
		Audience      string
		EventCodes    []struct {
			ValueName string
			Value     string
		}
		Effective   time.Time
		Onset       time.Time
		Expires     time.Time
		SenderName  string
		Headline    string
		Description string
		Instruction string
		Web         url.URL
		Contact     string
		Parameters  []struct {
			ValueName string
			Value     string
		}
		Resources []struct {
			ResourceDesc string
			MIMEType     string
			Size         int // approximate size in bytes
			URI          url.URL
			DerefURI     string // base-64 encoded binary
			Digest       string // SHA-1 hash
		}
		Areas []struct {
			AreaDesc string
			// because golang has no built in support for decimals, these values
			// are being left as `string` so the caller can handle as necessary.
			// The coordinate system used is WGS 84.
			Polygons []Polygon
			Circles  []Circle
			Geocodes []struct {
				ValueName string
				Value     string
			}
			Altitude string // decimal feet above mean sea level
			Ceiling  string // decimal feet above mean sea level
		}
	}
}

// isValidTimeString tests whether a time string is valid
func isValidTimeString(timeString string) bool {
	// not returning the error might be improper here, but I like that this
	// returns a simple boolean value.
	if _, err := time.Parse(timeFormat, timeString); err != nil {
		return false
	}
	return true
}

// isValidURLString tests whether a URL is valid
func isValidURLString(urlString string) bool {
	// not returning the error might be improper here, but I like that this
	// returns a simple boolean value.
	if _, err := url.Parse(urlString); err != nil {
		return false
	}
	return true
}

// Point represents a WGS 84 coordinate point on the earth
type Point struct {
	Latitude  string
	Longitude string
}

// Polygon
type Polygon []Point

// parsePolygon parses a polygon string
func parsePolygonString(polygonString string) (Polygon, error) {
	// 38.47,-120.14 38.34,-119.95 38.52,-119.74 38.62,-119.89 38.47,-120.14
	// TODO: test validity of numbers
	pointStrings := strings.Fields(polygonString)
	if len(pointStrings) < 4 {
		return nil, errors.New("a polygon must contain al least four points")
	}
	var polygon Polygon
	for _, ps := range pointStrings {
		vals := strings.Split(ps, ",")
		if len(vals) != 2 {
			return nil, errors.New("point must contain exactly two values")
		}
		polygon = append(polygon, Point{Latitude: vals[0], Longitude: vals[1]})

	}
	if polygon[0] != polygon[len(polygon)-1] {
		return nil, errors.New("first and last points must be equal")
	}
	return polygon, nil
}

// isValidPolygon tests whether a polygon string is valid
func isValidPolygonString(polygonString string) bool {
	if _, err := parsePolygonString(polygonString); err != nil {
		return false
	}
	return true
}

// Circle
type Circle struct {
	Point  Point
	Radius string
}

// parseCircleString
func parseCircleString(circleString string) (Circle, error) {
	// 32.9525,-115.5527 0
	// TODO: test validity of numbers
	pointRadiusStrings := strings.Fields(circleString)
	if len(pointRadiusStrings) != 2 {
		return Circle{}, errors.New("a circle must contain a central point and a radius")
	}
	pointStrings := strings.Split(pointRadiusStrings[0], ",")
	if len(pointStrings) != 2 {
		return Circle{}, errors.New("central point must contain exactly two values")
	}
	radiusString := pointRadiusStrings[1]
	circle := Circle{Point: Point{Latitude: pointStrings[0], Longitude: pointStrings[1]}, Radius: radiusString}
	return circle, nil
}

// isValidCircle tests whether a circle string is valid
func isValidCircleString(circleString string) bool {
	if _, err := parseCircleString(circleString); err != nil {
		return false
	}
	return true
}

// Reference
type Reference struct {
	Sender     string
	Identifier string
	Sent       time.Time
}

func parseReferencesString(referencesString string) ([]Reference, error) {
	refStrings := strings.Fields(referencesString)
	var refs []Reference
	for _, s := range refStrings {
		parts := strings.Split(s, ",")
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

func isValidReferencesString(referenceString string) bool {
	if _, err := parseCircleString(referenceString); err != nil {
		return false
	}
	return true
}

func splitSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings string) []string {
	// we use strings.SplitAfter to retain multiple whitespace in quoted
	// sections
	words := strings.SplitAfter(spaceDelimitedQuotedStrings, " ")
	var fields []string
	var currField string
	for _, word := range words {
		if strings.HasPrefix(word, `"`) {
			// first word of quoted section
			trimmed := strings.TrimPrefix(word, `"`)
			currField = trimmed
		} else if len(currField) == 0 {
			// this block handles words not in a quoted section
			fields = append(fields, strings.TrimSuffix(word, " "))
		} else if strings.HasSuffix(word, `" `) {
			// last word of quoted section
			trimmed := strings.TrimSuffix(word, `" `)
			currField += trimmed
			fields = append(fields, currField)
			currField = ""
		} else {
			// intermediate word of quoted section
			currField += word
		}
	}
	return fields
}

func parseAddressesString(addressesString string) []string {
	return splitSpaceDelimitedQuotedStrings(addressesString)
}

func parseIncidentsString(indidentsString string) []string {
	return splitSpaceDelimitedQuotedStrings(indidentsString)
}

func (a *alert) validate() error {
	var errorStrings []string
	var missingElements []string

	if len(a.Identifier) == 0 {
		missingElements = append(missingElements, "alert.identifier")
	} else if strings.ContainsAny(a.Identifier, restrictedCharacters) {
		errorStrings = append(errorStrings, "alert.identifier contains contains one or more restricted characters")
	}

	if len(a.Sender) == 0 {
		missingElements = append(missingElements, "alert.sender")
	} else if strings.ContainsAny(a.Sender, restrictedCharacters) {
		errorStrings = append(errorStrings, "alert.sender contains contains one or more restricted characters")
	}

	if len(a.Sent) == 0 {
		missingElements = append(missingElements, "alert.sent")
	} else if !isValidTimeString(a.Sent) {
		errorStrings = append(errorStrings, "invalid alert.sent time")
	}

	if len(a.Status) == 0 {
		missingElements = append(missingElements, "alert.status")
	} else if _, ok := AlertStatuses[a.Status]; !ok {
		errorStrings = append(errorStrings, "invalid alert.status")
	}

	if len(a.MsgType) == 0 {
		missingElements = append(missingElements, "alert.msgType")
	} else if _, ok := AlertMsgTypes[a.MsgType]; !ok {
		errorStrings = append(errorStrings, "invalid alert.msgType")
	}

	if len(a.Scope) == 0 {
		missingElements = append(missingElements, "alert.scope")
	} else if _, ok := AlertScopes[a.Scope]; !ok {
		errorStrings = append(errorStrings, "invalid alert.scope")
	} else if a.Scope == "Restricted" && len(a.Restriction) == 0 {
		errorStrings = append(errorStrings, "if alert.scope is Restricted must have alert.restriction")
	} else if a.Scope != "Restricted" && len(a.Restriction) > 0 {
		errorStrings = append(errorStrings, "if alert.scope is not Restricted must not have alert.restriction")
	} else if a.Scope == "Private" && len(a.Addresses) == 0 {
		errorStrings = append(errorStrings, "if alert.scope is Private must have alert.addresses")
	}

	if len(a.Addresses) > 0 {
		// "" and whitespace
	}

	if len(a.References) > 0 {
		if !isValidReferencesString(a.References) {
			errorStrings = append(errorStrings, "invalid alert.info[%d].area[%d].circle[%d]")
		}
	}

	if len(a.Incidents) > 0 {
		// "" and whitespace
	}

	for i, info := range a.Info {
		if len(info.Categories) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].category", i))
		} else {
			for j, cat := range info.Categories {
				if _, ok := AlertInfoCategories[cat]; !ok {
					errorStrings = append(errorStrings, fmt.Sprintf("invalid alert.info.category[%d]", j))
				}
			}
		}

		if len(info.Event) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].event", i))
		}

		for i, respType := range info.ResponseTypes {
			if _, ok := AlertInfoResponseTypes[respType]; !ok {
				errorStrings = append(errorStrings, fmt.Sprintf("invalid alert.info.responseType[%d]", i))
			}
		}

		if len(info.Urgency) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].urgency", i))
		} else if _, ok := AlertInfoUrgencies[info.Urgency]; !ok {
			errorStrings = append(errorStrings, "invalid alert.info.urgency")
		}

		if len(info.Severity) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].severity", i))
		} else if _, ok := AlertInfoSeverities[info.Severity]; !ok {
			errorStrings = append(errorStrings, "invalid alert.info.severity")
		}

		if len(info.Certainty) == 0 {
			missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].certainty", i))
		} else if _, ok := AlertInfoCertainties[info.Certainty]; !ok {
			errorStrings = append(errorStrings, "invalid alert.info.certainty")
		}

		if len(info.Effective) > 0 && !isValidTimeString(info.Effective) {
			errorStrings = append(errorStrings, "invalid alert.info.effective time")
		}

		if len(info.Onset) > 0 && !isValidTimeString(info.Onset) {
			errorStrings = append(errorStrings, "invalid alert.info.onset time")
		}

		if len(info.Expires) > 0 && !isValidTimeString(info.Expires) {
			errorStrings = append(errorStrings, "invalid alert.info.expires time")
		}

		if len(info.Web) > 0 && !isValidURLString(info.Web) {
			errorStrings = append(errorStrings, "invalid alert.info.web URL")
		}

		for j, resource := range info.Resources {
			if len(resource.ResourceDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].category", i, j))
			}

			if len(resource.MIMEType) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].resource[%d].mimeType", i, j))
			}

			if len(resource.URI) > 0 && !isValidURLString(resource.URI) {
				errorStrings = append(errorStrings, "invalid alert.info.resource.uri URL")
			}

			if !(len(resource.URI) > 0 || len(resource.DerefURI) > 0) {
				errorStrings = append(errorStrings, "invalid alert.info.resource.uri URL")
			}
		}

		for j, area := range info.Areas {
			if len(area.AreaDesc) == 0 {
				missingElements = append(missingElements, fmt.Sprintf("alert.info[%d].area[%d] must have uri or derefUri", i, j))
			}
			for _, p := range area.Polygons {
				if !isValidPolygonString(p) {
					errorStrings = append(errorStrings, fmt.Sprintf("invalid alert.info[%d].area[%d].polygon[%d]", i, j, p))
				}
			}
			for _, c := range area.Circles {
				if !isValidCircleString(c) {
					errorStrings = append(errorStrings, fmt.Sprintf("invalid alert.info[%d].area[%d].circle[%d]", i, j, c))
				}
			}
		}
	}

	if len(missingElements) > 0 {
		errorStrings = append(errorStrings, fmt.Sprintf("missing elements ", strings.Join(missingElements, ", ")))
	}
	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, "; "))
	}

	return nil
}

func (a *alert) convert() (*Alert, error) {
	return nil, nil
}

// ProcessAlertMessageXML takes an XML CAP alert message and returns an Alert struct
func ProcessAlertMessage(alertMessageXML []byte) (*Alert, error) {
	raw := &alert{}
	if err := xml.Unmarshal(alertMessageXML, raw); err != nil {
		return nil, fmt.Errorf("error unmarshalling alert message XML: %s", err)
	}
	if err := raw.validate(); err != nil {
		return nil, fmt.Errorf("error(s) validating alert: %s", err)
	}
	processedAlert, err := raw.convert()
	if err != nil {
		return nil, fmt.Errorf("error(s) converting alert: %s", err)
	}

	return processedAlert, nil
}
