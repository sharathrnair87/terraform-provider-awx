package awx

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	//"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	awx "github.com/sharathrnair87/goawx/client"
)

var _ provider.Provider = &AWXProvider{}

type AWXProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Token    types.String `tfsdk:"token"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

type Config struct {
	Hostname string
	Username string
	Password string
	Token    string
	Insecure bool
}

type AWXProvider struct {
	version string
}

func (p *AWXProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awx"
	resp.Version = p.version
}

func (p *AWXProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: "The AWX Host that we connect to. (defaults to AWX_HOSTNAME env variable if set)",
				Optional:    true,
			},
			"insecure": schema.BoolAttribute{
				Description: "If you are using a self signed certificate this should be set to true, default is false",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username to connect to the AWX host. (defaults to AWX_USERNAME env variable if set)",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to connect to the AWX host. (defaults to AWX_PASSWORD env variable if set)",
				Optional:    true,
				Sensitive:   true,
			},
			"token": schema.StringAttribute{
				Description: "The OAUTH2 connect to the AWX host. (defaults to AWX_TOKEN env variable if set)",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *AWXProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// parse the AWX Provider configuration
	var config AWXProviderModel

	// read values from provider awx {} block declaration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for misconfigured provider configuration values

	// Initialize Provider config with environment variables
	conf := Config{
		Hostname: os.Getenv("AWX_HOSTNAME"),
		Username: os.Getenv("AWX_USERNAME"),
		Password: os.Getenv("AWX_PASSWORD"),
		Token:    os.Getenv("AWX_TOKEN"),
		Insecure: true,
	}

	// Override config with TF Provider Block values if any
	if !config.Hostname.IsNull() {
		conf.Hostname = config.Hostname.ValueString()
	}

	if !config.Token.IsNull() {
		conf.Token = config.Token.ValueString()
	}

	if !config.Password.IsNull() {
		conf.Password = config.Password.ValueString()
	}

	if !config.Username.IsNull() {
		conf.Username = config.Username.ValueString()
	}

	if !config.Insecure.IsNull() {
		conf.Insecure = config.Insecure.ValueBool()
	}

	// Create TLS config to skip SSL verification
	c := http.DefaultClient
	if conf.Insecure {
		ct := http.DefaultTransport.(*http.Transport).Clone()
		ct.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		c.Transport = ct
	}

	var awxClient *awx.AWX
	var err error
	if conf.Token != "" {
		awxClient, err = awx.NewAWXToken(conf.Hostname, conf.Token, c)
	} else {
		awxClient, err = awx.NewAWX(conf.Hostname, conf.Username, conf.Password, c)
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to authenticate user against AWX API: check values in the provider initialization", fmt.Sprintf("Provider initialization failed: %+v", err))
	}

	if err == nil {
		//resp.Diagnostics.AddWarning("Auth Succeeded", "Auth Succeeded!")
		tflog.Trace(ctx, "Auth succeeded", map[string]interface{}{
			"auth": "success",
		})
	}

	resp.ResourceData = awxClient
	resp.DataSourceData = awxClient
}

func (p *AWXProvider) Resources(ctx context.Context) []func() resource.Resource {
	return allAWXResources
}

func (p *AWXProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AWXProvider{
			version: version,
		}
	}
}
