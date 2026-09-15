package resources

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

var (
	_ resource.ResourceWithConfigure   = (*fileResource)(nil)
	_ resource.ResourceWithImportState = (*fileResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*fileResource)(nil)
)

// Drift detection modes.
const (
	driftContent = "content" // download and hash on every refresh
	driftStat    = "stat"    // compare size and mtime, and only then download
	driftNone    = "none"    // never read the file back
)

// Destroy behaviours.
const (
	destroyLeave    = "leave"
	destroyTruncate = "truncate"
)

// NewFile returns the truenas_file resource.
func NewFile() resource.Resource { return &fileResource{} }

type fileResource struct {
	data *engine.ProviderData
}

type fileModel struct {
	ID             types.String `tfsdk:"id"`
	Path           types.String `tfsdk:"path"`
	Content        types.String `tfsdk:"content"`
	ContentSHA256  types.String `tfsdk:"content_sha256"`
	Mode           types.String `tfsdk:"mode"`
	Size           types.Int64  `tfsdk:"size"`
	Mtime          types.Int64  `tfsdk:"mtime"`
	DriftDetection types.String `tfsdk:"drift_detection"`
	OnDestroy      types.String `tfsdk:"on_destroy"`
}

func (r *fileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *fileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Writes a file on the TrueNAS host and keeps its contents converged.\n\n" +
			"This resource converges contents; it does not own the file's lifetime. TrueNAS has no API to " +
			"remove a file, so destroying the resource drops it from state and leaves the file in place " +
			"unless `on_destroy` says otherwise.\n\n" +
			"`content` is write-only: it never enters Terraform state. Only its SHA-256 is stored, which is " +
			"what lets a plan notice that the file changed. For a file whose contents must not be derived " +
			"into state at all, such as an environment file holding secrets, set `drift_detection = \"none\"`.\n\n" +
			"Reading and writing file contents requires an API key with the `FULL_ADMIN` role, because " +
			"`filesystem.get` and `filesystem.put` reach any path on the host and TrueNAS offers no narrower role.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The file's path.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Absolute path of the file on the TrueNAS host.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.RegexMatches(absolutePath, "must be an absolute path"),
				},
			},
			"content": schema.StringAttribute{
				Required:  true,
				WriteOnly: true,
				MarkdownDescription: "Contents to write. Write-only: the value never enters Terraform state, " +
					"only its SHA-256.",
			},
			"content_sha256": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "SHA-256 of the contents, in hex. Computed from `content`, and refreshed " +
					"from the file on the host unless `drift_detection = \"none\"`.",
			},
			"mode": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Permission bits in octal, e.g. `\"0644\"`. Left to TrueNAS when unset.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(octalMode, "must be octal permission bits, e.g. \"0644\""),
				},
			},
			"size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Size of the file on the host, in bytes.",
			},
			"mtime": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Modification time of the file on the host, in seconds since the epoch.",
			},
			"drift_detection": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(driftContent),
				MarkdownDescription: "How refresh decides the file changed on the host.\n\n" +
					"- `content` (default): download the file and hash it. Always correct, and transfers the " +
					"file on every refresh.\n" +
					"- `stat`: compare size and modification time, and download only when they move. Cheaper, " +
					"but it misses an edit that preserves both, such as a restored backup or `cp -p`.\n" +
					"- `none`: never read the file back. State holds nothing derived from the contents. Use it " +
					"for secrets, where even a hash confirms a guess.",
				Validators: []validator.String{
					stringvalidator.OneOf(driftContent, driftStat, driftNone),
				},
			},
			"on_destroy": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(destroyLeave),
				MarkdownDescription: "What destroying the resource does to the file. TrueNAS has no API to " +
					"delete one.\n\n" +
					"- `leave` (default): drop the resource from state and leave the file untouched.\n" +
					"- `truncate`: write an empty file. There is no undo, and an emptied file is quieter than " +
					"a missing one, so set it only where empty is safer than stale.",
				Validators: []validator.String{
					stringvalidator.OneOf(destroyLeave, destroyTruncate),
				},
			},
		},
	}
}

