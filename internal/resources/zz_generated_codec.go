// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// redacted is how the middleware masks secrets for callers without full admin rights.
const redacted = "********"

func tfType(a *node) tftypes.Type {
	return attrType(a).TerraformType(context.Background())
}

func nullOf(_ context.Context, a *node) tftypes.Value {
	return tftypes.NewValue(tfType(a), nil)
}

func objectType(attrs []*node) tftypes.Object {
	fields := make(map[string]tftypes.Type, len(attrs))
	for _, a := range attrs {
		fields[a.name] = tfType(a)
	}
	return tftypes.Object{AttributeTypes: fields}
}

// children returns a copy of an object value's attributes; null or unknown objects yield an
// empty map. tftypes hands out its internal map, so callers must never write to it directly.
func children(v tftypes.Value) map[string]tftypes.Value {
	out := map[string]tftypes.Value{}
	if !v.IsKnown() || v.IsNull() {
		return out
	}
	var inner map[string]tftypes.Value
	_ = v.As(&inner)
	for k, val := range inner {
		out[k] = val
	}
	return out
}

func elements(v tftypes.Value) []tftypes.Value {
	var out []tftypes.Value
	if !v.IsKnown() || v.IsNull() {
		return nil
	}
	_ = v.As(&out)
	return out
}

func mapElements(v tftypes.Value) map[string]tftypes.Value {
	out := map[string]tftypes.Value{}
	if !v.IsKnown() || v.IsNull() {
		return out
	}
	_ = v.As(&out)
	return out
}

// encodeValue converts a Terraform value to middleware JSON. ok is false when the value should be
// left out of the request: unknown values, and nulls unless sendNull is set on a nullable field.
func encodeValue(a *node, v tftypes.Value, sendNull bool) (out any, ok bool, err error) {
	if !v.IsKnown() {
		return nil, false, nil
	}
	if v.IsNull() {
		return nil, sendNull && a.nullable, nil
	}

	switch a.kind {
	case kindString:
		var s string
		if err = v.As(&s); err != nil {
			return nil, false, err
		}
		if a.intOrString {
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return i, true, nil
			}
		}
		return s, true, nil
	case kindInt:
		var f big.Float
		if err = v.As(&f); err != nil {
			return nil, false, err
		}
		if !f.IsInt() {
			return nil, false, fmt.Errorf("%s: %s is not an integer", a.path, f.Text('f', -1))
		}
		if i, acc := f.Int64(); acc == big.Exact {
			return i, true, nil
		}
		// Wider than int64: TrueNAS reports X.509 serials as JSON integers of up to 20 octets.
		// json.Number keeps every digit on the way back out.
		return json.Number(f.Text('f', -1)), true, nil
	case kindNumber:
		var f big.Float
		if err = v.As(&f); err != nil {
			return nil, false, err
		}
		n, _ := f.Float64()
		return n, true, nil
	case kindBool:
		var b bool
		err = v.As(&b)
		return b, err == nil, err
	case kindAny:
		var s string
		if err = v.As(&s); err != nil {
			return nil, false, err
		}
		dec := json.NewDecoder(bytes.NewReader([]byte(s)))
		dec.UseNumber()
		var decoded any
		if err = dec.Decode(&decoded); err != nil {
			return nil, false, fmt.Errorf("%s: invalid JSON: %w", a.path, err)
		}
		return decoded, true, nil
	case kindObject:
		obj, err := encodeObject(a.children, v, sendNull)
		return obj, err == nil, err
	case kindUnion:
		for name, cv := range children(v) {
			if cv.IsNull() || !cv.IsKnown() {
				continue
			}
			variant := a.child(name)
			obj, err := encodeObject(variant.children, cv, sendNull)
			if err != nil {
				return nil, false, err
			}
			obj[a.discriminator] = variant.api
			return obj, true, nil
		}
		return nil, sendNull && a.nullable, nil
	case kindList:
		elems := elements(v)
		list := make([]any, 0, len(elems))
		for _, ev := range elems {
			item, _, err := encodeValue(a.elem, ev, true)
			if err != nil {
				return nil, false, err
			}
			list = append(list, item)
		}
		return list, true, nil
	case kindMap:
		out := map[string]any{}
		for k, ev := range mapElements(v) {
			item, _, err := encodeValue(a.elem, ev, true)
			if err != nil {
				return nil, false, err
			}
			out[k] = item
		}
		return out, true, nil
	}
	return nil, false, fmt.Errorf("%s: cannot encode kind %s", a.path, a.kind)
}

func encodeObject(attrs []*node, v tftypes.Value, sendNull bool) (map[string]any, error) {
	vals := children(v)
	out := map[string]any{}
	for _, c := range attrs {
		if c.role == roleComputed {
			continue
		}
		enc, ok, err := encodeValue(c, vals[c.name], sendNull)
		if err != nil {
			return nil, err
		}
		if ok {
			out[c.api] = enc
		}
	}
	return out, nil
}

