// Code generated from the API schema. DO NOT EDIT.

package resources

// Kind, role and node mirror what the generator derived. They are data here: nothing in this
// package decides them, it only executes what the descriptors carry.

type kind uint8

const (
	kindAny kind = iota
	kindString
	kindInt
	kindNumber
	kindBool
	kindObject
	kindList
	kindMap
	kindUnion
)

func (k kind) String() string {
	return [...]string{"any", "string", "int", "number", "bool", "object", "list", "map", "union"}[k]
}

type role uint8

const (
	roleRequired role = iota
	roleOptional
	roleOptionalComputed
	roleComputed
)

// node is one Terraform attribute and how it maps to the API.
type node struct {
	name string
	api  string
	path string

	kind     kind
	role     role
	nullable bool

	replace     bool
	stable      bool
	sensitive   bool
	writeOnly   bool
	readable    bool
	identity    bool
	updatable   bool
	createOnly  bool
	property    bool
	rawProperty bool
	inherit     bool
	intOrString bool
	hasDefault  bool
	def         any

	description   string
	ref           string
	discriminator string
	canonical     func(string) string
	dialect       *dialect

	children []*node
	elem     *node
}

func (a *node) child(name string) *node {
	for _, c := range a.children {
		if c.name == name {
			return c
		}
	}
	return nil
}

// model is one resource: its attributes and the methods that act on them.
type model struct {
	typeName   string
	namespace  string
	primaryKey string
	idKind     kind

	createMethod string
	updateMethod string
	getMethod    string
	deleteMethod string

	singleton   bool
	adopt       bool
	adoptBy     string
	adoptByAttr string

	discriminator string
	variant       string

	listFilters [][]any
	listOptions map[string]any

	attrs    []*node
	versions []*node
}
