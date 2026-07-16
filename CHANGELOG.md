# Changelog

All notable changes follow [Keep a Changelog](https://keepachangelog.com/) and
Semantic Versioning.

## Unreleased

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
