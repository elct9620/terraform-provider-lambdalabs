package lambdalabs

import (
	"bytes"
	"context"
	"encoding/json"
)

type SshKey struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type ListSshKeysResponse struct {
	Data []SshKey `json:"data"`
}

func (c *Client) ListSshKeys(ctx context.Context) (*ListSshKeysResponse, error) {
	resp, err := c.Get(ctx, "/ssh-keys", nil)
	if err != nil {
		return nil, err
	}

	var res ListSshKeysResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type CreateSshKeyRequest struct {
	Name      string  `json:"name"`
	PublicKey *string `json:"public_key,omitempty"`
}

type CreateSshKeyResponse struct {
	Data SshKey `json:"data"`
}

func (c *Client) CreateSshKey(ctx context.Context, req *CreateSshKeyRequest) (*CreateSshKeyResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/ssh-keys", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var res CreateSshKeyResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type DeleteSshKeyRequest struct {
	Id string `json:"id"`
}

func (c *Client) DeleteSshKey(ctx context.Context, req *DeleteSshKeyRequest) error {
	_, err := c.Delete(ctx, "/ssh-keys/"+req.Id, nil)
	if err != nil {
		return err
	}

	return nil
}
