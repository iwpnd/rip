package rip

import (
	"fmt"
	"net/http"
	gourl "net/url"
	"strings"
)

type (
	RequestMiddleware  func(c *Client, r *Request) error
	ResponseMiddleware func(c *Client, r *Response) error
)

func parseRequestUrl(c *Client, r *Request) error {
	tempURL := r.RawURL
	if len(r.pathParams) > 0 {
		for p, v := range r.pathParams {
			tempURL = strings.ReplaceAll(tempURL, ":"+p, gourl.PathEscape(v))
		}
	}
	if len(c.pathParams) > 0 {
		for p, v := range c.pathParams {
			tempURL = strings.ReplaceAll(tempURL, ":"+p, gourl.PathEscape(v))
		}
	}
	url, err := gourl.Parse(tempURL)
	if err != nil {
		return err
	}

	// we go the url from the request method. so we merge with baseurl of client
	if !url.IsAbs() {
		tempURL = url.String()
		if tempURL != "" {
			tempURL = "/" + strings.Trim(tempURL, "/")
		}
		url, err = gourl.Parse(c.baseURL + tempURL)
		if err != nil {
			return err
		}
	}

	query := gourl.Values{}
	for k, v := range c.queryParams {
		for _, p := range v {
			query.Add(k, p)
		}
	}

	for k, v := range r.queryParams {
		query.Del(k)

		for _, p := range v {
			query.Add(k, p)
		}
	}

	if len(query) > 0 {
		url.RawQuery = query.Encode()
	}

	r.URL = url
	return nil
}

func parseResponseBody(c *Client, r *Response) error {
	if r.rawResponse == nil || r.IsError() {
		return nil
	}

	if r.rawResponse.StatusCode == http.StatusNoContent {
		return nil
	}

	if r.Request.Target != nil {
		err := Unmarshal(c, r, r.Request.Target)
		if err != nil {
			return err
		}
	}

	return nil
}

func debugResponse(client *Client, resp *Response) error {
	fmt.Printf("response from: %s\n", resp.Request.URL)
	fmt.Printf("response received code: %d\n", resp.StatusCode())
	if len(resp.Request.Body) != 0 {
		fmt.Println(resp.String())
	}

	fmt.Printf("isSuccess: %v\n", resp.IsSuccess())
	if resp.IsError() {
		fmt.Printf("isError: %s\n", resp.Err)
	}

	return nil
}
