package rip

import (
	"net/http"
	"net/http/httptest"
)

func createTestServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	return server
}

func createTestClient(url string) *Client {
	return C().SetBaseURL(url)
}
