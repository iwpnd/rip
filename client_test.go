package rip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

type tcase struct {
	Body          string
	Headers       Header
	Method        string
	Params        Params
	Path          string
	Query         Query
	expBody       interface{}
	expPath       string
	expStatusCode int
}

var (
	ts *httptest.Server
)

func setupTestServer() func() {
	ts = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				switch r.URL.Path {
				case "/test":
					{
						accept := r.Header.Get("Accept")
						switch accept {
						case contentTypeJSON:
							w.Header().Set("Content-Type", contentTypeJSON)
							w.WriteHeader(http.StatusOK)
							fmt.Fprint(w, fixture("response.json"))
						case contentTypeTEXT:
							w.WriteHeader(http.StatusOK)
							w.Header().Set("Content-Type", contentTypeTEXT)
							fmt.Fprint(w, "text response")
						default:
							w.WriteHeader(http.StatusNotAcceptable)
							fmt.Fprint(w, "nope")
						}
					}
				case "/test/1/2":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, fixture("response.json"))
				default:
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusNotFound)
				}
			case "POST":
				switch r.URL.Path {
				case "/test":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusCreated)
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					res, err := strconv.Unquote(string(body))
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					fmt.Fprint(w, res)
				case "/test/1/2":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusCreated)
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					res, err := strconv.Unquote(string(body))
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					fmt.Fprint(w, res)
				default:
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusNotFound)
				}
			case "PUT":
				switch r.URL.Path {
				case "/test":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusOK)
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					res, err := strconv.Unquote(string(body))
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					fmt.Fprint(w, res)
				case "/test/1/2":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusOK)
					body, err := ioutil.ReadAll(r.Body)
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					res, err := strconv.Unquote(string(body))
					if err != nil {
						http.Error(w, "can't read body", http.StatusBadRequest)
						return
					}
					fmt.Fprint(w, res)
				default:
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusNotFound)
				}
			case "DELETE":
				switch r.URL.Path {
				case "/test/1/2":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, string([]byte(`"{"ok":true}"`)))
				default:
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusNotFound)
				}
			}
		}))

	return func() {
		ts.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(b))
}

func TestClientWithOptions(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	c, err := NewClient(ts.URL,
		WithDefaultHeaders(map[string]string{
			"x-api-key": "api-key-test",
		}),
		WithTimeout(30),
	)

	if err != nil {
		t.Error("could not initialize client")
	}

	if c.options.Timeout == defaultTimeout {
		t.Error("should not be defaultTimeout")
	}

	if c.options.Header == nil {
		t.Error("should not be nil Header")
	}
}

func TestClientWithoutOptions(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	c, err := NewClient(ts.URL)

	if err != nil {
		t.Error("could not initialize client")
	}

	if c.options.Timeout != defaultTimeout {
		t.Errorf("should be defaultTimeout, got %v", c.options.Timeout)
	}

	if c.options.Header != nil {
		t.Error("should be nil Header")
	}
}

