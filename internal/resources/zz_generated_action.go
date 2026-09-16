// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Action declares a Terraform action that calls one middleware method, such as running a
// replication now. The action's configuration is derived from the method's arguments.

// methodAction invokes one API method. The generator resolved which method and which arguments;
// this only calls it.
type methodAction struct {
	typeName string
	method   string
	args     []*node
	schema   actionschema.Schema
	data     *ProviderData
}

func (a *methodAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = a.schema
}

func (a *methodAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + a.typeName
}

func (a *methodAction) Configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	data, ok := req.ProviderData.(*ProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *engine.ProviderData, got %T", req.ProviderData))
		return
	}
	a.data = data
}

func (a *methodAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {

	params, err := a.params(req.Config.Raw)
	if err != nil {
		resp.Diagnostics.AddError("Invalid action configuration", err.Error())
		return
	}
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Calling %s", a.method)})
	}
	if _, err := a.data.Client.Call(ctx, a.method, params...); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("%s failed", a.typeName), err.Error())
		return
	}
	if resp.SendProgress != nil {
		resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("%s finished", a.method)})
	}
}

// params encodes configured arguments positionally. Unset trailing arguments are left out; an
// unset argument before a set one is sent as its default (or null).
func (a *methodAction) params(config tftypes.Value) ([]any, error) {
	vals := children(config)
	params := make([]any, len(a.args))
	last := -1
	for i, arg := range a.args {
		enc, ok, err := encodeValue(arg, vals[arg.name], false)
		if err != nil {
			return nil, err
		}
		if ok {
			params[i], last = enc, i
		} else if arg.hasDefault {
			params[i] = arg.def
		}
	}
	return params[:last+1], nil
}
