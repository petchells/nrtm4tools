package rpsl

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

// ErrCannotParseRPSL is returned when the parser can't make sense of the RPSL
var ErrCannotParseRPSL = errors.New("invalid RPSL")

var trimHashChar = regexp.MustCompile("[^#]*")

// Rpsl is a data structure representing an NRTM object
type Rpsl struct {
	PrimaryKey string
	Source     string
	ObjectType string
	Payload    string
}

// ParseString parses a string and returns it as an RPSL object
func ParseString(str string) (Rpsl, error) {
	return parseString(str)
}

func parseString(str string) (Rpsl, error) {
	lines := strings.Split(str, "\n")
	var source, objectType string
	var primaryKey []string
	for _, rawLine := range lines {
		line := stripComment(rawLine)
		if len(line) == 0 {
			continue
		}
		if len(objectType) == 0 {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				log.Println("Cannot determine ObjectType")
				return Rpsl{}, ErrCannotParseRPSL
			}
			objectType = trimToUpper(parts[0])
			if isPrimaryKeyAttribute(objectType, objectType) {
				primaryKey = append(primaryKey, trimToUpper(parts[1]))
			}
		} else {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			attributeName := trimToLower(parts[0])
			if attributeName == "source" {
				source = trimToUpper(parts[1])
			} else if isPrimaryKeyAttribute(objectType, attributeName) {
				primaryKey = append(primaryKey, trimToUpper(parts[1]))
			}
		}
	}
	rpsl := Rpsl{PrimaryKey: strings.Join(primaryKey, ""), Source: source, ObjectType: objectType, Payload: str}
	if len(primaryKey) == 0 || len(source) == 0 || len(objectType) == 0 || ((objectType == "route" || objectType == "route6") && len(primaryKey) != 2) {
		return rpsl, ErrCannotParseRPSL
	}
	return rpsl, nil
}

func trimToLower(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

func trimToUpper(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

func isPrimaryKeyAttribute(objectType string, attributeName string) bool {
	if objectType == "PERSON" || objectType == "ROLE" {
		return attributeName == "nic-hdl"
	}
	if objectType == "ROUTE" || objectType == "ROUTE6" {
		return strings.EqualFold(attributeName, objectType) || attributeName == "origin"
	}
	return strings.EqualFold(attributeName, objectType)
}

func stripComment(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if ch != '#' {
			b.WriteRune(ch)
		} else {
			break
		}
	}
	return strings.TrimSpace(b.String())
}
