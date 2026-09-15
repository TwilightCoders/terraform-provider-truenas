package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

// File transfer does not go over the websocket. Reading a file means asking core.download to start
// a filesystem.get job and then fetching a one-shot URL over HTTP; writing means posting multipart
// to /_upload, which authenticates the HTTP request on its own rather than reusing the session.
//
// Both reuse the client's HTTP client, so they inherit its TLS settings and, under
// insecure_loopback, its refusal to speak plaintext to anything but a loopback address. Both also
// re-run the reverse-proxy check: the guard that protects the websocket handshake does not apply
// to these paths, and a proxy that terminates TLS in front of them is exactly what gets an API key
// revoked.
//
// The API key itself is never sent over HTTP. Uploads authenticate with a short-lived, single-use,
// origin-bound token minted over the authenticated websocket, so a proxy that logs or leaks the
// request cannot capture a durable credential.

// downloadTokenTTL bounds the life of the upload token. Middleware closes a download job's pipes
// 60 seconds after core.download returns, so neither side is worth a longer life.
const uploadTokenTTL = 60

// httpBase returns the scheme and host for the non-websocket endpoints.
func (c *Client) httpBase() string {
	if c.cfg.InsecureLoopback {
		return "http://" + c.cfg.Host
	}
	return "https://" + c.cfg.Host
}

// GetFile streams the contents of path on the TrueNAS host into w.
//
// The download URL carries its own one-shot token in the query string, which is how middleware
// requires it: the handler rejects the request unless auth_token is a query parameter. It is
// therefore visible to anything that proxies the request, so the URL is never logged and never
// appears in an error.
func (c *Client) GetFile(ctx context.Context, path string, w io.Writer) error {
	var result []json.RawMessage
	if err := c.CallInto(ctx, "core.download", &result, "filesystem.get", []any{path}, fileName(path)); err != nil {
		return fmt.Errorf("starting download of %s: %w", path, err)
	}
	if len(result) != 2 {
		return fmt.Errorf("core.download returned %d values, expected job id and URL", len(result))
	}
	var href string
	if err := json.Unmarshal(result[1], &href); err != nil {
		return fmt.Errorf("core.download returned an unreadable URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.httpBase()+href, http.NoBody)
	if err != nil {
		// The URL carries a credential, so report the path rather than the error's URL.
		return fmt.Errorf("downloading %s: malformed download URL", path)
	}
	resp, err := c.httpc.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", path, redactURL(err))
	}
	defer func() { _ = resp.Body.Close() }()

	if err := c.checkHTTPServer(resp); err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusGone:
		return fmt.Errorf("downloading %s: the download expired before it was fetched; "+
			"middleware closes the job after 60 seconds", path)
	case http.StatusUnauthorized:
		return fmt.Errorf("downloading %s: the download token was rejected", path)
	default:
		return fmt.Errorf("downloading %s: server returned %s", path, resp.Status)
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("downloading %s: %w", path, redactURL(err))
	}
	return nil
}

// PutFile writes content to path on the TrueNAS host, creating it if necessary. mode sets the
// file's permission bits when non-nil.
func (c *Client) PutFile(ctx context.Context, path string, mode *int64, content io.Reader) error {
	// Mint a token over the authenticated websocket rather than sending the API key over HTTP.
	// TrueNAS permanently revokes an API key that reaches it over a transport it considers
	// insecure, and /_upload authenticates separately from the session.
	var token string
	if err := c.CallInto(ctx, "auth.generate_token", &token, uploadTokenTTL, map[string]any{}, true, true); err != nil {
		return fmt.Errorf("authorizing upload of %s: %w", path, err)
	}

	options := map[string]any{"append": false}
	if mode != nil {
		options["mode"] = *mode
	}
	descriptor, err := json.Marshal(map[string]any{
		"method": "filesystem.put",
		"params": []any{path, options},
	})
	if err != nil {
		return err
	}

	// Stream the body: a compose file is small, but nothing here should assume that.
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)
	go func() {
		pw.CloseWithError(writeUploadBody(form, descriptor, content))
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.httpBase()+"/_upload", pr)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.Header.Set("Authorization", "Token "+token)

	resp, err := c.httpc.Do(req)
	if err != nil {
		return fmt.Errorf("uploading %s: %w", path, redactURL(err))
	}
	defer func() { _ = resp.Body.Close() }()

	if err := c.checkHTTPServer(resp); err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("uploading %s: server returned %s: %s", path, resp.Status, strings.TrimSpace(string(body)))
	}
	return nil
}

// writeUploadBody writes the two parts /_upload requires, in the order it requires them: a JSON
// "data" part naming the method, then the file contents.
func writeUploadBody(form *multipart.Writer, descriptor []byte, content io.Reader) error {
	data, err := form.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="data"`},
		"Content-Type":        {"application/json"},
	})
	if err != nil {
		return err
	}
	if _, err := data.Write(descriptor); err != nil {
		return err
	}
	file, err := form.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="file"; filename="file"`},
		"Content-Type":        {"application/octet-stream"},
	})
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, content); err != nil {
		return err
	}
	return form.Close()
}

// checkHTTPServer applies the reverse-proxy guard to a file-transfer response. The websocket
// handshake check does not cover these paths, and they are the ones that carry credentials.
func (c *Client) checkHTTPServer(resp *http.Response) error {
	if c.cfg.AllowReverseProxy {
		return nil
	}
	return checkServer(c.cfg.Host, resp)
}

// fileName returns the last path element, which middleware echoes in Content-Disposition.
func fileName(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 && i+1 < len(path) {
		return path[i+1:]
	}
	return path
}

// redactURL strips the query string from a URL error, so a download token never reaches a log or a
// Terraform diagnostic.
func redactURL(err error) error {
	var uerr *url.Error
	if !errors.As(err, &uerr) {
		return err
	}
	if u, perr := url.Parse(uerr.URL); perr == nil {
		u.RawQuery = ""
		return &url.Error{Op: uerr.Op, URL: u.String() + "?<redacted>", Err: uerr.Err}
	}
	return &url.Error{Op: uerr.Op, URL: "<redacted>", Err: uerr.Err}
}
