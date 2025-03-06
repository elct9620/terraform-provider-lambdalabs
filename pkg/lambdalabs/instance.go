package lambdalabs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

type RetrieveInstanceRequest struct {
	Id string `json:"id"`
}

type RetrieveInstanceResponse struct {
	Data Instance `json:"data"`
}

// RetrieveInstance to get the instance details
func (c *Client) RetrieveInstance(ctx context.Context, req *RetrieveInstanceRequest) (*RetrieveInstanceResponse, error) {
	resp, err := c.Get(ctx, "/instances/"+req.Id, nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res RetrieveInstanceResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

type ListInstanceTypesResponse struct {
	Data map[string]InstanceTypeInfo `json:"data"`
}

// ListInstanceTypes returns a list of available instance types
func (c *Client) ListInstanceTypes(ctx context.Context) (*ListInstanceTypesResponse, error) {
	resp, err := c.Get(ctx, "/instance-types", nil)
	if err != nil {
		return nil, err
	}

	var res ListInstanceTypesResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type LaunchInstanceRequest struct {
	Name             *string  `json:"name,omitempty"`
	RegionName       string   `json:"region_name"`
	InstanceTypeName string   `json:"instance_type_name"`
	SSHKeyNames      []string `json:"ssh_key_names"`
	FileSystemNames  []string `json:"file_system_names,omitempty"`
}

type LaunchInstanceResponse struct {
	Data struct {
		IDs []string `json:"instance_ids"`
	} `json:"data"`
}

// LaunchInstance creates a new instance with the specified configuration
func (c *Client) LaunchInstance(ctx context.Context, req *LaunchInstanceRequest) (*LaunchInstanceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/instance-operations/launch", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var res LaunchInstanceResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

type TerminateInstanceRequest struct {
	Ids []string `json:"instance_ids"`
}

type TerminateInstanceResponse struct {
	Data struct {
		TerminatedInstances []*Instance `json:"terminated_instances"`
	} `json:"data"`
}

// TerminateInstance terminates an existing instance by ID
func (c *Client) TerminateInstance(ctx context.Context, req *TerminateInstanceRequest) (*TerminateInstanceResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post(ctx, "/instance-operations/terminate", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var res TerminateInstanceResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
