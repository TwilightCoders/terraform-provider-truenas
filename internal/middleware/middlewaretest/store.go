package middlewaretest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/resources"
)

// Store is an in-memory CRUD service that validates requests against the API schema, the way
// middlewared's Pydantic models do: unknown fields, missing required fields and wrong JSON types
// are rejected, and defaults are applied.
type Store struct {
	namespace string
	argName   string
	shape     resources.Shape
	creatable bool
	updatable bool
	pk        string
	intPK     bool

	mu     sync.Mutex
	rows   map[string]map[string]any
	order  []string
	nextID int64

	// OnWrite, when set, runs after every create and update with the stored row, so tests can
	// fill read-only fields or simulate server-side normalization.
	OnWrite func(row map[string]any)
}

// ServeCRUD registers create, get_instance, update, delete and query for a resource type, backed
// by a Store shaped by what the provider believes that resource looks like.
//
// The shape comes from the provider's own descriptors rather than from a copy of the API's schema:
// a fake only ever receives calls the provider is capable of making, so the provider's belief is
// the right thing to hold it to, and it keeps a vendor's schema document out of this repository.
// More than one resource type may be given when they share a namespace: a dataset and a zvol are
// both pool.dataset, and one fake has to accept the fields of either.
func (s *Server) ServeCRUD(resourceTypes ...string) *Store {
	if len(resourceTypes) == 0 {
		panic("middlewaretest: ServeCRUD needs a resource type")
	}
	shape, ok := resources.ShapeFor(resourceTypes[0])
	if !ok {
		panic("middlewaretest: unknown resource type " + resourceTypes[0])
	}
	for _, extra := range resourceTypes[1:] {
		other, ok := resources.ShapeFor(extra)
		if !ok {
			panic("middlewaretest: unknown resource type " + extra)
		}
		if other.Namespace != shape.Namespace {
			panic("middlewaretest: " + extra + " is not in " + shape.Namespace)
		}
		shape.Fields = mergeFields(shape.Fields, other.Fields)
	}
	st := &Store{
		namespace: shape.Namespace,
		argName:   strings.ReplaceAll(shape.Namespace, ".", "_"),
		shape:     shape,
		creatable: !shape.Singleton,
		updatable: true,
		pk:        shape.PrimaryKey,
		intPK:     shape.IntPK,
		rows:      map[string]map[string]any{},
		nextID:    1,
	}
	if st.creatable {
		s.Handle(shape.Namespace+".create", st.handleCreate)
		s.Handle(shape.Namespace+".delete", st.handleDelete)
	}
	s.Handle(shape.Namespace+".get_instance", st.handleGet)
	s.Handle(shape.Namespace+".query", st.handleQuery)
	s.Handle(shape.Namespace+".update", st.handleUpdate)
	return st
}

// Rows returns copies of every stored row in creation order.
func (st *Store) Rows() []map[string]any {
	st.mu.Lock()
	defer st.mu.Unlock()
	out := make([]map[string]any, 0, len(st.order))
	for _, key := range st.order {
		out = append(out, deepCopy(st.rows[key]).(map[string]any))
	}
	return out
}

// Seed stores a row as if it already existed, keyed by its primary key. Like a real read, the
// row gets schema defaults and zero values for fields it leaves out.
func (st *Store) Seed(row map[string]any) {
	st.mu.Lock()
	defer st.mu.Unlock()
	stored := deepCopy(row).(map[string]any)
	if st.creatable {
		var ignored []middleware.FieldError
		checkObject("seed", st.shape.Fields, stored, true, &ignored)
	}
	fillReadOnly(st.shape.Fields, stored)
	key := fmt.Sprint(stored[st.pk])
	st.rows[key] = stored
	st.order = append(st.order, key)
}

// Mutate changes a stored row in place, simulating a change made outside Terraform.
func (st *Store) Mutate(id any, fn func(row map[string]any)) {
	st.mu.Lock()
	defer st.mu.Unlock()
	row, ok := st.rows[fmt.Sprint(id)]
	if !ok {
		panic(fmt.Sprintf("middlewaretest: no row %v", id))
	}
	fn(row)
}

// Remove deletes a stored row, simulating deletion outside Terraform.
func (st *Store) Remove(id any) {
	st.mu.Lock()
	defer st.mu.Unlock()
	key := fmt.Sprint(id)
	delete(st.rows, key)
	st.order = slices.DeleteFunc(st.order, func(k string) bool { return k == key })
}