func TestClientRequests(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	c, err := NewClient(ts.URL,
		WithDefaultHeaders(map[string]string{
			"x-api-key": "api-key-test",
		}),
		WithTimeout(30),
	)

	if err != nil {
		t.Error("could not initialize client")
	}

	fn := func(tc tcase) func(*testing.T) {
		return func(t *testing.T) {
			req := c.NR()

			if tc.Headers != nil {
				req.SetHeaders(tc.Headers)
			}

			if tc.Params != nil {
				req.SetParams(tc.Params)
			}

			if tc.Query != nil {
				req.SetQuery(tc.Query)
			}

			if tc.Body != "" {
				req.SetBody(tc.Body)
			}

			res, err := req.Execute(tc.Method, tc.Path)
			if err != nil {
				t.Error("failed to request")
			}
			defer res.RawResponse.Body.Close()

			if res.Request.URL != ts.URL+tc.expPath {
				t.Errorf("\n\n expected: %v, got: %v \n\n", ts.URL+tc.expPath, res.Request.URL)
			}

			if res.StatusCode() != tc.expStatusCode {
				t.Errorf("\n\n expected StatusCode %v, got: %v \n\n", tc.expStatusCode, res.StatusCode())
			}

			if tc.Query != nil {
				for k, v := range tc.Query {
					q := res.Request.RawRequest.URL.Query().Get(k)
					if q != fmt.Sprintf("%v", v) {
						t.Errorf("\n\n expected query param %v to be %v, got value %v \n\n", k, q, v)
					}
				}
			}

			if tc.Headers != nil {
				for k, v := range tc.Headers {
					h := res.Request.RawRequest.Header.Get(k)
					if h == "" {
						t.Errorf("\n\n Expected request header %v to be %v \n\n, got: %v", k, v, h)
					}

					if k == "x-api-key" {
						if h != tc.Headers[k] {
							t.Errorf("\n\n Expected default x-api-key to be overwritten \n\n")
						}
					}
				}
			}

			if tc.expBody != nil {
				if res.String() != tc.expBody {
					t.Errorf("failed. Response \n\n %+v \n\n does not match expected response \n\n %+v \n\n", res.String(), tc.expBody)
					return

				}
			}
		}
	}

	tests := map[string]tcase{
		"GET text": {
			Method: "GET",
			Path:   "/test",
			Headers: map[string]string{
				"Accept": contentTypeTEXT,
			},
			expPath:       "/test",
			expStatusCode: 200,
			expBody:       "text response",
		},
		"GET json": {
			Method: "GET",
			Path:   "/test",
			Headers: map[string]string{
				"Accept": contentTypeJSON,
			},
			Query:         Query{"test": 1},
			expPath:       "/test",
			expStatusCode: 200,
			expBody:       fixture("response.json"),
		},
		"GET json should not overwrite x-api-key": {
			Method: "GET",
			Path:   "/test",
			Headers: map[string]string{
				"Accept":    contentTypeJSON,
				"x-api-key": "should-overwrite-default",
			},
			expPath:       "/test",
			expStatusCode: 200,
			expBody:       fixture("response.json"),
		},
		"GET json resolve params": {
			Method: "GET",
			Path:   "/test/:id1/:id2",
			Params: Params{
				"id1": "1",
				"id2": "2",
			},
			Headers: map[string]string{
				"Accept": contentTypeJSON,
			},
			expPath:       "/test/1/2",
			expBody:       fixture("response.json"),
			expStatusCode: 200,
		},
		"GET fails": {
			Method:        "GET",
			Path:          "/test/fails",
			expPath:       "/test/fails",
			expStatusCode: 404,
		},
		"POST json": {
			Method: "POST",
			Path:   "/test",
			Headers: map[string]string{
				"Accept":       contentTypeJSON,
				"Content-Type": contentTypeJSON,
			},
			Body:          fixture("response.json"),
			expPath:       "/test",
			expBody:       fixture("response.json"),
			expStatusCode: 201,
		},
		"POST json resolve params": {
			Method: "POST",
			Path:   "/test/:id1/:id2",
			Params: Params{
				"id1": "1",
				"id2": "2",
			},
			Body: fixture("response.json"),
			Headers: map[string]string{
				"Accept":       contentTypeJSON,
				"Content-Type": contentTypeJSON,
			},
			expPath:       "/test/1/2",
			expBody:       fixture("response.json"),
			expStatusCode: 201,
		},
		"POST fails": {
			Method:        "POST",
			Path:          "/test/fails",
			Body:          fixture("response.json"),
			expPath:       "/test/fails",
			expStatusCode: 404,
		},
		"PUT json": {
			Method: "PUT",
			Path:   "/test",
			Headers: map[string]string{
				"Accept":       contentTypeJSON,
				"Content-Type": contentTypeJSON,
			},
			Body:          fixture("response.json"),
			expPath:       "/test",
			expBody:       fixture("response.json"),
			expStatusCode: 200,
		},
		"PUT json resolve params": {
			Method: "PUT",
			Path:   "/test/:id1/:id2",
			Params: Params{
				"id1": "1",
				"id2": "2",
			},
			Body: fixture("response.json"),
			Headers: map[string]string{
				"Accept":       contentTypeJSON,
				"Content-Type": contentTypeJSON,
			},
			expPath:       "/test/1/2",
			expBody:       fixture("response.json"),
			expStatusCode: 200,
		},
		"PUT fails": {
			Method:        "PUT",
			Path:          "/test/fails",
			Body:          fixture("response.json"),
			expPath:       "/test/fails",
			expStatusCode: 404,
		},
		"DELETE": {
			Method: "DELETE",
			Path:   "/test/:id1/:id2",
			Params: Params{
				"id1": "1",
				"id2": "2",
			},
			expPath:       "/test/1/2",
			expStatusCode: 200,
		},
		"DELETE fails": {
			Method:        "DELETE",
			Path:          "/test/fails",
			expPath:       "/test/fails",
			expStatusCode: 404,
		},
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}
}
