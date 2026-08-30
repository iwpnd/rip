package rip

import (
	"net/http"
	gourl "net/url"
	"time"
)

// TODO:
// - cookie jar

// Client wraps an http client.
type Client struct {
	baseURL      string
	pathParams   map[string]string
	queryParams  gourl.Values
	pinnedHeader http.Header

	before []RequestMiddleware
	after  []ResponseMiddleware

	httpClient *http.Client
}

func C() *Client {
	transport := T()
	params := map[string]string{}
	query := gourl.Values{}
	header := http.Header{}
	client := &Client{
		httpClient: &http.Client{
			Transport: transport,
		},
		before: []RequestMiddleware{
			parseRequestUrl,
		},
		after: []ResponseMiddleware{
			debugResponse,
			parseResponseBody,
			// TODO: add download middleware here
		},
		pathParams:   params,
		queryParams:  query,
		pinnedHeader: header,
	}

	return client
}

// NewClient creates a new Client
func NewClient() *Client {
	return C()
}

// NR creates a new request
func (c *Client) NR() *Request {
	return &Request{client: c, Header: http.Header{}}
}

func (c *Client) Get(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodGet
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Post(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodPost
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Put(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodPut
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Patch(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodPatch
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Delete(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodDelete
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Options(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodOptions
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) Head(url ...string) *Request {
	req := c.NR()
	req.Method = http.MethodHead
	if len(url) > 0 {
		req.RawURL = url[0]
	}
	return req
}

func (c *Client) WrapRoundTripper(mw RoundTripperMiddleware) *Client {
	c.httpClient.Transport = mw(c.httpClient.Transport)
	return c
}

func (c *Client) SetBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

func (c *Client) SetPathParams(params map[string]string) *Client {
	c.pathParams = params
	return c
}

func (c *Client) SetPathParam(param, value string) *Client {
	c.pathParams[param] = value
	return c
}

func (c *Client) SetQueryParams(params gourl.Values) *Client {
	for i, v := range params {
		for _, p := range v {
			c.queryParams.Add(i, p)
		}
	}
	return c
}

func (c *Client) SetQueryParm(param, value string) *Client {
	c.queryParams.Set(param, value)
	return c
}

func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

func (c *Client) SetPinnedHeaders(header http.Header) *Client {
	c.pinnedHeader = header
	return c
}

func (c *Client) SetPinnedHeader(key, value string) *Client {
	if c.pinnedHeader == nil {
		c.pinnedHeader = http.Header{}
	}

	c.pinnedHeader.Set(key, value)

	return c
}
