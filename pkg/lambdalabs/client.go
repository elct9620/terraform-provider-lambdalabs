package lambdalabs

import (
	"errors"
	"io"
	"net/http"
)

const BaseUrl = "https://cloud.lambdalabs.com/api/v1"

var (
	ErrUnauthorized = errors.New("Unauthorized")
	ErrForbidden    = errors.New("Forbidden")
	ErrBadRequest   = errors.New("Bad Request")
)

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

func (c *Client) Get(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
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
	req, err := http.NewRequest("POST", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
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

func (c *Client) Delete(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
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
