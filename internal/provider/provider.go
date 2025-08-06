package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	skprresource "github.com/skpr/terraform-provider-skpraws/internal/resource"
)

var _ provider.Provider = &SkprAwsProvider{}

type SkprAwsProvider struct {
	version string
}

// SkprAwsProviderModel describes the provider data model.
type SkprAwsProviderModel struct {
	Profile types.String `tfsdk:"profile"`
	Region  types.String `tfsdk:"region"`
}

func (p *SkprAwsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "skpraws"
	resp.Version = p.version
}

func (p *SkprAwsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "The AWS region",
				Required:            true,
			},
			"profile": schema.StringAttribute{
				MarkdownDescription: "The AWS profile",
				Optional:            true,
			},
		},
	}
}

func (p *SkprAwsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data SkprAwsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		awsConfig aws.Config
		err       error
	)
	if data.Profile.IsNull() {
		awsConfig, err = config.LoadDefaultConfig(context.TODO())
	} else {
		awsConfig, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(data.Profile.ValueString()))
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	awsConfig.Region = data.Region.ValueString()

	resp.DataSourceData = awsConfig
	resp.ResourceData = awsConfig
}

func (p *SkprAwsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		skprresource.NewManagedLoginBrandingResource,
	}
}

func (p *SkprAwsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *SkprAwsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *SkprAwsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SkprAwsProvider{
			version: version,
		}
	}
}
