// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	dsName    = "name"
	dsFilters = "query_filters"
)

type dataSource struct {
	crudResource
	attrs []*node // exposed attributes: readable, not sensitive, not write-only
}

// dataAttrs selects the attributes a data source exposes, all computed.
func dataAttrs(m *model) []*node {
	var out []*node
	for _, a := range m.attrs {
		if a.api == "" || !a.readable || a.writeOnly || a.sensitive || containsWriteOnly(a) || containsSensitive(a) {
			continue
		}
		// Deep copy: markComputed rewrites roles all the way down, and the descriptor these nodes
		// come from is shared with the resource. Copying only the top level would quietly turn every
		// nested attribute of the resource computed as well.
		c := cloneNode(a)
		markComputed(c)
		out = append(out, c)
	}
	return out
}

// cloneNode copies a node and everything under it.
func cloneNode(a *node) *node {
	if a == nil {
		return nil
	}
	c := *a
	c.children = nil
	for _, child := range a.children {
		c.children = append(c.children, cloneNode(child))
	}
	c.elem = cloneNode(a.elem)
	return &c
}

func containsSensitive(a *node) bool {
	for _, c := range a.children {
		if c.sensitive || containsSensitive(c) {
			return true
		}
	}
	return a.elem != nil && (a.elem.sensitive || containsSensitive(a.elem))
}

func (d *dataSource) hasName() bool {
	for _, a := range d.attrs {
		if a.name == dsName && a.kind == kindString {
			return true
		}
	}
	return false
}

func (d *dataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	var r resource.MetadataResponse
	d.crudResource.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: req.ProviderTypeName}, &r)
	resp.TypeName = r.TypeName
}

func (d *dataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	var r resource.ConfigureResponse
	d.crudResource.Configure(ctx, resource.ConfigureRequest{ProviderData: req.ProviderData}, &r)
	resp.Diagnostics.Append(r.Diagnostics...)
}

func (d *dataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if !d.ready(&resp.Diagnostics) {
		return
	}
	m := d.model
	typ := d.stateType()

	if m.singleton {
		raw, err := callIdempotent(ctx, d.data.Client, m.getMethod)
		obj, ok := d.object(raw, err, &resp.Diagnostics)
		if !ok {
			return
		}
		d.setState(resp, typ, obj, nil)
		return
	}

	config := children(req.Config.Raw)
	filters, err := m.queryFilters(jsonFromValue(config[dsFilters]))
	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root(dsFilters), "Invalid query_filters", err.Error())
		return
	}
	lookups := 0
	if v := config[idAttr]; v.IsKnown() && !v.IsNull() {
		id, err := m.idFromState(req.Config.Raw)
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root(idAttr), "Invalid id", err.Error())
			return
		}
		filters = append(filters, []any{m.primaryKey, "=", id})
		lookups++
	}
	if v, ok := config[dsName]; ok && v.IsKnown() && !v.IsNull() {
		var name string
		_ = v.As(&name)
		filters = append(filters, []any{dsName, "=", name})
		lookups++
	}
	if v := config[dsFilters]; v.IsKnown() && !v.IsNull() {
		lookups++
	}
	if lookups == 0 {
		resp.Diagnostics.AddError("Missing lookup", "Set id, name or query_filters to choose the object to read.")
		return
	}

	options := maps.Clone(m.listOptions)
	if options == nil {
		options = map[string]any{}
	}
	raw, err := callIdempotent(ctx, d.data.Client, m.namespace+".query", filters, options)
	if err != nil {
		m.addCallError(d.data.Client, &resp.Diagnostics, "read", err)
		return
	}
	decoded, err := decodeResult(raw)
	rows, ok := decoded.([]any)
	if err != nil || !ok {
		resp.Diagnostics.AddError("Unexpected read response", fmt.Sprintf("%s.query did not return a list", m.namespace))
		return
	}
	if len(rows) != 1 {
		resp.Diagnostics.AddError("Lookup did not match exactly one object",
			fmt.Sprintf("%s matched %d objects for filters %v.", m.namespace, len(rows), filters))
		return
	}
	obj, ok := d.object(nil, nil, &resp.Diagnostics, rows[0])
	if !ok {
		return
	}
	d.setState(resp, typ, obj, config)
}

// object decodes a single middleware object from raw, or takes an already decoded one.
func (d *dataSource) object(raw []byte, err error, diags *diag.Diagnostics, decoded ...any) (map[string]any, bool) {
	if err != nil {
		d.model.addCallError(d.data.Client, diags, "read", err)
		return nil, false
	}
	var v any
	if len(decoded) > 0 {
		v = decoded[0]
	} else if v, err = decodeResult(raw); err != nil {
		diags.AddError("Unexpected read response", err.Error())
		return nil, false
	}
	obj, ok := v.(map[string]any)
	if !ok {
		diags.AddError("Unexpected read response", fmt.Sprintf("expected an object, got %T", v))
		return nil, false
	}
	if m := d.model; m.discriminator != "" && obj[m.discriminator] != m.variant {
		diags.AddError("Wrong object type", fmt.Sprintf("%s has %s %v; this data source reads %s only.", m.namespace, m.discriminator, obj[m.discriminator], m.variant))
		return nil, false
	}
	return obj, true
}

