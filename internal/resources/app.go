package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

var (
	_ resource.ResourceWithConfigure        = (*app)(nil)
	_ resource.ResourceWithImportState      = (*app)(nil)
	_ resource.ResourceWithModifyPlan       = (*app)(nil)
	_ resource.ResourceWithConfigValidators = (*app)(nil)
)

var appNamePattern = regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`)

// NewApp returns the truenas_app resource: a TrueNAS custom app driven through app.*.
func NewApp() resource.Resource { return &app{} }

type app struct {
	data *engine.ProviderData
}

type appModel struct {
	Name            types.String `tfsdk:"name"`
	Include         types.List   `tfsdk:"include"`
	Compose         types.String `tfsdk:"compose"`
	Portals         types.Map    `tfsdk:"portals"`
	Notes           types.String `tfsdk:"notes"`
	DesiredState    types.String `tfsdk:"desired_state"`
	State           types.String `tfsdk:"state"`
	RedeployTrigger types.String `tfsdk:"redeploy_trigger"`
	Hold            types.Bool   `tfsdk:"hold"`
	DeleteImages    types.Bool   `tfsdk:"delete_images"`
	DeleteIXVolumes types.Bool   `tfsdk:"delete_ix_volumes"`
}

func (r *app) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *app) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a TrueNAS custom app: a Docker Compose project that TrueNAS owns, so the Apps UI, " +
			"lifecycle and status all work. Every change to a running app's compose restarts its containers. " +
			"Catalog apps are refused; convert them with `app.convert_to_custom` first.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "App name: lowercase letters, digits and dashes, up to 40 characters. Changing it forces a new app.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
					stringvalidator.RegexMatches(appNamePattern, "must start with a letter and contain only lowercase letters, digits and dashes"),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"include": schema.ListAttribute{
				Optional: true, ElementType: types.StringType,
				MarkdownDescription: "Absolute paths of compose files on TrueNAS to include. The files are not managed here; " +
					"`.env` interpolation and relative paths resolve against each file's directory.",
				Validators: []validator.List{listvalidator.SizeAtLeast(1)},
			},
			"compose": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Compose YAML stored inline in TrueNAS. Compared by meaning, not text. " +
					"Relative paths and `.env` resolve against TrueNAS's internal render directory; prefer `include` when they matter.",
			},
			"portals": schema.MapAttribute{
				Optional: true, Computed: true, ElementType: types.StringType,
				MarkdownDescription: "Links shown in the Apps UI, name to URL, e.g. `{ \"Web UI\" = \"http://plex.lan/\" }`. " +
					"TrueNAS reads them only when an app is created, so changing them replaces the app.",
				PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
			"notes": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "Markdown notes shown in the Apps UI. Read only when an app is created, so changing them replaces the app.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"desired_state": schema.StringAttribute{
				Optional: true, Computed: true, Default: stringdefault.StaticString("RUNNING"),
				MarkdownDescription: "`RUNNING` or `STOPPED`. Stopping removes the containers; volumes and images stay.",
				Validators:          []validator.String{stringvalidator.OneOf("RUNNING", "STOPPED")},
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "State reported by TrueNAS: `RUNNING`, `DEPLOYING`, `STOPPED`, `STOPPING` or `CRASHED`.",
			},
			"redeploy_trigger": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "Any value; changing it redeploys the app without changing its configuration. " +
					"Use `filesha256()` of included compose files that are shipped outside Terraform.",
			},
			"hold": schema.BoolAttribute{
				Optional: true, Computed: true, Default: booldefault.StaticBool(false),
				MarkdownDescription: "Refuse any plan that would change, redeploy, restart or destroy the app, while still " +
					"refreshing it. For stateful apps that are rolled out by hand. Setting it back to `false` is always allowed.",
			},
			"delete_images": schema.BoolAttribute{
				Optional: true, Computed: true, Default: booldefault.StaticBool(false),
				MarkdownDescription: "Delete the app's images when destroying it. Off by default so locally built images survive.",
			},
			"delete_ix_volumes": schema.BoolAttribute{
				Optional: true, Computed: true, Default: booldefault.StaticBool(false),
				MarkdownDescription: "Delete TrueNAS-managed ix volumes when destroying the app. Host paths are never touched.",
			},
		},
	}
}

func (r *app) ConfigValidators(context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(path.MatchRoot("include"), path.MatchRoot("compose")),
	}
}

func (r *app) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*engine.ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *engine.ProviderData, got %T", req.ProviderData))
		return
	}
	r.data = data
}

