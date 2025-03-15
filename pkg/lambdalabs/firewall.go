package lambdalabs

import (
	"context"
	"encoding/json"
)

// FirewallRule represents a firewall rule in the Lambda Labs API
type FirewallRule struct {
	Protocol      string   `json:"protocol"`
	PortRange     []int    `json:"port_range"`
	SourceNetwork string   `json:"source_network"`
	Description   string   `json:"description"`
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
