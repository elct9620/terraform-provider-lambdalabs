package provider_test

import (
	"github.com/elct9620/terraform-provider-lambdalabs/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"lambdalabs": providerserver.NewProtocol6WithError(provider.New("test")()),
}
