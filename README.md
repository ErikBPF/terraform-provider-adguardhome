# AdGuard Home Terraform/OpenTofu Provider

[![Test](https://github.com/ErikBPF/terraform-provider-adguardhome/actions/workflows/test.yaml/badge.svg)](https://github.com/ErikBPF/terraform-provider-adguardhome/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/v/release/ErikBPF/terraform-provider-adguardhome)](https://github.com/ErikBPF/terraform-provider-adguardhome/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE.md)

An independent Terraform and OpenTofu provider for the current AdGuard Home
API. The project prioritizes safe imports, convergent plans, secret-safe state,
and resources that have one clear owner.

This repository is a fork of Gustavo Michels' MIT-licensed
[`gmichels/terraform-provider-adguard`](https://github.com/gmichels/terraform-provider-adguard).
The fork preserves the upstream history and attribution while evolving under
the `ErikBPF/adguardhome` Registry address.

Registry installation begins with the first published fork release. Before
then, use a local development override; do not assume the address resolves.

## Install

```hcl
terraform {
  required_providers {
    adguardhome = {
      source  = "ErikBPF/adguardhome"
      version = "~> 0.1"
    }
  }
}

provider "adguardhome" {
  host     = "adguard.example.net"
  username = "operator"
  scheme   = "https"
  # Set ADGUARDHOME_PASSWORD in the environment.
}
```

Do not put API passwords, TLS private keys, or certificate bodies in committed
configuration. Treat Terraform/OpenTofu state as sensitive even when every
resource avoids secret output.

## Compatibility surface

The first fork releases retain the upstream resource suffixes and schemas under
the new `adguardhome_*` prefix. For example, upstream `adguard_config` becomes
`adguardhome_config`. Changing provider source and resource addresses requires
an explicit state migration or re-import; the provider never rewrites state
implicitly. The compatibility config resource is not the long-term interface.

The planned interface is endpoint-granular: DNS, filtering, query log,
statistics, DHCP, and TLS settings become independently owned singleton
resources; rewrites, filters, clients, and leases use stable collection
identities. See [the design and migration roadmap](docs/design.md).

## Development

```shell
go build ./...
go test ./...
go generate ./...
git diff --exit-code
```

Acceptance tests use only the repository's disposable AdGuard Home container.
Never point them at an existing server. Release procedure, required gates, and
artifact verification are documented in [RELEASING.md](RELEASING.md).
Current maintainer continuation state is recorded in
[docs/continue.md](docs/continue.md) until v0.1 adoption is complete.

## License and attribution

MIT. See [LICENSE.md](LICENSE.md). Copyright and attribution for the original
provider remain intact.
