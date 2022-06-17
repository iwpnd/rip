package rip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	defer ts.Close()

	c, err := NewClient(ts.URL, Options{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.Request("GET", "", RequestOptions{})
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	fmt.Println(res.Request.URL.String())

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}

func TestGetRequestWithParams(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
	defer ts.Close()

	c, err := NewClient(ts.URL, Options{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.Request("GET", "/:test1/:test2", RequestOptions{
		Params: map[string]interface{}{
			"test1": "test",
			"test2": 1,
		}})

	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	expectedURL := ts.URL + "/test/1"
	if res.Request.URL.String() != expectedURL {
		t.Errorf("expected: %v, got: %v", expectedURL, res.Request.URL.String())
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}

func TestParseParams(t *testing.T) {
	type tcase struct {
		path     string
		params   Params
		expected string
	}

	tests := map[string]tcase{
		"test no params": {
			path:     "/test",
			expected: "/test",
		},
		"test float": {
			path: "/test/:test",
			params: Params{
				"test": 1.1,
			},
			expected: "/test/1.100000",
		},
		"test int": {
			path: "/test/:test",
			params: Params{
				"test": 1,
			},
			expected: "/test/1",
		},
		"test string": {
			path: "/test/:test",
			params: Params{
				"test": "test",
			},
			expected: "/test/test",
		},
		"test multiple": {
			path: "/test/:test1/:test2/:test3",
			params: Params{
				"test1": "teststring",
				"test2": 1,
				"test3": 1.1,
			},
			expected: "/test/teststring/1/1.100000",
		},
		"test ignore object": {
			path: "/test/:test",
			params: Params{
				"test": map[string]interface{}{
					"bla": "bla",
				},
			},
			expected: "/test/:test",
		},
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			got := parseParams(tc.path, tc.params)

			if got != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, got)
			}
		}
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}

func TestParseQueryParams(t *testing.T) {
	type tcase struct {
		query    Query
		expected string
	}

	tests := map[string]tcase{
		"test one string query": {
			query: Query{
				"test1": "test1",
			},
			expected: "test1=test1",
		},
		"test multiple query": {
			query: Query{
				"test1": "test1",
				"test2": 1,
				"test3": 1.1,
				"test4": true,
			},
			expected: "test1=test1&test2=1&test3=1.100000&test4=true",
		},
		"test empty query": {
			expected: "",
		},
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			got := parseQuery(tc.query)

			if got != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, got)
			}
		}
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}