func (d *dataSource) stateType() tftypes.Object {
	t := objectType(d.attrs)
	if !d.model.singleton {
		t.AttributeTypes[dsFilters] = tftypes.String
	}
	return t
}

func (d *dataSource) setState(resp *datasource.ReadResponse, typ tftypes.Object, obj map[string]any, config map[string]tftypes.Value) {
	v, err := decodeObject(objectType(d.attrs), d.attrs, obj, tftypes.NewValue(objectType(d.attrs), nil))
	if err != nil {
		resp.Diagnostics.AddError("Unexpected read response", err.Error())
		return
	}
	vals := children(v)
	if !d.model.singleton {
		filters, ok := config[dsFilters]
		if !ok || !filters.IsKnown() {
			filters = tftypes.NewValue(tftypes.String, nil)
		}
		vals[dsFilters] = filters
	}
	resp.State.Raw = tftypes.NewValue(typ, vals)
}

func jsonFromValue(v tftypes.Value) jsontypes.Normalized {
	if !v.IsKnown() || v.IsNull() {
		return jsontypes.NewNormalizedNull()
	}
	var s string
	_ = v.As(&s)
	return jsontypes.NewNormalizedValue(s)
}

func (d *dataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs, err := dataSchemaAttributes(d.attrs)
	if err != nil {
		resp.Diagnostics.AddError("Invalid data source definition", err.Error())
		return
	}
	desc := fmt.Sprintf("Reads the current `%s` settings. Secrets are not exposed.", d.model.namespace)
	if !d.model.singleton {
		desc = fmt.Sprintf("Looks up exactly one existing `%s` object by `id`%s or `query_filters`. Secrets are not exposed.",
			d.model.namespace, map[bool]string{true: ", `name`", false: ""}[d.hasName()])
		id := attrs[idAttr]
		switch s := id.(type) {
		case dsschema.Int64Attribute:
			s.Optional, s.MarkdownDescription = true, "Identifier of the object to read."
			attrs[idAttr] = s
		case dsschema.StringAttribute:
			s.Optional, s.MarkdownDescription = true, "Identifier of the object to read."
			attrs[idAttr] = s
		}
		if d.hasName() {
			n := attrs[dsName].(dsschema.StringAttribute)
			n.Optional = true
			attrs[dsName] = n
		}
		attrs[dsFilters] = dsschema.StringAttribute{
			Optional:            true,
			CustomType:          jsontypes.NormalizedType{},
			MarkdownDescription: "Middleware query filters as JSON, e.g. `jsonencode([[\"path\", \"=\", \"/mnt/tank/x\"]])`. Must match exactly one object.",
		}
	}
	resp.Schema = dsschema.Schema{MarkdownDescription: desc, Attributes: attrs}
}
func dataSchemaAttributes(attrs []*node) (map[string]dsschema.Attribute, error) {
	out := make(map[string]dsschema.Attribute, len(attrs))
	for _, a := range attrs {
		built, err := dataSchemaAttribute(a)
		if err != nil {
			return nil, fmt.Errorf("attribute %s: %w", a.path, err)
		}
		out[a.name] = built
	}
	return out, nil
}
func dataSchemaAttribute(a *node) (dsschema.Attribute, error) {
	desc := a.description
	switch a.kind {
	case kindString:
		return dsschema.StringAttribute{Computed: true, MarkdownDescription: desc}, nil
	case kindInt:
		if a.identity {
			return dsschema.Int64Attribute{Computed: true, MarkdownDescription: desc}, nil
		}
		return dsschema.NumberAttribute{Computed: true, MarkdownDescription: desc}, nil
	case kindNumber:
		return dsschema.Float64Attribute{Computed: true, MarkdownDescription: desc}, nil
	case kindBool:
		return dsschema.BoolAttribute{Computed: true, MarkdownDescription: desc}, nil
	case kindAny:
		return dsschema.StringAttribute{Computed: true, CustomType: jsontypes.NormalizedType{}, MarkdownDescription: appendDesc(desc, "JSON-encoded.")}, nil
	case kindObject, kindUnion:
		children, err := dataSchemaAttributes(a.children)
		if err != nil {
			return nil, err
		}
		return dsschema.SingleNestedAttribute{Computed: true, MarkdownDescription: desc, Attributes: children}, nil
	case kindList:
		if len(a.elem.children) > 0 {
			children, err := dataSchemaAttributes(a.elem.children)
			if err != nil {
				return nil, err
			}
			return dsschema.ListNestedAttribute{Computed: true, MarkdownDescription: desc, NestedObject: dsschema.NestedAttributeObject{Attributes: children}}, nil
		}
		return dsschema.ListAttribute{Computed: true, MarkdownDescription: desc, ElementType: attrType(a.elem)}, nil
	case kindMap:
		if len(a.elem.children) > 0 {
			children, err := dataSchemaAttributes(a.elem.children)
			if err != nil {
				return nil, err
			}
			return dsschema.MapNestedAttribute{Computed: true, MarkdownDescription: desc, NestedObject: dsschema.NestedAttributeObject{Attributes: children}}, nil
		}
		return dsschema.MapAttribute{Computed: true, MarkdownDescription: desc, ElementType: attrType(a.elem)}, nil
	}
	return nil, fmt.Errorf("unsupported kind %s", a.kind)
}
