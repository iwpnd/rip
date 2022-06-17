package rip

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
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
	URL        *url.URL
	Options    Options
}

// NewClient creates a new Client
func NewClient(host string, options Options) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &Client{}, err
	}

	return &Client{httpClient: &http.Client{}, URL: u, Options: options}, nil
}

// Request performs a request on a resource
func (r Client) Request(method, path string, options RequestOptions) (*Response, error) {
	ppath := parseParams(path, options.Params)
	qs := parseQueryString(options.Query)

	reqURL := r.buildRequestURL(ppath, qs)
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return &Response{}, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &Response{RawResponse: resp, Request: req}, nil
}

func (r Client) buildRequestURL(path, qs string) string {
	url := r.URL.String() + path
	if qs != "" {
		url = url + "?" + qs
	}

	return url
}

func parseQueryString(query Query) string {
	if query == nil {
		return ""
	}

	var squery = make(map[string]string)
	for k, v := range query {
		var q string

		switch v.(type) {
		case float32:
		case float64:
			q = fmt.Sprintf("%.6f", v)
		case int32:
		case int64:
		case int:
			q = fmt.Sprintf("%d", v)
		case string:
			q = fmt.Sprintf("%s", v)
		default:
			q = ""
		}

		squery[k] = q
	}

	var keys []string
	for k := range squery {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf strings.Builder
	for _, k := range keys {
		escValue := url.QueryEscape(squery[k])
		escKey := url.QueryEscape(k)

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(escKey)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(escValue))
	}

	return buf.String()
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
