package provider_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func Test_FirewallData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.TrimSpace(r.URL.Path) == "/firewall-rules" {
			resBody := `
			{
				"data": [
					{
						"protocol": "tcp",
						"port_range": [22, 22],
						"source_network": "0.0.0.0/0",
						"description": "Allow SSH from anywhere"
					},
					{
						"protocol": "tcp",
						"port_range": [80, 80],
						"source_network": "0.0.0.0/0",
						"description": "Allow HTTP from anywhere"
					},
					{
						"protocol": "tcp",
						"port_range": [443, 443],
						"source_network": "0.0.0.0/0",
						"description": "Allow HTTPS from anywhere"
					}
				]
			}
			`
			w.Write([]byte(resBody)) //nolint:errcheck
		}
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig(server.URL) + `
				data "lambdalabs_firewall" "default" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "id", "firewall"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.#", "3"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.port_range.#", "2"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.port_range.0", "22"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.port_range.1", "22"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.source_network", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.0.description", "Allow SSH from anywhere"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.port_range.#", "2"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.port_range.0", "80"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.port_range.1", "80"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.source_network", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.1.description", "Allow HTTP from anywhere"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.port_range.#", "2"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.port_range.0", "443"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.port_range.1", "443"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.source_network", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.lambdalabs_firewall.default", "rules.2.description", "Allow HTTPS from anywhere"),
				),
			},
		},
	})
}