func (st *Store) handleCreate(_ context.Context, params []json.RawMessage) (any, error) {
	data, err := st.decodeArg(params, 0)
	if err != nil {
		return nil, err
	}
	var fields []middleware.FieldError
	create := st.shape.Fields
	if st.shape.Discriminator != "" {
		// The provider always sends the discriminator; it is fixed by the variant, not configured.
		create = append(append([]resources.ShapeField{}, create...),
			resources.ShapeField{Name: st.shape.Discriminator, Kind: "string"})
	}
	checkObject(st.argName+"_create", create, data, true, &fields)
	if len(fields) > 0 {
		return nil, validation(fields)
	}

	st.mu.Lock()
	defer st.mu.Unlock()
	var key string
	if st.intPK {
		data[st.pk] = st.nextID
		key = fmt.Sprint(st.nextID)
		st.nextID++
	} else {
		name, _ := data["name"].(string)
		if name == "" {
			return nil, validation([]middleware.FieldError{{Attribute: st.argName + "_create.name", Message: "Field required", Errno: 22}})
		}
		data[st.pk] = name
		key = name
	}
	if _, exists := st.rows[key]; exists {
		return nil, &middleware.Error{Errno: 17, Errname: "EEXIST", Reason: key + " already exists"}
	}
	fillReadOnly(st.shape.Fields, data)
	if st.OnWrite != nil {
		st.OnWrite(data)
	}
	st.rows[key] = data
	st.order = append(st.order, key)
	return deepCopy(data), nil
}

func (st *Store) handleUpdate(_ context.Context, params []json.RawMessage) (any, error) {
	key, err := st.decodeID(params)
	if err != nil {
		return nil, err
	}
	patch, err := st.decodeArg(params, 1)
	if err != nil {
		return nil, err
	}
	var fields []middleware.FieldError
	checkObject(st.argName+"_update", updatableFields(st.shape.Fields), patch, false, &fields)
	if len(fields) > 0 {
		return nil, validation(fields)
	}

	st.mu.Lock()
	defer st.mu.Unlock()
	row, ok := st.rows[key]
	if !ok {
		return nil, notFound(st.namespace, key)
	}
	for k, v := range patch {
		row[k] = v
	}
	if st.OnWrite != nil {
		st.OnWrite(row)
	}
	return deepCopy(row), nil
}

func (st *Store) handleGet(_ context.Context, params []json.RawMessage) (any, error) {
	key, err := st.decodeID(params)
	if err != nil {
		return nil, err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	row, ok := st.rows[key]
	if !ok {
		return nil, notFound(st.namespace, key)
	}
	return deepCopy(row), nil
}

func (st *Store) handleDelete(_ context.Context, params []json.RawMessage) (any, error) {
	key, err := st.decodeID(params)
	if err != nil {
		return nil, err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.rows[key]; !ok {
		return nil, notFound(st.namespace, key)
	}
	delete(st.rows, key)
	st.order = slices.DeleteFunc(st.order, func(k string) bool { return k == key })
	return true, nil
}

// handleQuery supports "=" and "!=" filters on top-level fields and the limit option.
func (st *Store) handleQuery(_ context.Context, params []json.RawMessage) (any, error) {
	var filters [][]any
	if len(params) > 0 {
		v, err := decode(params[0])
		if err != nil {
			return nil, err
		}
		list, _ := v.([]any)
		for _, f := range list {
			triple, ok := f.([]any)
			if !ok || len(triple) != 3 {
				return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: fmt.Sprintf("unsupported filter %v", f)}
			}
			filters = append(filters, triple)
		}
	}
	limit := 0
	if len(params) > 1 {
		var opts struct {
			Limit int `json:"limit"`
		}
		_ = json.Unmarshal(params[1], &opts)
		limit = opts.Limit
	}

	out := []map[string]any{}
	for _, row := range st.Rows() {
		if matches(row, filters) {
			out = append(out, row)
		}
		if limit > 0 && len(out) == limit {
			break
		}
	}
	return out, nil
}

func matches(row map[string]any, filters [][]any) bool {
	for _, f := range filters {
		field, _ := f[0].(string)
		op, _ := f[1].(string)
		equal := fmt.Sprint(row[field]) == fmt.Sprint(f[2])
		if (op == "=" && !equal) || (op == "!=" && equal) {
			return false
		}
	}
	return true
}

func (st *Store) decodeArg(params []json.RawMessage, index int) (map[string]any, error) {
	if len(params) <= index {
		return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: fmt.Sprintf("missing argument %d", index)}
	}
	v, err := decode(params[index])
	obj, ok := v.(map[string]any)
	if err != nil || !ok {
		return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: fmt.Sprintf("argument %d must be an object", index)}
	}
	return obj, nil
}

