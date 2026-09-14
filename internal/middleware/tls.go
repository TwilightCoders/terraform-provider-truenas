package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// TLSConfig controls how the server certificate is verified.
//
// With no fields set the system roots are used. TrueNAS ships a self-signed certificate, so most
// installations set CAPEM or pin Fingerprint.
type TLSConfig struct {
	// CAPEM adds PEM-encoded CA certificates to trust.
	CAPEM string
	// Fingerprint pins the server's leaf certificate by its SHA-256 digest (hex, colons optional).
	// When set, chain and hostname verification are replaced by the pin.
	Fingerprint string
	// InsecureSkipVerify disables verification entirely.
	InsecureSkipVerify bool
	// ServerName overrides the name used for verification and SNI.
	ServerName string
}

func (t TLSConfig) build() (*tls.Config, error) {
	cfg := &tls.Config{MinVersion: tls.VersionTLS12, ServerName: t.ServerName}

	if t.CAPEM != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(t.CAPEM)) {
			return nil, errors.New("tls: ca_pem contains no valid certificates")
		}
		cfg.RootCAs = pool
	}

	if t.Fingerprint != "" {
		want, err := parseFingerprint(t.Fingerprint)
		if err != nil {
			return nil, err
		}
		cfg.InsecureSkipVerify = true
		cfg.VerifyConnection = func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errors.New("tls: server presented no certificate")
			}
			got := sha256.Sum256(cs.PeerCertificates[0].Raw)
			if subtle.ConstantTimeCompare(got[:], want) != 1 {
				return fmt.Errorf("tls: certificate fingerprint %s does not match the pinned fingerprint", hex.EncodeToString(got[:]))
			}
			return nil
		}
		return cfg, nil
	}

	cfg.InsecureSkipVerify = t.InsecureSkipVerify
	return cfg, nil
}

func parseFingerprint(s string) ([]byte, error) {
	clean := strings.ToLower(strings.NewReplacer(":", "", " ", "").Replace(strings.TrimPrefix(strings.TrimSpace(s), "sha256:")))
	b, err := hex.DecodeString(clean)
	if err != nil || len(b) != sha256.Size {
		return nil, fmt.Errorf("tls: fingerprint must be a SHA-256 digest in hex, got %q", s)
	}
	return b, nil
}
