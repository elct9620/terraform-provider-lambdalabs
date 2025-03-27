package lambdalabs

import (
	"bytes"
	"context"
	"encoding/json"
)

// FirewallRule represents a firewall rule in the Lambda Labs API
type FirewallRule struct {
	Protocol      string `json:"protocol"`
	PortRange     []int  `json:"port_range"`
	SourceNetwork string `json:"source_network"`
	Description   string `json:"description"`
}

// ListFirewallRulesResponse represents the response from the List Firewall Rules API
type ListFirewallRulesResponse struct {
	Data []FirewallRule `json:"data"`
}

// ListFirewallRules retrieves all firewall rules for the authenticated user
func (c *Client) ListFirewallRules(ctx context.Context) (*ListFirewallRulesResponse, error) {
	resp, err := c.Get(ctx, "/firewall-rules", nil)
	if err != nil {
		return nil, err
	}

	var res ListFirewallRulesResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ReplaceFirewallRulesRequest represents the request to replace all firewall rules
type ReplaceFirewallRulesRequest struct {
	Data []FirewallRule `json:"data"`
}

// ReplaceFirewallRulesResponse represents the response from the Replace Firewall Rules API
type ReplaceFirewallRulesResponse struct {
	Data []FirewallRule `json:"data"`
}

// ReplaceFirewallRules replaces all inbound firewall rules with the provided rules
func (c *Client) ReplaceFirewallRules(ctx context.Context, req *ReplaceFirewallRulesRequest) (*ReplaceFirewallRulesResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.Put(ctx, "/firewall-rules", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var res ReplaceFirewallRulesResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
