package rip

import (
	"net/http"
	"net/url"
	"time"
)

const defaultTimeOut = 4

// ClientOptions ...
type ClientOptions struct {
	Header  Header
	Timeout int
}

// Client ...
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	Options    ClientOptions
}

// NewClient creates a new Client
func NewClient(host string, options ClientOptions) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &Client{}, err
	}

	timeout := time.Duration(defaultTimeOut) * time.Second
	if options.Timeout != 0 {
		timeout = time.Duration(options.Timeout) * time.Second
	}

	return &Client{httpClient: &http.Client{Timeout: timeout}, baseURL: u, Options: options}, nil
}

// NewRequest creates a new request from RequestOptions
func (c *Client) NewRequest(options RequestOptions) *Request {
	req := &Request{client: c, Options: options}

	req.parseQuery(options.Query)
	req.parseHeader(options.Header)
	req.parseBody(options.Body)

	return req
}

func (c *Client) execute(req *Request) (*Response, error) {
	// overwrite req header with host header
	if c.Options.Header != nil {
		for k, v := range c.Options.Header {
			req.RawRequest.Header.Set(k, v)
		}
	}
	resp, err := c.httpClient.Do(req.RawRequest)
	if err != nil {
		return &Response{}, err
	}

	response := &Response{
		Request: req, RawResponse: resp,
	}

	return response, nil
}
