package rip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type TestResponse struct {
	User string `json:"user"`
	Age  int    `json:"age"`
}

type TestResponseData struct {
	Data TestResponse `json:"data"`
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
				case "/text":
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", contentTypeTEXT)
					fmt.Fprint(w, "TestGet: text response")
				case "/json":
					w.Header().Set("Content-Type", contentTypeJSON)
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, fixture("response.json"))
				case "/test/1":
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
	return string(b)
}

func TestGetRequestText(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	path := "/text"
	url := ts.URL + path

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().Execute("GET", path)
	if err != nil {
		t.Error("failed to request")
	}
	defer res.RawResponse.Body.Close()

	if res.Request.URL != url {
		t.Errorf("expected: %v, got: %v", url, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}

func TestGetRequestJSON(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	path := "/json"
	url := ts.URL + path

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetHeader("Content-Type", contentTypeJSON).
		Execute("GET", path)

	if err != nil {
		t.Error("failed to request")
	}

	if res.Request.URL != url {
		t.Errorf("expected: %v, got: %v", url, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}

	if res.Request.RawRequest.Header.Get("Content-Type") != contentTypeJSON {
		t.Errorf("expected Content-Type: %v, got: %v",
			contentTypeJSON, res.Request.RawRequest.Header.Get("Content-Type"),
		)
	}

}

func TestGetRequestWithParams(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	path := "/test/1"
	params := map[string]interface{}{
		"test1": "test",
		"test2": 1,
	}
	url := ts.URL + path

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetParams(params).
		Execute("GET", "/:test1/:test2")
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	if res.Request.URL != url {
		t.Errorf("expected: %v, got: %v", url, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}

	r := &TestResponseData{}
	err = Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
	if err != nil {
		t.Error("failed to unmarshal response", err)
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

	r := &TestResponseData{}
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

	r := &TestResponseData{}
	err = Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
	if err != nil {
		t.Error("failed to unmarshal response", err)
	}
}