func (st *Store) decodeID(params []json.RawMessage) (string, error) {
	if len(params) == 0 {
		return "", &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: "missing id"}
	}
	v, err := decode(params[0])
	if err != nil {
		return "", err
	}
	if st.intPK {
		if _, ok := v.(json.Number); !ok {
			return "", &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: fmt.Sprintf("id must be an integer, got %T", v)}
		}
	}
	return fmt.Sprint(v), nil
}

func decode(raw json.RawMessage) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	err := dec.Decode(&v)
	return v, err
}

func notFound(namespace, key string) error {
	return &middleware.Error{Errno: 2, Errname: "ENOENT", Reason: fmt.Sprintf("%s %s does not exist", namespace, key)}
}

func validation(fields []middleware.FieldError) error {
	return &middleware.Error{Errno: 22, Errname: "EINVAL", Fields: fields}
}

// checkObject validates obj against fields and, when applyDefaults is set, fills missing defaults.
func checkObject(prefix string, fields []resources.ShapeField, obj map[string]any, applyDefaults bool, errs *[]middleware.FieldError) {
	known := map[string]resources.ShapeField{}
	for _, f := range fields {
		known[f.Name] = f
	}
	for key := range obj {
		if _, ok := known[key]; !ok {
			*errs = append(*errs, middleware.FieldError{Attribute: prefix + "." + key, Message: "Extra inputs are not permitted", Errno: 22})
		}
	}
	for _, f := range fields {
		path := prefix + "." + f.Name
		v, present := obj[f.Name]
		if !present {
			switch {
			case applyDefaults && f.HasDefault:
				obj[f.Name] = deepCopy(f.Default)
			case applyDefaults && !f.Nullable && objectDefaults(f) != nil:
				obj[f.Name] = objectDefaults(f)
			case applyDefaults && f.Required:
				*errs = append(*errs, middleware.FieldError{Attribute: path, Message: "Field required", Errno: 22})
			}
			continue
		}
		checkValue(path, f, v, errs)
	}
}

// mergeFields adds fields the first shape lacks, so a fake serving a shared namespace accepts
// what either resource sends.
func mergeFields(into, from []resources.ShapeField) []resources.ShapeField {
	have := map[string]bool{}
	for _, f := range into {
		have[f.Name] = true
	}
	for _, f := range from {
		if !have[f.Name] {
			into = append(into, f)
		}
	}
	return into
}

// updatableFields are those an update accepts.
func updatableFields(fields []resources.ShapeField) []resources.ShapeField {
	var out []resources.ShapeField
	for _, f := range fields {
		if f.Updatable {
			out = append(out, f)
		}
	}
	return out
}

func checkValue(path string, f resources.ShapeField, v any, errs *[]middleware.FieldError) {
	if v == nil {
		if !f.Nullable && f.Kind != "any" {
			*errs = append(*errs, middleware.FieldError{Attribute: path, Message: "Input should not be None", Errno: 22})
		}
		return
	}
	bad := func(want string) {
		*errs = append(*errs, middleware.FieldError{Attribute: path, Message: fmt.Sprintf("Input should be a valid %s", want), Errno: 22})
	}
	switch f.Kind {
	case "string":
		if _, isNumber := v.(json.Number); isNumber && f.IntOrString {
			return
		}
		if _, ok := v.(string); !ok {
			bad("string")
		}
	case "int":
		n, ok := v.(json.Number)
		// Wider than int64 is still an integer: TrueNAS reports certificate serials that way.
		if _, isInt := new(big.Int).SetString(n.String(), 10); !ok || !isInt {
			bad("integer")
		}
	case "number":
		if _, ok := v.(json.Number); !ok {
			bad("number")
		}
	case "bool":
		if _, ok := v.(bool); !ok {
			bad("boolean")
		}
	case "object":
		obj, ok := v.(map[string]any)
		if !ok {
			bad("dictionary")
			return
		}
		// Nested objects are whole models even inside partial updates, so their defaults apply.
		checkObject(path, f.Children, obj, true, errs)
	case "union":
		obj, ok := v.(map[string]any)
		if !ok {
			bad("dictionary")
			return
		}
		// Children are variants here. The value picks one by the discriminator.
		key, _ := obj[f.Discriminator].(string)
		for _, variant := range f.Children {
			if variant.Name == key {
				vf := append(append([]resources.ShapeField{}, variant.Children...),
					resources.ShapeField{Name: f.Discriminator, Kind: "string"})
				checkObject(path, vf, obj, true, errs)
				return
			}
		}
		*errs = append(*errs, middleware.FieldError{
			Attribute: path + "." + f.Discriminator,
			Message:   "Input tag does not match any expected tags",
			Errno:     22,
		})
	case "list":
		items, ok := v.([]any)
		if !ok {
			bad("list")
			return
		}
		if f.Elem == nil {
			return
		}
		for i, item := range items {
			checkValue(fmt.Sprintf("%s.%d", path, i), *f.Elem, item, errs)
		}
	case "map":
		if _, ok := v.(map[string]any); !ok {
			bad("dictionary")
		}
	}
}

