#!/usr/bin/env bash
# Writes a copy of the official catalog whose common pieces resolve to the
# commit under test, so CI exercises the real remote-fetch path rather than a
# local shortcut that would never fail the way a user's install can.
set -euo pipefail

out="$1"
sha="${GITHUB_SHA:-main}"

curl -fsSL https://raw.githubusercontent.com/Templetry/catalog/main/registry.json \
  | jq --arg sha "$sha" '.pieces |= map(if .repo == "Templetry/pieces" then .ref = $sha else . end)' \
  > "$out"

echo "pieces pinned to $sha"
jq -r '.pieces[] | "  \(.name) -> \(.repo)@\(.ref)"' "$out"
