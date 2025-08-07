package rip

import (
	"net/http"
	"net/url"
	"time"
)

// Option to use in option pattern.
type Option func(*Client)

// ClientOptions to configure the http client.
type ClientOptions struct {
	Header  Header
	Timeout time.Duration
}

// Client wraps an http client.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	options    *ClientOptions
	Header     Header
}

// WithTimeout sets timeout in seconds on rips httpClient.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.options.Timeout = timeout

		c.httpClient = &http.Client{Timeout: timeout}
	}
}

// WithDefaultHeaders sets client default headers (e.g. x-api-key)
func WithDefaultHeaders(headers Header) Option {
	return func(c *Client) {
		c.options.Header = headers
	}
}

// NewClient creates a new Client
func NewClient(host string, options ...Option) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &Client{}, err
	}

	client := &Client{
		baseURL:    u,
		options:    &ClientOptions{},
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client, nil
}

// NR creates a new request
func (c *Client) NR() *Request {
	h := http.Header{}

	// set default host header
	if c.options.Header != nil {
		for k, v := range c.options.Header {
			h.Set(k, v)
		}
	}

	return &Request{client: c, Header: h}
}

func (c *Client) execute(req *Request) (*Response, error) {
	// either caller is responsible to close the request
	// or Response methods do.
	resp, err := c.httpClient.Do(req.rawRequest) //nolint: bodyclose
	if err != nil {
		return &Response{Request: req, rawResponse: resp, Close: func() (err error) { return }}, err
	}

	response := &Response{
		Request: req, rawResponse: resp,
	}

	response.body = resp.Body

	response.Close = func() error {
		err := response.body.Close()
		if err != nil {
			return err
		}
		rErr := response.rawResponse.Body.Close()
		if rErr != nil {
			return rErr
		}
		return nil
	}

	return response, nil
}
