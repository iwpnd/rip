package rip

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = time.Duration(4) * time.Second

// Option ...
type Option func(*Client)

// ClientOptions ...
type ClientOptions struct {
	Header  Header
	Timeout time.Duration
}

// Client ...
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	options    *ClientOptions
	Header     Header
}

// WithTimeout sets timeout in seconds on rips httpClient
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.options.Timeout = time.Duration(timeout)
		timeout = time.Duration(timeout) * time.Second

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
		baseURL: u,
		options: &ClientOptions{},
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}

	client.options.Timeout = defaultTimeout

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
	resp, err := c.httpClient.Do(req.RawRequest)
	if err != nil {
		return &Response{}, err
	}
	defer resp.Body.Close()

	response := &Response{
		Request: req, RawResponse: resp,
	}

	response.body, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}
