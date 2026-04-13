REPO_PATH := /go/src/github.com/rancher/kontainer-driver-metadata
CI_IMAGE  := kdm-ci
ARCH      ?= $(shell go env GOHOSTARCH 2>/dev/null || echo amd64)

# Environment variables forwarded into the container
FORWARD_ENVS := \
    REPO TAG CI \
    LAST_COMMUNITY_RANCHER CATTLE_AGENT_IMAGE \
    PREV_COMMIT_PR_SHA PREV_COMMIT_PUSH_SHA \
    GITHUB_EVENT_NAME GITHUB_RUN_NUMBER GITHUB_REF_TYPE GITHUB_REF_NAME \
    STAGE_REGISTRY_ENDPOINT REGISTRY_ENDPOINT REGISTRY_USERNAME REGISTRY_PASSWORD \
    V2PROV_TEST_DIST V2PROV_TEST_RUN_REGEX KDM_TEST_K8S_MINOR DEBUG

# Only forward variables that are actually set in the environment
DOCKER_ENV_FLAGS = $(foreach v,$(FORWARD_ENVS),$(if $(value $(v)),--env '$(v)=$(value $(v))'))

TARGETS := $(shell ls scripts)

.ci-image:
	docker build \
	    -f Dockerfile.ci \
	    --build-arg TARGETARCH=$(ARCH) \
	    -t $(CI_IMAGE) \
	    .

$(TARGETS): .ci-image
	docker run --rm \
	    -v $$(pwd):$(REPO_PATH) \
	    -v /var/run/docker.sock:/var/run/docker.sock \
	    -v $$(pwd)/build:/tmp \
	    --privileged \
	    $(DOCKER_ENV_FLAGS) \
	    $(CI_IMAGE) $@

.DEFAULT_GOAL := ci

.PHONY: $(TARGETS) .ci-image
