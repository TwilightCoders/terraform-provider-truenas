# Terraform Provider for TrueNAS

Manage TrueNAS SCALE with Terraform through its versioned JSON-RPC 2.0 WebSocket API.

> **Status: pre-alpha.** Tested against an in-process fake of the TrueNAS middleware; not yet
> verified against a real TrueNAS box.

Resources are derived from TrueNAS's own API schema: each one is a short spec in
`internal/resources`, and a generic engine provides the schema, validation, create/read/update/delete,
import and resource identity.

| Resource | TrueNAS service |
|---|---|
| `truenas_cron_job` | `cronjob` |
| `truenas_group` | `group` |
| `truenas_init_script` | `initshutdownscript` |
| `truenas_smb_share` | `sharing.smb` |
| `truenas_snapshot_task` | `pool.snapshottask` |
| `truenas_user` | `user` |

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

Set `read_only = true` to import and plan against a production box with a guarantee that nothing
is changed: any plan that would create, update or destroy a resource fails.

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
