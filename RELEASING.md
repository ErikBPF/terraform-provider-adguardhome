# Releasing

Releases are tag-driven and require the protected `production-release`
environment.

## Before tagging

1. Update `CHANGELOG.md` with user-visible changes and migration notes.
2. Run `go test -race -count=1 ./...`, `go vet ./...`, `go generate ./...`, and
   require a clean diff.
3. Run acceptance tests only against the disposable container fixture.
4. Verify examples contain no credentials, private keys, certificate bodies,
   home addresses, or state files.
5. Run `goreleaser release --snapshot --clean` and inspect archive names,
   checksums, manifest, and SPDX SBOMs.

## Publish

Create a signed semantic-version tag from protected `main` and push the tag.
The release workflow reruns acceptance tests, waits for environment approval,
imports the GPG signing key, and runs GoReleaser. GitHub then attests the
published checksum file with keyless OIDC provenance.

Repository secrets required by the workflow are `GPG_PRIVATE_KEY` and
`GPG_PASSPHRASE`. Store only their values in GitHub; never place them in files,
workflow output, issue text, or release notes.

## Verify

1. Download the checksum file and detached GPG signature from the release.
2. Verify the signature with the published maintainer key.
3. Verify every downloaded archive against the checksum file.
4. Confirm each archive has an SPDX SBOM and the Registry manifest is attached.
5. Verify the checksum provenance with GitHub CLI artifact attestation support.
6. Confirm the Terraform Registry resolves `ErikBPF/adguardhome` at the tag.

If any verification fails, do not replace assets under the same tag. Mark the
release affected, fix forward, and publish a new patch version.
