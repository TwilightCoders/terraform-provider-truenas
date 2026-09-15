package resources

import "testing"

// The secret heuristic is deliberately narrow in both directions: a false positive hides a value
// the user expected to see in state, a false negative writes a credential into it.
func TestLooksSecret(t *testing.T) {
	for name, want := range map[string]bool{
		"password": true, "api_key": true, "secret_access_key": true, "private_key": true, "key": true,
		"pass": true, "token": true, "password_disabled": true, "host_rsa_key": true,
		"host_rsa_key_pub": false, "privatekey": true, "key_format": false, "monpwd": true,
		"v3_privpassphrase": true, "sed_passwd": true, "last_password_change": false,
		"keyboard": false, "username": false, "bypass": false,
	} {
		if looksSecret(name) != want {
			t.Errorf("looksSecret(%q) = %v", name, !want)
		}
	}
}
