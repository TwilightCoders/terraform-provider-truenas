// Code generated from the API schema. DO NOT EDIT.

package resources

// The runtime half of the API's conventions: how a wrapped value is shaped and what "inherit"
// means. The predicates that decided which fields these apply to ran at generation time and are
// not here.

type dialect struct {
	property         propertyFields
	inherit          string
	inheritedSources []string
}

// propertyFields names the parts of a wrapped value.
type propertyFields struct {
	value  string
	raw    string
	parsed string
	source string
}

// inheritedSource reports whether a wrapper's source means the value was not set here.
func (d *dialect) inheritedSource(src string) bool {
	if d == nil {
		return false
	}
	for _, s := range d.inheritedSources {
		if s == src {
			return true
		}
	}
	return false
}
