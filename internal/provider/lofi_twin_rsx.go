// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &T9LoFiTwinRsx{}
var _ resource.ResourceWithImportState = &T9LoFiTwinRsx{}

func NewT9LoFiTwinRsx() resource.Resource {
	return &T9LoFiTwinRsx{}
}

// T9LoFiTwinRsx defines the resource implementation.
type T9LoFiTwinRsx struct {
	client   *http.Client
	provider *Tensor9ProviderModel
}

// T9LoFiTwinRsxModel describes the resource data model.
type T9LoFiTwinRsxModel struct {
	Template     types.String `tfsdk:"template"`
	TemplateFmt  types.String `tfsdk:"template_fmt"`
	ProjectionId types.String `tfsdk:"projection_id"`
	InfraId      types.String `tfsdk:"infra_id"`
	Properties   types.Map    `tfsdk:"properties"`
	Id           types.String `tfsdk:"id"`
}

func (r *T9LoFiTwinRsx) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lofi_twin"
}

func (r *T9LoFiTwinRsx) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tensor9 Lo-Fidelity Digital Twin",

		Attributes: map[string]schema.Attribute{
			"template": schema.StringAttribute{
				MarkdownDescription: "The infra template that specifies the resource to create inside the appliance",
				Optional:            false,
				Required:            true,
			},
			"template_fmt": schema.StringAttribute{
				MarkdownDescription: "The format of the template",
				Optional:            false,
				Required:            true,
			},
			"projection_id": schema.StringAttribute{
				MarkdownDescription: "The id of the projection (and associated appliance) to create the resource in",
				Optional:            false,
				Required:            true,
			},
			"infra_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the infra that tracks the lifecycle and configuration of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"properties": schema.MapAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "A map of properties to configure the resource",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource - always set to the infra_id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *T9LoFiTwinRsx) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*Tensor9ProviderData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Tensor9ProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = providerData.Client
	r.provider = providerData.Model
}

func (r *T9LoFiTwinRsx) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data T9LoFiTwinRsxModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	println(fmt.Sprintf("Found provider endpoint: %s", r.provider.Endpoint))
	println(fmt.Sprintf("Found provider api_key: %s", r.provider.ApiKey))

	tflog.Debug(ctx, fmt.Sprintf("Found provider endpoint: %s", r.provider.Endpoint))
	tflog.Debug(ctx, fmt.Sprintf("Found provider api_key: %s", r.provider.ApiKey))

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.InfraId = types.StringValue("00000000000000000000000000000001") // TODO: actual infra id
	data.Id = data.InfraId

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *T9LoFiTwinRsx) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data T9LoFiTwinRsxModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *T9LoFiTwinRsx) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data T9LoFiTwinRsxModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *T9LoFiTwinRsx) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data T9LoFiTwinRsxModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *T9LoFiTwinRsx) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
