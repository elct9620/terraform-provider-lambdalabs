package provider

import (
	"context"
	"os"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &lambdalabsProvider{}
)

type lambdalabsProvider struct {
	version string
}

type lambdalabsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &lambdalabsProvider{version}
	}
}

func (p *lambdalabsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lambdalabs"
	resp.Version = p.version
}

func (p *lambdalabsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Lambdalabs Cloud GPU",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Lambdalabs API Endpoint",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API Key from Lambdalabs",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *lambdalabsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config lambdalabsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Lambdalabs API Endpoint",
			"The provider cannot create the Lambdalabs API client as there is an unknown configuration value for the Lambdalabs API Endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LAMBDALABS_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Lambdalabs API Key",
			"The provider cannot create the Lambdalabs API client as there is an unknown configuration value for the Lambdalabs API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LAMBDALABS_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("LAMBDALABS_ENDPOINT")
	apiKey := os.Getenv("LAMBDALABS_API_KEY")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if endpoint == "" {
		endpoint = api.Endpoint
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Lambdalabs API Key",
			"The provider cannot create the Lambdalabs API client as there is a missing or empty value for the Lambda API Key. "+
				"Set the api key value in the configuration or use the LAMBDALABS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := api.New(apiKey, api.WithEndpoint(endpoint))

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *lambdalabsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *lambdalabsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHKeyResource,
		NewInstanceResource,
	}
}
