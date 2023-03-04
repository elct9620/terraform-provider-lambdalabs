package lambdalabs

import (
	"errors"
	"io"
	"net/http"
)

const Endpoint = "https://cloud.lambdalabs.com/api/v1"

var (
	ErrUnauthorized = errors.New("Unauthorized")
	ErrForbidden    = errors.New("Forbidden")
	ErrBadRequest   = errors.New("Bad Request")
)

type Client struct {
	apiKey   string
	endpoint string
	http     *http.Client
}

type ClientOption = func(c *Client)

func New(apiKey string, options ...ClientOption) *Client {
	client := &Client{
		endpoint: Endpoint,
		apiKey:   apiKey,
		http:     &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) {
		c.endpoint = endpoint
	}
}

func (c *Client) Get(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.endpoint+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusForbidden:
		return nil, ErrForbidden
	}

	return resp, nil
}

func (c *Client) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.endpoint+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrUnauthorized
	case http.StatusForbidden:
		return nil, ErrForbidden
	case http.StatusBadRequest:
		return nil, ErrBadRequest
	}

	return resp, nil

}
