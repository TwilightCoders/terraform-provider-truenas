package middlewaretest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	st.create = mustMethod(snap, namespace+".create").Accepts[0].Type
	st.read = mustMethod(snap, namespace+".get_instance").Returns
	if snap.HasMethod(namespace + ".update") {
		st.update = mustMethod(snap, namespace+".update").Accepts[1].Type
	}

	s.Handle(namespace+".create", st.handleCreate)
	s.Handle(namespace+".get_instance", st.handleGet)
	s.Handle(namespace+".delete", st.handleDelete)
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

func (st *Store) handleQuery(context.Context, []json.RawMessage) (any, error) {
	return st.Rows(), nil
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
		if _, err := n.Int64(); !ok || err != nil {
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
