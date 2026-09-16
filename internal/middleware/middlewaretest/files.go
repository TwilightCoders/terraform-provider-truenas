package middlewaretest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// The file endpoints are modelled on middleware's own: a download is a job started over the
// websocket whose one-shot URL carries a token in the query string, and an upload is a multipart
// POST whose first part names the method. Both are reproduced closely enough that a client which
// works here works against a real box, including the parts that are easy to get wrong — the
// ordering of the multipart parts, the token living in the query string for downloads and in the
// Authorization header for uploads, and the job expiring before it is fetched.

// Files is an in-memory filesystem served by the download and upload endpoints.
type Files struct {
	mu      sync.Mutex
	content map[string][]byte
	jobs    map[int]*downloadJob
	tokens  map[string]bool
	nextJob int
	// ExpireDownloads makes every minted download URL report itself as already expired.
	ExpireDownloads bool
}

type downloadJob struct {
	path  string
	token string
}

func newFiles() *Files {
	return &Files{
		content: map[string][]byte{},
		jobs:    map[int]*downloadJob{},
		tokens:  map[string]bool{},
		nextJob: 1,
	}
}

// Put seeds a file.
func (f *Files) Put(path string, content []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.content[path] = append([]byte(nil), content...)
}

// Delete removes a file, simulating one removed outside Terraform.
func (f *Files) Delete(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.content, path)
}

// Get returns a file's contents and whether it exists.
func (f *Files) Get(path string) ([]byte, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, ok := f.content[path]
	return append([]byte(nil), b...), ok
}

// mintToken records a token the upload endpoint will accept once.
func (f *Files) mintToken(single bool) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	tok := fmt.Sprintf("token-%d", len(f.tokens)+1)
	f.tokens[tok] = single
	return tok
}

// useToken consumes a token, reporting whether it was valid.
func (f *Files) useToken(tok string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	single, ok := f.tokens[tok]
	if ok && single {
		delete(f.tokens, tok)
	}
	return ok
}

// startDownload registers a job for path and returns its id and URL.
func (f *Files) startDownload(path string) (int, string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	id := f.nextJob
	f.nextJob++
	tok := fmt.Sprintf("dl-%d", id)
	f.jobs[id] = &downloadJob{path: path, token: tok}
	return id, fmt.Sprintf("/_download/%d?auth_token=%s", id, tok)
}

// serveDownload implements GET /_download/{id}?auth_token=...
func (s *Server) serveDownload(w http.ResponseWriter, r *http.Request) {
	s.setServerHeader(w)
	id := 0
	if _, err := fmt.Sscanf(strings.TrimPrefix(r.URL.Path, "/_download/"), "%d", &id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	token := r.URL.Query().Get("auth_token")

	s.Files.mu.Lock()
	job, ok := s.Files.jobs[id]
	expired := s.Files.ExpireDownloads
	var content []byte
	if ok {
		content, ok = s.Files.content[job.path]
	}
	s.Files.mu.Unlock()

	switch {
	case token == "" || job == nil || token != job.token:
		// Middleware answers 401 for a missing or wrong token, before it looks at the job.
		w.WriteHeader(http.StatusUnauthorized)
	case expired:
		w.WriteHeader(http.StatusGone)
	case !ok:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(content)
	}
}

// serveUpload implements POST /_upload.
func (s *Server) serveUpload(w http.ResponseWriter, r *http.Request) {
	s.setServerHeader(w)

	auth := r.Header.Get("Authorization")
	token, isToken := strings.CutPrefix(auth, "Token ")
	switch {
	case strings.HasPrefix(auth, "Bearer "):
		// A real box accepts this, and permanently revokes the key when the transport looks
		// insecure. The provider must never take this path, so the fake refuses it outright.
		http.Error(w, "the API key must not be sent over HTTP", http.StatusForbidden)
		return
	case !isToken || !s.Files.useToken(token):
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	part, err := reader.NextPart()
	if err != nil || part.FormName() != "data" {
		http.Error(w, `"data" part must be the first on payload`, http.StatusMethodNotAllowed)
		return
	}
	var descriptor struct {
		Method string            `json:"method"`
		Params []json.RawMessage `json:"params"`
	}
	if err := json.NewDecoder(part).Decode(&descriptor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if descriptor.Method != "filesystem.put" || len(descriptor.Params) == 0 {
		http.Error(w, "unsupported method "+descriptor.Method, http.StatusUnprocessableEntity)
		return
	}
	var path string
	if err := json.Unmarshal(descriptor.Params[0], &path); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filePart, err := reader.NextPart()
	if err != nil || filePart.FormName() != "file" {
		http.Error(w, `"file" not found as second part on payload`, http.StatusMethodNotAllowed)
		return
	}
	content, err := io.ReadAll(filePart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Files.Put(path, content)
	w.WriteHeader(http.StatusOK)
}

// handleFileMethod serves the websocket half of file transfer: minting an upload token and
// starting a download job. A test that registers its own handler for either name wins.
func (s *Server) handleFileMethod(method string, params []json.RawMessage) (any, bool, error) {
	switch method {
	case "auth.generate_token":
		single := false
		if len(params) >= 4 {
			_ = json.Unmarshal(params[3], &single)
		}
		return s.Files.mintToken(single), true, nil
	case "core.download":
		if len(params) < 2 {
			return nil, true, fmt.Errorf("core.download needs a method and arguments")
		}
		var name string
		if json.Unmarshal(params[0], &name) != nil || name != "filesystem.get" {
			return nil, true, fmt.Errorf("unsupported download method %s", name)
		}
		var args []string
		if err := json.Unmarshal(params[1], &args); err != nil || len(args) == 0 {
			return nil, true, fmt.Errorf("core.download needs a path")
		}
		id, href := s.Files.startDownload(args[0])
		return []any{id, href}, true, nil
	}
	return nil, false, nil
}
