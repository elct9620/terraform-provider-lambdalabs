package lambdalabs

import (
	"bytes"
	"encoding/json"
	"io"
)

type Instance struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}

type LaunchInstancePayload struct {
	RegionName       string   `json:"region_name"`
	InstanceTypeName string   `json:"instance_type_name"`
	SSHKeyNames      []string `json:"ssh_key_names"`
}

type TerminateInstancePayload struct {
	IDs []string `json:"instance_ids"`
}

func (c *Client) GetInstance(id string) (*Instance, error) {
	resp, err := c.Get("/instances/"+id, nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data *Instance `json:"data"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}

func (c *Client) LaunchInstance(regionName, instanceTypeName string, sshKeyNames []string) (*Instance, error) {
	body, err := json.Marshal(LaunchInstancePayload{RegionName: regionName, InstanceTypeName: instanceTypeName, SSHKeyNames: sshKeyNames})
	if err != nil {
		return nil, err
	}

	resp, err := c.Post("/instance-operations/launch", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	var data struct {
		Data struct {
			IDs []string `json:"instance_ids"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)

	instance, err := c.GetInstance(data.Data.IDs[0])
	if err != nil {
		return nil, err
	}

	return instance, nil
}

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
	var data struct {
		Data struct {
			Instances []*Instance `json:"terminated_instances"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)

	return data.Data.Instances[0], nil
}
