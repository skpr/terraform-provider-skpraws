package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/document"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/skpr/terraform-provider-skpraws/internal/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ManagedLoginBrandingResource{}

func NewManagedLoginBrandingResource() resource.Resource {
	return &ManagedLoginBrandingResource{}
}

type ManagedLoginBrandingResource struct {
	config aws.Config
}

type ManagedLoginBrandingModel struct {
	Id         tftypes.String              `tfsdk:"id"`
	ClientId   tftypes.String              `tfsdk:"client_id"`
	UserPoolId tftypes.String              `tfsdk:"user_pool_id"`
	Settings   tftypes.String              `tfsdk:"settings"`
	Assets     []ManagedLoginBrandingAsset `tfsdk:"assets"`
}

type ManagedLoginBrandingAsset struct {
	Category  tftypes.String `tfsdk:"category"`
	ColorMode tftypes.String `tfsdk:"color_mode"`
	Bytes     tftypes.String `tfsdk:"bytes"`
}

func (r *ManagedLoginBrandingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_login_branding"
}

func (r *ManagedLoginBrandingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Managed Login Branding resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID",
				Required:            true,
			},
			"user_pool_id": schema.StringAttribute{
				MarkdownDescription: "Cognito user pool for resource",
				Required:            true,
			},
			"settings": schema.StringAttribute{
				MarkdownDescription: "Settings for branding",
				Optional:            true,
			},
			"assets": schema.ListNestedAttribute{
				MarkdownDescription: "Assets for branding.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{
							Description: "Category to fill for the asset",
							Required:    true,
						},
						"color_mode": schema.StringAttribute{
							Description: "light or dark color mode the asset will be used for",
							Required:    true,
						},
						"bytes": schema.StringAttribute{
							Description: "actual file data",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r *ManagedLoginBrandingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	awsConfig, ok := req.ProviderData.(aws.Config)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected aws.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.config = awsConfig
}

func (r *ManagedLoginBrandingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ManagedLoginBrandingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cognito := cognitoidentityprovider.NewFromConfig(r.config)
	input := cognitoidentityprovider.CreateManagedLoginBrandingInput{
		ClientId:                 data.ClientId.ValueStringPointer(),
		UserPoolId:               data.UserPoolId.ValueStringPointer(),
		UseCognitoProvidedValues: data.Settings.IsNull(),
	}
	if !data.Settings.IsNull() {
		var jsonData interface{}
		err := json.Unmarshal([]byte(data.Settings.ValueString()), &jsonData)
		if err != nil {
			resp.Diagnostics.AddError("Data Error", fmt.Sprintf("Invalid JSON for settings, got error: %s", err))
			return
		}
		input.Settings = document.NewLazyDocument(jsonData)
	}
	for _, asset := range data.Assets {
		cognitoAsset := awstypes.AssetType{
			Category:  types.AssetCategoryTypeFromString(asset.Category.ValueString()),
			ColorMode: types.ColorSchemeModeTypeFromString(asset.ColorMode.ValueString()),
			Bytes:     []byte(asset.Bytes.ValueString()),
			Extension: types.AssetExtensionTypeFromString("SVG"),
		}
		input.Assets = append(input.Assets, cognitoAsset)
	}
	out, err := cognito.CreateManagedLoginBranding(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource, got error: %s", err))
		return
	}

	data.Id = tftypes.StringValue(*out.ManagedLoginBranding.ManagedLoginBrandingId)
	for _, asset := range out.ManagedLoginBranding.Assets {
		modelAsset := ManagedLoginBrandingAsset{
			Bytes: tftypes.StringValue(string(asset.Bytes)),
		}
		data.Assets = append(data.Assets, modelAsset)
	}

	tflog.Trace(ctx, "created Managed Login Branding resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManagedLoginBrandingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ManagedLoginBrandingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cognito := cognitoidentityprovider.NewFromConfig(r.config)
	out, err := cognito.DescribeManagedLoginBranding(context.TODO(), &cognitoidentityprovider.DescribeManagedLoginBrandingInput{
		ManagedLoginBrandingId: data.Id.ValueStringPointer(),
		UserPoolId:             data.UserPoolId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read resource, got error: %s", err))
		return
	}

	data.Id = tftypes.StringValue(*out.ManagedLoginBranding.ManagedLoginBrandingId)
	data.UserPoolId = tftypes.StringValue(*out.ManagedLoginBranding.UserPoolId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManagedLoginBrandingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ManagedLoginBrandingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cognito := cognitoidentityprovider.NewFromConfig(r.config)
	input := cognitoidentityprovider.UpdateManagedLoginBrandingInput{
		ManagedLoginBrandingId:   data.Id.ValueStringPointer(),
		UserPoolId:               data.UserPoolId.ValueStringPointer(),
		UseCognitoProvidedValues: data.Settings.IsNull(),
	}
	if !data.Settings.IsNull() {
		var jsonData interface{}
		err := json.Unmarshal([]byte(data.Settings.ValueString()), &jsonData)
		if err != nil {
			resp.Diagnostics.AddError("Data Error", fmt.Sprintf("Invalid JSON for settings, got error: %s", err))
			return
		}
		input.Settings = document.NewLazyDocument(jsonData)
	}
	for _, asset := range data.Assets {
		cognitoAsset := awstypes.AssetType{
			Category:  types.AssetCategoryTypeFromString(asset.Category.ValueString()),
			ColorMode: types.ColorSchemeModeTypeFromString(asset.ColorMode.ValueString()),
			Bytes:     []byte(asset.Bytes.ValueString()),
			Extension: types.AssetExtensionTypeFromString("SVG"),
		}
		input.Assets = append(input.Assets, cognitoAsset)
	}
	_, err := cognito.UpdateManagedLoginBranding(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "updated Managed Login Branding resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManagedLoginBrandingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ManagedLoginBrandingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cognito := cognitoidentityprovider.NewFromConfig(r.config)
	_, err := cognito.DeleteManagedLoginBranding(ctx, &cognitoidentityprovider.DeleteManagedLoginBrandingInput{
		ManagedLoginBrandingId: data.Id.ValueStringPointer(),
		UserPoolId:             data.UserPoolId.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete resource, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted Managed Login Branding resource")
}
