package rip

import (
	"net/http"
	"net/http/cookiejar"
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

		c.httpClient.Timeout = timeout
	}
}

// WithCookieJar sets a cookie jar on rips httpClient.
func WithCookieJar(jar *cookiejar.Jar) Option {
	return func(c *Client) {
		c.httpClient.Jar = jar
	}
}

// WithDefaultHeaders sets client default headers (e.g. x-api-key)
func WithDefaultHeaders(headers Header) Option {
	return func(c *Client) {
		c.options.Header = headers
	}
}

// WithTransport sets a custom Transport.
func WithTransport(transport *http.Transport) Option {
	return func(c *Client) {
		c.httpClient.Transport = transport
	}
}

func defaultTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:        100,              // Maximum idle connections
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		IdleConnTimeout:     90 * time.Second, // Idle connection timeout
		DisableCompression:  false,            // Enable compression
		DisableKeepAlives:   false,            // Enable keep-alives
	}
}

// NewClient creates a new Client
func NewClient(host string, options ...Option) (*Client, error) {
	u, err := url.Parse(host)
	if err != nil {
		return &Client{}, err
	}

	transport := defaultTransport()

	client := &Client{
		baseURL: u,
		options: &ClientOptions{},
		httpClient: &http.Client{
			Transport: transport,
		},
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
	//nolint: bodyclose
	resp, err := c.httpClient.Do(req.rawRequest)
	if err != nil {
		return NewResponse(req, resp), err
	}

	response := &Response{
		Request: req, rawResponse: resp,
	}

	response.body = resp.Body
	response.Close = func() (err error) {
		if response.body != nil {
			err := response.body.Close()
			if err != nil {
				return err
			}
		}

		if response.rawResponse.Body != nil {
			rErr := response.rawResponse.Body.Close()
			if rErr != nil {
				return rErr
			}
		}

		return
	}

	return response, nil
}
