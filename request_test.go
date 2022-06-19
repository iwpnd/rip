package rip

import (
	"fmt"
	"net/http"
	"testing"
)

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
			r := &Request{}
			r.parsePath(tc.path, tc.params)

			if r.Path != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, r.Path)
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
			r := &Request{}
			r.parseQuery(tc.query)

			got := r.QueryParams.Encode()
			if got != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, got)
			}
		}
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}

func TestSetHeader(t *testing.T) {
	type tcase struct {
		options   RequestOptions
		key       string
		keys      []string
		value     string
		values    []string
		expHeader http.Header
	}

	tests := map[string]tcase{
		"test single header without request options": {
			options: RequestOptions{},
			key:     "test",
			value:   "test",
			expHeader: http.Header{
				"Test": []string{"test"},
			},
		},
		"test multiple header without request options": {
			options: RequestOptions{},
			keys:    []string{"test1", "test2"},
			values:  []string{"test", "test"},
			expHeader: http.Header{
				"Test1": []string{"test"},
				"Test2": []string{"test"},
			},
		},
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			r := &Request{Options: tc.options}

			if tc.key != "" && tc.value != "" {
				r.SetHeader(tc.key, tc.value)
			}

			if len(tc.keys) != 0 && len(tc.values) != 0 {
				for i, k := range tc.keys {
					r.SetHeader(k, tc.values[i])
				}
			}

			got := r.Header
			if fmt.Sprint(got) != fmt.Sprint(tc.expHeader) {
				t.Errorf("expected: %v, got: %v", fmt.Sprint(tc.expHeader), fmt.Sprint(got))
			}
		}
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}

func TestSetHeaders(t *testing.T) {
	type tcase struct {
		options     RequestOptions
		inputHeader Header
		expHeader   http.Header
	}

	tests := map[string]tcase{
		"test multiple header without request options": {
			options: RequestOptions{},
			inputHeader: Header{
				"test1": "test",
				"test2": "test",
			},
			expHeader: http.Header{
				"Test1": []string{"test"},
				"Test2": []string{"test"},
			},
		},
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			r := &Request{Options: tc.options}

			r.SetHeaders(tc.inputHeader)

			got := r.Header
			if fmt.Sprint(got) != fmt.Sprint(tc.expHeader) {
				t.Errorf("expected: %v, got: %v", fmt.Sprint(tc.expHeader), fmt.Sprint(got))
			}
		}
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}
