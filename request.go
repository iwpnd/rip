package rip

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"maps"
	"net/http"
	gourl "net/url"
)

// Request is the RIP request.
type Request struct {
	Header      http.Header
	pathParams  map[string]string
	queryParams gourl.Values

	Body    []byte
	body    io.ReadCloser // can only be read once
	GetBody func() (io.ReadCloser, error)
	Method  string
	URL     *gourl.URL
	RawURL  string
	Target  any

	contentLength int64
	close         bool

	client     *Client
	rawRequest *http.Request
}

func (r *Request) Get(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodGet, url)
}

func (r *Request) Post(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodPost, url)
}

func (r *Request) Put(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodPut, url)
}

func (r *Request) Patch(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodPatch, url)
}

func (r *Request) Head(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodHead, url)
}

func (r *Request) Option(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodOptions, url)
}

func (r *Request) Delete(ctx context.Context, url string) (*Response, error) {
	return r.Send(ctx, http.MethodDelete, url)
}

func (r *Request) Send(ctx context.Context, method, url string) (*Response, error) {
	r.Method = method
	r.RawURL = url

	resp, _ := r.execute(ctx) //nolint:errcheck // handled in request
	if resp.Err != nil {
		return resp, resp.Err
	}

	return resp, resp.Err
}

func (r *Request) Execute(ctx context.Context) (*Response, error) {
	resp, err := r.execute(ctx)
	return resp, err
}

// execute executes a given request using a method on a given path
func (r *Request) execute(ctx context.Context) (resp *Response, err error) {
	resp = &Response{Request: r}
	defer func() {
		if err != nil {
			resp.Err = err
		} else {
			err = resp.Err
		}
	}()

	for _, mw := range r.client.before {
		if err := mw(r.client, r); err != nil {
			return resp, err
		}
	}

	contentLength := int64(len(r.Body))
	if r.contentLength != 0 {
		contentLength = r.contentLength
	}

	var reqBody io.ReadCloser
	if r.GetBody != nil {
		reqBody, resp.Err = r.GetBody()
		if resp.Err != nil {
			return
		}
	}

	maps.Copy(r.Header, r.client.pinnedHeader)

	req := &http.Request{
		Method:        r.Method,
		Header:        r.Header.Clone(),
		URL:           r.URL,
		ContentLength: contentLength,
		Body:          reqBody,
		Close:         r.close,
	}

	req = req.WithContext(ctx)

	r.rawRequest = req
	var rawResponse *http.Response
	rawResponse, resp.Err = r.client.httpClient.Do(r.rawRequest) //nolint:bodyclose
	resp.rawResponse = rawResponse

	for _, mw := range r.client.after {
		if err := mw(r.client, resp); err != nil {
			return resp, err
		}
	}

	return resp, resp.Err
}

func (r *Request) SetPathParams(params map[string]string) *Request {
	if r.pathParams == nil {
		r.pathParams = map[string]string{}
	}

	r.pathParams = params
	return r
}

func (r *Request) SetPathParam(param, value string) *Request {
	if r.pathParams == nil {
		r.pathParams = map[string]string{}
	}

	r.pathParams[param] = value
	return r
}

func (r *Request) SetTarget(v any) *Request {
	if v == nil {
		return r
	}

	r.Target = v
	return r
}

func (r *Request) SetContentType(ct string) *Request {
	r.Header.Add("Content-Type", ct)
	return r
}

func (r *Request) SetContentLength(cl int64) *Request {
	r.rawRequest.ContentLength = cl
	return r
}

func (r *Request) SetBody(body any) *Request {
	if body == nil {
		return r
	}

	switch b := body.(type) {
	case io.ReadCloser:
		r.body = b
		r.GetBody = func() (io.ReadCloser, error) {
			return r.body, nil
		}
	case io.Reader:
		r.body = io.NopCloser(b)
		r.GetBody = func() (io.ReadCloser, error) {
			return r.body, nil
		}
	case []byte:
		r.Body = b
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(b)), nil
		}
	case string:
		r.Body = fmt.Append([]byte{}, body)
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(r.Body)), nil
		}
	default:
		mb, err := json.Marshal(b)
		if err != nil {
			r.Body = []byte{}
			r.GetBody = func() (io.ReadCloser, error) {
				return nil, ErrInvalidRequestBody
			}
			return r
		}
		r.Body = mb
		r.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(mb)), nil
		}
	}

	return r
}
