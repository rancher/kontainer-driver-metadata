# KDM Release Process 

## Introduction

All metadata related to provisioning kubernetes clusters for Rancher supported distros (RKE1/RKE2/k3s) is present in [data.json](https://github.com/rancher/kontainer-driver-metadata/blob/dev-v2.7/data/data.json). Releasing KDM is publishing this data.json to `releases.rancher.com/` which is then accessed by Rancher to provision clusters. Refer to [rancher docs](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/upgrade-kubernetes-without-upgrading-rancher#configuring-the-metadata-synchronization) on how data refresh is managed by Rancher.
RKE1 binary also independently [embeds this data](https://github.com/rancher/rke/blob/v1.4.3/codegen/codegen.go#L13) to be able to provision RKE1 standalone clusters.

Releases are merged to `releases.rancher.com/` using our automated [drone logic](https://github.com/rancher/kontainer-driver-metadata/blob/dev-v2.7/.drone.yml#L26) from `release-v2.x` branches. Data at `releases.rancher.com/` is accessed by both existing rancher server and new rancher installs. `release-v2.x` branches are thus critical and protected, only few people are granted to merge PRs to `release-v2.x` branches.

KDM is released on rancher minor version basis, so multiple releases could be required depending on the support matrix and product requirements. Currently we are
managing releases for 
- Rancher v2.6.x (`release-v2.6` branch)
- Rancher v2.7.x (`release-v2.7` branch) 

Terms used often: 
1. Dev branch - the branch changes are merged from, of the format `dev-v2.x` or `dev-v2.x-for-rancher-x.y.z`. Dev branches change often according to release lifecycle and we could be maintaining multiple dev branches simultaneously. 
2. Release branch - the branch changes are merged to, of the format `release-v2.x`. 
3. Release/KDM url - url where KDM data is published and publicly accessed, of the format `https://releases.rancher.com/kontainer-driver-metadata/branch_name/data.json`. 

Examples: 

- Dev branch: `dev-v2.6`, `dev-v2.7-for-2.6.11` 
- Release branch: `release-v2.6`, `release-v2.7`
- Release/KDM url: https://releases.rancher.com/kontainer-driver-metadata/release-v2.6/data.json, https://releases.rancher.com/kontainer-driver-metadata/dev-v2.7/data.json

## Prepare PRs 

In order to trigger release process, PRs are directly merged from target dev branch to release branch. The first step is to prepare pull requests for release.

1. Open pull request using https://github.com/rancher/kontainer-driver-metadata/compare (base: `release-v2.x`, compare: target dev branch). 
2. Often pull requests can't be automatically merged by Github and requires rebasing off of `release-v2.x` branch. This is because there could be changes merged in `release-v2.x` branch from a different dev branch and we're now trying to merge from a different target dev branch. We need to ensure all required changes are present in our release PR, so this would require rebasing or resolving conflicts in github UI or opening a new PR to target dev branch. 
A couple of example PRs to help: 
- https://github.com/rancher/kontainer-driver-metadata/pull/1104
- https://github.com/rancher/kontainer-driver-metadata/pull/1056 
- https://github.com/rancher/kontainer-driver-metadata/pull/1027
- https://github.com/rancher/kontainer-driver-metadata/pull/1084 

### Prepare release notes for KDM (TBD)

## Review PRs 

There are a few common things to check while reviewing PRs for release. Most of this should already have been reviewed during the developer workflow. We do end up catching last minute bugs when reviewing PRs at release time so it's often a good idea to double check things.

RKE1: (Reviewers: @snasovich @kinarashah @jiaqiluo)
1. Kubernetes versions are correct according to the release plan
2. Rancher server information is correct
3. Default kubernetes version is correct 
4. Only one kubernetes version per minor version is released
5. Addon version information is correct
6. Service options are correct if updated for a patch version

RKE2/K3s: (Reviewers: @snasovich @rancher-max @kinarashah @jiaqiluo)
1. Kubernetes versions are correct according to the release plan (no RC versions)
2. Rancher server information is correct
3. Default kubernetes version is correct
4. Server and Agent arg information is correct
5. Chart information is correct
6. Only one kubernetes version per minor version is released

## Prepare draft release notes for RKE1

1. Draft a release note under https://github.com/rancher/rke/releases, it should have no tag attached and title would be "Draft Release vX.Y.Z" 
2. Copy paste the existing release notes from the latest tag 
3. Add major bug fixes and enhancements based on PRs present in diff between the last release and latest HEAD commit. Contact developers of PRs if more context is required. If the diff only includes KDM data, update enhancement section with the latest k8s versions about to be released.
4. List of kubernetes versions to be released is populated automatically on our RC tags under `RKE Kubernetes versions` (Example: https://github.com/rancher/rke/releases/tag/v1.4.4-rc7)

Note: This is our current process, but we also have https://github.com/rancherlabs/release-notes/tree/main/rke/ (which seems outdated at the moment). Skip this section if docs team is creating draft release notes.

## Release KDM

1. PRs to be merged by @kinarashah or @snasovich 
2. QA to perform post release checks for KDM after drone publish tasks are successfully completed (monitor at https://drone-publish.rancher.io/rancher/kontainer-driver-metadata)
3. Release RKE1 (refer to Release RKE1 section)
4. QA to perform post release checks for RKE after RKE tags are available 
5. Announce in [rancher forums](https://forums.rancher.com/c/announcements/) and slack. Example forums posts for reference: 
- https://forums.rancher.com/t/kubernetes-v1-24-10-and-v1-23-16/40218
- https://forums.rancher.com/t/kubernetes-v1-24-8-v1-23-14-and-v1-22-16/39559

## Release RKE1

We release RKE after KDM is released. Currently we are managing releases for 
- RKE v1.3.x (`release/v1.3`) corresponding to Rancher v2.6.x (https://github.com/rancher/rancher/blob/v2.6.11/go.mod#L121)
- RKE v1.4.x (`release/v1.4`) corresponding to Rancher v2.7.x (https://github.com/rancher/rancher/blob/v2.7.1/go.mod#L119)

1. Checkout at the right branch for RKE (`release/v1.x`). 
2. Update the target release branch in codegen/codegen.go (https://github.com/rancher/rke/blob/v1.4.3/codegen/codegen.go#L13)
3. Run `go generate`
4. Generated data/data.json should match with the latest changes. Open a pull request with the changes, one review should suffice. (Reviewers: @kinarashah @snasovich @jiaqiluo @HarrisonWAffel) 

Example PRs: 
- https://github.com/rancher/rke/pull/3193 
- https://github.com/rancher/rke/pull/3194 
5. After the PR is merged, checkout the branch at latest `release/v1.x` 
6. Tag and push,

   ```
      git remote add rancher https://github.com/rancher/rke.git
      git tag vX.Y.Z
      git push rancher vX.Y.Z
   ```
7. Monitor drone status for tag (https://drone-publish.rancher.io/rancher/rke). Drone uploads the release manifests to the release tag under Assets (https://github.com/rancher/rke/releases) 
8. Copy release notes from our draft release notes to the released tag 
9. Mark the patch version for the latest minor version (`v1.4.x` if releasing both `v1.3.x` and `v1.4.x` as the latest version)
