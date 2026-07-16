# Continue — AdGuard Home provider v0.1

## Last action

The `ErikBPF/terraform-provider-adguardhome` fork was rebranded to the
`ErikBPF/adguardhome` Registry address and `adguardhome_*` type prefix. Unit,
race, vet, build, generated-doc, workflow-lint, and unsigned snapshot-release
gates pass. The snapshot produced all archives, checksums, and SPDX SBOMs.

## Next action

Configure the GitHub `production-release` environment plus
`GPG_PRIVATE_KEY`/`GPG_PASSPHRASE`, verify the public key is documented, then
tag the reviewed commit `v0.1.0`. Verify the GitHub release and register
`ErikBPF/adguardhome` in the Terraform Registry before changing homelab-IaC.

## Why

The upstream provider sends an invalid empty DHCP body during an unrelated
singleton update. This fork contains the tested disabled-DHCP short circuit,
canonical refresh behavior, TLS timestamp stabilization, computed-field plan
modifiers, and diagnostic redaction needed for the imported Discovery config.

## After publication

1. Exact-pin `ErikBPF/adguardhome` `0.1.0` in homelab-IaC and regenerate its
   lock file.
2. Migrate the encrypted state provider address from `gmichels/adguard` to
   `ErikBPF/adguardhome`, then move `adguard_config.this` to
   `adguardhome_config.this`; inspect state addresses without printing values.
3. Generate a fresh saved plan. Apply only the imported config singleton, then
   run DNS/API/exporter smoke tests and require a second no-op plan.
4. Rotate the previously exposed state passphrase after the provider migration
   and re-encrypt/verify every affected state object.
5. Resume the endpoint-granular v1 roadmap only after v0.1 adoption is stable.

## Open threads

- Real release signing is blocked only by the GitHub GPG secrets/environment.
- Terraform Registry onboarding may require a manual namespace/provider step.
- The state passphrase appeared once in prior command output and must rotate.
- Servarr YAML remains authoritative until the new provider applies cleanly and
  unsupported bootstrap/TLS/user fields have an explicit retained owner.

## Do not

- Do not reuse the failed stock-provider saved plan.
- Do not use a local development override against Discovery.
- Do not remove YAML, volumes, credentials, TLS material, or runtime state.
- Do not claim v1 granular resources exist; `docs/design.md` is a roadmap.
