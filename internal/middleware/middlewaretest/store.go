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

	"github.com/TwilightCoders/terraform-provider-truenas/internal/apischema"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

// Store is an in-memory CRUD service that validates requests against the API schema, the way
// middlewared's Pydantic models do: unknown fields, missing required fields and wrong JSON types
// are rejected, and defaults are applied.
type Store struct {
	namespace string
	argName   string
	create    *apischema.Type
	update    *apischema.Type
	read      *apischema.Type
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

// ServeCRUD registers create, get_instance, update, delete and query for namespace, backed by
// a Store shaped by snap.
func (s *Server) ServeCRUD(snap *apischema.Snapshot, namespace string) *Store {
	svc, ok := snap.Service(namespace)
	if !ok {
		panic("middlewaretest: unknown service " + namespace)
	}
	st := &Store{
		namespace: namespace,
		argName:   strings.ReplaceAll(namespace, ".", "_"),
		pk:        svc.PrimaryKey,
		intPK:     svc.PrimaryKeyType == "integer",
		rows:      map[string]map[string]any{},
		nextID:    1,
	}
	if snap.HasMethod(namespace + ".create") {
		st.create = mustMethod(snap, namespace+".create").Accepts[0].Type
		s.Handle(namespace+".create", st.handleCreate)
	}
	st.read = mustMethod(snap, namespace+".get_instance").Returns
	if snap.HasMethod(namespace + ".update") {
		st.update = mustMethod(snap, namespace+".update").Accepts[1].Type
	}

	s.Handle(namespace+".get_instance", st.handleGet)
	if snap.HasMethod(namespace + ".delete") {
		s.Handle(namespace+".delete", st.handleDelete)
	}
	s.Handle(namespace+".query", st.handleQuery)
	if st.update != nil {
		s.Handle(namespace+".update", st.handleUpdate)
	}
	return st
}

func mustMethod(snap *apischema.Snapshot, name string) *apischema.Method {
	m, err := snap.Method(name)
	if err != nil {
		panic(err)
	}
	return m
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
	if st.create != nil && st.create.Kind == apischema.KindObject {
		var ignored []middleware.FieldError
		checkObject("seed", st.create, stored, true, &ignored)
	}
	fillReadOnly(st.read, stored)
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
	if st.create.Kind == apischema.KindUnion {
		checkValue(st.argName+"_create", st.create, data, &fields)
	} else {
		checkObject(st.argName+"_create", st.create, data, true, &fields)
	}
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
	fillReadOnly(st.read, data)
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
	checkObject(st.argName+"_update", st.update, patch, false, &fields)
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

// checkObject validates obj against t and, when applyDefaults is set, fills missing defaults.
func checkObject(prefix string, t *apischema.Type, obj map[string]any, applyDefaults bool, errs *[]middleware.FieldError) {
	for key := range obj {
		if t.Field(key) == nil {
			*errs = append(*errs, middleware.FieldError{Attribute: prefix + "." + key, Message: "Extra inputs are not permitted", Errno: 22})
		}
	}
	for _, f := range t.Fields {
		path := prefix + "." + f.Name
		v, present := obj[f.Name]
		if !present {
			switch {
			case applyDefaults && f.Type.HasDefault:
				obj[f.Name] = deepCopy(f.Type.Default)
			case applyDefaults && !f.Type.Nullable && objectDefaults(f.Type) != nil:
				obj[f.Name] = objectDefaults(f.Type)
			case applyDefaults && f.Required:
				*errs = append(*errs, middleware.FieldError{Attribute: path, Message: "Field required", Errno: 22})
			}
			continue
		}
		checkValue(path, f.Type, v, errs)
	}
}

func checkValue(path string, t *apischema.Type, v any, errs *[]middleware.FieldError) {
	if v == nil {
		if !t.Nullable && t.Kind != apischema.KindAny {
			*errs = append(*errs, middleware.FieldError{Attribute: path, Message: "Input should not be None", Errno: 22})
		}
		return
	}
	bad := func(want string) {
		*errs = append(*errs, middleware.FieldError{Attribute: path, Message: fmt.Sprintf("Input should be a valid %s", want), Errno: 22})
	}
	switch t.Kind {
	case apischema.KindString:
		s, ok := v.(string)
		if _, isNumber := v.(json.Number); isNumber && t.IntOrString {
			return
		}
		if !ok {
			bad("string")
			return
		}
		if len(t.Enum) > 0 && !slices.Contains(t.Enum, any(s)) {
			*errs = append(*errs, middleware.FieldError{Attribute: path, Message: fmt.Sprintf("Input should be one of %v", t.Enum), Errno: 22})
		}
	case apischema.KindInt:
		n, ok := v.(json.Number)
		// Wider than int64 is still an integer: TrueNAS reports certificate serials that way.
		if _, isInt := new(big.Int).SetString(n.String(), 10); !ok || !isInt {
			bad("integer")
		}
	case apischema.KindNumber:
		if _, ok := v.(json.Number); !ok {
			bad("number")
		}
	case apischema.KindBool:
		if _, ok := v.(bool); !ok {
			bad("boolean")
		}
	case apischema.KindObject:
		obj, ok := v.(map[string]any)
		if !ok {
			bad("dictionary")
			return
		}
		// Nested objects are whole models even inside partial updates, so their defaults apply.
		checkObject(path, t, obj, true, errs)
	case apischema.KindList:
		items, ok := v.([]any)
		if !ok {
			bad("list")
			return
		}
		for i, item := range items {
			checkValue(fmt.Sprintf("%s.%d", path, i), t.Elem, item, errs)
		}
	case apischema.KindUnion:
		obj, ok := v.(map[string]any)
		if !ok {
			bad("dictionary")
			return
		}
		key, _ := obj[t.Discriminator].(string)
		for _, variant := range t.Variants {
			if k, _ := variant.VariantKey(t.Discriminator); k == key {
				checkObject(path, variant, obj, true, errs)
				return
			}
		}
		*errs = append(*errs, middleware.FieldError{Attribute: path + "." + t.Discriminator, Message: "Input tag does not match any expected tags", Errno: 22})
	}
}

// objectDefaults returns the default object Pydantic's default_factory would build: every field
// of t has a default. It returns nil otherwise.
func objectDefaults(t *apischema.Type) map[string]any {
	if t.Kind != apischema.KindObject || len(t.Fields) == 0 {
		return nil
	}
	out := map[string]any{}
	for _, f := range t.Fields {
		if !f.Type.HasDefault {
			return nil
		}
		out[f.Name] = deepCopy(f.Type.Default)
	}
	return out
}

// fillReadOnly adds zero values for read-shape fields the stored row lacks.
func fillReadOnly(read *apischema.Type, row map[string]any) {
	if read == nil {
		return
	}
	for _, f := range read.Fields {
		if _, ok := row[f.Name]; !ok {
			row[f.Name] = zero(f.Type)
		}
	}
}

func zero(t *apischema.Type) any {
	if t.Nullable || t.Kind == apischema.KindAny || t.Kind == apischema.KindUnion {
		return nil
	}
	switch t.Kind {
	case apischema.KindString:
		return ""
	case apischema.KindInt, apischema.KindNumber:
		return json.Number("0")
	case apischema.KindBool:
		return false
	case apischema.KindList:
		return []any{}
	case apischema.KindMap:
		return map[string]any{}
	default:
		obj := map[string]any{}
		fillReadOnly(t, obj)
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
	update    *apischema.Type

	mu   sync.Mutex
	data map[string]any
}

// ServeConfig registers <namespace>.config and <namespace>.update, starting from initial.
func (s *Server) ServeConfig(snap *apischema.Snapshot, namespace string, initial map[string]any) *ConfigStore {
	cs := &ConfigStore{
		namespace: namespace,
		update:    mustMethod(snap, namespace+".update").Accepts[0].Type,
		data:      deepCopy(initial).(map[string]any),
	}
	fillReadOnly(mustMethod(snap, namespace+".config").Returns, cs.data)
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
		checkObject(strings.ReplaceAll(namespace, ".", "_")+"_update", cs.update, patch, false, &fields)
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
