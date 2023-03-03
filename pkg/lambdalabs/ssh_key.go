package lambdalabs

import "errors"

var (
	ErrSSHKeyNotFound = errors.New("SSH Key not exists")
)

type SSHKey struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (c *Client) ListSSHKeys() ([]*SSHKey, error) {
	return []*SSHKey{}, nil
}

func (c *Client) GetSSHKey(id string) (*SSHKey, error) {
	keys, err := c.ListSSHKeys()
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

func (c *Client) CreateSSHKey(name string) (*SSHKey, error) {
	return &SSHKey{}, nil
}

func (c *Client) CreateSSHKeyWithPublicKey(name, publicKey string) (*SSHKey, error) {
	return &SSHKey{}, nil
}
