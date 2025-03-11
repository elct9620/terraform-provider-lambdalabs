package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &imageData{}
	_ datasource.DataSourceWithConfigure = &imageData{}
)

type imageData struct {
	client *api.Client
}

type imageRegionModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type imageDataModel struct {
	Id           types.String     `tfsdk:"id"`
	Name         types.String     `tfsdk:"name"`
	Description  types.String     `tfsdk:"description"`
	Family       types.String     `tfsdk:"family"`
	Version      types.String     `tfsdk:"version"`
	Architecture types.String     `tfsdk:"architecture"`
	CreatedTime  types.String     `tfsdk:"created_time"`
	UpdatedTime  types.String     `tfsdk:"updated_time"`
	Region       *imageRegionModel `tfsdk:"region"`
}

type imagesFilterModel struct {
	Region       types.String `tfsdk:"region"`
	Family       types.String `tfsdk:"family"`
	Architecture types.String `tfsdk:"architecture"`
}

type imagesDataModel struct {
	Id     types.String        `tfsdk:"id"`
	Filter *imagesFilterModel  `tfsdk:"filter"`
	Images []imageDataModel    `tfsdk:"images"`
}

func NewImageData() datasource.DataSource {
	return &imageData{}
}

func (d *imageData) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (d *imageData) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Available Images Data",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"filter": schema.SingleNestedAttribute{
				Description: "Filter the images",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						Description: "Filter by region name",
						Optional:    true,
					},
					"family": schema.StringAttribute{
						Description: "Filter by image family",
						Optional:    true,
					},
					"architecture": schema.StringAttribute{
						Description: "Filter by architecture",
						Optional:    true,
					},
				},
			},
			"images": schema.ListNestedAttribute{
				Description: "List of available images",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Image ID",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Image name",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Image description",
							Computed:    true,
						},
						"family": schema.StringAttribute{
							Description: "Image family",
							Computed:    true,
						},
						"version": schema.StringAttribute{
							Description: "Image version",
							Computed:    true,
						},
						"architecture": schema.StringAttribute{
							Description: "Image architecture",
							Computed:    true,
						},
						"created_time": schema.StringAttribute{
							Description: "Image creation time",
							Computed:    true,
						},
						"updated_time": schema.StringAttribute{
							Description: "Image last update time",
							Computed:    true,
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

func (d *imageData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *imageData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model imagesDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListImages(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list images", err.Error())
		return
	}

	// Apply filters if provided
	filteredImages := res.Data
	if model.Filter != nil {
		filteredImages = []api.Image{}
		for _, image := range res.Data {
			// Filter by region if specified
			if !model.Filter.Region.IsNull() && model.Filter.Region.ValueString() != "" {
				if image.Region.Name != model.Filter.Region.ValueString() {
					continue
				}
			}
			
			// Filter by family if specified
			if !model.Filter.Family.IsNull() && model.Filter.Family.ValueString() != "" {
				if image.Family != model.Filter.Family.ValueString() {
					continue
				}
			}
			
			// Filter by architecture if specified
			if !model.Filter.Architecture.IsNull() && model.Filter.Architecture.ValueString() != "" {
				if image.Architecture != model.Filter.Architecture.ValueString() {
					continue
				}
			}
			
			filteredImages = append(filteredImages, image)
		}
	}

	images := make([]imageDataModel, 0, len(filteredImages))
	for _, image := range filteredImages {
		images = append(images, imageDataModel{
			Id:           types.StringValue(image.ID),
			Name:         types.StringValue(image.Name),
			Description:  types.StringValue(image.Description),
			Family:       types.StringValue(image.Family),
			Version:      types.StringValue(image.Version),
			Architecture: types.StringValue(image.Architecture),
			CreatedTime:  types.StringValue(image.CreatedTime),
			UpdatedTime:  types.StringValue(image.UpdatedTime),
			Region: &imageRegionModel{
				Name:        types.StringValue(image.Region.Name),
				Description: types.StringValue(image.Region.Description),
			},
		})
	}

	model.Id = types.StringValue("images")
	model.Images = images

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
