package lambdalabs

import (
	"context"
	"encoding/json"
)

// ListFileSystemsResponse represents the response from the List File Systems API
type ListFileSystemsResponse struct {
	Data []FileSystem `json:"data"`
}

// ListFileSystems retrieves all file systems for the authenticated user
func (c *Client) ListFileSystems(ctx context.Context) (*ListFileSystemsResponse, error) {
	resp, err := c.Get(ctx, "/file-systems", nil)
	if err != nil {
		return nil, err
	}

	var res ListFileSystemsResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
