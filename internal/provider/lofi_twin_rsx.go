// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	client *http.Client
}

// T9LoFiTwinRsxModel describes the resource data model.
type T9LoFiTwinRsxModel struct {
	Template     types.String `tfsdk:"template"`
	TemplateFmt  types.String `tfsdk:"template_fmt"`
	ProjectionId types.String `tfsdk:"projection_id"`
	Defaulted    types.String `tfsdk:"defaulted"`
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
				MarkdownDescription: "The template that specifies the resource to create inside the appliance",
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
			"defaulted": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute with default value",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("example value when not configured"),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
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

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *T9LoFiTwinRsx) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

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
