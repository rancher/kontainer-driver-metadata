#!/usr/bin/env bash
RELEASES="release-v2.5"

for RELEASE in $RELEASES; do
    echo "Images check for released Rancher k8s versions for ${RELEASE}"

    RELEASE_URL="https://releases.rancher.com/kontainer-driver-metadata/${RELEASE}/data.json"

    # Released data file
    RELEASEDATAFILE=$(curl --connect-timeout 60 --max-time 60 -s ${RELEASE_URL})

    # Build released versions checksums file
    echo "Creating kdm-released-${RELEASE}.txt from https://releases.rancher.com/kontainer-driver-metadata/${RELEASE}/data.json"
    for K8SVERSION in $(echo $RELEASEDATAFILE | jq -r '.K8sVersionRKESystemImages | keys[]'); do
        CHECKSUM=$(echo $RELEASEDATAFILE | jq -r '.K8sVersionRKESystemImages["'"$K8SVERSION"'"]' | sha256sum | awk '{ print $1 }')
        echo "${K8SVERSION} ${CHECKSUM}" >> "kdm-released-${RELEASE}.txt"
    done

    # Compare to current data
    echo "Comparing images for released versions in kdm-released-${RELEASE}.txt to data/data.json"
    DEVDATAFILE=$(cat data/data.json)

    while read -r K8SVERSION DATACHECKSUM; do
        DEVCHECKSUM=$(echo $DEVDATAFILE | jq -r '.K8sVersionRKESystemImages["'"$K8SVERSION"'"]' | sha256sum | awk '{ print $1 }')
        if [ "${DEVCHECKSUM}" != "${DATACHECKSUM}" ]; then
            echo "Images checksum for released k8s version ${K8SVERSION} not equal (data.json: ${DEVCHECKSUM}) vs (kdm-released-${RELEASE}.txt: ${DATACHECKSUM})"
            echo "We cannot change images in a released Rancher k8s version, please create a new Rancher k8s version to make the change (e.g. create v1.20.9-rancher1-2 with new images instead of changing v1.20.9-rancher1-1)"
            exit 1
        fi
    done < "kdm-released-${RELEASE}.txt"
done
