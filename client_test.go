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

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Payload struct {
	Data User `json:"data"`
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
					return
				}
			case "POST":
				switch r.URL.Path {
				case "/test/2":
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
				}
			case "PUT":
				switch r.URL.Path {
				case "/test/3":
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

func TestClientGET(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	type TCase struct {
		Method        string
		Path          string
		Params        Params
		Headers       Header
		Body          []byte
		expPath       string
		expStatusCode int
		expBody       interface{}
	}

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("could not initialize client")
	}

	fn := func(tc TCase) func(*testing.T) {
		return func(t *testing.T) {
			req := c.NR()

			if tc.Headers != nil {
				req.SetHeaders(tc.Headers)
			}

			if tc.Params != nil {
				req.SetParams(tc.Params)
			}

			res, err := req.Execute(tc.Method, tc.Path)
			if err != nil {
				t.Error("failed to request")
			}
			defer res.RawResponse.Body.Close()

			if res.Request.URL != ts.URL+tc.expPath {
				t.Errorf("expected: %v, got: %v", ts.URL+tc.expPath, res.Request.URL)
			}

			if res.StatusCode() != tc.expStatusCode {
				t.Errorf("expected StatusCode %v, got: %v", tc.expStatusCode, res.StatusCode())
			}

			if tc.expBody != nil {
				if res.String() != tc.expBody {
					t.Errorf("failed. Response \n\n %+v \n\n does not match expected response \n\n %+v \n\n", res.String(), tc.expBody)
					return

				}
			}
		}
	}

	tests := map[string]TCase{
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
	}

	for name, tc := range tests {
		t.Run(name, fn(tc))
	}

}

func TestPostWithBody(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	path := "/:path/:id"
	params := map[string]interface{}{
		"path": "test",
		"id":   2,
	}
	url := ts.URL + "/test/2"
	body := fixture("response.json")

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetParams(params).
		SetBody(body).
		Execute("POST", path)
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	if res.Request.URL != url {
		t.Errorf("expected url: %v, got: %v", url, res.Request.URL)
	}

	if res.StatusCode() != 201 {
		t.Errorf("expected status code 201, got: %v", res.StatusCode())
	}

	if res.Request.RawRequest.Header.Get("Content-Type") != contentTypeJSON {
		t.Errorf("expected Content-Type: %v, got: %v",
			contentTypeJSON, res.Request.RawRequest.Header.Get("Content-Type"),
		)
	}

	r := &Payload{}
	err = Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
	if err != nil {
		t.Error("failed to unmarshal response", err)
	}
}

func TestPutWithBody(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	path := "/:path/:id"
	params := map[string]interface{}{
		"path": "test",
		"id":   3,
	}
	url := ts.URL + "/test/3"
	body := fixture("response.json")

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetParams(params).
		SetBody(body).
		Execute("PUT", path)
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	if res.Request.URL != url {
		t.Errorf("expected url: %v, got: %v", url, res.Request.URL)
	}

	if res.StatusCode() != 201 {
		t.Errorf("expected status code 201, got: %v", res.StatusCode())
	}

	if res.Request.RawRequest.Header.Get("Content-Type") != contentTypeJSON {
		t.Errorf("expected Content-Type: %v, got: %v",
			contentTypeJSON, res.Request.RawRequest.Header.Get("Content-Type"),
		)
	}

	r := &Payload{}
	err = Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
	if err != nil {
		t.Error("failed to unmarshal response", err)
	}
}
