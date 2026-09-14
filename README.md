# Terraform Provider for TrueNAS

Manage TrueNAS SCALE with Terraform through its versioned JSON-RPC 2.0 WebSocket API.

> **Status: pre-alpha.** Under active construction. No resources are available yet.

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
