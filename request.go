package rip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ErrClientMissing occurs when Request is instantiated without Client.NR()
var ErrClientMissing = errors.New("use .NR() to create a new request instead")

// Header ...
type Header = map[string]string

// Params ...
type Params = map[string]interface{}

// Query ...
type Query = map[string]interface{}

// Request ...
type Request struct {
	Body       interface{}
	Header     http.Header
	Params     Params
	Path       string
	Query      url.Values
	Result     interface{} // NOTE: can I pass struct here to unmarshal resp body to?
	URL        string
	client     *Client
	rawRequest *http.Request
}

// Execute executes a given request using a method on a given path
func (r *Request) Execute(method, path string) (*Response, error) {
	if r.client == nil {
		return &Response{}, ErrClientMissing
	}

	var err error

	r.parsePath(path, r.Params)
	r.parseURL()

	if rd, ok := r.Body.(io.Reader); ok {
		r.rawRequest, err = http.NewRequest(method, r.URL, rd)
	} else {
		r.rawRequest, err = http.NewRequest(method, r.URL, nil)
	}

	if err != nil {
		return &Response{}, err
	}

	r.rawRequest.Header = r.Header

	if r.Query != nil {
		r.rawRequest.URL.RawQuery = r.Query.Encode()
	}

	resp, err := r.client.execute(r)
	if err != nil {
		return &Response{}, err
	}
	resp.Close = func() {
		resp.body.Close()
		resp.rawResponse.Body.Close()
	}

	return resp, err
}

// SetQuery to set query parameters
func (r *Request) SetQuery(query Query) *Request {
	r.parseQuery(query)

	return r
}

func (r *Request) parseQuery(query Query) {
	r.Query = url.Values{}
	for k, v := range query {
		switch v.(type) {
		case float32:
		case float64:
			r.Query.Add(k, fmt.Sprintf("%.6f", v))
		case int32:
		case int64:
		case int:
			r.Query.Add(k, fmt.Sprintf("%d", v))
		case string:
			r.Query.Add(k, fmt.Sprintf("%s", v))
		case bool:
			r.Query.Add(k, fmt.Sprintf("%t", v))
		default:
			break
		}
	}
}

// SetParams to replace in path
func (r *Request) SetParams(params Params) *Request {
	r.Params = params

	return r
}

func (r *Request) parsePath(path string, params Params) {
	r.Path = path

	for k, v := range params {
		var p string

		switch v.(type) {
		case float32:
		case float64:
			p = fmt.Sprintf("%.6f", v)
		case int32:
		case int64:
		case int:
			p = fmt.Sprintf("%d", v)
		case string:
			p = fmt.Sprintf("%s", v)
		default:
			p = ""
		}

		if p != "" {
			r.Path = strings.Replace(r.Path, fmt.Sprintf(":%s", k), p, 1)
		}
	}
}

// SetHeader to set a single header
func (r *Request) SetHeader(key, value string) *Request {
	if r.Header == nil {
		r.Header = http.Header{}
	}

	r.Header.Add(key, value)

	return r
}

// SetHeaders to set multiple header as map[string]string
func (r *Request) SetHeaders(header Header) *Request {
	if r.Header == nil {
		r.Header = http.Header{}
	}

	for k, v := range header {
		r.Header.Set(k, v)
	}

	return r
}

// SetBody to set a request body
func (r *Request) SetBody(body interface{}) *Request {
	if body == nil {
		return r
	}

	b := r.parseBody(body)
	r.Body = b

	return r
}

// NOTE: rn expected json only
func (r *Request) parseBody(body interface{}) io.Reader {
	if body == nil {
		return nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	b := bytes.NewBuffer(jsonBody)

	contentType := r.Header.Get("Content-Type")
	if !IsJSON(contentType) {
		r.Header.Set("Content-Type", "application/json")
	}

	return b
}

func (r *Request) parseURL() {
	r.URL = r.client.baseURL.String() + r.Path
}