func (r *fileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// files returns the client, which must be able to transfer file contents.
func (r *fileResource) files() (engine.FileClient, error) {
	fc, ok := r.data.Client.(engine.FileClient)
	if !ok {
		return nil, fmt.Errorf("the configured client cannot transfer files")
	}
	return fc, nil
}

var (
	absolutePath = regexp.MustCompile(`^/.+`)
	octalMode    = regexp.MustCompile(`^0?[0-7]{3,4}$`)
)

// ModifyPlan computes the planned content hash from the configuration and refuses changes when the
// provider is read-only.
//
// content is write-only, so it is absent from both plan and state but present in the configuration.
// Hashing it here is what makes an edited file show up as a change without the contents ever being
// written to state.
func (r *fileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.data == nil {
		return
	}
	if req.Plan.Raw.IsNull() { // destroy
		if r.data.ReadOnly {
			var state fileModel
			resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
			if !resp.Diagnostics.HasError() && state.OnDestroy.ValueString() == destroyTruncate {
				resp.Diagnostics.AddError("Provider is read-only",
					fmt.Sprintf("This plan would truncate %s, but the provider has read_only = true.", state.Path.ValueString()))
			}
		}
		return
	}

	var config fileModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.Content.IsUnknown() {
		// The contents are not known until apply, so neither is the hash.
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("content_sha256"), types.StringUnknown())...)
		return
	}
	planned := hashString(config.Content.ValueString())
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("content_sha256"), types.StringValue(planned))...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.State.Raw.IsNull() {
		if r.data.ReadOnly {
			resp.Diagnostics.AddError("Provider is read-only",
				"This plan would create a file, but the provider has read_only = true.")
		}
		return
	}
	var state fileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.ContentSHA256.ValueString() == planned && state.Mode.Equal(config.Mode) {
		return // nothing will be written, so size and mtime keep their values
	}
	if r.data.ReadOnly {
		resp.Diagnostics.AddError("Provider is read-only",
			fmt.Sprintf("This plan would write %s, but the provider has read_only = true.", config.Path.ValueString()))
		return
	}
	// Writing the file changes both, and the prior values Terraform carries into the plan for a
	// computed attribute would not match what the host reports afterwards.
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("size"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("mtime"), types.Int64Unknown())...)
}

