#!/usr/bin/env bash
# Regenerates app/secrets.properties (gitignored) from the bound Doppler config.
# Doppler names are UPPER_SNAKE, the two signing keys Gradle reads are not, so
# they are mapped back; everything else passes through unchanged.
set -euo pipefail

cd "$(dirname "$0")/.."
out="app/secrets.properties"
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT

doppler secrets download --no-file --format env-no-quotes \
  | grep -v '^DOPPLER_' \
  | sed -e 's/^SIGNING_KEY_ALIAS=/signing_key_alias=/' \
        -e 's/^SIGNING_KEYSTORE_PASSWORD=/signing_keystore_password=/' \
  > "$tmp"

[ -s "$tmp" ] || { echo "doppler returned no secrets; $out left untouched" >&2; exit 1; }
mv "$tmp" "$out"
trap - EXIT
echo "wrote $out ($(wc -l < "$out" | tr -d ' ') entries)"