func (r *app) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	readOnly := r.data != nil && r.data.ReadOnly
	if req.State.Raw.IsNull() {
		if readOnly && !req.Plan.Raw.IsNull() {
			resp.Diagnostics.AddError("Provider is read-only", "This plan would create an app, but the provider has read_only = true.")
		}
		return
	}
	var state appModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.Plan.Raw.IsNull() {
		if state.Hold.ValueBool() {
			resp.Diagnostics.AddError("App is held", fmt.Sprintf("%s has hold = true; set hold = false before destroying it.", state.Name.ValueString()))
		}
		if readOnly {
			resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would destroy %s, but the provider has read_only = true.", state.Name.ValueString()))
		}
		return
	}

	var plan appModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Portals and notes are read only on create.
	if known(plan.Portals) && !state.Portals.IsNull() && !portalsEqual(plan.Portals, state.Portals) {
		resp.RequiresReplace.Append(path.Root("portals"))
	}
	if known(plan.Notes) && !state.Notes.IsNull() && strings.TrimSpace(plan.Notes.ValueString()) != strings.TrimSpace(state.Notes.ValueString()) {
		resp.RequiresReplace.Append(path.Root("notes"))
	}

	if readOnly {
		switch calls := plannedCalls(plan, state); {
		case len(resp.RequiresReplace) > 0:
			resp.Diagnostics.AddError("Provider is read-only", fmt.Sprintf("This plan would replace %s, but the provider has read_only = true.", state.Name.ValueString()))
		case len(calls) > 0:
			resp.Diagnostics.AddError("Provider is read-only",
				fmt.Sprintf("This plan would call %s on %s, but the provider has read_only = true.", strings.Join(calls, ", "), state.Name.ValueString()))
		}
	}

	if state.Hold.ValueBool() && plan.Hold.ValueBool() && appChanges(plan, state) {
		resp.Diagnostics.AddError("App is held",
			fmt.Sprintf("%s has hold = true, and this plan would change, redeploy or restart it. Roll it out by hand, or set hold = false.", state.Name.ValueString()))
	}
}

func known(v interface {
	IsNull() bool
	IsUnknown() bool
},
) bool {
	return !v.IsNull() && !v.IsUnknown()
}

// appChanges reports whether plan differs from state in anything that touches the running app.
func appChanges(plan, state appModel) bool {
	return !plan.Include.Equal(state.Include) || !plan.Compose.Equal(state.Compose) ||
		(known(plan.Portals) && !portalsEqual(plan.Portals, state.Portals)) ||
		(known(plan.Notes) && !plan.Notes.Equal(state.Notes)) ||
		!plan.DesiredState.Equal(state.DesiredState) || !plan.RedeployTrigger.Equal(state.RedeployTrigger)
}

func (r *app) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan appModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || !r.writable(&resp.Diagnostics) {
		return
	}
	config, diags := composeConfig(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := plan.Name.ValueString()
	if _, err := r.data.Client.Call(ctx, "app.create", map[string]any{
		"app_name": name, "custom_app": true, "custom_compose_config": config,
	}); err != nil {
		resp.Diagnostics.AddError("Unable to create app", err.Error())
		return
	}
	if plan.DesiredState.ValueString() == "STOPPED" {
		if _, err := r.data.Client.Call(ctx, "app.stop", name); err != nil {
			resp.Diagnostics.AddError("Unable to stop app", err.Error())
			return
		}
	}
	r.refresh(ctx, &plan, &resp.Diagnostics, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *app) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state appModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.data == nil {
		resp.Diagnostics.AddError("Provider not configured", "The TrueNAS provider has not been configured with a connection.")
		return
	}
	if !r.refresh(ctx, &state, &resp.Diagnostics, false) {
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// plannedCalls lists the app.* methods an update from state to plan would call, in order. Attributes
// that only configure the provider (hold, delete_images, delete_ix_volumes) call nothing.
func plannedCalls(plan, state appModel) []string {
	var calls []string
	switch {
	case !plan.Include.Equal(state.Include) || !plan.Compose.Equal(state.Compose):
		calls = append(calls, "app.update")
	case !plan.RedeployTrigger.Equal(state.RedeployTrigger) && plan.DesiredState.ValueString() == "RUNNING":
		calls = append(calls, "app.redeploy")
	}
	switch {
	case plan.DesiredState.ValueString() == "STOPPED" && state.State.ValueString() != "STOPPED":
		calls = append(calls, "app.stop")
	case plan.DesiredState.ValueString() == "RUNNING" && state.State.ValueString() == "STOPPED":
		calls = append(calls, "app.start")
	}
	return calls
}

func (r *app) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state appModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.data == nil {
		resp.Diagnostics.AddError("Provider not configured", "The TrueNAS provider has not been configured with a connection.")
		return
	}
	calls := plannedCalls(plan, state)
	if len(calls) > 0 && !r.writable(&resp.Diagnostics) {
		return
	}
	name := plan.Name.ValueString()
	for _, method := range calls {
		var err error
		switch method {
		case "app.update":
			config, diags := composeConfig(ctx, plan)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			_, err = r.data.Client.Call(ctx, method, name, map[string]any{"custom_compose_config": config})
		default:
			_, err = r.data.Client.Call(ctx, method, name)
		}
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("%s failed", method), err.Error())
			return
		}
	}
	r.refresh(ctx, &plan, &resp.Diagnostics, true)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *app) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state appModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() || !r.writable(&resp.Diagnostics) {
		return
	}
	_, err := r.data.Client.Call(ctx, "app.delete", state.Name.ValueString(), map[string]any{
		"remove_images":     state.DeleteImages.ValueBool(),
		"remove_ix_volumes": state.DeleteIXVolumes.ValueBool(),
	})
	if err != nil && !middleware.IsNotFound(err) {
		resp.Diagnostics.AddError("Unable to delete app", err.Error())
	}
}

