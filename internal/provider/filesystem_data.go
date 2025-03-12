package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &filesystemData{}
	_ datasource.DataSourceWithConfigure = &filesystemData{}
)

type filesystemData struct {
	client *api.Client
}

type filesystemRegionModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type filesystemUserModel struct {
	ID     types.String `tfsdk:"id"`
	Email  types.String `tfsdk:"email"`
	Status types.String `tfsdk:"status"`
}

type filesystemDataModel struct {
	ID         types.String           `tfsdk:"id"`
	Name       types.String           `tfsdk:"name"`
	MountPoint types.String           `tfsdk:"mount_point"`
	Created    types.String           `tfsdk:"created"`
	CreatedBy  *filesystemUserModel   `tfsdk:"created_by"`
	IsInUse    types.Bool             `tfsdk:"is_in_use"`
	Region     *filesystemRegionModel `tfsdk:"region"`
	BytesUsed  types.Int64            `tfsdk:"bytes_used"`
}

type filesystemsFilterModel struct {
	Region types.String `tfsdk:"region"`
}

type filesystemsDataModel struct {
	ID          types.String            `tfsdk:"id"`
	Filter      *filesystemsFilterModel `tfsdk:"filter"`
	FileSystems []filesystemDataModel   `tfsdk:"filesystems"`
}

func NewFilesystemData() datasource.DataSource {
	return &filesystemData{}
}

func (d *filesystemData) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystems"
}

func (d *filesystemData) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Available File Systems Data",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filter the file systems",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Description: "Filter by region name",
						Optional:    true,
					},
				},
			},
			"filesystems": schema.ListNestedAttribute{
				Description: "List of available file systems",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "File System ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "File System name",
							Computed:    true,
						},
						"mount_point": schema.StringAttribute{
							Description: "File System mount point",
							Computed:    true,
						},
						"created": schema.StringAttribute{
							Description: "File System creation time",
							Computed:    true,
						},
						"is_in_use": schema.BoolAttribute{
							Description: "Whether the file system is in use",
							Computed:    true,
						},
						"bytes_used": schema.Int64Attribute{
							Description: "Bytes used in the file system",
							Computed:    true,
						},
						"created_by": schema.SingleNestedAttribute{
							Description: "User who created the file system",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "User ID",
									Computed:    true,
								},
								"email": schema.StringAttribute{
									Description: "User email",
									Computed:    true,
								},
								"status": schema.StringAttribute{
									Description: "User status",
									Computed:    true,
								},
							},
						},
						"region": schema.SingleNestedAttribute{
							Description: "Region information",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Region name",
									Computed:    true,
								},
								"description": schema.StringAttribute{
									Description: "Region description",
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

func (d *filesystemData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *filesystemData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model filesystemsDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListFileSystems(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list file systems", err.Error())
		return
	}

	// Apply filters if provided
	filteredFileSystems := res.Data
	if model.Filter != nil {
		filteredFileSystems = []api.FileSystem{}
		for _, fs := range res.Data {
			// Filter by region if specified
			if !model.Filter.Region.IsNull() && model.Filter.Region.ValueString() != "" {
				if fs.Region.Name != model.Filter.Region.ValueString() {
					continue
				}
			}

			filteredFileSystems = append(filteredFileSystems, fs)
		}
	}

	filesystems := make([]filesystemDataModel, 0, len(filteredFileSystems))
	for _, fs := range filteredFileSystems {
		filesystems = append(filesystems, filesystemDataModel{
			ID:         types.StringValue(fs.ID),
			Name:       types.StringValue(fs.Name),
			MountPoint: types.StringValue(fs.MountPoint),
			Created:    types.StringValue(fs.Created),
			IsInUse:    types.BoolValue(fs.IsInUse),
			BytesUsed:  types.Int64Value(fs.BytesUsed),
			CreatedBy: &filesystemUserModel{
				ID:     types.StringValue(fs.CreatedBy.ID),
				Email:  types.StringValue(fs.CreatedBy.Email),
				Status: types.StringValue(fs.CreatedBy.Status),
			},
			Region: &filesystemRegionModel{
				Name:        types.StringValue(fs.Region.Name),
				Description: types.StringValue(fs.Region.Description),
			},
		})
	}

	model.Id = types.StringValue("filesystems")
	model.FileSystems = filesystems

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
