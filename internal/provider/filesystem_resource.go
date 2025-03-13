package provider

import (
	"context"

	"github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &filesystemResource{}
	_ resource.ResourceWithConfigure = &filesystemResource{}
)

type filesystemResource struct {
	client *api.Client
}

func NewFilesystemResource() resource.Resource {
	return &filesystemResource{}
}

// Metadata returns the resource type name.
func (r *filesystemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem"
}

// Schema defines the schema for the resource.
func (r *filesystemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage File Systems",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "File System ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The File System name",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where the file system will be created",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *filesystemResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *filesystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var fs filesystemResourceModel
	diags := req.Plan.Get(ctx, &fs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := lambdalabs.CreateFileSystemRequest{
		Name:   fs.Name.ValueString(),
		Region: fs.Region.ValueString(),
	}

	res, err := r.client.CreateFileSystem(ctx, &payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating File System",
			"Could not create File System, unexpected error: "+err.Error(),
		)
		return
	}

	fs.ID = types.StringValue(res.Data.ID)
	fs.Name = types.StringValue(res.Data.Name)
	fs.Region = types.StringValue(res.Data.Region.Name)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, fs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *filesystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state filesystemResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filesystems, err := r.client.ListFileSystems(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lambdalabs File System",
			"Could not list Lambdalabs File Systems: "+err.Error(),
		)
		return
	}

	var filesystem *lambdalabs.FileSystem
	filesystemId := state.ID.ValueString()
	for _, fs := range filesystems.Data {
		if fs.ID == filesystemId {
			filesystem = &fs
			break
		}
	}

	if filesystem == nil {
		resp.Diagnostics.AddError(
			"Error Reading Lambdalabs File System",
			"Could not find Lambdalabs File System ID "+filesystemId,
		)
		return
	}

	state.ID = types.StringValue(filesystem.ID)
	state.Name = types.StringValue(filesystem.Name)
	state.Region = types.StringValue(filesystem.Region.Name)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *filesystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update Lambdalabs File System",
		"Unsupported Method",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *filesystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state filesystemResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.DeleteFileSystem(ctx, &lambdalabs.DeleteFileSystemRequest{
		ID: state.ID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Delete Lambdalabs File System",
			"Could not delete Lambdalabs File System ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Check if the file system was actually deleted
	if len(res.Data.DeletedIDs) == 0 {
		resp.Diagnostics.AddError(
			"Error Delete Lambdalabs File System",
			"File System ID "+state.ID.ValueString()+" was not deleted",
		)
		return
	}
}
