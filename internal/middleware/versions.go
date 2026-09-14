package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// apiVersion is a parsed "vMAJOR.MINOR.PATCH" middleware API version.
type apiVersion struct {
	major, minor, patch int
}

func parseAPIVersion(s string) (apiVersion, bool) {
	parts := strings.Split(strings.TrimPrefix(s, "v"), ".")
	if !strings.HasPrefix(s, "v") || len(parts) != 3 {
		return apiVersion{}, false
	}
	var v apiVersion
	for i, dst := range []*int{&v.major, &v.minor, &v.patch} {
		n, err := strconv.Atoi(parts[i])
		if err != nil || n < 0 {
			return apiVersion{}, false
		}
		*dst = n
	}
	return v, true
}

// SelectVersion picks the highest available API version in want's major.minor series that is not
// newer than want.
func SelectVersion(available []string, want string) (string, error) {
	w, ok := parseAPIVersion(want)
	if !ok {
		return "", fmt.Errorf("invalid API version %q: expected vMAJOR.MINOR.PATCH", want)
	}
	best, found := "", apiVersion{patch: -1}
	for _, s := range available {
		v, ok := parseAPIVersion(s)
		if !ok || v.major != w.major || v.minor != w.minor || v.patch > w.patch {
			continue
		}
		if v.patch > found.patch {
			best, found = s, v
		}
	}
	if best == "" {
		return "", fmt.Errorf("server offers API versions %v; this provider needs v%d.%d.x (at most %s)",
			available, w.major, w.minor, want)
	}
	return best, nil
}

// ProxyError reports that the host is served by something other than TrueNAS's own web server.
type ProxyError struct {
	Host   string
	Server string
}

func (e *ProxyError) Error() string {
	return fmt.Sprintf("%s is answered by %q, not TrueNAS's own web server. TrueNAS permanently revokes an API key "+
		"the first time it arrives over plaintext, which is what happens behind a reverse proxy that terminates TLS and "+
		"forwards over HTTP. Point host at TrueNAS's HTTPS port directly, or set allow_reverse_proxy if the proxy "+
		"re-encrypts to TrueNAS.", e.Host, e.Server)
}

// checkServer rejects responses from a web server other than TrueNAS's nginx.
func checkServer(host string, resp *http.Response) error {
	server := resp.Header.Get("Server")
	if server == "" || strings.HasPrefix(strings.ToLower(server), "nginx") {
		return nil
	}
	return &ProxyError{Host: host, Server: server}
}

func fetchVersions(ctx context.Context, httpc *http.Client, base string, allowProxy bool) ([]string, error) {
	url := base + "/api/versions"
	host := strings.TrimPrefix(strings.TrimPrefix(base, "https://"), "http://")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discovering API versions: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if !allowProxy {
		if err := checkServer(host, resp); err != nil {
			return nil, err
		}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovering API versions: GET %s returned %s (TrueNAS 25.04 or later is required)", url, resp.Status)
	}
	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("discovering API versions: %w", err)
	}
	return versions, nil
}