func (r *fileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config fileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !r.write(ctx, plan, config.Content.ValueString(), &resp.Diagnostics) {
		return
	}
	plan.ID = plan.Path
	plan.ContentSHA256 = types.StringValue(hashString(config.Content.ValueString()))
	r.stat(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *fileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, config fileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	planned := hashString(config.Content.ValueString())
	// Only transfer when something actually differs; a change to drift_detection or on_destroy is
	// state-only and must not touch the host.
	if state.ContentSHA256.ValueString() != planned || !state.Mode.Equal(plan.Mode) {
		if !r.write(ctx, plan, config.Content.ValueString(), &resp.Diagnostics) {
			return
		}
	}
	plan.ID = plan.Path
	plan.ContentSHA256 = types.StringValue(planned)
	r.stat(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *fileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	priorSize, priorMtime := state.Size, state.Mtime
	if !r.statAllowMissing(ctx, &state, resp, &resp.Diagnostics) || resp.Diagnostics.HasError() {
		return
	}

	switch state.DriftDetection.ValueString() {
	case driftNone:
		// Nothing derived from the contents may enter state, so keep what we last wrote.
	case driftStat:
		if state.Size.Equal(priorSize) && state.Mtime.Equal(priorMtime) {
			break
		}
		r.refreshHash(ctx, &state, &resp.Diagnostics)
	default:
		r.refreshHash(ctx, &state, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	state.ID = state.Path
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *fileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.OnDestroy.ValueString() != destroyTruncate {
		return // TrueNAS cannot remove a file; the resource only leaves state.
	}
	if r.data.ReadOnly {
		resp.Diagnostics.AddError("Provider is read-only",
			fmt.Sprintf("Destroying this resource would truncate %s, but the provider has read_only = true.",
				state.Path.ValueString()))
		return
	}
	r.write(ctx, state, "", &resp.Diagnostics)
}

func (r *fileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("drift_detection"), types.StringValue(driftContent))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("on_destroy"), types.StringValue(destroyLeave))...)
}

// write uploads content to the file's path.
func (r *fileResource) write(ctx context.Context, m fileModel, content string, diags *diag.Diagnostics) bool {
	if r.data.ReadOnly {
		diags.AddError("Provider is read-only", "The provider has read_only = true and will not change TrueNAS.")
		return false
	}
	fc, err := r.files()
	if err != nil {
		diags.AddError("Cannot transfer files", err.Error())
		return false
	}
	mode, err := parseMode(m.Mode)
	if err != nil {
		diags.AddAttributeError(path.Root("mode"), "Invalid mode", err.Error())
		return false
	}
	if err := fc.PutFile(ctx, m.Path.ValueString(), mode, strings.NewReader(content)); err != nil {
		diags.AddError("Cannot write file", err.Error())
		return false
	}
	return true
}

// refreshHash downloads the file and records its hash.
func (r *fileResource) refreshHash(ctx context.Context, m *fileModel, diags *diag.Diagnostics) {
	fc, err := r.files()
	if err != nil {
		diags.AddError("Cannot transfer files", err.Error())
		return
	}
	var buf bytes.Buffer
	if err := fc.GetFile(ctx, m.Path.ValueString(), &buf); err != nil {
		diags.AddError("Cannot read file", err.Error())
		return
	}
	m.ContentSHA256 = types.StringValue(hashString(buf.String()))
}

// stat records the file's size, mtime and mode, reporting an error if it is missing.
func (r *fileResource) stat(ctx context.Context, m *fileModel, diags *diag.Diagnostics) {
	info, err := r.statCall(ctx, m.Path.ValueString())
	if err != nil {
		diags.AddError("Cannot stat file", err.Error())
		return
	}
	applyStat(m, info)
}

// statAllowMissing records stat information, removing the resource from state when the file is gone.
func (r *fileResource) statAllowMissing(ctx context.Context, m *fileModel, resp *resource.ReadResponse, diags *diag.Diagnostics) bool {
	info, err := r.statCall(ctx, m.Path.ValueString())
	if middleware.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return false
	}
	if err != nil {
		diags.AddError("Cannot stat file", err.Error())
		return false
	}
	applyStat(m, info)
	return true
}

type statInfo struct {
	Size  int64 `json:"size"`
	Mtime any   `json:"mtime"`
	Mode  int64 `json:"mode"`
}

func (r *fileResource) statCall(ctx context.Context, p string) (*statInfo, error) {
	raw, err := r.data.Client.Call(ctx, "filesystem.stat", p)
	if err != nil {
		return nil, err
	}
	var info statInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func applyStat(m *fileModel, info *statInfo) {
	m.Size = types.Int64Value(info.Size)
	m.Mtime = types.Int64Value(mtimeSeconds(info.Mtime))
	if !m.Mode.IsNull() {
		// Report the bits TrueNAS actually has, without the file-type bits stat includes.
		m.Mode = types.StringValue(fmt.Sprintf("0%o", info.Mode&0o7777))
	}
}

// mtimeSeconds reads the modification time, which middleware reports either as a number or as a
// {"$date": milliseconds} wrapper.
func mtimeSeconds(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case map[string]any:
		if ms, ok := t["$date"].(float64); ok {
			return int64(ms / 1000)
		}
	}
	return 0
}

func parseMode(m types.String) (*int64, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	v, err := strconv.ParseInt(m.ValueString(), 8, 32)
	if err != nil {
		return nil, fmt.Errorf("%q is not octal permission bits", m.ValueString())
	}
	return &v, nil
}

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
