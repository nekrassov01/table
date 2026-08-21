#!/usr/bin/env bash
set -euo pipefail

# This script combines testdata for the specified package
# and copies it to the clipboard. This is useful for inspecting
# package output during manual testing.

pkg="${1:-}"
if [[ -z "$pkg" ]]; then
  echo "usage: $0 {text|markdown|backlog|csv|html}" >&2
  exit 1
fi

dir="$(dirname "$0")/${pkg}/testdata"
if [[ ! -d "$dir" ]]; then
  echo "not found: $dir" >&2
  exit 1
fi

{
  first=true
  for f in "$dir"/*.txt; do
    [[ ! -s "$f" ]] && continue
    perl -MEncode=decode,FB_CROAK -0777 -ne '
      exit 1 if /\x00/;
      decode("UTF-8", $_, FB_CROAK);
    ' "$f" 2>/dev/null || continue
    if "$first"; then
      first=false
    else
      printf '\n'
    fi
    cat "$f"
  done
} | pbcopy
