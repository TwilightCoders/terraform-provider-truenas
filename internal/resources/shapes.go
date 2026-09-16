package resources

// Shape is what a call to one namespace looks like, exposed so a fake server can accept the same
// calls the provider makes and answer with the same fields it reads.
//
// It is derived from the generated descriptors rather than from a copy of the API's own schema,
// which is deliberate: this repository holds what is needed to use the provider, and a vendor's
// schema document is neither ours to license nor needed at runtime. A fake built from the
// descriptors validates exactly the calls the provider is capable of making, which is all a fake
// is ever asked to judge.
type Shape struct {
	// Namespace is the API namespace, e.g. "pool.dataset".
	Namespace string
	// PrimaryKey is the field identifying an object, and IntPK reports whether it is an integer.
	PrimaryKey string
	IntPK      bool
	// Singleton is set for settings, which are read and updated but never created.
	Singleton bool
	// Discriminator and Variant are set when the resource is one branch of a union create: the
	// provider always sends Discriminator = Variant, which is not otherwise an attribute.
	Discriminator string
	Variant       string
	// Fields are the attributes of one object, keyed by the name the API uses.
	Fields []ShapeField
}

// ShapeField is one field of a call or a result.
type ShapeField struct {
	// Name is the field as the API spells it.
	Name string
	// Kind is "string", "int", "number", "bool", "object", "list", "map", "union" or "any".
	Kind string
	// Required reports whether a create must supply it.
	Required bool
	// Readable reports whether the API reports it back, so a fake knows what to answer with.
	Readable bool
	// Updatable reports whether an update accepts it.
	Updatable bool
	// Nullable reports whether the field accepts null.
	Nullable bool
	// Property reports whether the API wraps this value in an object rather than reporting it
	// directly, as ZFS properties are.
	Property bool
	// Discriminator names the field that tells a union's variants apart. When set, Children are
	// the variants rather than fields, each named by the value it is chosen by.
	Discriminator string
	// IntOrString marks a field the API accepts as either.
	IntOrString bool
	// Default is the value the server chooses when a create omits the field.
	Default    any
	HasDefault bool
	// Children are an object's fields; Elem is a list or map's element.
	Children []ShapeField
	Elem     *ShapeField
}

// ShapeFor returns the shape of a resource type, and whether the provider serves it. Keyed by
// type rather than namespace: dataset and zvol are variants of one namespace and differ in shape.
func ShapeFor(resourceType string) (Shape, bool) {
	m, ok := descriptors[resourceType]
	if !ok {
		return Shape{}, false
	}
	s := Shape{
		Namespace:     m.namespace,
		PrimaryKey:    m.primaryKey,
		IntPK:         m.idKind == kindInt,
		Singleton:     m.singleton,
		Discriminator: m.discriminator,
		Variant:       m.variant,
	}
	for _, a := range m.attrs {
		if a.api == "" || a.name == idAttr {
			continue
		}
		s.Fields = append(s.Fields, shapeField(a))
	}
	return s, true
}

func shapeField(a *node) ShapeField {
	f := ShapeField{
		Name:          a.api,
		Kind:          a.kind.String(),
		Required:      a.role == roleRequired,
		Readable:      a.readable,
		Updatable:     a.updatable,
		Nullable:      a.nullable,
		Property:      a.property,
		Discriminator: a.discriminator,
		IntOrString:   a.intOrString,
		Default:       a.def,
		HasDefault:    a.hasDefault,
	}
	for _, c := range a.children {
		f.Children = append(f.Children, shapeField(c))
	}
	if a.elem != nil {
		e := shapeField(a.elem)
		f.Elem = &e
	}
	return f
}
