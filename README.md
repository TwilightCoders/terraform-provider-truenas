# Terraform Provider for TrueNAS

Manage TrueNAS SCALE with Terraform through its versioned JSON-RPC 2.0 WebSocket API.

> **Status: pre-alpha.** Tested against an in-process fake of the TrueNAS middleware; not yet
> verified against a real TrueNAS box.

Resources are derived from TrueNAS's own API schema: each one is a short spec in
`internal/resources`, and a generic engine provides the schema, validation, create/read/update/delete,
import and resource identity.

| Resource | TrueNAS service |
|---|---|
| `truenas_acme_dns_authenticator` | `acme.dns.authenticator` |
| `truenas_certificate` | `certificate` |
| `truenas_cloudsync_credentials` | `cloudsync.credentials` |
| `truenas_cloudsync_task` | `cloudsync` |
| `truenas_cron_job` | `cronjob` |
| `truenas_dataset` | `pool.dataset` (filesystems) |
| `truenas_general_config` | `system.general` (settings) |
| `truenas_group` | `group` |
| `truenas_init_script` | `initshutdownscript` |
| `truenas_network_config` | `network.configuration` (settings) |
| `truenas_nfs_config` | `nfs` (settings) |
| `truenas_replication` | `replication` |
| `truenas_smb_config` | `smb` (settings) |
| `truenas_smb_share` | `sharing.smb` |
| `truenas_snapshot_task` | `pool.snapshottask` |
| `truenas_ssh_config` | `ssh` (settings) |
| `truenas_user` | `user` |
| `truenas_zvol` | `pool.dataset` (volumes) |

## Usage

```hcl
terraform {
  required_providers {
    truenas = {
      source = "twilightcoders/truenas"
    }
  }
}

provider "truenas" {
  host     = "nas.lan"
  username = "truenas_admin"
  # api_key from TRUENAS_API_KEY
  tls = {
    fingerprint = "AA:BB:CC:..." # TrueNAS ships a self-signed certificate
  }
}

resource "truenas_snapshot_task" "home" {
  dataset        = "tank/home"
  recursive      = true
  lifetime_value = 14
  lifetime_unit  = "DAY"
  schedule = {
    minute = "45"
    hour   = "2"
  }
}
```

Settings resources (`*_config`) manage values that always exist: creating one applies only the
attributes you set, and destroying it removes it from state without changing TrueNAS.

Set `read_only = true` to import and plan against a production box with a guarantee that nothing
is changed: any plan that would create, update or destroy a resource fails.

> **Connect to TrueNAS directly.** TrueNAS permanently revokes an API key the first time it arrives
> over plaintext. A reverse proxy that terminates TLS and forwards to TrueNAS over HTTP does exactly
> that, even though your side of the connection is encrypted. The provider refuses to send the key
> when `host` is answered by a web server other than TrueNAS's own; `allow_reverse_proxy = true`
> overrides this for proxies that re-encrypt.
>
> Without HTTPS on TrueNAS, tunnel to it over SSH and set `insecure_loopback = true`: TrueNAS accepts
> keys over plaintext from loopback, and the provider refuses plaintext to anything that is not
> loopback.
>
> ```hcl
> # ssh -f -N -L 18080:127.0.0.1:80 nas
> provider "truenas" {
>   host              = "127.0.0.1:18080"
>   insecure_loopback = true
> }
> ```

### Importing existing objects

Import blocks for this provider need `provider = truenas`; without it, Terraform's config generation
looks for `hashicorp/truenas`.

```hcl
import {
  to       = truenas_snapshot_task.home
  provider = truenas
  identity = { id = 5 }
}
```

`terraform plan -generate-config-out=generated.tf` then writes the configuration for everything imported.

With Terraform 1.14 or later, `terraform query` finds existing objects without knowing their IDs.
Put `list` blocks in a `.tfquery.hcl` file:

```hcl
list "truenas_dataset" "tank" {
  provider = truenas
  config {
    query_filters = jsonencode([["pool", "=", "tank"]])
  }
}
```

`terraform query -generate-config-out=datasets.tf` writes import blocks and configuration for every
match. Users and groups list only non-built-in accounts.

## Requirements

| Component | Version |
|---|---|
| TrueNAS SCALE | 25.10 or later |
| Terraform | 1.14 or later |
| Go (development) | pinned in `go.mod`; fetched automatically by the `go` command |

## Development

Everything runs through `make`. Run `make help` for the full list.

```bash
make build        # bin/terraform-provider-truenas
make test         # unit tests
make lint         # golangci-lint (pinned in tools/go.mod)
make docs         # regenerate docs/ with tfplugindocs
make schema-golden # accept resource schema changes
make install      # install into the local Terraform plugin mirror
make dev-override # print a ~/.terraformrc dev_overrides block
```

### Exploring the TrueNAS API

```bash
TRUENAS_HOST=nas.lan make api-docs DOC=api_methods
TRUENAS_SSH=nas.lan  make api-method METHOD=pool.snapshottask.create
```

Both are read-only.

## License

[MIT](LICENSE)
