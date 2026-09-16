// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Client is what the engine needs from a transport: a way to call a method, and enough error
// classification to tell a missing object from a failure, a dropped connection from a refusal, and
// a validation failure from either. The transport owns those semantics because only it knows what
// its errors mean; the engine only asks the questions.
type Client interface {
	Call(ctx context.Context, method string, params ...any) (json.RawMessage, error)
	// IsNotFound reports whether err means the object does not exist.
	IsNotFound(err error) bool
	// IsRetryable reports whether err is a dropped connection, which an idempotent call may retry.
	IsRetryable(err error) bool
	// ValidationFields returns the per-attribute validation failures err carries, if any.
	ValidationFields(err error) []FieldError
}

// FieldError is one per-attribute validation failure reported by the server. Attribute is the
// server's own dotted path, which the engine maps onto the Terraform attribute it came from.
type FieldError struct {
	Attribute string
	Message   string
}

// FileClient transfers file contents over the HTTP endpoints, which the websocket API does not
// carry. Resources that need it assert for it on the configured client.
type FileClient interface {
	Client
	GetFile(ctx context.Context, path string, w io.Writer) error
	PutFile(ctx context.Context, path string, mode *int64, content io.Reader) error
}

// ProviderData is what the provider hands every resource in Configure.
type ProviderData struct {
	Client Client
	// ReadOnly makes any plan that would create, update or delete fail.
	ReadOnly bool
}

// crudResource serves one generated descriptor. Every generated resource embeds it and supplies
// its own Schema; the descriptor is already resolved, so there is nothing here to derive or fail.
type crudResource struct {
	model *model
	data  *ProviderData
}

func (r *crudResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.model.typeName
}

func (r *crudResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	var id identityschema.Attribute = identityschema.StringAttribute{RequiredForImport: true, Description: "Identifier assigned by TrueNAS."}
	if r.model != nil && r.model.idKind == kindInt {
		id = identityschema.Int64Attribute{RequiredForImport: true, Description: "Identifier assigned by TrueNAS."}
	}
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{idAttr: id}}
}

func (r *crudResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *engine.ProviderData, got %T", req.ProviderData))
		return
	}
	r.data = data
}

func (r *crudResource) ModifyPlan(_ context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.model != nil && !req.State.Raw.IsNull() && !req.Plan.Raw.IsNull() {
		resp.RequiresReplace.Append(r.model.createOnlyChanges(req.Plan.Raw, req.State.Raw)...)
	}
	if r.data == nil || !r.data.ReadOnly {
		return
	}
	switch {
	case req.State.Raw.IsNull():
		resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would create a %s, but the provider has read_only = true.", r.model.typeName))
	case req.Plan.Raw.IsNull():
		resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would destroy a %s, but the provider has read_only = true.", r.model.typeName))
	case len(resp.RequiresReplace) > 0:
		resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would replace a %s, but the provider has read_only = true.", r.model.typeName))
	case !req.Plan.Raw.Equal(req.State.Raw) && r.model != nil:
		// Changes that only update state (such as create-only values recorded after an import) do
		// not touch TrueNAS, so read-only mode allows them.
		params, err := r.model.updateParams(req.Plan.Raw, req.State.Raw, req.Config.Raw)
		if err != nil || len(params) > 0 {
			resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would update a %s, but the provider has read_only = true.", r.model.typeName))
		}
	}
}

func (r *crudResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.writable(&resp.Diagnostics) {
		return
	}
	m := r.model
	params, err := m.createParams(req.Plan.Raw, req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}
	if m.singleton {
		if len(params) > 0 {
			if _, err := r.data.Client.Call(ctx, m.updateMethod, params); err != nil {
				m.addCallError(r.data.Client, &resp.Diagnostics, "configure", err)
				return
			}
		}
		state, ok := r.readInto(ctx, nil, req.Plan.Raw, &resp.Diagnostics)
		if ok {
			resp.State.Raw = m.reconcileState(req.Plan.Raw, state)
			setIdentity(ctx, resp.Identity, m.singletonID(state), &resp.Diagnostics)
		}
		return
	}
	if m.adopt {
		r.adoptCreate(ctx, req, resp, params)
		return
	}
	raw, err := r.data.Client.Call(ctx, m.createMethod, params)
	if err != nil {
		m.addCallError(r.data.Client, &resp.Diagnostics, "create", err)
		return
	}
	id, err := m.idFromResult(raw)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected create response", err.Error())
		return
	}
	state, ok := r.readInto(ctx, id, req.Plan.Raw, &resp.Diagnostics)
	if !ok {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("Created resource disappeared", fmt.Sprintf("%s %v was created but could not be read back.", r.model.typeName, id))
		}
		return
	}
	resp.State.Raw = m.reconcileState(req.Plan.Raw, state)
	setIdentity(ctx, resp.Identity, id, &resp.Diagnostics)
}

