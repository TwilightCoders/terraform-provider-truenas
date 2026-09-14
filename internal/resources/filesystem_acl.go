package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

var (
	_ resource.ResourceWithConfigure   = (*filesystemACL)(nil)
	_ resource.ResourceWithImportState = (*filesystemACL)(nil)
)

// NFS4 permission and flag names, as filesystem.getacl reports them.
var (
	nfs4Perms = []string{
		"READ_DATA", "WRITE_DATA", "APPEND_DATA", "READ_NAMED_ATTRS", "WRITE_NAMED_ATTRS", "EXECUTE",
		"DELETE", "DELETE_CHILD", "READ_ATTRIBUTES", "WRITE_ATTRIBUTES", "READ_ACL", "WRITE_ACL", "WRITE_OWNER", "SYNCHRONIZE",
	}
	nfs4Flags        = []string{"FILE_INHERIT", "DIRECTORY_INHERIT", "NO_PROPAGATE_INHERIT", "INHERIT_ONLY", "INHERITED"}
	nfs4BasicPerms   = []string{"FULL_CONTROL", "MODIFY", "READ", "TRAVERSE"}
	nfs4BasicFlags   = []string{"INHERIT", "NOINHERIT"}
	nfs4Tags         = []string{"owner@", "group@", "everyone@", "USER", "GROUP"}
	posixTags        = []string{"USER_OBJ", "GROUP_OBJ", "OTHER", "MASK", "USER", "GROUP"}
	aclEntryAttrType = map[string]attr.Type{
		"tag": types.StringType, "type": types.StringType, "id": types.Int64Type,
		"perms": types.StringType, "advanced_perms": types.SetType{ElemType: types.StringType},
		"flags": types.StringType, "advanced_flags": types.SetType{ElemType: types.StringType},
		"read": types.BoolType, "write": types.BoolType, "execute": types.BoolType, "default": types.BoolType,
	}
)

// NewFilesystemACL returns the truenas_filesystem_acl resource.
func NewFilesystemACL() resource.Resource { return &filesystemACL{} }

type filesystemACL struct {
	data *engine.ProviderData
}

type aclModel struct {
	Path      types.String `tfsdk:"path"`
	ACLType   types.String `tfsdk:"acltype"`
	UID       types.Int64  `tfsdk:"uid"`
	GID       types.Int64  `tfsdk:"gid"`
	Entries   types.Set    `tfsdk:"entries"`
	Recursive types.Bool   `tfsdk:"recursive"`
	Traverse  types.Bool   `tfsdk:"traverse"`
}

type aclEntry struct {
	Tag           types.String `tfsdk:"tag"`
	Type          types.String `tfsdk:"type"`
	ID            types.Int64  `tfsdk:"id"`
	Perms         types.String `tfsdk:"perms"`
	AdvancedPerms types.Set    `tfsdk:"advanced_perms"`
	Flags         types.String `tfsdk:"flags"`
	AdvancedFlags types.Set    `tfsdk:"advanced_flags"`
	Read          types.Bool   `tfsdk:"read"`
	Write         types.Bool   `tfsdk:"write"`
	Execute       types.Bool   `tfsdk:"execute"`
	Default       types.Bool   `tfsdk:"default"`
}

func (r *filesystemACL) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_filesystem_acl"
}

