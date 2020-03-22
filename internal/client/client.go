package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	timeout = 240
	agent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"
	noCache = "no-cache"
)

type CustomClient struct {
	http.Client
}

// NewClient will return a new pre configured AgentHTTP
func NewClient() CustomClient {
	return CustomClient{
		Client: http.Client{
			Timeout: time.Second * timeout,
		},
	}
}

// Get receives an string url and will perform an GET HTTP Request on that URL and return the response body
func (c *CustomClient) Get(target string) (body io.ReadCloser, err error) {
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return body, fmt.Errorf("client.HTTP.Get.NewRequest: %w", err)
	}
	addHeaders(req)
	body, err = c.execute(req)
	if err != nil {
		return body, fmt.Errorf("client.HTTP.Get.Do: %w", err)
	}
	return body, nil
}

// execute will receive an http.Request, performs that call and returns the body content
func (c *CustomClient) execute(req *http.Request) (io.ReadCloser, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.execute.Do: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client.execute.Status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// addHeaders will add the default HTTP headers to the request
func addHeaders(req *http.Request) {
	req.Header.Add("Cache-Control", noCache)
	req.Header.Add("Pragma", noCache)
	req.Header.Add("User-Agent", agent)
	req.Header.Add("Accept", "*/*")
}
