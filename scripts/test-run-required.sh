#!/bin/bash
set -ex

echo "Checking if rancher integration testing is required"
echo "Environment variable DRONE_BUILD_EVENT is ${DRONE_BUILD_EVENT}"

if [ -z "$CI" ]; then
  echo "Not running in CI, rancher integration testing is required"
  exit 0
fi

if [ -z "$KDM_TEST_K8S_MINOR" ]; then
  echo "Error: KDM_TEST_K8S_MINOR not defined. This should not be happening in CI"
  exit 1
fi

if [ -z "$DRONE_COMMIT_BEFORE" ]; then
  echo "Error: DRONE_COMMIT_BEFORE not defined. This should not be happening in CI"
  exit 1
fi

# Only run check if Drone build event is 'push' or 'pull_request'
if [ "${DRONE_BUILD_EVENT}" = "push" ] || [ "${DRONE_BUILD_EVENT}" = "pull_request" ]; then
  # Check if the channels file contains changes to versions from the minor version
  if [ "$(git --no-pager diff --no-color -G "^  - version:" $DRONE_COMMIT_BEFORE -- "$CHANNELS_FILE" | grep -c -P "(^\+\s+- version: v1.$KDM_TEST_K8S_MINOR)")" -ne 0 ]; then
    exit 0
  fi
fi

echo "Skipping CI, no changes detected for relevant minor version"
exit 1
