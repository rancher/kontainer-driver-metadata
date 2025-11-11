#!/bin/bash

# This script compares two Rancher metadata branches and checks if there are any
# RKE2 or K3s versions present in the release branch but missing in the dev branch.
#
# Usage:
#   ./check-versions.sh <dev-branch> <release-branch>
#
# Example:
#   ./check-versions.sh dev-v2.12 release-v2.12

set -euo pipefail

# Check dependencies
for cmd in jq curl; do
  if ! command -v "$cmd" &> /dev/null; then
    echo "Error: '$cmd' is not installed or not in PATH. Please install it."
    exit 1
  fi
done

if [[ $# -ne 2 ]]; then
      echo "Check if any released RKE2/K3s version from the release branch is missing in the dev branch"
      echo "Usage:"
      echo "   $0 <dev-branch> <release-branch>"
      exit 1
fi

DEV_BRANCH="$1"
REL_BRANCH="$2"

# Check if branch name is provided
if [[ -z "$DEV_BRANCH" || -z "$REL_BRANCH" ]]; then
  echo "Error: Both <dev-branch> and <release-branch> must be non-empty."
  echo "Usage: $0 <dev-branch> <release-branch>"
  exit 1
fi

FILE_DEV=$(mktemp)
FILE_REL=$(mktemp)
trap 'rm -f "$FILE_DEV" "$FILE_REL"' EXIT

BASE_URL="https://releases.rancher.com/kontainer-driver-metadata"

echo "Downloading data.json from $DEV_BRANCH..."
curl -s -f -o "$FILE_DEV" "${BASE_URL}/${DEV_BRANCH}/data.json" || {
  echo "Error: failed to download data.json from $DEV_BRANCH"
  exit 1
}

echo "Downloading data.json from $REL_BRANCH..."
curl -s -f -o "$FILE_REL" "${BASE_URL}/${REL_BRANCH}/data.json" || {
  echo "Error: failed to download data.json from $REL_BRANCH"
  exit 1
}

for f in "$FILE_DEV" "$FILE_REL"; do
  for key in ".rke2.releases" ".k3s.releases"; do
    if ! jq -e "$key" "$f" > /dev/null; then
      echo "Error: Key $key not found in $f"
      exit 1
    fi
  done
done

echo "Comparing release branch ($REL_BRANCH) against dev branch ($DEV_BRANCH)..."

VERSIONS_RKE2_DEV=$(jq -r '.rke2.releases[].version' "$FILE_DEV" | sort)
VERSIONS_RKE2_REL=$(jq -r '.rke2.releases[].version' "$FILE_REL" | sort)
VERSIONS_K3S_DEV=$(jq -r '.k3s.releases[].version' "$FILE_DEV" | sort)
VERSIONS_K3S_REL=$(jq -r '.k3s.releases[].version' "$FILE_REL" | sort)

# Compare versions and find missing ones
MISSING_RKE2=$(comm -13 <(echo "$VERSIONS_RKE2_DEV") <(echo "$VERSIONS_RKE2_REL"))
MISSING_K3S=$(comm -13 <(echo "$VERSIONS_K3S_DEV") <(echo "$VERSIONS_K3S_REL"))

# Output results
FOUND=0
echo "Comparison results:"

if [[ -n "$MISSING_RKE2" ]]; then
  FOUND=1
  echo "[FAIL] The following RKE2 versions exist in $REL_BRANCH but not in $DEV_BRANCH:"
  echo "$MISSING_RKE2"
else
  echo "[PASS] All RKE2 versions from $REL_BRANCH exist in $DEV_BRANCH."
fi

if [[ -n "$MISSING_K3S" ]]; then
  FOUND=1
  echo "[FAIL] The following K3s versions exist in $REL_BRANCH but not in $DEV_BRANCH:"
  echo "$MISSING_K3S"
else
  echo "[PASS] All K3s versions from $REL_BRANCH exist in $DEV_BRANCH."
fi

exit $FOUND
