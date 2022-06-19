package rip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Header ...
type Header = map[string]string

// Params ...
type Params = map[string]interface{}

// Options ...
type Options struct {
	Header  Header
	Timeout int
}

// RequestOptions ...
type RequestOptions struct {
	Header Header
	Params Params
	Query  Query
	Body   interface{}
}

// Query ...
type Query = map[string]interface{}

// Request ...
type Request struct {
	Path        string
	RawRequest  *http.Request
	Header      http.Header
	QueryParams url.Values
	Result      interface{}
	Body        io.Reader
	URL         string
	Options     RequestOptions
	client      *Client
}

// Execute executes a given request using a method on a given path
func (r *Request) Execute(method, path string) (*Response, error) {
	var err error

	r.parsePath(path, r.Options.Params)
	r.parseURL()

	fmt.Print(r.URL)

	r.RawRequest, err = http.NewRequest(method, r.URL, r.Body)
	if err != nil {
		return &Response{}, err
	}

	r.RawRequest.Header = r.Header

	resp, err := r.client.httpClient.Do(r.RawRequest)
	if err != nil {
		return &Response{}, err
	}

	response := &Response{
		Request:     r,
		RawResponse: resp,
	}

	return response, err
}

func (r *Request) parseQuery(query Query) {
	r.QueryParams = url.Values{}
	for k, v := range query {
		switch v.(type) {
		case float32:
		case float64:
			r.QueryParams.Add(k, fmt.Sprintf("%.6f", v))
		case int32:
		case int64:
		case int:
			r.QueryParams.Add(k, fmt.Sprintf("%d", v))
		case string:
			r.QueryParams.Add(k, fmt.Sprintf("%s", v))
		case bool:
			r.QueryParams.Add(k, fmt.Sprintf("%t", v))
		default:
			break
		}
	}
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

func (r *Request) parseHeader(header Header) {
	r.Header = http.Header{}
	for k, v := range header {
		r.Header.Set(k, v)
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

func (r *Request) parseBody(body interface{}) error {
	if body == nil {
		r.Body = nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	r.Body = bytes.NewBuffer(jsonBody)

	contentType := r.Header.Get("Content-Type")
	if !IsJSON(contentType) {
		r.Header.Set("Content-Type", "application/json")
	}

	return nil
}

func (r *Request) parseURL() {
	r.URL = r.client.baseURL.String() + r.Path

	if r.QueryParams.Encode() != "" {
		r.URL = r.URL + "?" + r.QueryParams.Encode()
	}
}
