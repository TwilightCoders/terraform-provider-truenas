// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// addCallError reports a failed middleware call, attaching validation failures to the attributes
// they concern.
func (m *model) addCallError(c Client, diags *diag.Diagnostics, op string, err error) {
	summary := fmt.Sprintf("Unable to %s %s", op, m.typeName)
	fields := c.ValidationFields(err)
	if len(fields) == 0 {
		diags.AddError(summary, err.Error())
		return
	}
	for _, f := range fields {
		detail := f.Message
		if p, ok := m.attributePath(f.Attribute); ok {
			diags.AddAttributeError(p, summary, detail)
			continue
		}
		diags.AddError(summary, fmt.Sprintf("%s: %s", f.Attribute, detail))
	}
}

// attributePath maps a middleware validation attribute such as "cronjob_create.schedule.minute"
// to the Terraform attribute path. It stops at unions and at the deepest attribute it can follow.
func (m *model) attributePath(apiAttr string) (path.Path, bool) {
	segments := strings.Split(apiAttr, ".")
	if len(segments) < 2 {
		return path.Empty(), false
	}
	segments = segments[1:] // method argument name

	attrs := m.attrs
	var p path.Path
	var current *node
	for i, seg := range segments {
		if current != nil && current.kind == kindList {
			idx, err := strconv.Atoi(seg)
			if err != nil {
				return p, true
			}
			p = p.AtListIndex(idx)
			current = current.elem
			attrs = current.children
			continue
		}
		next := findByAPI(attrs, seg)
		if next == nil {
			return p, i > 0
		}
		if i == 0 {
			p = path.Root(next.name)
		} else {
			p = p.AtName(next.name)
		}
		current = next
		if next.kind == kindUnion {
			return p, true
		}
		attrs = next.children
	}
	return p, true
}

func findByAPI(attrs []*node, api string) *node {
	for _, a := range attrs {
		if a.api == api {
			return a
		}
	}
	return nil
}
