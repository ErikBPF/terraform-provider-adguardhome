# Changelog

All notable changes follow [Keep a Changelog](https://keepachangelog.com/) and
Semantic Versioning.

## Unreleased

### Added

- State move support from all five legacy `adguard_*` resource types after
  migrating from `gmichels/adguard` to `ErikBPF/adguardhome`.

### Fixed

- Imported or provider-migrated singleton config now converges when
  `blocked_services_pause_schedule` is explicitly null. Version 0.1.5 is
  affected by a perpetual diff when AdGuard Home normalizes the unset schedule
  time zone.

- Config creation now skips AdGuard Home's rejected empty DHCP payload when
  DHCP is absent or disabled.  TLS validation fields and the bookkeeping
  timestamp now converge across unrelated updates and no-op plans.  Version
  0.1.4 is affected by both lifecycle defects; upgrade to 0.1.5 or newer.

- Registry checksums now cover only the 13 provider archives and manifest;
  SPDX SBOMs remain separate release assets. Versions 0.1.1 and 0.1.2 are
  affected by checksums that incorrectly reference SBOM assets.
- Release 0.1.3 has valid Registry checksums and signatures, but is affected by
  missing GitHub provenance: its post-release validation looked for the
  renamed Registry manifest in `dist/` instead of its GoReleaser source path,
  so the attestation step did not run.

## [0.1.0] - 2026-07-16

### Added

- Independent `ErikBPF/adguardhome` provider identity.
- Disabled-DHCP and plan-convergence regression coverage.
- Signed checksums, SPDX SBOMs, and GitHub build-provenance attestations in the
  release contract.
- Endpoint-granular v1 interface roadmap.

### Changed

- CI is read-only on pull requests and rejects stale generated documentation.
- Release publication requires disposable acceptance tests and environment
  approval.

### Compatibility

- Existing `adguard_*` resources and data sources remain the v0 compatibility
  surface. Granular resource names are not yet implemented.
