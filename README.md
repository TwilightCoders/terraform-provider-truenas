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
| `truenas_app` | `app` custom apps (hand-written) |
| `truenas_certificate` | `certificate` |
| `truenas_cloudsync_credentials` | `cloudsync.credentials` |
| `truenas_cloudsync_task` | `cloudsync` |
| `truenas_cron_job` | `cronjob` |
| `truenas_dataset` | `pool.dataset` (filesystems) |
| `truenas_filesystem_acl` | `filesystem.getacl` / `setacl` (hand-written) |
| `truenas_general_config` | `system.general` (settings) |
| `truenas_group` | `group` |
| `truenas_init_script` | `initshutdownscript` |
| `truenas_iscsi_auth`, `_extent`, `_initiator`, `_portal`, `_target`, `_target_extent` | `iscsi.*` |
| `truenas_network_config` | `network.configuration` (settings) |
| `truenas_nfs_config` | `nfs` (settings) |
| `truenas_nfs_share` | `sharing.nfs` |
| `truenas_replication` | `replication` |
| `truenas_service` | `service` (adopts existing services) |
| `truenas_smb_config` | `smb` (settings) |
| `truenas_smb_share` | `sharing.smb` |
| `truenas_snapshot_task` | `pool.snapshottask` |
| `truenas_ssh_config` | `ssh` (settings) |
| `truenas_static_route` | `staticroute` |
| `truenas_tunable` | `tunable` |
| `truenas_user` | `user` |
| `truenas_vm` | `vm` |
| `truenas_vm_device` | `vm.device` |
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
  # A publicly trusted certificate needs no tls block. For a private CA:
  # tls = { ca_pem = file("ca.pem") }
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

Every resource has a data source of the same name. Data sources for objects look up exactly one by
`id`, `name` or `query_filters`; settings data sources read the current values. Data sources never
expose sensitive or write-only attributes.

```hcl
data "truenas_acme_dns_authenticator" "cloudflare" {
  name = "cloudflare"
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

## Actions

With Terraform 1.14 or later, actions run one-off operations: `truenas_replication_run`,
`truenas_cloudsync_sync`, `truenas_snapshot_task_run`, `truenas_scrub_run`, `truenas_app_redeploy`,
`truenas_app_start`, `truenas_app_stop`, `truenas_service_control` and `truenas_ui_restart`. Job-backed
operations wait for completion.

```hcl
action "truenas_replication_run" "archive" {
  config { id = truenas_replication.archive.id }
}
```

Run one with `terraform apply -invoke=action.truenas_replication_run.archive`, or attach it to a
resource with `lifecycle { action_trigger { ... } }`.

## Functions

`provider::truenas::size_bytes("1.5T")` converts binary sizes to bytes for attributes such as
`quota` and `volsize`.

## Secrets and Terraform state

Terraform state is a plain file. The provider keeps credentials out of it with write-only attributes:
Terraform sends them to TrueNAS and never records them. Changing a write-only value alone plans
nothing; bump the matching `*_wo_version` attribute to send it again.

| Resource | Write-only | Resend with |
|---|---|---|
| `truenas_user` | `password` | `password_wo_version` |
| `truenas_certificate` | `privatekey`, `passphrase` | `privatekey_wo_version`, `passphrase_wo_version` |
| `truenas_cloudsync_credentials` | every secret under `storage` | `storage_wo_version` |
| `truenas_cloudsync_task` | `encryption_password`, `encryption_salt` | their `*_wo_version` |
| `truenas_acme_dns_authenticator` | every secret under `authenticator` | `authenticator_wo_version` |
| `truenas_replication` | `encryption_key` | `encryption_key_wo_version` |
| `truenas_dataset`, `truenas_zvol` | `encryption_options.passphrase`, `encryption_options.key` | `encryption_options_wo_version` |
| `truenas_vm_device` | display `password` under `attributes` | `attributes_wo_version` |
| `truenas_iscsi_auth` | `secret`, `peersecret` | `secret_wo_version`, `peersecret_wo_version` |

Values the provider cannot recognize as secrets are stored normally: commands in `truenas_cron_job`
and `truenas_init_script`, and compose files, are recorded as written. Keep secrets in `.env` files
or TrueNAS itself rather than inline.

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