func (r *app) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	for attr, v := range map[string]any{"hold": false, "delete_images": false, "delete_ix_volumes": false, "desired_state": "RUNNING"} {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(attr), v)...)
	}
}

func (r *app) writable(diags *diag.Diagnostics) bool {
	switch {
	case r.data == nil:
		diags.AddError("Provider not configured", "The TrueNAS provider has not been configured with a connection.")
	case r.data.ReadOnly:
		diags.AddError("Provider is read-only", "The provider has read_only = true and will not change TrueNAS.")
	default:
		return true
	}
	return false
}

type appEntry struct {
	Name      string            `json:"name"`
	State     string            `json:"state"`
	CustomApp bool              `json:"custom_app"`
	Portals   map[string]string `json:"portals"`
	Notes     *string           `json:"notes"`
	Config    map[string]any    `json:"config"`
}

// refresh reads the app into m. It keeps configured values that are equivalent to what TrueNAS
// reports, so formatting differences never plan changes. afterWrite keeps planned values that
// TrueNAS normalizes. It returns false when the app does not exist.
func (r *app) refresh(ctx context.Context, m *appModel, diags *diag.Diagnostics, afterWrite bool) bool {
	name := m.Name.ValueString()
	raw, err := r.data.Client.Call(ctx, "app.query", []any{[]any{"name", "=", name}}, map[string]any{"extra": map[string]any{"retrieve_config": true}})
	if err != nil {
		diags.AddError("Unable to read app", err.Error())
		return false
	}
	var entries []appEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		diags.AddError("Unexpected app response", err.Error())
		return false
	}
	if len(entries) == 0 {
		if afterWrite {
			diags.AddError("App disappeared", fmt.Sprintf("%s could not be read back.", name))
		}
		return false
	}
	e := entries[0]
	if !e.CustomApp {
		diags.AddError("Catalog app",
			fmt.Sprintf("%s is a catalog app. Its configuration comes from the catalog chart, so a compose file here would be silently ignored. "+
				"Convert it with app.convert_to_custom before managing it.", name))
		return false
	}

	m.State = types.StringValue(e.State)
	delete(e.Config, "x-portals")
	delete(e.Config, "x-notes")
	if includes, ok := onlyIncludes(e.Config); ok {
		list, d := types.ListValueFrom(ctx, types.StringType, includes)
		diags.Append(d...)
		if !m.Include.Equal(list) {
			m.Include = list
		}
		if !afterWrite || m.Compose.IsUnknown() {
			m.Compose = types.StringNull()
		}
	} else {
		if !yamlEqual(m.Compose.ValueString(), e.Config) {
			out, err := yaml.Marshal(e.Config)
			if err != nil {
				diags.AddError("Unexpected app config", err.Error())
				return false
			}
			m.Compose = types.StringValue(string(out))
		}
		if !afterWrite || m.Include.IsUnknown() {
			m.Include = types.ListNull(types.StringType)
		}
	}

	portals, d := types.MapValueFrom(ctx, types.StringType, e.Portals)
	diags.Append(d...)
	if !portalsEqual(m.Portals, portals) {
		m.Portals = portals
	}
	notes := types.StringNull()
	if e.Notes != nil {
		notes = types.StringValue(*e.Notes)
	}
	if strings.TrimSpace(m.Notes.ValueString()) != strings.TrimSpace(notes.ValueString()) || m.Notes.IsUnknown() || m.Notes.IsNull() != notes.IsNull() {
		m.Notes = notes
	}
	return true
}

