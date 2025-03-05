package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &sshKeyData{}
	_ datasource.DataSourceWithConfigure = &sshKeyData{}
)

type sshKeyData struct {
	client *api.Client
}

type sshKeyDataModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PublicKey types.String `tfsdk:"public_key"`
}

func NewSshKeyData() datasource.DataSource {
	return &sshKeyData{}
}

func (d *sshKeyData) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (d *sshKeyData) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SSH Key Data",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "SSH Key ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The SSH Key name",
				Optional:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "The public key to install into instance",
				Computed:    true,
			},
		},
	}
}

func (d *sshKeyData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *sshKeyData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model sshKeyDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListSshKeys(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list ssh keys", err.Error())
		return
	}

	keyName := model.Name.ValueString()

	for _, key := range res.Data {
		if key.Name == keyName {
			model.Id = types.StringValue(key.Id)
			model.PublicKey = types.StringValue(key.PublicKey)
			resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
			return
		}
	}

	resp.Diagnostics.AddError("ssh key not found", "The ssh key with name "+keyName+" not found")
}
