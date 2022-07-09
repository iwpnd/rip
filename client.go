package rip

import (
	"io/ioutil"
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

// NR creates a new request
func (c *Client) NR() *Request {
	h := http.Header{}

	// set default host header
	if c.Options.Header != nil {
		for k, v := range c.Options.Header {
			h.Set(k, v)
		}
	}

	return &Request{client: c, Header: h}
}

func (c *Client) execute(req *Request) (*Response, error) {
	resp, err := c.httpClient.Do(req.RawRequest)
	if err != nil {
		return &Response{}, err
	}

	response := &Response{
		Request: req, RawResponse: resp,
	}

	response.body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}
