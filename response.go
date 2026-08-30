package rip

import (
	"io"
	"net/http"
)

// Response the rip response wrapping the original request and response.
type Response struct {
	Request     *Request
	Err         error
	rawResponse *http.Response
	body        []byte
}

func (r *Response) Close() error {
	return r.rawResponse.Body.Close()
}

// Status returns the response status.
func (r *Response) Status() string {
	if r.rawResponse == nil {
		return ""
	}

	return r.rawResponse.Status
}

// StatusCode returns the response status code.
func (r *Response) StatusCode() int {
	if r.rawResponse == nil {
		return 0
	}

	return r.rawResponse.StatusCode
}

func (r *Response) GetHeader(key string) string {
	if r.rawResponse == nil {
		return ""
	}

	return r.rawResponse.Header.Get(key)
}

func (r *Response) ContentType() string {
	return r.GetHeader("Content-Type")
}

func (r *Response) String() (string, error) {
	b, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (r *Response) Bytes() (body []byte, err error) {
	if r.Err != nil {
		return []byte{}, r.Err
	}
	if r.body != nil { // already read once
		return r.body, nil
	}
	if r.rawResponse == nil || r.rawResponse.Body == nil {
		return []byte{}, nil
	}

	defer func() {
		err := r.rawResponse.Body.Close()
		if err != nil {
			r.Err = err
		}
		r.body = body
	}()

	b, err := io.ReadAll(r.rawResponse.Body)
	if err != nil {
		return nil, err
	}
	r.body = b

	return b, nil
}

type ResultState int

const (
	SuccessState ResultState = iota + 1
	ErrorState
	UnknownState
)

func (r *Response) responseState() ResultState {
	if r == nil {
		return UnknownState
	}

	if r.Err != nil {
		return ErrorState
	}

	if code := r.StatusCode(); code > 199 && code < 300 {
		return SuccessState
	} else if code > 399 {
		return ErrorState
	} else {
		return UnknownState
	}
}

// IsSuccess returns true when response state is success.
func (r *Response) IsSuccess() bool {
	return r.responseState() == SuccessState
}

// IsError returns true if response state is error state.
func (r *Response) IsError() bool {
	return r.responseState() == ErrorState
}
