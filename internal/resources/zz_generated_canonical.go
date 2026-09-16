// Code generated from the API schema. DO NOT EDIT.

package resources

import (
	"net/netip"
	"strings"
)

// CanonicalSAN returns the subject alternative name TrueNAS will report after being sent s, so an
// equivalent spelling in configuration is not read back as drift.
//
// TrueNAS writes and reads SANs through different vocabularies. On write, normalize_san
// (truenas_crypto_utils/generate_utils.py:56) splits on the first colon, or infers "IP" for a bare
// address and "DNS" for anything else. The builder beneath it (generate_utils.py:70) then honours
// only "IP" — every other type, including "email" and "URI", becomes a DNSName. On read, the SAN
// extension is rendered by OpenSSL (truenas_crypto_utils/read.py:133), which labels those two
// cases "DNS:" and "IP Address:".
//
// So only two forms can ever come back, and "email:" or "URI:" in configuration is not a way to
// request those name types: it produces a DNS name whose value keeps the rest of the string.
func CanonicalSAN(s string) string {
	typ, value := "", s
	if t, v, found := strings.Cut(s, ":"); found {
		// An unprefixed IPv6 address is cut here too, exactly as TrueNAS cuts it. Reproducing that
		// keeps the comparison honest: "2001:db8::1" really is stored as the DNS name "db8::1".
		typ, value = t, v
	} else if _, err := netip.ParseAddr(s); err == nil {
		typ = "IP"
	}
	if typ == "IP" {
		return "IP Address:" + value
	}
	return "DNS:" + value
}
