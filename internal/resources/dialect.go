package resources

import (
	"slices"
	"strings"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/apischema"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
)

// Dialect describes the conventions TrueNAS layers on top of its own schema. None of this is
// expressible in JSON Schema, so the engine cannot infer it and this provider has to say it.
var Dialect = &engine.Dialect{
	PropertyWrapper:  isPropertyWrapper,
	Property:         engine.PropertyFields{Value: "value", Raw: "rawvalue", Parsed: "parsed", Source: "source"},
	AcceptsInherit:   acceptsInherit,
	Inherit:          "INHERIT",
	InheritedSources: []string{"INHERITED", "DEFAULT"},
	LooksSecret:      looksSecret,
}

// isPropertyWrapper reports whether t is the {value, rawvalue, parsed, source} object pool.dataset
// reports ZFS properties as, rather than the value itself.
func isPropertyWrapper(t *apischema.Type) bool {
	if t == nil || t.Kind != apischema.KindObject {
		return false
	}
	return t.Field("value") != nil && t.Field("rawvalue") != nil && t.Field("source") != nil
}

// acceptsInherit reports whether t takes "INHERIT", meaning take the value from the parent dataset.
func acceptsInherit(t *apischema.Type) bool {
	if t.Default == "INHERIT" || t.IntOrString {
		return true
	}
	for _, e := range t.Enum {
		if e == "INHERIT" {
			return true
		}
	}
	return false
}

// TrueNAS does not mark which fields hold credentials, so they are recognised by name. The rules
// are deliberately narrow: a false positive hides a value from state that the user expected to
// see, and a false negative writes a credential into it.
var secretWords = []string{"password", "passphrase", "token", "private_key", "api_key", "encryption_key", "pass"}

// notSecret lists names that mention passwords without holding one.
var notSecret = []string{"last_password_change", "password_age", "password_history", "password_change_required"}

func looksSecret(name string) bool {
	n := strings.ToLower(name)
	if slices.Contains(notSecret, n) {
		return false
	}
	for _, part := range []string{"secret", "passw", "passphrase", "pwd"} {
		if strings.Contains(n, part) {
			return true
		}
	}
	if n == "privatekey" || strings.HasSuffix(n, "_key") {
		return true
	}
	for _, w := range secretWords {
		if n == w || strings.HasSuffix(n, "_"+w) {
			return true
		}
	}
	return n == "key"
}
