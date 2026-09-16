// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var pathQueryFilters = path.Root("query_filters")

type listResource struct {
	crudResource
}

type listConfig struct {
	QueryFilters jsontypes.Normalized `tfsdk:"query_filters"`
}

func (r *listResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.crudResource.Metadata(ctx, req, resp)
}

func (r *listResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.crudResource.Configure(ctx, req, resp)
}

func (r *listResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var diags diag.Diagnostics
	if !r.ready(&diags) {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	var cfg listConfig
	if diags = req.Config.Get(ctx, &cfg); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	m := r.model
	filters, err := m.queryFilters(cfg.QueryFilters)
	if err != nil {
		diags.AddAttributeError(pathQueryFilters, "Invalid query_filters", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	options := maps.Clone(m.listOptions)
	if options == nil {
		options = map[string]any{}
	}
	if req.Limit > 0 {
		options["limit"] = req.Limit
	}

	raw, err := callIdempotent(ctx, r.data.Client, m.namespace+".query", filters, options)
	if err != nil {
		m.addCallError(r.data.Client, &diags, "list", err)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	decoded, err := decodeResult(raw)
	rows, ok := decoded.([]any)
	if err != nil || !ok {
		diags.AddError("Unexpected list response", fmt.Sprintf("%s.query did not return a list", m.namespace))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, item := range rows {
			row, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if !push(m.listResult(ctx, req, row)) {
				return
			}
		}
	}
}

func (m *model) listResult(ctx context.Context, req list.ListRequest, row map[string]any) list.ListResult {
	result := req.NewListResult(ctx)
	id, err := m.normalizeID(row[m.primaryKey])
	if err != nil {
		result.Diagnostics.AddError("Unexpected list response", err.Error())
		return result
	}
	setIdentity(ctx, result.Identity, id, &result.Diagnostics)
	result.DisplayName = displayName(row, id)
	if req.IncludeResource {
		v, err := decodeObject(objectType(m.attrs), m.attrs, row, tftypes.NewValue(objectType(m.attrs), nil))
		if err != nil {
			result.Diagnostics.AddError("Unexpected list response", err.Error())
			return result
		}
		result.Resource.Raw = v
	}
	return result
}

// listFilters combines the spec's filters, the variant discriminator and the user's filters.
func (m *model) queryFilters(user jsontypes.Normalized) ([]any, error) {
	filters := []any{}
	for _, f := range m.listFilters {
		filters = append(filters, f)
	}
	if m.discriminator != "" {
		filters = append(filters, []any{m.discriminator, "=", m.variant})
	}
	if user.IsNull() || user.IsUnknown() {
		return filters, nil
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(user.ValueString())))
	dec.UseNumber()
	var extra []any
	if err := dec.Decode(&extra); err != nil {
		return nil, fmt.Errorf("must be a JSON list of filters: %w", err)
	}
	return append(filters, extra...), nil
}

func displayName(row map[string]any, id any) string {
	for _, key := range []string{"name", "username", "path", "dataset", "comment", "description"} {
		if s, ok := row[key].(string); ok && s != "" {
			return s
		}
	}
	return fmt.Sprint(id)
}

func (r *listResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: fmt.Sprintf("Lists existing `%s` objects for `terraform query`.", r.model.namespace),
		Attributes: map[string]listschema.Attribute{
			"query_filters": listschema.StringAttribute{
				Optional:   true,
				CustomType: jsontypes.NormalizedType{},
				MarkdownDescription: "Middleware query filters as JSON, e.g. " +
					"`jsonencode([[\"pool\", \"=\", \"tank\"]])`. Combined with any filters the resource always applies.",
			},
		},
	}
}
