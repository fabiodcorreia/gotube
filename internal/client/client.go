package client

import (
	"fmt"
	"io"
	"net/http"
)

const timeout = 240

// Client is responsible for holding a pre configured http.Client
type Client struct {
	client *http.Client
}

// NewClient will return a new pre configured AgentHTTP
func NewClient() Client {
	return Client{
		client: &http.Client{
			//Timeout: time.Second * timeout,
		},
	}
}

// Get receives an string url and will perform an GET HTTP Request on that URL and return the response body
func (c *Client) Get(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client.HTTP.Get.NewRequest: %w", err)
	}
	addHeaders(req)
	body, err := c.execute(req)
	if err != nil {
		return nil, fmt.Errorf("client.HTTP.Get.Do: %w", err)
	}
	return body, nil
}

// execute will receive an http.Request, performs that call and returns the body content
func (c *Client) execute(req *http.Request) (io.ReadCloser, error) {
	resp, err := c.client.Do(req)
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
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36")
	req.Header.Add("Accept", "*/*")
}
