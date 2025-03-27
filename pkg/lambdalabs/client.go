package lambdalabs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const BaseUrl = "https://cloud.lambdalabs.com/api/v1"

type Client struct {
	baseUrl string
	*http.Client
}

type ClientOption = func(c *Client)

func New(apiKey string, options ...ClientOption) *Client {
	client := &Client{
		baseUrl: BaseUrl,
		Client: &http.Client{
			Transport: &Transport{
				apiKey: apiKey,
			},
		},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func WithBaseUrl(baseUrl string) ClientOption {
	return func(c *Client) {
		c.baseUrl = baseUrl
	}
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

func assertError(resp *http.Response) (*http.Response, error) {
	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}

	var errorResponse ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return nil, err
	}

	return nil, &errorResponse.Error
}

func (c *Client) Get(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return assertError(resp)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return assertError(resp)
}

func (c *Client) Delete(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return assertError(resp)
}

func (c *Client) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return assertError(resp)
}