func (r *crudResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if !r.ready(&resp.Diagnostics) {
		return
	}
	var id any
	if !r.model.singleton {
		var err error
		if id, err = r.model.idFromState(req.State.Raw); err != nil {
			resp.Diagnostics.AddError("Invalid state", err.Error())
			return
		}
	}
	state, ok := r.readInto(ctx, id, req.State.Raw, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.model.singleton {
		id = r.model.singletonID(state)
	}
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.State.Raw = state
	setIdentity(ctx, resp.Identity, id, &resp.Diagnostics)
}

func (r *crudResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if !r.ready(&resp.Diagnostics) {
		return
	}
	m := r.model
	var id any
	if !m.singleton {
		var err error
		if id, err = m.idFromState(req.State.Raw); err != nil {
			resp.Diagnostics.AddError("Invalid state", err.Error())
			return
		}
	}
	params, err := m.updateParams(req.Plan.Raw, req.State.Raw, req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}
	if len(params) > 0 {
		if !r.writable(&resp.Diagnostics) {
			return
		}
		args := []any{id, params}
		if m.singleton {
			args = []any{params}
		}
		if _, err := r.data.Client.Call(ctx, m.updateMethod, args...); err != nil {
			m.addCallError(r.data.Client, &resp.Diagnostics, "update", err)
			return
		}
	}
	state, ok := r.readInto(ctx, id, req.Plan.Raw, &resp.Diagnostics)
	if !ok {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("Resource disappeared", fmt.Sprintf("%s %v no longer exists.", r.model.typeName, id))
		}
		return
	}
	resp.State.Raw = m.reconcileState(req.Plan.Raw, state)
	if m.singleton {
		id = m.singletonID(state)
	}
	setIdentity(ctx, resp.Identity, id, &resp.Diagnostics)
}

func (r *crudResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.model != nil && (r.model.singleton || r.model.adopt) {
		resp.Diagnostics.AddWarning("Settings left unchanged",
			fmt.Sprintf("Removed %s from Terraform state. TrueNAS keeps its current settings.", r.model.typeName))
		return
	}
	if !r.writable(&resp.Diagnostics) {
		return
	}
	id, err := r.model.idFromState(req.State.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Invalid state", err.Error())
		return
	}
	if _, err := r.data.Client.Call(ctx, r.model.deleteMethod, id); err != nil && !r.data.Client.IsNotFound(err) {
		r.model.addCallError(r.data.Client, &resp.Diagnostics, "delete", err)
	}
}

func (r *crudResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var id any
	switch {
	case r.model.singleton:
		// There is one instance; any import ID selects it. Read fills in the real id.
		id = r.model.placeholderID()
	case req.ID != "":
		parsed, err := r.model.parseID(req.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid import ID", err.Error())
			return
		}
		id = parsed
	default:
		v, diags := identityID(ctx, req.Identity, r.model.idKind)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		id = v
	}
	if r.model.idKind == kindInt {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(idAttr), types.Int64Value(id.(int64)))...)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(idAttr), types.StringValue(id.(string)))...)
	}
	setIdentity(ctx, resp.Identity, id, &resp.Diagnostics)
}

func (r *crudResource) ready(diags *diag.Diagnostics) bool {
	if r.data == nil || r.data.Client == nil {
		diags.AddError("Provider not configured", "The provider has not been configured with a connection.")
		return false
	}
	return true
}

// writable is ready plus the read-only guard, for operations that change TrueNAS.
func (r *crudResource) writable(diags *diag.Diagnostics) bool {
	if !r.ready(diags) {
		return false
	}
	if r.data.ReadOnly {
		diags.AddError("Provider is read-only", "The provider has read_only = true and will not change TrueNAS.")
		return false
	}
	return true
}

