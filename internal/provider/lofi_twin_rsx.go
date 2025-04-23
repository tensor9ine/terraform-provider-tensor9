// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
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

type T9LoFiTwinRsxModel struct {
	Template     types.String `tfsdk:"template"`
	TemplateFmt  types.String `tfsdk:"template_fmt"`
	ProjectionId types.String `tfsdk:"projection_id"`
	Params       types.Map    `tfsdk:"params"`
	Schema       types.Map    `tfsdk:"schema"`
	Outputs      types.Map    `tfsdk:"outputs"`
	RsxId        types.String `tfsdk:"rsx_id"`
	InfraId      types.String `tfsdk:"infra_id"`
	Id           types.String `tfsdk:"id"`
}

type TfLoFiTemplate struct {
	Raw string `json:"raw"`
	Fmt string `json:"fmt"`
}

type TfLoFiTwinRsx struct {
	RsxId        *string                   `json:"rsxId"`
	Template     *TfLoFiTemplate           `json:"template"`
	ProjectionId *string                   `json:"projectionId"`
	Params       *map[string]string        `json:"params"`
	Schema       *map[string]TfRsxPropType `json:"schema"`
	Outputs      *map[string]string        `json:"outputs"`
	InfraId      *string                   `json:"infraId"`
}

type TfRsxEvt struct {
	ApiKey      string         `json:"apiKey"`
	RsxType     string         `json:"rsxType"`
	EvtType     string         `json:"evtType"`
	LoFiTwinRsx *TfLoFiTwinRsx `json:"loFiTwinRsx"`
}

type TfRsxEvtResult struct {
	EvtType     string                `json:"evtType"`
	RsxType     string                `json:"rsxType"`
	ResultType  string                `json:"resultType"`
	LoFiTwinRsx *Delta[TfLoFiTwinRsx] `json:"loFiTwinRsx"`
	Reason      *string               `json:"reason"`
}

type Delta[T any] struct {
	Before *T `json:"before"`
	After  *T `json:"after"`
}

type TfRsxPropType string

const (
	Bool   TfRsxPropType = "Bool"
	I32    TfRsxPropType = "I32"
	I64    TfRsxPropType = "I64"
	F32    TfRsxPropType = "F32"
	F64    TfRsxPropType = "F64"
	Str    TfRsxPropType = "Str"
	Secret TfRsxPropType = "Secret"
)

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
			"params": schema.MapAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "A map of parameters to pass to the template that specified the resource",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"schema": schema.MapAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The schema describing the properties of the resource",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"outputs": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "A map of outputs published by the resource upon create/update",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"rsx_id": schema.StringAttribute{
				Optional:            false,
				Required:            true,
				MarkdownDescription: "The rsx id of the twin rsx in the compiled twin stack",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"infra_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The id of the infra that tracks the lifecycle and configuration of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
	var rsxModel T9LoFiTwinRsxModel

	// Read Terraform plan rsxModel into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &rsxModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	println(fmt.Sprintf("Found provider endpoint: %s", r.provider.Endpoint))
	println(fmt.Sprintf("Found provider api_key: %s", r.provider.ApiKey))

	tflog.Debug(ctx, fmt.Sprintf("Found provider endpoint: %s", r.provider.Endpoint))
	tflog.Debug(ctx, fmt.Sprintf("Found provider api_key: %s", r.provider.ApiKey))

	params := mapToStringMap(rsxModel.Params)
	schema := schematize(rsxModel.Schema)
	var evt = TfRsxEvt{
		ApiKey:  r.provider.ApiKey.ValueString(),
		RsxType: "LoFiTwin",
		EvtType: "Create",
		LoFiTwinRsx: &TfLoFiTwinRsx{
			RsxId: rsxModel.RsxId.ValueStringPointer(),
			Template: &TfLoFiTemplate{
				Raw: rsxModel.Template.ValueString(),
				Fmt: rsxModel.TemplateFmt.ValueString(),
			},
			ProjectionId: rsxModel.ProjectionId.ValueStringPointer(),
			Params:       &params,
			Schema:       &schema,
		},
	}
	evtJson, err := json.Marshal(evt)
	if err != nil {
		resp.Diagnostics.AddError("JSON Encoding Error", fmt.Sprintf("Failed to encode request body: %s", err))
		return
	}

	evtReq, err := http.NewRequest("POST", r.provider.Endpoint.ValueString()+"/stack/tf/react", bytes.NewReader(evtJson))
	evtResultResp, err := r.client.Do(evtReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create rsx, got error: %s", err))
		return
	}

	evtResultBytes, err := io.ReadAll(evtResultResp.Body)

	if err != nil {
		resp.Diagnostics.AddError("Read Response Error", err.Error())
		return
	}

	evtResultStr := string(evtResultBytes)
	//println("Received evt result", evtResultStr)
	tflog.Debug(ctx, fmt.Sprintf("Create response body: %s", evtResultStr))

	err = evtResultResp.Body.Close()
	if err != nil {
		resp.Diagnostics.AddError("Close Response Body Error", err.Error())
		return
	}

	var evtResult TfRsxEvtResult

	err = json.Unmarshal(evtResultBytes, &evtResult)
	if err != nil {
		resp.Diagnostics.AddError("JSON decode error", fmt.Sprintf("Failed to decode response JSON: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Got tf evt result: %s", evtResultStr))
	//println(fmt.Sprintf("Got tf evt result for tf stack reactor: %s", evtResultStr))

	rsxModel.InfraId = types.StringValue(*evtResult.LoFiTwinRsx.After.InfraId)
	rsxModel.Id = rsxModel.InfraId

	outputs, diag := types.MapValueFrom(ctx, types.StringType, evtResult.LoFiTwinRsx.After.Outputs)
	resp.Diagnostics.Append(diag...)
	rsxModel.Outputs = outputs

	tflog.Debug(ctx, fmt.Sprintf("created an lo fi twin resource; infra_id=%s", rsxModel.InfraId.ValueString()))
	//println(fmt.Sprintf("created lo fi twin resource; infra_id=%s; rsx_id=%s; outputs=%s", rsxModel.InfraId.ValueString(), rsxModel.RsxId.ValueString(), rsxModel.Outputs.String()))

	// Save rsxModel into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &rsxModel)...)
}

func (r *T9LoFiTwinRsx) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var rsxModel T9LoFiTwinRsxModel

	// Read Terraform prior state rsxModel into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &rsxModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: handle read

	// Save updated rsxModel into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &rsxModel)...)
}

func (r *T9LoFiTwinRsx) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var rsxModel T9LoFiTwinRsxModel

	// Read Terraform plan rsxModel into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &rsxModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: handle update

	// Save updated rsxModel into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &rsxModel)...)
}

func (r *T9LoFiTwinRsx) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var rsxModel T9LoFiTwinRsxModel

	// Read Terraform prior state rsxModel into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &rsxModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: handle delete
}

func (r *T9LoFiTwinRsx) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapToStringMap(attrMap types.Map) map[string]string {
	result := make(map[string]string)
	for k, v := range attrMap.Elements() {
		if strVal, ok := v.(types.String); ok {
			result[k] = strVal.ValueString()
		}
	}
	return result
}

func schematize(attrMap types.Map) map[string]TfRsxPropType {
	result := make(map[string]TfRsxPropType)
	for k, v := range attrMap.Elements() {
		if strVal, ok := v.(types.String); ok {
			result[k] = TfRsxPropType(strVal.ValueString())
		}
	}
	return result
}
