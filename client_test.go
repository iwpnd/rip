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
	mux    *http.ServeMux
	ts     *httptest.Server
	client *Client
)

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func setup() func() {
	mux = http.NewServeMux()
	ts = httptest.NewServer(mux)

	return func() {
		ts.Close()
	}
}

func TestGetRequestText(t *testing.T) {
	teardown := setup()
	defer teardown()

	testPath := "/text"

	mux.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("get: text response"))
	})

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NewRequest(RequestOptions{}).Execute("GET", testPath)
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
	teardown := setup()
	defer teardown()

	testPath := "/json"

	mux.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"test":"test"}}`))
	})

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NewRequest(RequestOptions{}).SetHeader("Content-Type", "application/json").Execute("GET", testPath)
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
	teardown := setup()
	defer teardown()

	testPath := "/test/1"

	mux.HandleFunc(testPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"test":"test"}}`))
	})

	c, err := NewClient(ts.URL, ClientOptions{})
	if err != nil {
		t.Error("Cannot initialize client")
	}

	res, err := c.NewRequest(RequestOptions{
		Params: map[string]interface{}{
			"test1": "test",
			"test2": 1,
		}}).Execute("GET", "/:test1/:test2")
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
