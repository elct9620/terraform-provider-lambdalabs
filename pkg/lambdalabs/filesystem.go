package lambdalabs

import (
	"bytes"
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

// CreateFileSystemRequest represents the request to create a new file system
type CreateFileSystemRequest struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// CreateFileSystemResponse represents the response from the Create File System API
type CreateFileSystemResponse struct {
	Data FileSystem `json:"data"`
}

// CreateFileSystem creates a new file system
func (c *Client) CreateFileSystem(ctx context.Context, req *CreateFileSystemRequest) (*CreateFileSystemResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/file-systems", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var res CreateFileSystemResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteFileSystemRequest represents the request to delete a file system
type DeleteFileSystemRequest struct {
	ID string `json:"id"`
}

// DeleteFileSystemResponse represents the response from the Delete File System API
type DeleteFileSystemResponse struct {
	Data struct {
		DeletedIDs []string `json:"deleted_ids"`
	} `json:"data"`
}

// DeleteFileSystem deletes a file system by ID
func (c *Client) DeleteFileSystem(ctx context.Context, req *DeleteFileSystemRequest) (*DeleteFileSystemResponse, error) {
	resp, err := c.Delete(ctx, "/file-systems/"+req.ID, nil)
	if err != nil {
		return nil, err
	}

	var res DeleteFileSystemResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
