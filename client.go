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

// Headers ...
type Headers = map[string]string

// Params ...
type Params = map[string]interface{}

// Query ...
type Query = map[string]interface{}

// Options ...
type Options struct {
	Headers Headers
}

// RequestOptions ...
type RequestOptions struct {
	Headers Headers
	Params  Params
	Query   Query
	Body    interface{}
}

// Client ...
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	Options    Options
}

// NewClient creates a new Client
func NewClient(host string, options Options) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &Client{}, err
	}

	return &Client{httpClient: &http.Client{}, baseURL: u, Options: options}, nil
}

// Request performs a request on a resource
func (r Client) Request(method, path string, options RequestOptions) (*Response, error) {
	ppath := parseParams(path, options.Params)
	qs := parseQuery(options.Query)
	headers := parseHeader(r.Options.Headers, options.Headers)

	body, err := parseBody(options.Body)
	if err != nil {
		return &Response{}, err
	}

	if body != nil {
		headers.Set("Accept", "application/json")
	}

	reqURL := r.buildRequestURL(ppath, qs)
	req, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return &Response{}, err
	}

	req.Header = headers

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return &Response{}, err
	}

	return &Response{RawResponse: resp, Request: req}, nil
}

func (r Client) buildRequestURL(path, qs string) string {
	url := r.baseURL.String() + path
	if qs != "" {
		url = url + "?" + qs
	}

	return url
}

func parseBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBody), nil
}

func parseHeader(defaults, overwrites Headers) http.Header {
	h := http.Header{}

	for k, v := range defaults {
		h.Set(k, v)
	}

	for k, v := range overwrites {
		h.Set(k, v)
	}

	return h
}

func parseQuery(query Query) string {
	if query == nil {
		return ""
	}

	q := url.Values{}
	for k, v := range query {
		switch v.(type) {
		case float32:
		case float64:
			q.Add(k, fmt.Sprintf("%.6f", v))
		case int32:
		case int64:
		case int:
			q.Add(k, fmt.Sprintf("%d", v))
		case string:
			q.Add(k, fmt.Sprintf("%s", v))
		case bool:
			q.Add(k, fmt.Sprintf("%t", v))
		default:
			break
		}
	}

	return q.Encode()
}

func parseParams(path string, params Params) string {
	if params == nil {
		return path
	}

	var ppath = path

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
			ppath = strings.Replace(ppath, fmt.Sprintf(":%s", k), p, 1)
		}
	}

	return ppath
}
