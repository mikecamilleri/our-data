package ouralerts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ValidateMessageXML validates an XML CAP alert message
func ValidateMessageXML(messageXML []byte) error {
	a, err := unmarshallMessageXML(messageXML)
	if err != nil {
		return fmt.Errorf("error unmarshalling message XML: %s", err)
	}
	if err := a.validate(); err != nil {
		return fmt.Errorf("error validating alert message: %s", err)
	}
	return nil
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

// isValidAddressesString tests whether an addresses string is valid
func isValidAddressesString(addressesString string) bool {
	return isValidSpaceDelimitedQuotedStrings(addressesString)
}

// isValidIncidentsString tests whether an incidents string is valid
func isValidIncidentsString(incidentsString string) bool {
	return isValidSpaceDelimitedQuotedStrings(incidentsString)
}

// isValidSpaceDelimitedQuotedStrings tests whether a space delimited quoted
// string is valid
func isValidSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings string) bool {
	if _, err := splitSpaceDelimitedQuotedStrings(spaceDelimitedQuotedStrings); err != nil {
		return false
	}
	return true
}

// isValidTimeString tests whether a time string is valid
func isValidTimeString(timeString string) bool {
	if _, err := parseTimeString(timeString); err != nil {
		return false
	}
	return true
}

// isValidURLString tests whether a URL string is valid
func isValidURLString(urlString string) bool {
	if _, err := parseURLString(urlString); err != nil {
		return false
	}
	return true
}

// isValidReferencesString tests whether a references string is valid
func isValidReferencesString(referencesString string) bool {
	if _, err := parseReferencesString(referencesString); err != nil {
		return false
	}
	return true
}

// isValidPolygonString tests whether a polygon string is valid
func isValidPolygonString(polygonString string) bool {
	if _, err := parsePolygonString(polygonString); err != nil {
		return false
	}
	return true
}

// isValidCircleString tests whether a circle string is valid
func isValidCircleString(circleString string) bool {
	if _, err := parseCircleString(circleString); err != nil {
		return false
	}
	return true
}
