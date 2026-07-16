# Provider interface design

**Status:** Accepted direction; granular resources are a roadmap, not the v0 API.

## Decision

The provider will converge on one resource owner per AdGuard Home API concern.
Singleton API documents stay grouped by endpoint; splitting fields that the API
updates together would create lost-update races. Collections use stable API IDs
or canonical natural keys.

The target singleton resources are DNS settings, filtering settings, query-log
settings, statistics settings, DHCP settings, and TLS settings. Rewrites, list
filters, persistent clients, and DHCP leases are collection resources. User
rules remain one ordered-list resource because the API exposes them as one
ordered document.

Matching read-only data sources expose non-secret configuration and import
identities. They do not create a second ownership mechanism.

## Safety semantics

- Imports resolve exactly one identity and never update the server.
- Read normalizes API defaults, duration units, IP/MAC formatting, ordering,
  and null-versus-empty representations before writing state.
- Updates read the current endpoint, change only fields owned by that resource,
  and verify the post-update fingerprint.
- Collection deletion removes exactly one bound identity.
- Singleton deletion relinquishes management by default. Destructive reset
  requires an explicit option and a plan-visible warning.
- API passwords and private keys never appear in diagnostics or data sources.
  TLS file inputs store paths, not key bodies.
- Unsupported server versions or fields fail before mutation.

## Compatibility roadmap

### v0: fork compatibility

The fork keeps the existing resource suffixes and schemas under the new
`adguardhome_*` prefix: config, rewrite, list filter, user rules, and client,
plus their data sources. Moving from upstream requires an explicit provider and
resource-address state migration or re-import. Fixes in this line focus on idempotency, disabled-DHCP
behavior, import safety, current AdGuard Home compatibility, and redaction.

### v1: endpoint-granular interface

Add the granular resources and inventory data sources behind acceptance tests
covering create/read/update/delete, import, a second no-op plan, concurrent UI
drift, and multiple supported AdGuard Home versions. Publish state migration
guidance before deprecating `adguardhome_config`.

### Later majors

Remove compatibility names only after at least one stable deprecation cycle.
Write-only secret attributes may be added only when both Terraform and OpenTofu
offer compatible semantics; ordinary sensitive strings are insufficient because
they remain in state.

## Deliberate tradeoffs

Granular resources increase import count and cannot make changes across several
endpoints atomic. In return, ownership is explicit, plans are smaller, and one
resource cannot silently overwrite an unrelated concern. Endpoint grouping is
less granular than one field per resource, but it matches the API's true update
boundary and prevents shallow, conflicting abstractions.
