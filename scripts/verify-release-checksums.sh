#!/usr/bin/env bash
set -euo pipefail

dist_dir=${1:-dist}
mapfile -t checksum_files < <(find "$dist_dir" -maxdepth 1 -type f -name '*_SHA256SUMS')
if test "${#checksum_files[@]}" -ne 1; then
  echo "expected exactly one *_SHA256SUMS in $dist_dir; found ${#checksum_files[@]}" >&2
  exit 1
fi
checksum_file=${checksum_files[0]}

mapfile -t archives < <(find "$dist_dir" -maxdepth 1 -type f -name '*.zip' -printf '%f\n' | sort)
if test "${#archives[@]}" -ne 13; then
  echo "expected exactly 13 provider archives in $dist_dir; found ${#archives[@]}" >&2
  exit 1
fi

checksum_name=${checksum_file##*/}
manifest_name=${checksum_name%_SHA256SUMS}_manifest.json
mapfile -t manifests < <(find "$dist_dir" -maxdepth 1 -type f -name '*_manifest.json' -printf '%f\n')
if test "${#manifests[@]}" -ne 1 || test "${manifests[0]}" != "$manifest_name"; then
  echo "expected published Registry manifest $manifest_name in $dist_dir" >&2
  exit 1
fi

expected=$(mktemp)
actual=$(mktemp)
trap 'rm -f "$expected" "$actual"' EXIT
printf '%s\n' "${archives[@]}" "$manifest_name" | sort >"$expected"
awk '{print $2}' "$checksum_file" | sort >"$actual"

if grep -Eq '\.sbom\.|\.spdx\.' "$actual"; then
  echo "checksum must not reference SBOM assets" >&2
  exit 1
fi

if ! diff -u "$expected" "$actual"; then
  echo "checksum must reference exactly 13 archives and the Registry manifest" >&2
  exit 1
fi

if ! (cd "$dist_dir" && sha256sum --check --strict "$checksum_name" >/dev/null); then
  echo "one or more Registry release checksums do not match local subjects" >&2
  exit 1
fi

printf '%s\n' "$checksum_file"
