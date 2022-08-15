# rip

rip is a minimal api client framework i just want use to REST in peace.

## Installation

```
go get -u github.com/iwpnd/rip
```

## Usage

```go
package main

import (
  "encoding/json"

  "github.com/iwpnd/rip"
  )

type BlogPost struct {
    Id string
    Content string
  }

type BlogApiClient struct {
    *rip.Client
}

func NewBlogApiClient(host string, options ...rip.Option) (*BlogApiClient, error) {
    c, err := rip.NewClient(host, options)
    if err != nil {
        return &BlogApiClient{}, err
    }
    return &BlogApiClient{c}, nil
}

func (c *BlogApiClient) GetById(id string) (*BlogPost, error) {
    req, err := c.NR().
        SetHeader("Accept", "application/json").
        SetParams(Params{"id": "1"})

    res, err := req.Execute("GET", "/blog/:id")
    if err != nil {
        return &BlogPost{}, err
    }

    b := &BlogPost{}
    err = rip.Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
    if err != nil {
        return &BlogPost, err
    }

    return b
}

func (c *BlogApiClient) Create(post BlogPost) (*BlogPost, error) {
    b, err := json.Marshal(post)
    if err != nil {
        return &BlogPost, err
    }

    req := c.NR().
        SetHeaders(rip.Header{
          "Content-Type": "application/json",
          "Accept":       "application/json",
          }).
        SetBody(b)


    res, err := req.Execute("POST", "/blog")
    if err != nil {
        return &Response{}, err
    }

    b := &BlogPost{}
    err = rip.Unmarshal(res.Header().Get("Content-Type"), res.Body(), r)
    if err != nil {
        return &BlogPost, err
    }

    return b
}

func main() {
    c, err := NewBlogApiClient(
      "https://myblog.io",
      rip.WithDefaultHeaders(rip.Header{"x-api-key": os.Getenv("API_KEY_BLOGAPI")}),
      rip.WithTimeout(30)
    )
    if err != nil {
        panic("AAAH!")
    }

    b, err := c.GetById(id)
    if err != nil {
        t.Errorf("could not get blogpost for id %v :(", id)
    }

    fmt.Printf("blogpost: \n %v\n\n", b)

}
```

## License

MIT

## Maintainer

Benjamin Ramser - [@iwpnd](https://github.com/iwpnd)

Project Link: [https://github.com/iwpnd/rip](https://github.com/iwpnd/rip)
