package rip

import (
	"net/http"
	"time"
)

func T() *http.Transport {
	return defaultTransport()
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

type RoundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

type RoundTripperMiddleware = func(base http.RoundTripper) http.RoundTripper

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func WithAuth(key string) func(base http.RoundTripper) http.RoundTripper {
	return func(base http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("Authorization", "Bearer "+key)
			return base.RoundTrip(req)
		})
	}
}

func WithAPIKey(key string) func(base http.RoundTripper) http.RoundTripper {
	return func(base http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			req = req.Clone(req.Context())
			req.Header.Set("x-api-key", key)
			return base.RoundTrip(req)
		})
	}
}
