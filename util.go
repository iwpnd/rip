package rip

import (
	"encoding/json"
	"regexp"
)

var jsonCheck = regexp.MustCompile(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`)

// IsJSON helper to determine content type
func IsJSON(ct string) bool {
	return jsonCheck.MatchString(ct)
}

func Unmarshal(c *Client, r *Response, v any) error {
	ct := r.ContentType()

	b, err := r.Bytes()
	if err != nil {
		return err
	}

	return unmarshal(ct, b, v)
}

// Unmarshal helper
func unmarshal(ct string, b []byte, d any) error {
	if IsJSON(ct) {
		err := json.Unmarshal(b, d)
		if err != nil {
			return err
		}
	}

	return nil
}
