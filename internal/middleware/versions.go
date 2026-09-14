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

func fetchVersions(ctx context.Context, httpc *http.Client, host string) ([]string, error) {
	url := fmt.Sprintf("https://%s/api/versions", host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discovering API versions: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovering API versions: GET %s returned %s (TrueNAS 25.04 or later is required)", url, resp.Status)
	}
	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("discovering API versions: %w", err)
	}
	return versions, nil
}
