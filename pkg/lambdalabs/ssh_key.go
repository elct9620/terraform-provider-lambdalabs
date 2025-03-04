package lambdalabs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrSSHKeyNotFound = errors.New("SSH Key not exists")
)

type SSHKey struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type CreateSSHKeyPayload struct {
	Name string `json:"name"`
}

type CreateSSHKeyWithPKeyPayload struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

func (c *Client) ListSSHKeys(ctx context.Context) ([]*SSHKey, error) {
	resp, err := c.Get(ctx, "/ssh-keys", nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data []*SSHKey `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (c *Client) GetSSHKey(ctx context.Context, id string) (*SSHKey, error) {
	keys, err := c.ListSSHKeys(ctx)
	if err != nil {
		return nil, err
	}

	for idx := range keys {
		if keys[idx].ID == id {
			return keys[idx], nil
		}
	}

	return nil, ErrSSHKeyNotFound
}

func (c *Client) CreateSSHKey(ctx context.Context, name string) (*SSHKey, error) {
	body, err := json.Marshal(CreateSSHKeyPayload{Name: name})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/ssh-keys", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data *SSHKey `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (c *Client) CreateSSHKeyWithPublicKey(ctx context.Context, name, publicKey string) (*SSHKey, error) {
	body, err := json.Marshal(CreateSSHKeyWithPKeyPayload{Name: name, PublicKey: publicKey})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/ssh-keys", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data *SSHKey `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (c *Client) DeleteSSHKey(id string) error {
	_, err := c.Delete("/ssh-keys/"+id, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}

	return nil
}
