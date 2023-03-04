package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &sshKeyResource{}
	_ resource.ResourceWithConfigure = &sshKeyResource{}
)

type sshKeyResource struct {
	client *api.Client
}

type sshKeyModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func NewSSHKeyResource() resource.Resource {
	return &sshKeyResource{}
}

// Metadata returns the resource type name.
func (r *sshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

// Schema defines the schema for the resource.
func (r *sshKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "SSH Key ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The SSH Key name",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "The public key to install into instance",
				Optional:            true,
				Computed:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "If public key not given the Lambdalabs will generated one and return in this field",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *sshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *sshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var key sshKeyModel
	diags := req.Plan.Get(ctx, &key)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var createdKey *api.SSHKey
	var err error
	if key.PublicKey.ValueString() == "" {
		createdKey, err = r.client.CreateSSHKey(key.Name.ValueString())
	} else {
		createdKey, err = r.client.CreateSSHKeyWithPublicKey(key.Name.ValueString(), key.PublicKey.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SSH Key",
			"Could not create SSH Key, unexpected error: "+err.Error(),
		)
		return
	}

	key.ID = types.StringValue(createdKey.ID)
	key.Name = types.StringValue(createdKey.Name)
	key.PublicKey = types.StringValue(createdKey.PublicKey)
	key.PrivateKey = types.StringValue(createdKey.PrivateKey)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, key)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sshKeyModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.GetSSHKey(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lambdalabs SSH Key",
			"Could not find Lambdalabs SSH Key ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(key.ID)
	state.Name = types.StringValue(key.Name)
	state.PublicKey = types.StringValue(key.PublicKey)
	state.PrivateKey = types.StringValue(key.PrivateKey)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Lambdalabs SSH Key",
		"Unsupported Method",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Error Delete Lambdalabs SSH Key",
		"Unsupported Method",
	)
}
