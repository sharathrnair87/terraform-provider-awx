package awx

import (
	"context"
	"fmt"
	//"log"
	"strconv"

	//"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	//"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	//"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	awx "github.com/sharathrnair87/goawx/client"
)

var (
	_ resource.Resource                = &awxOrganizationResource{}
	_ resource.ResourceWithConfigure   = &awxOrganizationResource{}
	_ resource.ResourceWithImportState = &awxOrganizationResource{}
)

func init() {
	registerResource(NewAWXOrganizationResource)
}

func NewAWXOrganizationResource() resource.Resource {
	return &awxOrganizationResource{}
}

func (r *awxOrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_organization"
}

type awxOrganizationResource struct {
	client *awx.AWX
}

type awxOrganizationResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	MaxHosts           types.Int64  `tfsdk:"max_hosts"`
	CustomVirtualenv   types.String `tfsdk:"custom_virtualenv"`
	DefaultEnvironment types.Int64  `tfsdk:"default_environment"`
}

func (r *awxOrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"max_hosts": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"custom_virtualenv": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"default_environment": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				//PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			//"default_environment": schema.StringAttribute{
			//	Optional:      true,
			//	Computed:      true,
			//	PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			//},
		},
	}
}

func (r *awxOrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*awx.AWX)
}

func (r *awxOrganizationResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {

	var data *awxOrganizationResourceModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating AWX Organization", map[string]interface{}{
		"name": data.Name.String(),
	})

	org := make(map[string]interface{})

	org["name"] = data.Name.ValueString()
	org["description"] = data.Description.ValueString()
	org["max_hosts"] = data.MaxHosts.ValueInt64()

	if !data.CustomVirtualenv.IsNull() && !data.CustomVirtualenv.IsUnknown() {
		org["custom_virtualenv"] = data.CustomVirtualenv.ValueString()
	}

	if !data.DefaultEnvironment.IsNull() && !data.DefaultEnvironment.IsUnknown() {
		org["default_environment"] = data.DefaultEnvironment.ValueInt64()
		//org["default_environment"] = data.DefaultEnvironment.ValueString()
	}

	organization, err := r.client.OrganizationsService.CreateOrganization(org, map[string]string{})

	if err != nil {
		res.Diagnostics.AddError("AWX API Error!", fmt.Sprintf("Unable to create AWX Organization: %s", err.Error()))
		return
	}

	r.organizationModelToState(organization, data)

	tflog.Debug(ctx, "created AWX Organization", map[string]interface{}{
		"name": data.Name.String(),
	})

	res.Diagnostics.Append(res.State.Set(ctx, &data)...)
}

func (r *awxOrganizationResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {

	var data *awxOrganizationResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if res.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError("TF State Error!", fmt.Sprintf("ID value in state is not numeric: %s", err.Error()))
		return
	}

	/*
	  result, err := r.client.OrganizationsService.ListOrganizations(map[string]string{})

	  if err != nil {
	    res.Diagnostics.AddError("Unable to Auth!", fmt.Sprintf("Unable to Auth: %+v", err))
	    return
	  }

	  tflog.Trace(ctx, "ping succeeded", map[string]interface{}{
	    "result": result,
	  })
	*/

	organization, err := r.client.OrganizationsService.GetOrganizationsByID(id, make(map[string]string))
	if err != nil {
		res.Diagnostics.AddError("AWX API Error!", fmt.Sprintf("Unable to find AWX Organization: %s", err.Error()))
		return
	}

	tflog.Trace(ctx, "AWX Organization found", map[string]interface{}{
		"name":                organization.Name,
		"default_environment": organization.DefaultEnvironment,
	})

	r.organizationModelToState(organization, data)
	// Save updated data into Terraform state
	res.Diagnostics.Append(res.State.Set(ctx, &data)...)
}

func (r *awxOrganizationResource) organizationModelToState(organization *awx.Organization, data *awxOrganizationResourceModel) {
	data.ID = types.StringValue(strconv.Itoa(organization.ID))
	data.Name = types.StringValue(organization.Name)
	data.Description = types.StringValue(organization.Description)
	data.MaxHosts = types.Int64Value(int64(organization.MaxHosts))
	data.CustomVirtualenv = types.StringValue(organization.CustomVirtualenv)
	//data.DefaultEnvironment = types.StringValue(organization.DefaultEnvironment)
	data.DefaultEnvironment = types.Int64Value(int64(organization.DefaultEnvironment))
}

func (r *awxOrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {

	var data *awxOrganizationResourceModel
	res.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if res.Diagnostics.HasError() {
		return
	}

	err := r.update(ctx, data, &res.Diagnostics)
	if err != nil {
		res.Diagnostics.AddError("Unable to update AWX Organization", err.Error())
	}

	res.Diagnostics.Append(res.State.Set(ctx, &data)...)
}

func (r *awxOrganizationResource) update(ctx context.Context, data *awxOrganizationResourceModel, diags *diag.Diagnostics) error {

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		diags.AddError("TF State Error!", fmt.Sprintf("ID value in state is not numeric: %s", err.Error()))
		return err
	}

	_, err = r.client.OrganizationsService.GetOrganizationsByID(id, map[string]string{})
	if err != nil {
		diags.AddError("AWX API Error!", fmt.Sprintf("Unable to find AWX Organization: %s", err.Error()))
		return err
	}

	org := make(map[string]interface{})

	org["name"] = data.Name.ValueString()
	org["description"] = data.Description.ValueString()
	org["max_hosts"] = data.MaxHosts.ValueInt64()

	if !data.CustomVirtualenv.IsNull() && !data.CustomVirtualenv.IsUnknown() {
		org["custom_virtualenv"] = data.CustomVirtualenv.ValueString()
	}

	if !data.DefaultEnvironment.IsNull() && !data.DefaultEnvironment.IsUnknown() {
		org["default_environment"] = data.DefaultEnvironment.ValueInt64()
		//org["default_environment"] = data.DefaultEnvironment.ValueString()
	}

	organization, err := r.client.OrganizationsService.UpdateOrganization(id, org, map[string]string{})

	if err != nil {
		diags.AddError("AWX API Error!", fmt.Sprintf("Unable to Update Organization: %s", err.Error()))
		return err
	}

	r.organizationModelToState(organization, data)

	return nil
}

func (r *awxOrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (r *awxOrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {

	var data *awxOrganizationResourceModel
	res.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if res.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		res.Diagnostics.AddError("TF State Error!", fmt.Sprintf("ID value in state is not numeric: %s", err.Error()))
		return
	}

	_, err = r.client.OrganizationsService.DeleteOrganization(id)
	if err != nil {
		res.Diagnostics.AddError("AWX API Error!", fmt.Sprintf("Unable to delete Organization: %s", err.Error()))
		return
	}

}