// readInto fetches the resource and decodes it over prior. ok is false when it does not exist.
func (r *crudResource) readInto(ctx context.Context, id any, prior tftypes.Value, diags *diag.Diagnostics) (tftypes.Value, bool) {
	m := r.model
	args := []any{id}
	if m.singleton {
		args = nil
	}
	raw, err := callIdempotent(ctx, r.data.Client, m.getMethod, args...)
	if r.data.Client.IsNotFound(err) {
		return tftypes.Value{}, false
	}
	if err != nil {
		m.addCallError(r.data.Client, diags, "read", err)
		return tftypes.Value{}, false
	}
	obj, err := decodeResult(raw)
	if err != nil {
		diags.AddError("Unexpected read response", fmt.Sprintf("%s: %v", m.getMethod, err))
		return tftypes.Value{}, false
	}
	m2, ok := obj.(map[string]any)
	if !ok {
		diags.AddError("Unexpected read response", fmt.Sprintf("%s returned %T, expected an object", m.getMethod, obj))
		return tftypes.Value{}, false
	}
	if m.discriminator != "" && m2[m.discriminator] != m.variant {
		diags.AddError("Wrong resource type",
			fmt.Sprintf("%s %v has %s %v; this resource manages %s only.", m.namespace, id, m.discriminator, m2[m.discriminator], m.variant))
		return tftypes.Value{}, false
	}
	state, err := decodeObject(objectType(m.attrs), m.attrs, m2, prior)
	if err != nil {
		diags.AddError("Unexpected read response", err.Error())
		return tftypes.Value{}, false
	}
	return state, true
}

func (m *model) createParams(plan, config tftypes.Value) (map[string]any, error) {
	planned, configured := children(plan), children(config)
	out := map[string]any{}
	if m.discriminator != "" {
		out[m.discriminator] = m.variant
	}
	for _, a := range m.attrs {
		if a.name == idAttr || a.role == roleComputed || a.api == "" {
			continue
		}
		v := withWriteOnly(a, planned[a.name], configured[a.name])
		enc, ok, err := encodeValue(a, v, false)
		if err != nil {
			return nil, err
		}
		if ok {
			out[a.api] = enc
		}
	}
	return out, nil
}

func (m *model) updateParams(plan, state, config tftypes.Value) (map[string]any, error) {
	planned, prior, configured := children(plan), children(state), children(config)
	out := map[string]any{}
	for _, a := range m.attrs {
		if a.name == idAttr || a.role == roleComputed || a.api == "" || a.replace || !a.updatable {
			continue
		}
		v := planned[a.name]
		version := a.name + "_wo_version"
		switch {
		case a.writeOnly:
			if planned[version].Equal(prior[version]) {
				continue
			}
			v = configured[a.name]
		case containsWriteOnly(a):
			if v.Equal(prior[a.name]) && planned[version].Equal(prior[version]) {
				continue
			}
			v = withWriteOnly(a, v, configured[a.name])
		case !a.readable:
			// Only send what the configuration explicitly sets; state cannot tell us what the
			// server holds.
			if v = configured[a.name]; v.IsNull() || v.Equal(prior[a.name]) {
				continue
			}
		case v.Equal(prior[a.name]):
			continue
		}
		enc, ok, err := encodeValue(a, v, true)
		if err != nil {
			return nil, err
		}
		if ok {
			out[a.api] = enc
		}
	}
	return out, nil
}

func (m *model) reconcileState(plan, applied tftypes.Value) tftypes.Value {
	planned, got := children(plan), children(applied)
	vals := make(map[string]tftypes.Value, len(m.attrs))
	for _, a := range m.attrs {
		vals[a.name] = reconcile(a, planned[a.name], got[a.name])
	}
	return tftypes.NewValue(objectType(m.attrs), vals)
}

// adoptCreate finds an existing object by the adopt-by key and applies the configured values.
func (r *crudResource) adoptCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, params map[string]any) {
	m := r.model
	key := params[m.adoptBy]
	delete(params, m.adoptBy)
	raw, err := callIdempotent(ctx, r.data.Client, m.namespace+".query", []any{[]any{m.adoptBy, "=", key}}, map[string]any{})
	if err != nil {
		m.addCallError(r.data.Client, &resp.Diagnostics, "find", err)
		return
	}
	decoded, err := decodeResult(raw)
	rows, ok := decoded.([]any)
	if err != nil || !ok || len(rows) != 1 {
		resp.Diagnostics.AddAttributeError(path.Root(m.adoptByAttr), "Not found",
			fmt.Sprintf("No single %s has %s = %v; this resource manages objects that already exist.", m.namespace, m.adoptBy, key))
		return
	}
	row, _ := rows[0].(map[string]any)
	id, err := m.normalizeID(row[m.primaryKey])
	if err != nil {
		resp.Diagnostics.AddError("Unexpected response", err.Error())
		return
	}
	if len(params) > 0 {
		if _, err := r.data.Client.Call(ctx, m.updateMethod, id, params); err != nil {
			m.addCallError(r.data.Client, &resp.Diagnostics, "update", err)
			return
		}
	}
	state, ok := r.readInto(ctx, id, req.Plan.Raw, &resp.Diagnostics)
	if !ok {
		return
	}
	resp.State.Raw = m.reconcileState(req.Plan.Raw, state)
	setIdentity(ctx, resp.Identity, id, &resp.Diagnostics)
}

