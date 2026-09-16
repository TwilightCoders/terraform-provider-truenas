// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// idAttr is the Terraform attribute carrying the resource identity.
const idAttr = "id"

// attrType returns the framework type of an attribute's value.
func attrType(a *node) attr.Type {
	switch a.kind {
	case kindString:
		return types.StringType
	case kindInt:
		if a.identity {
			return types.Int64Type
		}
		return types.NumberType
	case kindNumber:
		return types.Float64Type
	case kindBool:
		return types.BoolType
	case kindObject, kindUnion:
		fields := make(map[string]attr.Type, len(a.children))
		for _, c := range a.children {
			fields[c.name] = attrType(c)
		}
		return types.ObjectType{AttrTypes: fields}
	case kindList:
		return types.ListType{ElemType: attrType(a.elem)}
	case kindMap:
		return types.MapType{ElemType: attrType(a.elem)}
	default:
		return jsontypes.NormalizedType{}
	}
}

// containsWriteOnly reports whether any attribute nested under a is write-only.
func containsWriteOnly(a *node) bool {
	for _, c := range a.children {
		if c.writeOnly || containsWriteOnly(c) {
			return true
		}
	}
	return false
}

func markComputed(a *node) {
	a.role = roleComputed
	a.hasDefault, a.def = false, nil
	a.replace = false
	for _, c := range a.children {
		markComputed(c)
	}
	if a.elem != nil {
		markComputed(a.elem)
	}
}

func appendDesc(desc, extra string) string {
	if desc == "" {
		return extra
	}
	return desc + " " + extra
}

// nilOf is a typed null for an attribute, which is what setting a planned value to "nothing"
// requires: the framework needs the type as well as the absence.
func nilOf(a *node) attr.Value {
	t := attrType(a)
	v, err := t.ValueFromTerraform(context.Background(), tftypes.NewValue(t.TerraformType(context.Background()), nil))
	if err != nil {
		return nil
	}
	return v
}