func onlyIncludes(config map[string]any) ([]string, bool) {
	if len(config) != 1 {
		return nil, false
	}
	list, ok := config["include"].([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		s, ok := v.(string)
		if !ok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// composeConfig builds the stored compose: the includes or parsed compose, plus x-portals and x-notes.
func composeConfig(ctx context.Context, m appModel) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics
	config := map[string]any{}
	switch {
	case known(m.Include):
		var includes []string
		diags.Append(m.Include.ElementsAs(ctx, &includes, false)...)
		config["include"] = includes
	case known(m.Compose):
		if err := yaml.Unmarshal([]byte(m.Compose.ValueString()), &config); err != nil {
			diags.AddAttributeError(path.Root("compose"), "Invalid compose YAML", err.Error())
			return nil, diags
		}
		if config == nil {
			diags.AddAttributeError(path.Root("compose"), "Invalid compose YAML", "compose must be a YAML mapping")
			return nil, diags
		}
		delete(config, "x-portals")
		delete(config, "x-notes")
	}
	if known(m.Portals) {
		portals := map[string]string{}
		diags.Append(m.Portals.ElementsAs(ctx, &portals, false)...)
		entries, err := xPortals(portals)
		if err != nil {
			diags.AddAttributeError(path.Root("portals"), "Invalid portal", err.Error())
			return nil, diags
		}
		if len(entries) > 0 {
			config["x-portals"] = entries
		}
	}
	if known(m.Notes) {
		config["x-notes"] = m.Notes.ValueString()
	}
	return config, diags
}

// xPortals converts name→URL portals into the x-portals entries TrueNAS validates.
func xPortals(portals map[string]string) ([]map[string]any, error) {
	names := make([]string, 0, len(portals))
	for name := range portals {
		names = append(names, name)
	}
	sort.Strings(names)
	out := make([]map[string]any, 0, len(names))
	for _, name := range names {
		u, err := url.Parse(portals[name])
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
			return nil, fmt.Errorf("%s: %q must be an http or https URL with a host", name, portals[name])
		}
		port := defaultPort(u.Scheme)
		if p := u.Port(); p != "" {
			if port, err = strconv.Atoi(p); err != nil {
				return nil, fmt.Errorf("%s: invalid port in %q", name, portals[name])
			}
		}
		portalPath := u.EscapedPath()
		if portalPath == "" {
			portalPath = "/"
		}
		out = append(out, map[string]any{"name": name, "scheme": u.Scheme, "host": u.Hostname(), "port": port, "path": portalPath})
	}
	return out, nil
}

func defaultPort(scheme string) int {
	if scheme == "https" {
		return 443
	}
	return 80
}

// portalsEqual compares portal maps by normalized URL: default ports and an empty path are implied.
func portalsEqual(a, b types.Map) bool {
	if a.IsNull() || b.IsNull() || a.IsUnknown() || b.IsUnknown() {
		return a.IsNull() == b.IsNull() && a.IsUnknown() == b.IsUnknown()
	}
	ae, be := a.Elements(), b.Elements()
	if len(ae) != len(be) {
		return false
	}
	for name, av := range ae {
		bv, ok := be[name]
		if !ok || normalizeURL(av.(types.String).ValueString()) != normalizeURL(bv.(types.String).ValueString()) {
			return false
		}
	}
	return true
}

func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if port == "" || port == strconv.Itoa(defaultPort(u.Scheme)) {
		port = ""
	} else {
		port = ":" + port
	}
	p := u.EscapedPath()
	if p == "" {
		p = "/"
	}
	return strings.ToLower(u.Scheme) + "://" + host + port + p
}

// yamlEqual reports whether text parses to the same document as want.
func yamlEqual(text string, want map[string]any) bool {
	if text == "" {
		return false
	}
	var got map[string]any
	if yaml.Unmarshal([]byte(text), &got) != nil {
		return false
	}
	delete(got, "x-portals")
	delete(got, "x-notes")
	a, errA := json.Marshal(got)
	b, errB := json.Marshal(want)
	if errA != nil || errB != nil {
		return reflect.DeepEqual(got, want)
	}
	var an, bn any
	_ = json.NewDecoder(bytes.NewReader(a)).Decode(&an)
	_ = json.NewDecoder(bytes.NewReader(b)).Decode(&bn)
	return reflect.DeepEqual(an, bn)
}
