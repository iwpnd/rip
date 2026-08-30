package rip

import (
	"errors"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
)

func TestRequestBody(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			defer r.Body.Close()

			w.Header().Add("method", r.Method)
			io.Copy(w, r.Body)
		default:
			w.Header().Add("method", r.Method)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	s := createTestServer(testHandler)
	defer s.Close()

	c := createTestClient(s.URL)

	type tcase struct {
		fn             func(c *Client, body any) *Response
		body           string
		toBodyFn       func(b string) any
		expectedMethod string
		expectedCode   int
		expectErr      bool
		expectStruct   bool
	}
	tests := map[string]tcase{
		"POST success struct": {
			body: `{"foo":"bar"}`,
			toBodyFn: func(b string) any {
				return b
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().
					SetBody(struct {
						Foo string `json:"foo"`
					}{Foo: "bar"}).
					Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
			expectStruct:   true,
		},
		"POST error": {
			body: `{"foo":"bar"}`,
			toBodyFn: func(b string) any {
				return b
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().
					SetBody(math.Inf(1)).
					Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      true,
			expectStruct:   true,
		},
		"POST success string": {
			body: "foobar",
			toBodyFn: func(b string) any {
				return b
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().SetBody(body).Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST success bytes": {
			body: "foobar",
			toBodyFn: func(b string) any {
				return []byte(b)
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().SetBody(body).Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST success reader": {
			body: "foobar",
			toBodyFn: func(b string) any {
				return strings.NewReader(b)
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().SetBody(body).Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST success readcloser": {
			body: "foobar",
			toBodyFn: func(b string) any {
				return io.NopCloser(strings.NewReader(b))
			},
			fn: func(c *Client, body any) *Response {
				resp, _ := c.NR().SetBody(body).Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp := tc.fn(c, tc.toBodyFn(tc.body))
			if resp.Err != nil && !tc.expectErr {
				t.Fatalf("no transport error expected, got: %s", resp.Err)
			}
			if tc.expectErr {
				if resp.Err == nil {
					t.Errorf("expected error got nil")
				}

				_, err := resp.Request.GetBody()
				if err == nil {
					t.Fatalf("expected error got nil")
				}

				if !errors.Is(err, ErrInvalidRequestBody) {
					t.Fatalf("expected error to be ErrInvalidRequestBody, got: %s", err)
				}

				return
			}

			receivedBody, err := resp.String()
			if err != nil {
				t.Fatal("should not err")
			}
			if tc.expectErr {
				_, err := resp.Request.GetBody()
				if err == nil {
					t.Errorf("expected error got nil")
					return
				}
			}

			if receivedBody != tc.body {
				t.Fatalf("expected send %s to equal got %s", tc.body, receivedBody)
			}

			gotMethod := resp.GetHeader("method")
			if gotMethod != tc.expectedMethod {
				t.Fatalf("expected method POST, got '%s'", gotMethod)
			}
		})
	}
}

func TestRequestMethods(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.Header().Add("method", r.Method)
			w.WriteHeader(http.StatusOK)
		default:
			w.Header().Add("method", r.Method)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	s := createTestServer(testHandler)
	defer s.Close()

	c := createTestClient(s.URL)

	type tcase struct {
		fn             func(c *Client) *Response
		expectedMethod string
		expectedCode   int
		expectErr      bool
	}
	tests := map[string]tcase{
		"GET success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodGet, "/success")
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   200,
			expectErr:      false,
		},
		"GET success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodGet, "/:path")
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   200,
			expectErr:      false,
		},
		"GET success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Get("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   200,
			expectErr:      false,
		},
		"GET on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Get(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   200,
			expectErr:      false,
		},
		"GET error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodGet, "/")
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   500,
			expectErr:      true,
		},
		"GET error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Get("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodGet,
			expectedCode:   500,
			expectErr:      true,
		},
		"POST success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPost, "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodPost, "/:path")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Post("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Post(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   200,
			expectErr:      false,
		},
		"POST error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPost, "/")
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   500,
			expectErr:      true,
		},
		"POST error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Post("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPost,
			expectedCode:   500,
			expectErr:      true,
		},
		"PUT success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPut, "/success")
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   200,
			expectErr:      false,
		},
		"PUT success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodPut, "/:path")
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   200,
			expectErr:      false,
		},
		"PUT success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Put("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   200,
			expectErr:      false,
		},
		"PUT on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Put(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   200,
			expectErr:      false,
		},
		"PUT error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPut, "/")
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   500,
			expectErr:      true,
		},
		"PUT error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Put("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPut,
			expectedCode:   500,
			expectErr:      true,
		},
		"PATCH success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPatch, "/success")
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   200,
			expectErr:      false,
		},
		"PATCH success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodPatch, "/:path")
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   200,
			expectErr:      false,
		},
		"PATCH success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Patch("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   200,
			expectErr:      false,
		},
		"PATCH on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Patch(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   200,
			expectErr:      false,
		},
		"PATCH error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodPatch, "/")
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   500,
			expectErr:      true,
		},
		"PATCH error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Patch("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodPatch,
			expectedCode:   500,
			expectErr:      true,
		},
		"DELETE success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodDelete, "/success")
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   200,
			expectErr:      false,
		},
		"DELETE success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodDelete, "/:path")
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   200,
			expectErr:      false,
		},
		"DELETE success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Delete("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   200,
			expectErr:      false,
		},
		"DELETE on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Delete(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   200,
			expectErr:      false,
		},
		"DELETE error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodDelete, "/")
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   500,
			expectErr:      true,
		},
		"DELETE error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Delete("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodDelete,
			expectedCode:   500,
			expectErr:      true,
		},
		"OPTION success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodHead, "/success")
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   200,
			expectErr:      false,
		},
		"OPTION success params": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().
					SetPathParam("path", "success").
					Send(t.Context(), http.MethodHead, "/:path")
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   200,
			expectErr:      false,
		},
		"OPTION success execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Head("/success").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   200,
			expectErr:      false,
		},
		"OPTION on request success": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Head(t.Context(), "/success")
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   200,
			expectErr:      false,
		},
		"OPTION error": {
			fn: func(c *Client) *Response {
				resp, _ := c.NR().Send(t.Context(), http.MethodHead, "/")
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   500,
			expectErr:      true,
		},
		"OPTION error execute": {
			fn: func(c *Client) *Response {
				resp, _ := c.Head("/").Execute(t.Context())
				return resp
			},
			expectedMethod: http.MethodHead,
			expectedCode:   500,
			expectErr:      true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp := tc.fn(c)
			if tc.expectErr && !resp.IsError() {
				t.Fatal("expected err bot got nil")
			}

			gotMethod := resp.GetHeader("method")
			if gotMethod != tc.expectedMethod {
				t.Fatalf("expected method GET, got '%s'", gotMethod)
			}

			gotStatusCode := resp.StatusCode()
			if gotStatusCode != resp.StatusCode() {
				t.Fatalf("expected status code %d, got '%d'", tc.expectedCode, gotStatusCode)
			}
		})
	}
}