// decodeValue converts middleware JSON to a Terraform value. prior supplies values the API cannot
// return: redacted secrets and fields missing from the read shape.
func decodeValue(a *node, api any, prior tftypes.Value) (tftypes.Value, error) {
	typ := tfType(a)
	if ref, ok := api.(map[string]any); ok && a.ref != "" && a.kind != kindList {
		api = ref[a.ref]
	}
	if a.property {
		api = unwrapProperty(a, api)
	}
	if a.writeOnly || api == nil {
		return tftypes.NewValue(typ, nil), nil
	}
	if a.sensitive && api == redacted && prior.IsKnown() {
		return prior, nil
	}

	switch a.kind {
	case kindString:
		s, ok := api.(string)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		// The server rewrites some values on write; keep the configured spelling when it means
		// the same thing, so the rewrite is not reported as drift on every plan.
		if a.canonical != nil && prior.IsKnown() && !prior.IsNull() {
			var p string
			if prior.As(&p) == nil && a.canonical(p) == s {
				return prior, nil
			}
		}
		return tftypes.NewValue(typ, s), nil
	case kindInt, kindNumber:
		f, err := bigFloat(api)
		if err != nil {
			return tftypes.Value{}, mismatch(a, api)
		}
		if a.kind == kindInt && !f.IsInt() {
			return tftypes.Value{}, mismatch(a, api)
		}
		return tftypes.NewValue(typ, f), nil
	case kindBool:
		b, ok := api.(bool)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		return tftypes.NewValue(typ, b), nil
	case kindAny:
		b, err := json.Marshal(api)
		if err != nil {
			return tftypes.Value{}, err
		}
		if equalJSON(prior, b) {
			return prior, nil
		}
		return tftypes.NewValue(typ, string(b)), nil
	case kindObject:
		m, ok := api.(map[string]any)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		return decodeObject(typ, a.children, m, prior)
	case kindUnion:
		m, ok := api.(map[string]any)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		key, _ := m[a.discriminator].(string)
		priorVariants := children(prior)
		vals := map[string]tftypes.Value{}
		matched := false
		for _, variant := range a.children {
			if variant.api != key {
				vals[variant.name] = tftypes.NewValue(tfType(variant), nil)
				continue
			}
			v, err := decodeObject(tfType(variant), variant.children, m, priorOrNull(priorVariants, variant))
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[variant.name], matched = v, true
		}
		if !matched {
			return tftypes.Value{}, fmt.Errorf("%s: unknown %s %q", a.path, a.discriminator, key)
		}
		return tftypes.NewValue(typ, vals), nil
	case kindList:
		items, ok := api.([]any)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		priorItems := elements(prior)
		vals := make([]tftypes.Value, 0, len(items))
		for i, item := range items {
			p := tftypes.NewValue(tfType(a.elem), nil)
			if i < len(priorItems) {
				p = priorItems[i]
			}
			v, err := decodeValue(a.elem, item, p)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals = append(vals, v)
		}
		return tftypes.NewValue(typ, vals), nil
	case kindMap:
		items, ok := api.(map[string]any)
		if !ok {
			return tftypes.Value{}, mismatch(a, api)
		}
		priorItems := mapElements(prior)
		vals := make(map[string]tftypes.Value, len(items))
		for k, item := range items {
			p, has := priorItems[k]
			if !has {
				p = tftypes.NewValue(tfType(a.elem), nil)
			}
			v, err := decodeValue(a.elem, item, p)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[k] = v
		}
		return tftypes.NewValue(typ, vals), nil
	}
	return tftypes.Value{}, fmt.Errorf("%s: cannot decode kind %s", a.path, a.kind)
}

func decodeObject(typ tftypes.Type, attrs []*node, m map[string]any, prior tftypes.Value) (tftypes.Value, error) {
	priorVals := children(prior)
	vals := make(map[string]tftypes.Value, len(attrs))
	for _, c := range attrs {
		p := priorOrNull(priorVals, c)
		if !c.readable {
			vals[c.name] = p
			continue
		}
		raw := m[c.api]
		if c.kind == kindUnion {
			raw = withDiscriminator(raw, c.discriminator, m)
		}
		v, err := decodeValue(c, raw, p)
		if err != nil {
			return tftypes.Value{}, err
		}
		vals[c.name] = v
	}
	return tftypes.NewValue(typ, vals), nil
}

