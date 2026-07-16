#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
verifier=$repo_root/scripts/verify-release-checksums.sh
fixture=$(mktemp -d)
trap 'rm -rf "$fixture"' EXIT

for platform in \
  darwin_amd64 darwin_arm64 \
  freebsd_386 freebsd_amd64 freebsd_arm freebsd_arm64 \
  linux_386 linux_amd64 linux_arm linux_arm64 \
  windows_386 windows_amd64 windows_arm64; do
  : >"$fixture/terraform-provider-adguardhome_9.8.7_${platform}.zip"
done

checksum=$fixture/terraform-provider-adguardhome_9.8.7_SHA256SUMS
manifest=$fixture/terraform-provider-adguardhome_9.8.7_manifest.json
: >"$manifest"
(cd "$fixture" && sha256sum -- *.zip "${manifest##*/}") >"$checksum"

test "$($verifier "$fixture")" = "$checksum"
expected_subjects=$(mktemp)
actual_subjects=$(mktemp)
find "$fixture" -maxdepth 1 -type f -name '*.zip' -printf '%f\n' | sort >"$expected_subjects"
printf '%s\n' "${manifest##*/}" >>"$expected_subjects"
sort -o "$expected_subjects" "$expected_subjects"
awk '{print $2}' "$checksum" | sort >"$actual_subjects"
diff -u "$expected_subjects" "$actual_subjects"
rm -f "$expected_subjects" "$actual_subjects"

printf '%064d  terraform-provider-adguardhome_9.8.7_linux_amd64.zip.spdx.sbom.json\n' 0 >>"$checksum"
if "$verifier" "$fixture" >"$fixture/stdout" 2>"$fixture/stderr"; then
  echo "expected SBOM checksum fixture to fail" >&2
  exit 1
fi
grep -Fq 'checksum must not reference SBOM assets' "$fixture/stderr"
