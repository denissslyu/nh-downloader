package lrr_rpc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"nh-downloader/internal/config"
)

type HttpClient struct {
	url     string
	method  string
	headers map[string]string
	params  map[string]string
	body    []byte
}

func (c *HttpClient) URL(url string) *HttpClient {
	c.url = config.Lanraragi().Url + url
	return c
}

func newHttpClient() *HttpClient {
	client := &HttpClient{
		headers: map[string]string{
			"Accept": "application/json",
		},
	}
	if config.Lanraragi().KeyBase64 != "" {
		client.headers["Authorization"] = fmt.Sprintf("Bearer %s", config.Lanraragi().KeyBase64)
	}
	return client
}

func (c *HttpClient) Method(method string) *HttpClient {
	c.method = method
	return c
}

func (c *HttpClient) Params(params map[string]string) *HttpClient {
	c.params = params
	return c
}

func (c *HttpClient) Body(body []byte) *HttpClient {
	c.body = body
	return c
}

func (c *HttpClient) send() ([]byte, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	queryParams := u.Query()
	for key, value := range c.params {
		queryParams.Add(key, value)
	}

	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest(c.method, u.String(), bytes.NewReader(c.body))
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
