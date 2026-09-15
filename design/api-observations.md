# Where the API schema and the server disagree

`core.get_methods` is generated from middleware's Pydantic models. It describes what the models
declare, not what the server does. Every entry below is a case where the two differ, found by a
real apply against real hardware rather than by reading the schema.

Each one is recorded with how it was proven, because "the schema says X" turned out not to be
evidence.

## certificate.serial — an integer that does not fit an integer

The schema types it `integer`. X.509 serials are up to 20 octets, and a public CA issues ~140-bit
random values, so `Int64Attribute` refused to read the resource back out of state and the whole
resource became unusable — a plan that could only propose destroying a live certificate.

Worse, and invisible: `big.Float` defaults to a 64-bit mantissa, so parsing the serial silently
rounded it before storing.

```
want 586444785138067375905931086053142391461381
got  586444785138067375900000000000000000000000
```

**Proven by** a round-trip test that fails on the rounded value; the wire protocol itself is
arbitrary-precision, so only the provider's own conversions were lossy.

**Now:** every integer but the primary key is arbitrary-precision, and conversions go through
`big.Int`.

## certificate.san — written one way, stored another

`normalize_san` splits on the first colon, or infers `IP` for an address and `DNS` for anything
else. The builder beneath it honours **only** `IP`; every other type, `email:` and `URI:` included,
becomes a DNS name. On read, OpenSSL renders the extension, labelling the two cases `DNS:` and
`IP Address:`.

So a bare name written to the API reads back prefixed, and `san` is readable and forces
replacement: every plan proposed replacing a live certificate, and the second plan was the first to
show it.

Consequences that do not follow from the schema at all:

| Written | Stored as | Reported |
|---|---|---|
| `example.com` | DNS name | `DNS:example.com` |
| `IP:192.0.2.1` | IP name | `IP Address:192.0.2.1` |
| `email:a@example.com` | **DNS name** `a@example.com` | `DNS:a@example.com` |
| `IP Address:192.0.2.1` | **DNS name** `192.0.2.1` | `DNS:192.0.2.1` |

The reported form is not valid input. A bare IPv6 address is mangled, because the colon split
happens before anything looks at the value.

**Proven by** reading `truenas_crypto_utils/generate_utils.py` and by rendering a certificate with
all four name types through pyOpenSSL on a live box, rather than guessing the labels.

**Now:** comparison canonicalises, so both spellings mean the same name.

## csr, certificate, and ten subject fields — the server owns them

`csr` is declared as something the operator supplies. For `CERTIFICATE_CREATE_CSR` TrueNAS
generates it and reports it in `CSR`. The schema's read type is nullable, which is true but
useless: it is nullable because it is empty for *other* create types, not because the server may
leave it unset for this one.

Anything typed as user-supplied but filled by the server plans as `null` against a real value, and
on a replace-forcing attribute that proposes destroying the resource forever. Ten attributes on
this resource were in that state; only two were noticed, because the configuration happened to set
the rest.

**Proven by** `certificate.get_instance` on a real CSR showing the field populated.

**Now:** a field that forces replacement and is reported by the server is optional and computed,
and carries no schema default.

## certificate.digest_algorithm — applied, and never reported

The most misleading of the four, because both the schema and the API agree, and both are wrong.

TrueNAS accepts it, honours it, and then reports `null` for it forever. Setting it explicitly
compared the configured value against that null on every plan, and the attribute forces
replacement, so the certificate was reissued on every apply.

```
CSR on disk:  Signature Algorithm: sha256WithRSAEncryption
API reports:  digest_algorithm = None
```

**Proven by** reading the generated CSR with `openssl req -noout -text` — a second observation
point. No comparison of request to response could have detected it: the request succeeded, the
response was well-formed, and the schema described the field correctly.

**Now:** `engine.Field.Unreported` records such a field instead of comparing it.

## certificate ACME dns_mapping — keys must match a computed list verbatim

`issue_cert.py` rejects any key not in `certificate.get_domain_names(csr_id)`, which returns
`[common] + san` **raw** — the common name bare, the SANs `DNS:`-prefixed. So a wildcard given as
the common name is accepted bare while the apex must be written `DNS:apex.example.com`.

Wildcard plus apex in one certificate is buildable; the validator is the only picky part, and its
error message names the domain without saying what it was compared against.

**Proven by** reading the validator and calling `certificate.get_domain_names` directly.

## API keys are revoked, not rejected, over plaintext

Middleware permanently revokes a plaintext API key the first time it arrives over a transport it
considers insecure. A reverse proxy that terminates TLS and forwards HTTP is such a transport, and
the resulting error says the key is `EXPIRED`, which sounds like something that will resolve.

The first key used against a real box died this way, before a single resource was managed.

**Now:** the provider refuses to send a key when the `Server` header is not TrueNAS's own web
server, on both the websocket and the file-transfer endpoints, and file uploads authenticate with a
short-lived token minted over the already-authenticated session so the key never touches HTTP.

## The lesson under all of them

Four of these were found in one day, by one person running `terraform plan` against production
hardware. None were found by reading the schema, and the test suite agreed with the schema every
time — twice because the fake server was wrong in the same direction as the assumption.

When the schema and the server disagree, the server wins, and the only way to know is to look at
the server. Where possible, look at a second point: the certificate on disk rather than the API's
opinion of it.
