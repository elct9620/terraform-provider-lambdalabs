package main

import (
	"context"

	"github.com/elct9620/terraform-provider-lambdalabs/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), lambdalabs.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/elct9620/lambdalabs",
	})
}
