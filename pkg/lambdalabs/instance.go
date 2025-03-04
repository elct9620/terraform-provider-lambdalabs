package lambdalabs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
)

type LaunchInstancePayload struct {
	Name             *string  `json:"name"`
	RegionName       string   `json:"region_name"`
	InstanceTypeName string   `json:"instance_type_name"`
	SSHKeyNames      []string `json:"ssh_key_names"`
	FileSystemNames  []string `json:"file_system_names,omitempty"`
}

type TerminateInstancePayload struct {
	IDs []string `json:"instance_ids"`
}

type InstanceOption = func(payload *LaunchInstancePayload)

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

// LaunchInstance creates a new instance with the specified configuration
func (c *Client) LaunchInstance(ctx context.Context, regionName, instanceTypeName string, sshKeyNames []string, options ...InstanceOption) (*Instance, error) {
	payload := LaunchInstancePayload{
		RegionName:       regionName,
		InstanceTypeName: instanceTypeName,
		SSHKeyNames:      sshKeyNames,
	}

	for _, option := range options {
		option(&payload)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.Post("/instance-operations/launch", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			IDs []string `json:"instance_ids"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	res, err := c.RetrieveInstance(ctx, &RetrieveInstanceRequest{Id: data.Data.IDs[0]})
	if err != nil {
		return nil, err
	}

	return &res.Data, nil
}

// WithInstanceName sets the instance name in the payload
func WithInstanceName(name string) InstanceOption {
	return func(payload *LaunchInstancePayload) {
		payload.Name = &name
	}
}

// WithFileSystemNames sets the file system names in the payload
func WithFileSystemNames(names []string) InstanceOption {
	return func(payload *LaunchInstancePayload) {
		payload.FileSystemNames = names
	}
}

// TerminateInstance terminates an existing instance by ID
func (c *Client) TerminateInstance(id string) (*Instance, error) {
	body, err := json.Marshal(TerminateInstancePayload{IDs: []string{id}})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post("/instance-operations/terminate", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Instances []*Instance `json:"terminated_instances"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Instances) == 0 {
		return nil, nil
	}

	return data.Data.Instances[0], nil
}
