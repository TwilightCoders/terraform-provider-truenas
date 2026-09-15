# `truenas_file` — design note

Status: proposal. No code yet.

## Correction

An earlier design note for `truenas_app` said TrueNAS has no API to read a file back, and
recommended leaving file shipping to `apply.sh` on the strength of that. **The premise was wrong**,
and the recommendation built on it should not be trusted. The read exists, and the write method is
`filesystem.put`, not `filesystem.file_receive`.

That claim was asserted without being tested. It is the expensive kind of wrong: it closed off a
resource that turns out to be buildable, and it did so in a document that reads as settled.

## The mechanism, as verified

Reading a file is a two-step dance across two transports:

1. `core.download("filesystem.get", ["/path"], "<filename>")` → `[job_id, "/_download/<job_id>?auth_token=…"]`
2. `GET https://host/_download/<job_id>?auth_token=…` → the bytes, chunked `application/octet-stream`

Writing is a multipart POST to `/_upload`: first part `data` is JSON
`{"method": "filesystem.put", "params": ["/path", {"mode": …, "append": false}]}`, second part
`file` carries the bytes.

`filesystem.stat` gives `size`, `mtime`, `mode`, `uid`, `gid`, `inode` cheaply, with no transfer.

## Four constraints that shape the design

These come from reading middlewared rather than from probing, and three of them are not visible
from a successful manual download.

### 1. `filesystem.get` and `filesystem.put` require `FULL_ADMIN`

Every other method this provider calls is reachable with a narrower role. `filesystem.stat` needs
only `READONLY_ADMIN`. Adding `truenas_file` therefore **escalates the privilege of the API key the
whole provider runs as** — one resource type forces every other resource in the same configuration
to run as full admin.

This is the single biggest cost of the feature and it should be a deliberate decision, not a
side effect. It argues for keeping `truenas_file` in a separate Terraform configuration with its
own key, rather than mixing it into the configuration that manages datasets and shares.

### 2. The download URL expires in 60 seconds

`FileApplication.register_job` (`middlewared/apps/file_app.py`) schedules cleanup at **60 seconds**
for an unbuffered job, 3600 if buffered. After that the job's pipes close and the URL returns
`410 Gone`; a bad or missing token is `401`.

So the download is not a URL that can be minted, queued, and fetched later. Mint and fetch must be
adjacent, and a refresh of many files is a sequence of short-lived fetches, each of which can expire
under a slow link. Retry means minting a fresh job, not retrying the URL.

### 3. The auth token travels in the query string

`?auth_token=…` lands in the access log of anything that proxies the request. Where a reverse proxy
fronts the API, this token is written to its access log on every file read, in cleartext, by
default. It is single-use and job-scoped, which limits the damage, but it is still a
credential in a log file.

The provider must never log the URL, and the token must be redacted from every error message that
quotes it. Worth being explicit because the natural thing to do when a fetch fails is to log the URL.

### 4. `/_upload` authenticates over HTTP, not over the websocket session

`upload()` calls `parse_credentials(request)` and authenticates the HTTP request on its own. The
established JSON-RPC session does not carry over. So writing a file means **sending the API key as
an HTTP header** — a second transport with its own authentication.

This is the exact shape of the failure that burned the first API key: middlewared revokes a
plaintext API key presented over a transport it considers insecure, and a TLS-terminating reverse
proxy forwarding cleartext to `:8080` is such a transport. The provider already guards the
websocket path (`checkServer`, the loopback-only plaintext dialer). **The upload path needs the same
guard, and it does not inherit it.** Shipping `truenas_file` without that is an invitation to burn
another key.

## No server-side hashing

There is no checksum method — `filesystem.*` has nothing, and the only hash-shaped methods in the
API belong to network interfaces and NVMe-oF. So content drift cannot be detected without
transferring the file and hashing locally.

`filesystem.stat` is the only cheap signal, and it is advisory:

- An edit that changes content always moves `mtime`, so stat-gating catches the realistic cases.
- An edit that restores `mtime` and preserves `size` defeats it. Rare, but a restored backup or a
  `cp -p` does exactly that.

So stat-gating is a performance optimisation that can miss drift, not a correctness mechanism. It
should be opt-in per resource (`drift_detection = "stat"` vs `"content"`), default `"content"`, and
the documentation must say plainly what `"stat"` can miss. Silently defaulting to the fast path
would reintroduce the invisible-drift problem this resource exists to solve.

## Where the line is on reading into state

This is the part worth arguing about, and the ask was specifically for it.

**Content never enters state.** `content` is a write-only attribute. State holds:

- `sha256` of the content the provider last wrote or read
- `size`, `mtime`, `mode`, `uid`, `gid` from `stat`
- the path

**One resource, one explicit path. No globs, no directory recursion.** A glob turns "manage these
files" into "read whatever matches into my state file", and the set changes without the
configuration changing. Every path that gets read is a path someone wrote down.

**A hash is not nothing.** For a file with low entropy or guessable contents, a stored SHA-256
confirms a guess. `/mnt/…/apps/docker/.env` is exactly that: an attacker who can read state and can
guess the shape of the file can verify a candidate offline. So for secret files the answer is not
"hash it instead of storing it" — it is **do not read it at all**:

- `drift_detection = "none"` — the provider writes on change of `content_wo_version` and never
  reads the file back. State holds the path and nothing derived from the contents.
- This is the mode `.env` files use. It gives up drift detection deliberately, and says so.

**Refuse by default rather than ask forgiveness.** A file the provider did not write should not be
adopted into state silently. Importing an existing file is an explicit `terraform import`, the same
as everything else.

## Sketch

```hcl
resource "truenas_file" "compose" {
  path    = "/mnt/<pool>/apps/docker/example.yml"
  content = file("${path.module}/stacks/example.yml")   # write-only
  mode    = "0644"

  # "content" (default) downloads and hashes on refresh
  # "stat"    compares size and mtime only; can miss a same-size, mtime-preserving edit
  # "none"    never reads the file back; for files whose contents must not be derived into state
  drift_detection = "content"
}
```

- **Create / Update** → `/_upload` multipart with `filesystem.put`, then `filesystem.setperm` /
  `filesystem.chown` if `mode`/`uid`/`gid` are set and `put`'s `mode` option is insufficient.
- **Read** → `filesystem.stat`; then, per `drift_detection`, mint `core.download` and hash the
  stream without buffering the whole file in memory.
- **Delete** → there is no `filesystem.unlink`. Deleting the resource removes it from state and
  leaves the file. This must be documented loudly; it is surprising and it is not fixable from the
  API surface as it stands.

The missing delete is worth weighing before building. A resource that cannot remove what it created
is a leak by construction, and `stacks/apply.sh` can `rm`.

## Open questions for Dale

1. **Is the `FULL_ADMIN` escalation acceptable?** If `truenas_file` shares a configuration with the
   dataset and share resources, all of them run as full admin. Separate configuration and key, or
   accept it?
2. **No delete.** Is "Terraform writes it and never removes it" acceptable, or does that keep a
   bash reconciler in the loop anyway for teardown?
3. **Which files are in scope?** Compose files, yes. The recommendation here is that `.env` uses
   `drift_detection = "none"` and is otherwise left to whatever manages secrets today.

## What this replaces

With read and drift working, `redeploy_trigger = filesha256(…)` stops being a local-only proxy —
it can compare against what is actually on the box. That is the difference between Terraform
orchestrating `stacks/apply.sh` and replacing its file-shipping half.
