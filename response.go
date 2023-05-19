package rip

import (
	"io"
	"net/http"
	"strings"
)

// Response ...
type Response struct {
	Request     *Request
	rawResponse *http.Response
	body        io.ReadCloser
	Close       func()
}

// ContentLength returns the content-length
func (r *Response) ContentLength() int64 {
	if r.rawResponse == nil {
		return 0
	}

	return r.rawResponse.ContentLength
}

// Status returns the response status
func (r *Response) Status() string {
	if r.rawResponse == nil {
		return ""
	}

	return r.rawResponse.Status
}

// StatusCode returns the response status code
func (r *Response) StatusCode() int {
	if r.rawResponse == nil {
		return 0
	}

	return r.rawResponse.StatusCode
}

// Header method returns the response headers
func (r *Response) Header() http.Header {
	if r.rawResponse == nil {
		return http.Header{}
	}
	return r.rawResponse.Header
}

// String method returns the body of the server response as String.
func (r *Response) String() string {
	if r.body == nil {
		return ""
	}
	defer r.body.Close()

	body, err := io.ReadAll(r.body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(body))
}

// Body returns Body as byte array
func (r *Response) Body() []byte {
	if r.body == nil {
		return []byte{}
	}
	defer r.body.Close()

	body, err := io.ReadAll(r.body)
	if err != nil {
		return []byte{}
	}

	return body
}

// RawBody returns raw response body. be sure to close
func (r *Response) RawBody() io.ReadCloser {
	if r.rawResponse == nil {
		return nil
	}
	return r.rawResponse.Body
}

// IsSuccess returns true if 199 < StatusCode < 300
func (r *Response) IsSuccess() bool {
	return r.StatusCode() > 199 && r.StatusCode() < 300
}

// IsError returns true if StatusCode > 399
func (r *Response) IsError() bool {
	return r.StatusCode() > 399
}
