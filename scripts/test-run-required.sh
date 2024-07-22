#!/bin/bash
set -ex

echo "Checking if rancher integration testing is required"
echo "Environment variable GITHUB_EVENT_NAME is ${GITHUB_EVENT_NAME}"

if [ -z "$CI" ]; then
  echo "Not running in CI, rancher integration testing is required"
  exit 0
fi

if [ -z "$KDM_TEST_K8S_MINOR" ]; then
  echo "Error: KDM_TEST_K8S_MINOR not defined. This should not be happening in CI"
  exit 1
fi

if [ -z "$PREV_COMMIT_SHA" ]; then
  echo "Error: PREV_COMMIT_SHA not defined. This should not be happening in CI"
  exit 1
fi

# Only run check if Github build event is 'push' or 'pull_request'
if [ "${GITHUB_EVENT_NAME}" = "push" ] || [ "${GITHUB_EVENT_NAME}" = "pull_request" ]; then
  # Check if the channels file contains changes to versions from the minor version
  if [ "$(git --no-pager diff --no-color -G "^ - version:" $PREV_COMMIT_SHA -- "$CHANNELS_FILE" | grep -c -P "(^\+\s+- version: v1.$KDM_TEST_K8S_MINOR)")" -ne 0 ]; then
    exit 0
  fi
fi

echo "Skipping CI, no changes detected for relevant minor version"
exit 1