// withDiscriminator copies the discriminator from the parent object into a union value that
// omits it; sharing.smb returns options without the purpose that selects them.
func withDiscriminator(raw any, discriminator string, parent map[string]any) any {
	obj, ok := raw.(map[string]any)
	if !ok {
		return raw
	}
	if _, has := obj[discriminator]; has {
		return raw
	}
	key, ok := parent[discriminator]
	if !ok {
		return raw
	}
	out := make(map[string]any, len(obj)+1)
	for k, v := range obj {
		out[k] = v
	}
	out[discriminator] = key
	return out
}

// priorOrNull returns the known prior value of a, or null.
func priorOrNull(vals map[string]tftypes.Value, a *node) tftypes.Value {
	if v, ok := vals[a.name]; ok && v.IsFullyKnown() {
		return v
	}
	return tftypes.NewValue(tfType(a), nil)
}

// reconcile returns the value to store after an apply. Terraform rejects state that contradicts a
// known planned value, so planned values win; the API's value fills in anything unknown. A real
// difference then surfaces as drift on the next refresh.
func reconcile(a *node, planned, applied tftypes.Value) tftypes.Value {
	if !planned.IsKnown() {
		return applied
	}
	if planned.IsFullyKnown() {
		return planned
	}
	switch a.kind {
	case kindObject, kindUnion:
		if planned.IsNull() || applied.IsNull() {
			return applied
		}
		p, ap := children(planned), children(applied)
		vals := make(map[string]tftypes.Value, len(a.children))
		for _, c := range a.children {
			vals[c.name] = reconcile(c, p[c.name], ap[c.name])
		}
		return tftypes.NewValue(tfType(a), vals)
	case kindList:
		p, ap := elements(planned), elements(applied)
		if len(p) != len(ap) {
			return applied
		}
		vals := make([]tftypes.Value, len(p))
		for i := range p {
			vals[i] = reconcile(a.elem, p[i], ap[i])
		}
		return tftypes.NewValue(tfType(a), vals)
	}
	return applied
}

// unwrapProperty turns a ZFS property wrapper into the plain value create accepts.
func unwrapProperty(a *node, api any) any {
	f := a.dialect.property
	w, ok := api.(map[string]any)
	if !ok {
		if api == nil && a.inherit {
			// A field backed by a user property is omitted entirely when it was never set.
			return a.dialect.inherit
		}
		return api
	}
	if src, _ := w[f.source].(string); a.inherit && a.dialect.inheritedSource(src) {
		return a.dialect.inherit
	}
	switch {
	case a.rawProperty:
		return w[f.raw]
	case a.kind == kindInt || a.kind == kindNumber || a.kind == kindBool:
		return w[f.parsed]
	case a.intOrString:
		if _, isString := w[f.parsed].(string); !isString && w[f.parsed] != nil {
			return fmt.Sprint(w[f.parsed])
		}
		return w[f.value]
	default:
		return w[f.value]
	}
}

func bigFloat(v any) (*big.Float, error) {
	switch n := v.(type) {
	case json.Number:
		// Integers are converted exactly; big.Float's default 64-bit mantissa would silently
		// round values such as a 140-bit certificate serial.
		if i, ok := new(big.Int).SetString(n.String(), 10); ok {
			return new(big.Float).SetInt(i), nil
		}
		f, _, err := big.ParseFloat(n.String(), 10, 512, big.ToNearestEven)
		return f, err
	case float64:
		return big.NewFloat(n), nil
	case int64:
		return new(big.Float).SetInt64(n), nil
	case int:
		return new(big.Float).SetInt64(int64(n)), nil
	}
	return nil, fmt.Errorf("%v is not a number", v)
}

func equalJSON(prior tftypes.Value, current []byte) bool {
	if !prior.IsKnown() || prior.IsNull() {
		return false
	}
	var s string
	if prior.As(&s) != nil {
		return false
	}
	var a, b any
	if json.Unmarshal([]byte(s), &a) != nil || json.Unmarshal(current, &b) != nil {
		return false
	}
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	return bytes.Equal(ja, jb)
}

func mismatch(a *node, api any) error {
	return fmt.Errorf("%s: API returned %T, expected %s", a.path, api, a.kind)
}

// withWriteOnly fills write-only attributes nested in a planned value from configuration, where
// Terraform keeps them.
func withWriteOnly(a *node, planned, configured tftypes.Value) tftypes.Value {
	if a.writeOnly {
		return configured
	}
	if (a.kind != kindObject && a.kind != kindUnion) || !planned.IsKnown() || planned.IsNull() || !containsWriteOnly(a) {
		return planned
	}
	p, c := children(planned), children(configured)
	vals := make(map[string]tftypes.Value, len(a.children))
	for _, child := range a.children {
		cv, ok := c[child.name]
		if !ok {
			cv = tftypes.NewValue(tfType(child), nil)
		}
		vals[child.name] = withWriteOnly(child, p[child.name], cv)
	}
	return tftypes.NewValue(tfType(a), vals)
}
