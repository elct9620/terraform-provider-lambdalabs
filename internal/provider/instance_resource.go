package provider

import (
	"context"
	"time"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	helper "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	InstanceStateBooting     string = "booting"
	InstanceStateActive      string = "active"
	InstanceStateContactable string = "contactable"
)

var (
	_                            resource.Resource              = &instanceResource{}
	_                            resource.ResourceWithConfigure = &instanceResource{}
	defaultInstanceCreateTimeout                                = 10 * time.Minute
	instanceCreateDelay                                         = 10 * time.Second
)

type instanceResource struct {
	client *api.Client
}

type instanceModel struct {
	ID               types.String   `tfsdk:"id"`
	Name             types.String   `tfsdk:"name"`
	IP               types.String   `tfsdk:"ip"`
	RegionName       types.String   `tfsdk:"region_name"`
	InstanceTypeName types.String   `tfsdk:"instance_type_name"`
	SSHKeyNames      types.List     `tfsdk:"ssh_key_names"`
	FileSystemNames  types.List     `tfsdk:"file_system_names"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

// Metadata returns the resource type name.
func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema defines the schema for the resource.
func (r *instanceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage instances",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The instance ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The instance name",
				Optional:            true,
			},
			"ip": schema.StringAttribute{
				MarkdownDescription: "The public IP address",
				Computed:            true,
			},
			"region_name": schema.StringAttribute{
				MarkdownDescription: "The instance region name",
				Required:            true,
			},
			"instance_type_name": schema.StringAttribute{
				MarkdownDescription: "The instance type name",
				Required:            true,
			},
			"ssh_key_names": schema.ListAttribute{
				MarkdownDescription: "The SSH Key names to install into instance",
				Required:            true,
				ElementType:         types.StringType,
			},
			"file_system_names": schema.ListAttribute{
				MarkdownDescription: "Optional list of file system names to attach to the instance",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*api.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var instance instanceModel
	diags := req.Plan.Get(ctx, &instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	keyNames := []string{}
	diags = instance.SSHKeyNames.ElementsAs(ctx, &keyNames, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileSystemNames := []string{}
	if !instance.FileSystemNames.IsNull() {
		diags = instance.FileSystemNames.ElementsAs(ctx, &fileSystemNames, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createTimeout, diags := instance.Timeouts.Create(ctx, defaultInstanceCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	options := []api.InstanceOption{}
	if !instance.Name.IsNull() {
		options = append(options, api.WithInstanceName(instance.Name.ValueString()))
	}

	createdInstance, err := r.client.LaunchInstance(
		instance.RegionName.ValueString(),
		instance.InstanceTypeName.ValueString(),
		keyNames,
		options...,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating instance",
			"Could not create instance, unexpected error: "+err.Error(),
		)
		return
	}

	latestInstance, err := r.waitInstanceCreated(ctx, createdInstance.ID, createTimeout)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating instance",
			"Could not create instance, unexpected error: "+err.Error(),
		)
		return
	}

	instance.ID = types.StringValue(latestInstance.ID)
	instance.IP = types.StringValue(latestInstance.IP)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, instance)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	latestInstance, err := r.client.GetInstance(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Lambdalabs SSH Key",
			"Could not find Lambdalabs SSH Key ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.IP = types.StringValue(latestInstance.IP)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Error Update Lambdalabs instance",
		"Unsupported Method",
	)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state instanceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	_, err := r.client.TerminateInstance(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting instance",
			"Could not delete instance, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *instanceResource) waitInstanceCreated(ctx context.Context, id string, createTimeout time.Duration) (*api.Instance, error) {
	changeConfig := &helper.StateChangeConf{
		Pending: []string{
			InstanceStateBooting,
		},
		Target: []string{
			InstanceStateActive,
			InstanceStateContactable,
		},
		Refresh: func() (any, string, error) {
			resp, err := r.client.GetInstance(id)
			if err != nil {
				return nil, "", err
			}
			return resp, resp.Status, nil
		},
		Timeout: createTimeout,
		Delay:   instanceCreateDelay,
	}
	raw, err := changeConfig.WaitForStateContext(ctx)

	if v, ok := raw.(*api.Instance); ok {
		return v, err
	}

	return nil, err
}
