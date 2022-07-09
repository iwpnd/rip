package rip

import (
	"encoding/json"
	"regexp"
)

var (
	contentTypeTEXT = "text/plain; charset=utf-8"
	contentTypeJSON = "application/json"

	jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)
)

// IsJSON helper to determine content type
func IsJSON(ct string) bool {
	return jsonCheck.MatchString(ct)
}

// Unmarshal helper
func Unmarshal(ct string, b []byte, d interface{}) error {
	if IsJSON(ct) {
		err := json.Unmarshal(b, d)
		if err != nil {
			return err
		}
	}

	return nil
}