// createOnlyChanges returns create-only attributes whose known prior value the plan changes.
func (m *model) createOnlyChanges(plan, state tftypes.Value) []path.Path {
	planned, prior := children(plan), children(state)
	var paths []path.Path
	for _, a := range m.attrs {
		p, s := planned[a.name], prior[a.name]
		if a.createOnly && !s.IsNull() && p.IsKnown() && !p.Equal(s) {
			paths = append(paths, path.Root(a.name))
		}
	}
	return paths
}

// readRetries and readBackoff bound retries of reads across dropped connections.
var (
	readRetries = 5
	readBackoff = 500 * time.Millisecond
)

// callIdempotent retries a read whose connection dropped, which happens when a UI settings change
// restarts TrueNAS's web server. Only use it for methods without side effects.
func callIdempotent(ctx context.Context, c Client, method string, params ...any) (json.RawMessage, error) {
	backoff := readBackoff
	for attempt := 0; ; attempt++ {
		raw, err := c.Call(ctx, method, params...)
		if err == nil || !c.IsRetryable(err) || attempt == readRetries {
			return raw, err
		}
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		backoff *= 2
	}
}

// singletonID returns the id a config service reported in state.
func (m *model) singletonID(state tftypes.Value) any {
	id, err := m.idFromState(state)
	if err != nil {
		return m.placeholderID()
	}
	return id
}

func (m *model) placeholderID() any {
	if m.idKind == kindInt {
		return int64(1)
	}
	return m.namespace
}

func (m *model) idFromResult(raw json.RawMessage) (any, error) {
	v, err := decodeResult(raw)
	if err != nil {
		return nil, err
	}
	if obj, ok := v.(map[string]any); ok {
		v = obj[m.primaryKey]
	}
	return m.normalizeID(v)
}

func (m *model) idFromState(state tftypes.Value) (any, error) {
	v := children(state)[idAttr]
	if !v.IsKnown() || v.IsNull() {
		return nil, errors.New("resource state has no id")
	}
	if m.idKind == kindInt {
		var f big.Float
		if err := v.As(&f); err != nil {
			return nil, err
		}
		i, _ := f.Int64()
		return i, nil
	}
	var s string
	err := v.As(&s)
	return s, err
}

func (m *model) normalizeID(v any) (any, error) {
	switch m.idKind {
	case kindInt:
		f, err := bigFloat(v)
		if err != nil || !f.IsInt() {
			return nil, fmt.Errorf("expected an integer id, got %v", v)
		}
		i, _ := f.Int64()
		return i, nil
	default:
		s, ok := v.(string)
		if !ok || s == "" {
			return nil, fmt.Errorf("expected a string id, got %v", v)
		}
		return s, nil
	}
}

func (m *model) parseID(s string) (any, error) {
	if m.idKind == kindInt {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a numeric %s id", s, m.typeName)
		}
		return i, nil
	}
	return s, nil
}

func decodeResult(raw json.RawMessage) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	err := dec.Decode(&v)
	return v, err
}

func setIdentity(ctx context.Context, identity *tfsdk.ResourceIdentity, id any, diags *diag.Diagnostics) {
	if identity == nil {
		return
	}
	switch v := id.(type) {
	case int64:
		diags.Append(identity.SetAttribute(ctx, path.Root(idAttr), types.Int64Value(v))...)
	case string:
		diags.Append(identity.SetAttribute(ctx, path.Root(idAttr), types.StringValue(v))...)
	}
}

func identityID(ctx context.Context, identity *tfsdk.ResourceIdentity, kind kind) (any, diag.Diagnostics) {
	var diags diag.Diagnostics
	if identity == nil {
		diags.AddError("Missing import ID", "Provide an import ID or an identity.")
		return nil, diags
	}
	if kind == kindInt {
		var v types.Int64
		diags.Append(identity.GetAttribute(ctx, path.Root(idAttr), &v)...)
		if diags.HasError() || v.IsNull() || v.IsUnknown() {
			diags.AddError("Missing import identity", "identity.id is required.")
			return nil, diags
		}
		return v.ValueInt64(), diags
	}
	var v types.String
	diags.Append(identity.GetAttribute(ctx, path.Root(idAttr), &v)...)
	if diags.HasError() || v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		diags.AddError("Missing import identity", "identity.id is required.")
		return nil, diags
	}
	return v.ValueString(), diags
}
