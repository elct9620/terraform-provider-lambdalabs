package lambdalabs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFirewallRules(t *testing.T) {
	cases := []struct {
		name     string
		response string
		expected *ListFirewallRulesResponse
		err      bool
	}{
		{
			name: "success",
			response: `{
				"data": [
					{
						"protocol": "tcp",
						"port_range": [22, 22],
						"source_network": "0.0.0.0/0",
						"description": "Allow SSH from anywhere"
					}
				]
			}`,
			expected: &ListFirewallRulesResponse{
				Data: []FirewallRule{
					{
						Protocol:      "tcp",
						PortRange:     []int{22, 22},
						SourceNetwork: "0.0.0.0/0",
						Description:   "Allow SSH from anywhere",
					},
				},
			},
			err: false,
		},
		{
			name: "unauthorized",
			response: `{
				"error": {
					"code": "global/invalid-api-key",
					"message": "API key was invalid, expired, or deleted.",
					"suggestion": "Check your API key or create a new one, then try again."
				}
			}`,
			expected: nil,
			err:      true,
		},
		{
			name: "forbidden",
			response: `{
				"error": {
					"code": "global/account-inactive",
					"message": "Your account is inactive.",
					"suggestion": "Make sure you have verified your email address and have a valid payment method. Contact Support if problems continue."
				}
			}`,
			expected: nil,
			err:      true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/firewall-rules" {
					t.Errorf("Expected path %q, got %q", "/firewall-rules", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Errorf("Expected method %q, got %q", http.MethodGet, r.Method)
				}

				if c.name == "unauthorized" {
					w.WriteHeader(http.StatusUnauthorized)
				} else if c.name == "forbidden" {
					w.WriteHeader(http.StatusForbidden)
				} else {
					w.WriteHeader(http.StatusOK)
				}
				_, _ = w.Write([]byte(c.response))
			}))
			defer server.Close()

			client := New("test-api-key", WithBaseUrl(server.URL))
			res, err := client.ListFirewallRules(context.Background())

			if c.err {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}

				if len(res.Data) != len(c.expected.Data) {
					t.Errorf("Expected %d firewall rules, got %d", len(c.expected.Data), len(res.Data))
				}

				for i, rule := range res.Data {
					expectedRule := c.expected.Data[i]
					if rule.Protocol != expectedRule.Protocol {
						t.Errorf("Expected protocol %q, got %q", expectedRule.Protocol, rule.Protocol)
					}
					if rule.SourceNetwork != expectedRule.SourceNetwork {
						t.Errorf("Expected source network %q, got %q", expectedRule.SourceNetwork, rule.SourceNetwork)
					}
					if rule.Description != expectedRule.Description {
						t.Errorf("Expected description %q, got %q", expectedRule.Description, rule.Description)
					}
					if len(rule.PortRange) != len(expectedRule.PortRange) {
						t.Errorf("Expected port range length %d, got %d", len(expectedRule.PortRange), len(rule.PortRange))
					} else {
						for j, port := range rule.PortRange {
							if port != expectedRule.PortRange[j] {
								t.Errorf("Expected port range[%d] %d, got %d", j, expectedRule.PortRange[j], port)
							}
						}
					}
				}
			}
		})
	}
}
