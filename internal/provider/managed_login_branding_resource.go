package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/document"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Id         types.String `tfsdk:"id"`
	ClientId   types.String `tfsdk:"client_id"`
	UserPoolId types.String `tfsdk:"user_pool_id"`
	Settings   types.String `tfsdk:"settings"`
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
	out, err := cognito.CreateManagedLoginBranding(ctx, &input)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create resource, got error: %s", err))
		return
	}

	data.Id = types.StringValue(*out.ManagedLoginBranding.ManagedLoginBrandingId)

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

	data.Id = types.StringValue(*out.ManagedLoginBranding.ManagedLoginBrandingId)
	data.UserPoolId = types.StringValue(*out.ManagedLoginBranding.UserPoolId)

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
