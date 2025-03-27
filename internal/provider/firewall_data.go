package provider

import (
	"context"

	api "github.com/elct9620/terraform-provider-lambdalabs/pkg/lambdalabs"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &firewallData{}
	_ datasource.DataSourceWithConfigure = &firewallData{}
)

type firewallData struct {
	client *api.Client
}

type firewallRuleModel struct {
	Protocol      types.String `tfsdk:"protocol"`
	PortRange     types.List   `tfsdk:"port_range"`
	SourceNetwork types.String `tfsdk:"source_network"`
	Description   types.String `tfsdk:"description"`
}

type firewallDataModel struct {
	Id    types.String       `tfsdk:"id"`
	Rules []firewallRuleModel `tfsdk:"rules"`
}

func NewFirewallData() datasource.DataSource {
	return &firewallData{}
}

func (d *firewallData) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall"
}

func (d *firewallData) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Firewall Rules Data",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"rules": schema.ListNestedAttribute{
				Description: "List of firewall rules",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							Description: "The protocol (tcp, udp)",
							Computed:    true,
						},
						"port_range": schema.ListAttribute{
							Description: "The port range [start, end]",
							Computed:    true,
							ElementType: types.Int64Type,
						},
						"source_network": schema.StringAttribute{
							Description: "The source network in CIDR notation",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the rule",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *firewallData) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*api.Client)
}

func (d *firewallData) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model firewallDataModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := d.client.ListFirewallRules(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to list firewall rules", err.Error())
		return
	}

	model.Id = types.StringValue("firewall")
	model.Rules = make([]firewallRuleModel, 0, len(res.Data))

	for _, rule := range res.Data {
		var portRange []int64
		for _, p := range rule.PortRange {
			portRange = append(portRange, int64(p))
		}

		portRangeList, diags := types.ListValueFrom(ctx, types.Int64Type, portRange)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		model.Rules = append(model.Rules, firewallRuleModel{
			Protocol:      types.StringValue(rule.Protocol),
			PortRange:     portRangeList,
			SourceNetwork: types.StringValue(rule.SourceNetwork),
			Description:   types.StringValue(rule.Description),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
