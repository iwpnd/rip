package rip

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
			if r.Method == "GET" {
				switch r.URL.Path {
				case "/text":
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte("TestGet: text response"))
				case "/json":
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"data":{"test":"test"}}`))
				case "/test/1":
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"data":{"test":"test"}}`))
				default:
					return
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

	testPath := "/text"

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().Execute("GET", testPath)
	if err != nil {
		t.Error("failed to request")
	}
	defer res.RawResponse.Body.Close()

	expectedURL := ts.URL + testPath
	if res.Request.URL != expectedURL {
		t.Errorf("expected: %v, got: %v", expectedURL, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}

func TestGetRequestJSON(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	testPath := "/json"

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetHeader("Content-Type", "application/json").
		Execute("GET", testPath)

	if err != nil {
		t.Error("failed to request")
	}
	defer res.RawResponse.Body.Close()

	expectedURL := ts.URL + testPath
	if res.Request.URL != expectedURL {
		t.Errorf("expected: %v, got: %v", expectedURL, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}

func TestGetRequestWithParams(t *testing.T) {
	teardown := setupTestServer()
	defer teardown()

	testPath := "/test/1"
	testParams := map[string]interface{}{
		"test1": "test",
		"test2": 1,
	}

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NR().
		SetParams(testParams).
		Execute("GET", "/:test1/:test2")
	if err != nil {
		t.Errorf("expected err to be nil got: %v", err)
	}

	expectedURL := ts.URL + testPath
	if res.Request.URL != expectedURL {
		t.Errorf("expected: %v, got: %v", expectedURL, res.Request.URL)
	}

	if res.StatusCode() != 200 {
		t.Errorf("expected StatusCode 200, got: %v", res.StatusCode())
	}
}
