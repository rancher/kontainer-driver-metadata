#!/bin/bash
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 branch"
  exit 1
fi

BRANCH=$1

DEV_URL="https://releases.rancher.com/kontainer-driver-metadata/${BRANCH}/data.json"
echo "Using URL ${DEV_URL} to compare k8s version(s)"

DEVDATAFILE=$(curl --connect-timeout 60 --max-time 60 -s ${DEV_URL})

PRDATAFILE=$(cat data/data.json)

NEWK8SVERSIONS=$(diff --changed-group-format='%<%>' --unchanged-group-format='' <(echo $PRDATAFILE | jq -r '.K8sVersionRKESystemImages | keys[]') <(echo $DEVDATAFILE | jq -r '.K8sVersionRKESystemImages | keys[]'))

if [[ -z $NEWK8SVERSIONS ]]; then
  echo "No new k8s version found in PR's data/data.json"
else
  echo "New k8s version(s) found in PR's data/data.json:"
  echo $NEWK8SVERSIONS

  for NEWK8SVERSION in $NEWK8SVERSIONS; do
    echo "Getting previous version for $NEWK8SVERSION"
    PREVK8SVERSION=$(echo $PRDATAFILE | jq -r '.K8sVersionRKESystemImages | keys[]' | sort -V | grep -B1 $NEWK8SVERSION | head -1)
    echo "Previous version for $NEWK8SVERSION is: $PREVK8SVERSION"
    echo "Diffing previous ($PREVK8SVERSION) and new version ($NEWK8SVERSION)"
    diff <(echo $PRDATAFILE| jq -r --arg PREVK8SVERSION "$PREVK8SVERSION" '.K8sVersionRKESystemImages[$PREVK8SVERSION]') <(echo $PRDATAFILE | jq -r --arg NEWK8SVERSION "$NEWK8SVERSION" '.K8sVersionRKESystemImages[$NEWK8SVERSION]')
  done
fi
