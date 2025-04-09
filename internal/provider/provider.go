// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure Tensor9Provider satisfies various provider interfaces.
var _ provider.Provider = &Tensor9Provider{}
var _ provider.ProviderWithFunctions = &Tensor9Provider{}
var _ provider.ProviderWithEphemeralResources = &Tensor9Provider{}

// Tensor9Provider defines the provider implementation.
type Tensor9Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Tensor9ProviderModel describes the provider data model.
type Tensor9ProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

type Tensor9ProviderData struct {
	Client *http.Client
	Model  *Tensor9ProviderModel
}

func (p *Tensor9Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tensor9"
	resp.Version = p.version
}

func (p *Tensor9Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The endpoint of the vctrl's terraform stack reactor to send CRUD requests to",
				Optional:            false,
				Required:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The api key used to authenticate with the vctrl's terraform stack reactor",
				Optional:            false,
				Required:            true,
			},
		},
	}
}

func (p *Tensor9Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Tensor9ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = &Tensor9ProviderData{
		Client: client,
		Model:  &data,
	}
}

func (p *Tensor9Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewT9LoFiTwinRsx,
	}
}

func (p *Tensor9Provider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *Tensor9Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *Tensor9Provider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Tensor9Provider{
			version: version,
		}
	}
}
