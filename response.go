package rip

import (
	"io"
	"net/http"
	"strings"
)

// Response ...
type Response struct {
	Request     *Request
	RawResponse *http.Response
	body        []byte
}

// Status returns the response status
func (r *Response) Status() string {
	if r.RawResponse == nil {
		return ""
	}

	return r.RawResponse.Status
}

// StatusCode returns the response status code
func (r *Response) StatusCode() int {
	if r.RawResponse == nil {
		return 0
	}

	return r.RawResponse.StatusCode
}

// Header method returns the response headers
func (r *Response) Header() http.Header {
	if r.RawResponse == nil {
		return http.Header{}
	}
	return r.RawResponse.Header
}

// String method returns the body of the server response as String.
func (r *Response) String() string {
	if r.body == nil {
		return ""
	}
	return strings.TrimSpace(string(r.body))
}

// Body returns Body as byte array
func (r *Response) Body() []byte {
	if r.body == nil {
		return []byte{}
	}

	return r.body
}

// RawBody returns raw response body. be sure to close
func (r *Response) RawBody() io.ReadCloser {
	if r.RawResponse == nil {
		return nil
	}
	return r.RawResponse.Body
}

// IsSuccess returns true if 199 < StatusCode < 300
func (r *Response) IsSuccess() bool {
	return r.StatusCode() > 199 && r.StatusCode() < 300
}

// IsError returns true if StatusCode > 399
func (r *Response) IsError() bool {
	return r.StatusCode() > 399
}
