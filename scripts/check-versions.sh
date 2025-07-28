#!/bin/bash

# This Bash script checks whether any released RKE, RKE2, or K3s versions
# from a specified Rancher metadata branch are missing in a local data.json file.
#
# Usage: ./scripts/check-versions.sh <branch-name>
#
# Example: ./scripts/check-versions.sh release-v2.11


# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: jq is not installed or not in PATH. Please install jq."
    exit 1
fi

if ! command -v curl &> /dev/null; then
    echo "Error: curl is not installed or not in PATH. Please install curl."
    exit 1
fi

# Check if branch name is provided
if [[ -z "$1" ]]; then
    echo "Check if any released RKE/RKE2/K3s version from the target branch is missing in the local data.json file"
    echo "Usage:"
    echo "  $0 <branch-name>"
    exit 1
fi

# Local file
FILE_LOCAL="$(dirname $0)/../data/data.json"
# Check if local file exists and is readable
if [[ ! -f "$FILE_LOCAL" ]]; then
    echo "Error: Local file not found: $FILE_LOCAL"
    exit 1
fi

# Create a secure temporary file
FILE_DOWNLOADED=$(mktemp)
# Ensure cleanup even on error/exit
trap 'rm -f "$FILE_DOWNLOADED"' EXIT

BRANCH="$1"
URL="https://releases.rancher.com/kontainer-driver-metadata/${BRANCH}/data.json"
echo "Downloading data from $URL..."
if curl -s -o "$FILE_DOWNLOADED" "$URL"; then
  echo "Download complete: $FILE_DOWNLOADED"
else
  echo "Error: failed to download the data.json file"
  exit 1
fi

# Check for expected keys in downloaded file
if ! jq -e '.K8sVersionServiceOptions' "$FILE_DOWNLOADED" > /dev/null; then
  echo "Error: Key '.K8sVersionServiceOptions' not found in downloaded file $FILE_DOWNLOADED."
  exit 1
fi
if ! jq -e '.rke2.releases' "$FILE_DOWNLOADED" > /dev/null; then
   echo "Error: Key '.rke2.releases' not found in downloaded file $FILE_DOWNLOADED."
   exit 1
fi
if ! jq -e '.k3s.releases' "$FILE_DOWNLOADED" > /dev/null; then
   echo "Error: Key '.k3s.releases' not found in downloaded file $FILE_DOWNLOADED."
   exit 1
fi

echo "Starting comparison with branch: $BRANCH"

# Extract RKE1 versions
echo "Extracting RKE1 versions..."
VERSIONS_LOCAL=$(jq -r '.K8sVersionServiceOptions | keys[]' "$FILE_LOCAL" | sort)
VERSIONS_DOWNLOADED=$(jq -r '.K8sVersionServiceOptions | keys[]' "$FILE_DOWNLOADED" | sort)

# Extract RKE2 versions
echo "Extracting RKE2 versions..."
VERSIONS_RKE2_LOCAL=$(jq -r '.rke2.releases | .[] | .version' "$FILE_LOCAL" | sort)
VERSIONS_RKE2_DOWNLOADED=$(jq -r '.rke2.releases | .[] | .version' "$FILE_DOWNLOADED" | sort)

# Extract K3s versions
echo "Extracting K3s versions..."
VERSIONS_K3S_LOCAL=$(jq -r '.k3s.releases | .[] | .version' "$FILE_LOCAL" | sort)
VERSIONS_K3S_DOWNLOADED=$(jq -r '.k3s.releases | .[] | .version' "$FILE_DOWNLOADED" | sort)

# Compare versions and find missing ones
echo "Comparing versions..."
MISSING_RKE1_FROM_LOCAL=$(comm -13 <(echo "$VERSIONS_LOCAL") <(echo "$VERSIONS_DOWNLOADED"))
MISSING_RKE2_FROM_LOCAL=$(comm -13 <(echo "$VERSIONS_RKE2_LOCAL") <(echo "$VERSIONS_RKE2_DOWNLOADED"))
MISSING_K3S_FROM_LOCAL=$(comm -13 <(echo "$VERSIONS_K3S_LOCAL") <(echo "$VERSIONS_K3S_DOWNLOADED"))

# Output results
FOUND=0
echo "Comparison results:"
if [[ -n "$MISSING_RKE1_FROM_LOCAL" ]]; then
  FOUND=1
  echo "[Fail] The following RKE1 versions are in $BRANCH but missing from local data.json:"
  echo "$MISSING_RKE1_FROM_LOCAL"
else
  echo "[PASS] No RKE1 versions are missing in local data.json compared to $BRANCH."
fi

if [[ -n "$MISSING_RKE2_FROM_LOCAL" ]]; then
  FOUND=1
  echo "[Fail] The following RKE2 versions are in $BRANCH but missing from local data.json:"
  echo "$MISSING_RKE2_FROM_LOCAL"
else
  echo "[PASS] No RKE2 versions are missing in local data.json compared to $BRANCH."
fi

if [[ -n "$MISSING_K3S_FROM_LOCAL" ]]; then
  FOUND=1
  echo "[Fail] The following K3s versions are in $BRANCH but missing from local data.json:"
  echo "$MISSING_K3S_FROM_LOCAL"
else
  echo "[PASS] No K3s versions are missing in local data.json compared to $BRANCH."
fi

if [[ $FOUND -eq 1 ]]; then
  exit 1
else
  exit 0
fi
