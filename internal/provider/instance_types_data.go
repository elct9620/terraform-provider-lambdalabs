package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &instanceTypesData{}
	_ datasource.DataSourceWithConfigure = &instanceTypesData{}
)

type instanceTypesData struct {
	client *api.Client
}

type instanceTypeSpecsModel struct {
	VCPUs      types.Int64 `tfsdk:"vcpus"`
	MemoryGiB  types.Int64 `tfsdk:"memory_gib"`
	StorageGiB types.Int64 `tfsdk:"storage_gib"`
	GPUs       types.Int64 `tfsdk:"gpus"`
}

type instanceTypeModel struct {
	Name              types.String          `tfsdk:"name"`
	Description       types.String          `tfsdk:"description"`
	GPUDescription    types.String          `tfsdk:"gpu_description"`
	PriceCentsPerHour types.Int64          `tfsdk:"price_cents_per_hour"`
	Specs             *instanceTypeSpecsModel `tfsdk:"specs"`
}

type instanceTypesDataModel struct {
	ID            types.String                     `tfsdk:"id"`
	InstanceTypes map[string]*instanceTypeModel    `tfsdk:"instance_types"`
}

func NewInstanceTypesData() datasource.DataSource {
	return &instanceTypesData{}
}

func (d *instanceTypesData) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_types"
}

func (d *instanceTypesData) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Instance Types Data",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"instance_types": schema.MapNestedAttribute{
				Description: "Available instance types",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Instance type name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Instance type description",
							Computed:    true,
						},
						"gpu_description": schema.StringAttribute{
							Description: "GPU description",
							Computed:    true,
						},
						"price_cents_per_hour": schema.Int64Attribute{
							Description: "Price in cents per hour",
							Computed:    true,
						},
						"specs": schema.SingleNestedAttribute{
							Description: "Instance specifications",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"vcpus": schema.Int64Attribute{
									Description: "Number of virtual CPUs",
									Computed:    true,
								},
								"memory_gib": schema.Int64Attribute{
									Description: "Memory in GiB",
									Computed:    true,
								},
								"storage_gib": schema.Int64Attribute{
									Description: "Storage in GiB",
									Computed:    true,
								},
								"gpus": schema.Int64Attribute{
									Description: "Number of GPUs",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *instanceTypesData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *instanceTypesData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model instanceTypesDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListInstanceTypes(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list instance types", err.Error())
		return
	}

	model.ID = types.StringValue("instance_types")
	model.InstanceTypes = make(map[string]*instanceTypeModel)

	for name, info := range res.Data {
		instanceType := info.InstanceType
		model.InstanceTypes[name] = &instanceTypeModel{
			Name:              types.StringValue(instanceType.Name),
			Description:       types.StringValue(instanceType.Description),
			GPUDescription:    types.StringValue(instanceType.GPUDescription),
			PriceCentsPerHour: types.Int64Value(int64(instanceType.PriceCentsPerHour)),
			Specs: &instanceTypeSpecsModel{
				VCPUs:      types.Int64Value(int64(instanceType.Specs.VCPUs)),
				MemoryGiB:  types.Int64Value(int64(instanceType.Specs.MemoryGiB)),
				StorageGiB: types.Int64Value(int64(instanceType.Specs.StorageGiB)),
				GPUs:       types.Int64Value(int64(instanceType.Specs.GPUs)),
			},
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
