#!/usr/bin/env bash
set -euo pipefail

dist_dir=${1:-dist}
mapfile -t checksum_files < <(find "$dist_dir" -maxdepth 1 -type f -name '*_SHA256SUMS')
test "${#checksum_files[@]}" -eq 1
checksum_file=${checksum_files[0]}

mapfile -t archives < <(find "$dist_dir" -maxdepth 1 -type f -name '*.zip' -printf '%f\n' | sort)
test "${#archives[@]}" -eq 13

mapfile -t manifests < <(find "$dist_dir" -maxdepth 1 -type f -name '*_manifest.json' -printf '%f\n')
test "${#manifests[@]}" -eq 1

expected=$(mktemp)
actual=$(mktemp)
trap 'rm -f "$expected" "$actual"' EXIT
printf '%s\n' "${archives[@]}" "${manifests[0]}" | sort >"$expected"
awk '{print $2}' "$checksum_file" | sort >"$actual"

if grep -Eq '\.sbom\.|\.spdx\.' "$actual"; then
  echo "checksum must not reference SBOM assets" >&2
  exit 1
fi

if ! diff -u "$expected" "$actual"; then
  echo "checksum must reference exactly 13 archives and the Registry manifest" >&2
  exit 1
fi

printf '%s\n' "$checksum_file"