func (r *filesystemACL) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the access control list of a path (filesystem.getacl / filesystem.setacl). " +
			"TrueNAS orders entries canonically, so `entries` is a set. Destroying the resource leaves the ACL in place.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required: true, MarkdownDescription: "Absolute path, e.g. `/mnt/tank/share`. Changing it forces a new resource.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"acltype": schema.StringAttribute{
				Optional: true, Computed: true,
				MarkdownDescription: "`NFS4` or `POSIX1E`. Defaults to what the dataset uses.",
				Validators:          []validator.String{stringvalidator.OneOf("NFS4", "POSIX1E")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"uid": schema.Int64Attribute{
				Optional: true, Computed: true, MarkdownDescription: "Owning user id.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"gid": schema.Int64Attribute{
				Optional: true, Computed: true, MarkdownDescription: "Owning group id.",
				PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"recursive": schema.BoolAttribute{
				Optional: true, Computed: true, Default: booldefault.StaticBool(false),
				MarkdownDescription: "Apply the ACL to everything below `path` as well. Only affects writes; not detected on read.",
			},
			"traverse": schema.BoolAttribute{
				Optional: true, Computed: true, Default: booldefault.StaticBool(false),
				MarkdownDescription: "With `recursive`, cross into child datasets.",
			},
			"entries": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "Access control entries.",
				Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"tag": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "NFS4: `owner@`, `group@`, `everyone@`, `USER`, `GROUP`. POSIX1E: `USER_OBJ`, `GROUP_OBJ`, `OTHER`, `MASK`, `USER`, `GROUP`.",
						Validators:          []validator.String{stringvalidator.OneOf(append(append([]string{}, nfs4Tags...), posixTags[:4]...)...)},
					},
					"id":   schema.Int64Attribute{Optional: true, MarkdownDescription: "User or group id for `USER` and `GROUP` entries."},
					"type": schema.StringAttribute{Optional: true, MarkdownDescription: "NFS4 only: `ALLOW` or `DENY`.", Validators: []validator.String{stringvalidator.OneOf("ALLOW", "DENY")}},
					"perms": schema.StringAttribute{
						Optional: true, MarkdownDescription: "NFS4 basic permission set: `FULL_CONTROL`, `MODIFY`, `READ` or `TRAVERSE`. Conflicts with `advanced_perms`.",
						Validators: []validator.String{stringvalidator.OneOf(nfs4BasicPerms...)},
					},
					"advanced_perms": schema.SetAttribute{
						Optional: true, ElementType: types.StringType,
						MarkdownDescription: "NFS4 individual permissions that are granted, e.g. `READ_DATA`.",
						Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.OneOf(nfs4Perms...))},
					},
					"flags": schema.StringAttribute{
						Optional: true, MarkdownDescription: "NFS4 basic inheritance: `INHERIT` or `NOINHERIT`. Conflicts with `advanced_flags`.",
						Validators: []validator.String{stringvalidator.OneOf(nfs4BasicFlags...)},
					},
					"advanced_flags": schema.SetAttribute{
						Optional: true, ElementType: types.StringType,
						MarkdownDescription: "NFS4 individual inheritance flags that are set, e.g. `FILE_INHERIT`.",
						Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.OneOf(nfs4Flags...))},
					},
					"read":    schema.BoolAttribute{Optional: true, MarkdownDescription: "POSIX1E read permission."},
					"write":   schema.BoolAttribute{Optional: true, MarkdownDescription: "POSIX1E write permission."},
					"execute": schema.BoolAttribute{Optional: true, MarkdownDescription: "POSIX1E execute permission."},
					"default": schema.BoolAttribute{Optional: true, MarkdownDescription: "POSIX1E: a default (inherited) entry."},
				}},
			},
		},
	}
}

func (r *filesystemACL) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *filesystemACL) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aclModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || !r.writable(&resp.Diagnostics) {
		return
	}
	resp.Diagnostics.Append(r.apply(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *filesystemACL) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aclModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || !r.writable(&resp.Diagnostics) {
		return
	}
	resp.Diagnostics.Append(r.apply(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *filesystemACL) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aclModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if r.data == nil {
		resp.Diagnostics.AddError("Provider not configured", "The TrueNAS provider has not been configured with a connection.")
		return
	}
	found, diags := r.read(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *filesystemACL) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("ACL left in place",
		"Removed the ACL from Terraform state. TrueNAS keeps it: stripping ACLs from data is never done implicitly.")
}

func (r *filesystemACL) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("path"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("recursive"), false)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("traverse"), false)...)
}