// objectDefaults returns the default object every-field-has-a-default implies, or nil.
func objectDefaults(f resources.ShapeField) map[string]any {
	if f.Kind != "object" || len(f.Children) == 0 {
		return nil
	}
	out := map[string]any{}
	for _, c := range f.Children {
		if !c.HasDefault {
			return nil
		}
		out[c.Name] = deepCopy(c.Default)
	}
	return out
}

// fillReadOnly adds zero values for readable fields the stored row lacks, so a read answers with
// the same fields the provider expects to find.
func fillReadOnly(fields []resources.ShapeField, row map[string]any) {
	for _, f := range fields {
		if !f.Readable {
			continue
		}
		if _, ok := row[f.Name]; !ok {
			row[f.Name] = zero(f)
		}
	}
}

func zero(f resources.ShapeField) any {
	// A wrapped value is reported as an object, not as the value, and a row that never set one
	// has no wrapper at all. Tests rebuild the wrapper from nil.
	if f.Property || f.Nullable || f.Kind == "any" || f.Kind == "union" {
		return nil
	}
	switch f.Kind {
	case "string":
		return ""
	case "int", "number":
		return json.Number("0")
	case "bool":
		return false
	case "list":
		return []any{}
	case "map":
		return map[string]any{}
	default:
		obj := map[string]any{}
		fillReadOnly(f.Children, obj)
		return obj
	}
}

func deepCopy(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			out[k] = deepCopy(val)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = deepCopy(val)
		}
		return out
	default:
		return v
	}
}

// ConfigStore is an in-memory config service: one settings object read with <ns>.config and
// changed with <ns>.update.
type ConfigStore struct {
	namespace string
	fields    []resources.ShapeField

	mu   sync.Mutex
	data map[string]any
}

// ServeConfig registers <namespace>.config and <namespace>.update, starting from initial.
func (s *Server) ServeConfig(resourceType string, initial map[string]any) *ConfigStore {
	shape, ok := resources.ShapeFor(resourceType)
	if !ok {
		panic("middlewaretest: unknown resource type " + resourceType)
	}
	namespace := shape.Namespace
	cs := &ConfigStore{
		namespace: namespace,
		fields:    shape.Fields,
		data:      deepCopy(initial).(map[string]any),
	}
	fillReadOnly(shape.Fields, cs.data)
	s.Handle(namespace+".config", func(context.Context, []json.RawMessage) (any, error) {
		return cs.Data(), nil
	})
	s.Handle(namespace+".update", func(_ context.Context, params []json.RawMessage) (any, error) {
		if len(params) != 1 {
			return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: "update takes one argument"}
		}
		v, err := decode(params[0])
		patch, ok := v.(map[string]any)
		if err != nil || !ok {
			return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errno: 22, Errname: "EINVAL", Reason: "argument must be an object"}
		}
		var fields []middleware.FieldError
		checkObject(strings.ReplaceAll(namespace, ".", "_")+"_update", updatableFields(cs.fields), patch, false, &fields)
		if len(fields) > 0 {
			return nil, validation(fields)
		}
		cs.mu.Lock()
		for k, val := range patch {
			cs.data[k] = val
		}
		cs.mu.Unlock()
		return cs.Data(), nil
	})
	return cs
}

// Data returns a copy of the current settings.
func (cs *ConfigStore) Data() map[string]any {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return deepCopy(cs.data).(map[string]any)
}
