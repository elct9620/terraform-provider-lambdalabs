package lambdalabs

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
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"public_key": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"private_key": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
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
	if key.PublicKey.String() == "" {
		createdKey, err = r.client.CreateSSHKey(key.Name.String())
	} else {
		createdKey, err = r.client.CreateSSHKeyWithPublicKey(key.Name.String(), key.PublicKey.String())
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

	key, err := r.client.GetSSHKey(state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lambdalabs SSH Key",
			"Could not find Lambdalabs SSH Key ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	newState := sshKeyModel{
		ID:         types.StringValue(key.ID),
		Name:       types.StringValue(key.Name),
		PublicKey:  types.StringValue(key.PublicKey),
		PrivateKey: types.StringValue(key.PrivateKey),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Error Update Lambdalabs SSH Key",
		"Unsupported Method",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Error Delete Lambdalabs SSH Key",
		"Unsupported Method",
	)
}