func (r *filesystemACL) writable(diags *diag.Diagnostics) bool {
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

// apply sets the ACL and refreshes the model from what TrueNAS stored.
func (r *filesystemACL) apply(ctx context.Context, m *aclModel) diag.Diagnostics {
	var diags diag.Diagnostics
	var entries []aclEntry
	diags.Append(m.Entries.ElementsAs(ctx, &entries, false)...)
	if diags.HasError() {
		return diags
	}
	acltype := m.ACLType.ValueString()
	dacl := make([]map[string]any, 0, len(entries))
	for i, e := range entries {
		wire, err := entryToAPI(e, acltype)
		if err != nil {
			diags.AddAttributeError(path.Root("entries"), "Invalid ACL entry", fmt.Sprintf("entry %d: %v", i+1, err))
			continue
		}
		dacl = append(dacl, wire)
	}
	if diags.HasError() {
		return diags
	}

	params := map[string]any{
		"path": m.Path.ValueString(),
		"dacl": dacl,
		"options": map[string]any{
			"recursive": m.Recursive.ValueBool(),
			"traverse":  m.Traverse.ValueBool(),
		},
	}
	if acltype != "" {
		params["acltype"] = acltype
	}
	if !m.UID.IsNull() && !m.UID.IsUnknown() {
		params["uid"] = m.UID.ValueInt64()
	}
	if !m.GID.IsNull() && !m.GID.IsUnknown() {
		params["gid"] = m.GID.ValueInt64()
	}
	if _, err := r.data.Client.Call(ctx, "filesystem.setacl", params); err != nil {
		diags.AddError("Unable to set ACL", err.Error())
		return diags
	}
	planned := *m
	found, readDiags := r.read(ctx, m)
	diags.Append(readDiags...)
	if !found && !diags.HasError() {
		diags.AddError("Path disappeared", fmt.Sprintf("%s no longer exists.", m.Path.ValueString()))
	}
	// Terraform requires known planned values back; a server-side difference shows as drift on refresh.
	m.Entries = planned.Entries
	if !planned.ACLType.IsUnknown() {
		m.ACLType = planned.ACLType
	}
	if !planned.UID.IsUnknown() && !planned.UID.IsNull() {
		m.UID = planned.UID
	}
	if !planned.GID.IsUnknown() && !planned.GID.IsNull() {
		m.GID = planned.GID
	}
	return diags
}

type wireACL struct {
	ACLType string      `json:"acltype"`
	UID     *int64      `json:"uid"`
	GID     *int64      `json:"gid"`
	ACL     []wireEntry `json:"acl"`
}

type wireEntry struct {
	Tag     string          `json:"tag"`
	Type    string          `json:"type"`
	ID      *int64          `json:"id"`
	Perms   json.RawMessage `json:"perms"`
	Flags   json.RawMessage `json:"flags"`
	Default bool            `json:"default"`
}

// read refreshes m from filesystem.getacl. found is false when the path does not exist.
func (r *filesystemACL) read(ctx context.Context, m *aclModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	raw, err := r.data.Client.Call(ctx, "filesystem.getacl", m.Path.ValueString(), true, false)
	if err != nil {
		if middleware.IsNotFound(err) {
			return false, diags
		}
		diags.AddError("Unable to read ACL", err.Error())
		return false, diags
	}
	var acl wireACL
	if err := json.Unmarshal(raw, &acl); err != nil {
		diags.AddError("Unexpected ACL response", err.Error())
		return false, diags
	}
	if acl.ACLType == "DISABLED" {
		diags.AddError("ACLs disabled", fmt.Sprintf("%s has ACL support disabled; set the dataset's acltype first.", m.Path.ValueString()))
		return false, diags
	}
	m.ACLType = types.StringValue(acl.ACLType)
	m.UID = int64Ptr(acl.UID)
	m.GID = int64Ptr(acl.GID)

	values := make([]attr.Value, 0, len(acl.ACL))
	for _, e := range acl.ACL {
		entry, err := entryFromAPI(e, acl.ACLType)
		if err != nil {
			diags.AddError("Unexpected ACL entry", err.Error())
			return false, diags
		}
		obj, d := types.ObjectValueFrom(ctx, aclEntryAttrType, entry)
		diags.Append(d...)
		values = append(values, obj)
	}
	set, d := types.SetValue(types.ObjectType{AttrTypes: aclEntryAttrType}, values)
	diags.Append(d...)
	m.Entries = set
	return true, diags
}

func entryToAPI(e aclEntry, acltype string) (map[string]any, error) {
	out := map[string]any{"tag": e.Tag.ValueString()}
	if !e.ID.IsNull() {
		out["id"] = e.ID.ValueInt64()
	} else {
		out["id"] = -1
	}
	if acltype == "POSIX1E" {
		if e.Read.IsNull() || e.Write.IsNull() || e.Execute.IsNull() || e.Default.IsNull() {
			return nil, fmt.Errorf("POSIX1E entries need read, write, execute and default")
		}
		out["default"] = e.Default.ValueBool()
		out["perms"] = map[string]bool{"READ": e.Read.ValueBool(), "WRITE": e.Write.ValueBool(), "EXECUTE": e.Execute.ValueBool()}
		return out, nil
	}

	if e.Type.IsNull() {
		return nil, fmt.Errorf("NFS4 entries need type (ALLOW or DENY)")
	}
	if e.Perms.IsNull() && e.AdvancedPerms.IsNull() {
		return nil, fmt.Errorf("NFS4 entries need perms or advanced_perms")
	}
	if e.Flags.IsNull() && e.AdvancedFlags.IsNull() {
		return nil, fmt.Errorf("NFS4 entries need flags or advanced_flags")
	}
	out["type"] = e.Type.ValueString()
	perms, err := basicOrAdvanced(e.Perms, e.AdvancedPerms, nfs4Perms, "perms")
	if err != nil {
		return nil, err
	}
	flags, err := basicOrAdvanced(e.Flags, e.AdvancedFlags, nfs4Flags, "flags")
	if err != nil {
		return nil, err
	}
	out["perms"], out["flags"] = perms, flags
	return out, nil
}

func basicOrAdvanced(basic types.String, advanced types.Set, names []string, attrName string) (map[string]any, error) {
	hasBasic, hasAdvanced := !basic.IsNull(), !advanced.IsNull()
	switch {
	case hasBasic && hasAdvanced:
		return nil, fmt.Errorf("set %s or advanced_%s, not both", attrName, attrName)
	case hasBasic:
		return map[string]any{"BASIC": basic.ValueString()}, nil
	}
	granted := map[string]bool{}
	if hasAdvanced {
		for _, v := range advanced.Elements() {
			granted[v.(types.String).ValueString()] = true
		}
	}
	out := make(map[string]any, len(names))
	for _, n := range names {
		out[n] = granted[n]
	}
	return out, nil
}

func entryFromAPI(e wireEntry, acltype string) (aclEntry, error) {
	entry := aclEntry{
		Tag: types.StringValue(e.Tag), ID: types.Int64Null(), Type: types.StringNull(),
		Perms: types.StringNull(), AdvancedPerms: types.SetNull(types.StringType),
		Flags: types.StringNull(), AdvancedFlags: types.SetNull(types.StringType),
		Read: types.BoolNull(), Write: types.BoolNull(), Execute: types.BoolNull(), Default: types.BoolNull(),
	}
	if e.ID != nil && *e.ID >= 0 {
		entry.ID = types.Int64Value(*e.ID)
	}
	if acltype == "POSIX1E" {
		var perms map[string]bool
		if err := json.Unmarshal(e.Perms, &perms); err != nil {
			return entry, err
		}
		entry.Read, entry.Write, entry.Execute = types.BoolValue(perms["READ"]), types.BoolValue(perms["WRITE"]), types.BoolValue(perms["EXECUTE"])
		entry.Default = types.BoolValue(e.Default)
		return entry, nil
	}

	entry.Type = types.StringValue(e.Type)
	var err error
	if entry.Perms, entry.AdvancedPerms, err = basicOrAdvancedFromAPI(e.Perms); err != nil {
		return entry, err
	}
	entry.Flags, entry.AdvancedFlags, err = basicOrAdvancedFromAPI(e.Flags)
	return entry, err
}

func basicOrAdvancedFromAPI(raw json.RawMessage) (types.String, types.Set, error) {
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		return types.StringNull(), types.SetNull(types.StringType), err
	}
	if basic, ok := m["BASIC"].(string); ok {
		return types.StringValue(basic), types.SetNull(types.StringType), nil
	}
	var granted []string
	for name, v := range m {
		if b, ok := v.(bool); ok && b {
			granted = append(granted, name)
		}
	}
	sort.Strings(granted)
	elems := make([]attr.Value, len(granted))
	for i, g := range granted {
		elems[i] = types.StringValue(g)
	}
	set, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		return types.StringNull(), types.SetNull(types.StringType), fmt.Errorf("%v", diags)
	}
	return types.StringNull(), set, nil
}

func int64Ptr(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}
