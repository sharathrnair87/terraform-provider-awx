package awx

import (
	"context"
	"fmt"
	//"log"
	"strconv"

	//"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	//"github.com/hashicorp/terraform-plugin-framework/diag"
	//"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	awx "github.com/sharathrnair87/goawx/client"
)

var (
	_ datasource.DataSource              = &awxOrganizationDataSource{}
	_ datasource.DataSourceWithConfigure = &awxOrganizationDataSource{}
)

func init() {
	registerDataSource(NewAWXOrganizationDataSource)
}

func NewAWXOrganizationDataSource() datasource.DataSource {
	return &awxOrganizationDataSource{}
}

type awxOrganizationDataSource struct {
	client *awx.AWX
}

func (d *awxOrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_organization"
}

type awxOrganizationDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	MaxHosts         types.Int64  `tfsdk:"max_hosts"`
	CustomVirtualenv types.String `tfsdk:"custom_virtualenv"`
	//DefaultEnvironment types.String `tfsdk:"default_environment"`
	DefaultEnvironment types.Int64 `tfsdk:"default_environment"`
}

func (d *awxOrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
				// ConflictsWith Name
			},
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// Conflicts with ID
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"max_hosts": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"custom_virtualenv": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"default_environment": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (d *awxOrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*awx.AWX)
}

func (d *awxOrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var state awxOrganizationDataSourceModel

	// Read TF plan into DataSourceModel
	res.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if res.Diagnostics.HasError() {
		return
	}

	params := make(map[string]string)

	if !state.ID.IsNull() {
		params["id"] = state.ID.ValueString()
	}

	if !state.Name.IsNull() {
		params["name"] = state.Name.ValueString()
	}

	organizations, err := d.client.OrganizationsService.ListOrganizations(params)

	tflog.Debug(ctx, "Organizations", map[string]interface{}{
		"orgs": organizations,
	})

	if err != nil {
		res.Diagnostics.AddError("Unable to find Organization", fmt.Sprintf("Unable to find Organization with parameters %+v, got %s", params, err.Error()))
		return
	}

	if len(organizations) == 0 {
		res.Diagnostics.AddError("No Organizations found", fmt.Sprintf("No Organization found with parameters %+v", params))
		return
	}

	if len(organizations) > 1 {
		res.Diagnostics.AddError("Multiple entries!", fmt.Sprintf("Search with parameters %+v, returns multiple Organizations", params))
		return
	}

	organization := organizations[0]

	state.MaxHosts = types.Int64Value(int64(organization.MaxHosts))
	//state.DefaultEnvironment = types.StringValue(organization.DefaultEnvironment)
	state.DefaultEnvironment = types.Int64Value(int64(organization.DefaultEnvironment))
	state.CustomVirtualenv = types.StringValue(organization.CustomVirtualenv)
	state.Name = types.StringValue(organization.Name)
	state.Description = types.StringValue(organization.Description)
	state.ID = types.StringValue(strconv.Itoa(organization.ID))

	// Set State
	diags := res.State.Set(ctx, &state)
	res.Diagnostics.Append(diags...)
}
