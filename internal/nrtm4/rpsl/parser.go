package rpsl

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

var ErrCannotParseRPSL = errors.New("invalid RPSL")

var trimHashChar = regexp.MustCompile("[^#]*")

type Rpsl struct {
	PrimaryKey string
	Source     string
	ObjectType string
	Payload    string
}

func ParseString(str string) (Rpsl, error) {
	return parseString(str)
}

func parseString(str string) (Rpsl, error) {
	lines := strings.Split(str, "\n")
	var source, objectType string
	var primaryKey []string
	for _, rawLine := range lines {
		line := string(trimHashChar.Find([]byte(strings.TrimSpace(rawLine))))
		if len(line) == 0 {
			continue
		}
		if len(objectType) == 0 {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				log.Println("Cannot determine ObjectType")
				return Rpsl{}, ErrCannotParseRPSL
			}
			objectType = trimToLower(parts[0])
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
	if objectType == "person" || objectType == "role" {
		return attributeName == "nic-hdl"
	}
	if objectType == "route" || objectType == "route6" {
		return attributeName == objectType || attributeName == "origin"
	}
	return attributeName == objectType
}
