# KDM Release Process 

## Introduction

All metadata related to provisioning kubernetes clusters for Rancher supported distros (RKE2/k3s) is present in [data.json](https://github.com/rancher/kontainer-driver-metadata/blob/dev-v2.7/data/data.json). Releasing KDM is publishing this data.json to `releases.rancher.com/` which is then accessed by Rancher to provision clusters. Refer to [rancher docs](https://ranchermanager.docs.rancher.com/getting-started/installation-and-upgrade/upgrade-kubernetes-without-upgrading-rancher#configuring-the-metadata-synchronization) on how data refresh is managed by Rancher.
RKE1 binary also independently [embeds this data](https://github.com/rancher/rke/blob/v1.4.3/codegen/codegen.go#L13) to be able to provision RKE1 standalone clusters.

Releases are merged to `releases.rancher.com/` using our automated [GHA logic](https://github.com/rancher/kontainer-driver-metadata/blob/dev-v2.10/.github/workflows/workflow.yaml#L79) from `release-v2.x` branches. Data at `releases.rancher.com/` is accessed by both existing rancher server and new rancher installs. `release-v2.x` branches are thus critical and protected, only few people are granted to merge PRs to `release-v2.x` branches.

KDM is released on a rancher minor version basis, so multiple releases could be required depending on the support matrix and product requirements. Currently, we are
managing releases for
- Rancher v2.9.x (`release-v2.9` branch)
- Rancher v2.10.x (`release-v2.10` branch)
- Rancher v2.11.x (`release-v2.11` branch)
- Rancher v2.12.x (`release-v2.12` branch)

Terms used often: 
1. Dev branch - the branch changes are merged from, of the format `dev-v2.x` or `dev-v2.x-for-rancher-x.y.z`. Dev branches change often according to release lifecycle, and we could be maintaining multiple dev branches simultaneously. 
2. Release branch - the branch changes are merged to, of the format `release-v2.x`. 
3. Release/KDM url - url where KDM data is published and publicly accessed, of the format `https://releases.rancher.com/kontainer-driver-metadata/branch_name/data.json`. 
4. Release Milestone - a GitHub milestone within the [rancher/rancher repository](https://github.com/rancher/rancher) which is used to track issues that are included in KDM releases. All issues in a KDM milestone must be resolved (Tested, Closed out, etc) before a release.

Examples: 

- Dev branch: `dev-v2.6`, `dev-v2.7-for-2.6.11` 
- Release branch: `release-v2.6`, `release-v2.7`
- Release/KDM url: https://releases.rancher.com/kontainer-driver-metadata/release-v2.6/data.json, https://releases.rancher.com/kontainer-driver-metadata/dev-v2.7/data.json
- Release Milestone: `v2.7-KDM-Aug-2023-patches`, `v2.6-KDM-Aug-2023-patches`

## Review Issues in Current Milestone

Before any pull requests are opened, be sure to double-check the relevant release milestone for the upcoming KDM release. Ensure that all issues in the current milestone have been resolved. The release captain should raise concerns regarding any open issues or bugs, and ensure they are closed out before the release process begins.

## Prepare PRs 

### Preparing The Release For a New Rancher Minor Version

When KDM is being released for a new minor version of Rancher (e.g. `v2.10.0`, `v2.11.0`, etc) a few additional checks must be made.

+ Ensure that all provisioning tests are utilizing the correct version of Rancher and the correct Kubernetes version range. Ideally this will have been done during the creation of the dev branch, however this should be double checked to ensure that CI functions properly upon release
   + Example: https://github.com/rancher/kontainer-driver-metadata/pull/1553
+ The `release-v2.x` branch will not exist yet. To ensure that all versions are correctly included, raise a preview PR against the prior release branch (e.g `dev-v2.10` against `release-v2.9`). **Do not merge this PR**
   + Example: https://github.com/rancher/kontainer-driver-metadata/pull/1552

Once these checks are complete, creating the KDM release for the new minor version is as simple as creating the new `release-v2.x` branch based off of the corresponding `dev-v2.x` branch and pushing it to Github.

## Opening The Release PRs

In order to trigger the release process, PRs are directly merged from the target dev branch to release branch. PRs raised from forks should not be merged into the release branches. The first step is to prepare pull requests for release.

1. Open pull request using https://github.com/rancher/kontainer-driver-metadata/compare (base: `release-v2.x`, compare: target dev branch). 
2. Often pull requests can't be automatically merged by Github and requires rebasing off of `release-v2.x` branch. This is because there could be changes merged in `release-v2.x` branch from a different dev branch, and we're now trying to merge from a different target dev branch. We need to ensure all required changes are present in our release PR, so this would require rebasing or resolving conflicts in github UI or opening a new PR to target dev branch. 
A couple of example PRs to help: 
- https://github.com/rancher/kontainer-driver-metadata/pull/1104
- https://github.com/rancher/kontainer-driver-metadata/pull/1056 
- https://github.com/rancher/kontainer-driver-metadata/pull/1027
- https://github.com/rancher/kontainer-driver-metadata/pull/1084

## Prime configuration variables

When preparing a new KDM release line (or advancing an existing one), update the prime-related variables used by CI/test scripts to correctly detect prime and select a prime-published agent image:

```yaml
# .github/workflows/provisioning-tests.yaml
- name: Provisioning Operations tests
  run: dapper provisioning-tests
  env:
    LAST_COMMUNITY_RANCHER: v2.9.3
    # ...existing env...
```

**What to update each release:**

* **`LAST_COMMUNITY_RANCHER`** — Set to the cutoff Rancher server version for *community* builds on this line. This is compared to a channel entry’s `minChannelServerVersion` to decide if a K8s release is **prime-only**.

**Where these are used:** KDM CI/test scripts (e.g., `scripts/prime-route`, `scripts/provisioning-tests`) referenced by Rancher provisioning tests.

**When creating a new KDM branch:** if there is no prime release yet for this Rancher minor, set `LAST_COMMUNITY_RANCHER` to a high placeholder similar to how `maxChannelServerVersion` is used in `data.json` (for example `v2.X.99`). Replace this placeholder with the agreed cutoff once the prime-only point is defined for the line.

**When to bump:**

* At the start of a new release line or when the prime cutoff moves (bump `LAST_COMMUNITY_RANCHER`).

### During prime cutoff: pin distro versions in rancher/rancher

While a release line is at or near the prime cutoff, CI in **rancher/rancher** should **pin** the Kubernetes versions per distro to avoid unintentionally selecting a prime-only “latest”. Add the per-distro pins in the Rancher workflow’s **Run tests** step:

```yaml
# rancher/rancher: .github/workflows/provisioning-tests.yml (excerpt)
- name: Run tests
  run: ./.dapper provisioning-tests
  env:
    # per-distro pins used by the provisioning script if SOME_K8S_VERSION is unset
    K3S_PINNED_VERSION: "v1.30.14+k3s2"
    RKE2_PINNED_VERSION: "v1.30.14+rke2r2"
    # ...existing env...
```

[Example reference: the change in this commit shows the pins added in the Rancher workflow](https://github.com/rancher/rancher/pull/52176/commits/dcf3704ca888037cdf37513a04eaeb1d07cb3a70#diff-3b900f349b5944ff0314a1b0597aff3c4b531758380e00bd0e3ece5844fead00R57).

## Prime images update

The KDM release process will automatically sync the newly released images into the rancher prime registry once a release occurs.

### Prepare release notes for KDM (TBD)

## Review KDM PRs 

There are a few common things to check while reviewing PRs for release. Most of this should already have been reviewed during the developer workflow. We do end up catching last minute bugs when reviewing PRs at release time so it's often a good idea to double check things.

RKE2/K3s: (Reviewers: @ShylajaDevadiga @snasovich)
1. Kubernetes versions are correct according to the release plan (no RC versions)
2. Rancher server information is correct
3. Default kubernetes version is correct
4. Server and Agent arg information is correct
5. Chart information is correct

## Release KDM

1. PRs to be merged by @kinarashah or @snasovich 
2. Once merged, look at the Github Action pipeline and ensure that all prime images have been successfully pushed to the trusted registry (monitor at https://github.com/rancher/kontainer-driver-metadata/actions/workflows/workflow.yaml) 
3. Ping QA in Slack to perform post release checks for KDM once the GHA tasks are successfully completed 
4. QA to perform post sync checks for Rancher prime after GHA publish tasks are successfully completed
5. Release RKE1 (refer to Release RKE1 section)
6. QA to perform post release checks for RKE after RKE tags are available 
7. Announce in [rancher forums](https://forums.rancher.com/c/announcements/) and slack (use the `team-rancher-release-coordination` channel). Example forums posts for reference: 
- https://forums.rancher.com/t/kubernetes-v1-24-10-and-v1-23-16/40218
- https://forums.rancher.com/t/kubernetes-v1-24-8-v1-23-14-and-v1-22-16/39559


## Post OOB Release Tasks for KDM

Some OOB releases utilize branches other than the standard dev-v2.x branches so that others are not blocked prior to a KDM release. In these cases, the primary dev branch needs to be updated post release. 
This can be done by either merging the secondary dev branch into the primary dev branch, or by merging the release branch into the primary dev branch. Ensure that commits are fully sync'd and that no merge commits exist in the release branch that are not included in the standard dev branches. 

Here are a few example PR's for this process. While typically one review will suffice, double check with @snasovich or @kinarashah before merging, as this may impact the work of others.  
+ https://github.com/rancher/kontainer-driver-metadata/pull/1167
+ https://github.com/rancher/kontainer-driver-metadata/pull/1131

## Prepare KDM For A Rancher/Rancher Release 

All officially released versions of Rancher are expected to use the KDM release branches. Rancher references a particular KDM branch in a number of places within its code base. A majority of the time, Rancher will be referencing a KDM development branch so that developers do not have to update any environment variables to work with the latest KDM data. However, when preparing for a Rancher release, this reference **must** be updated to point at a release branch - otherwise development work done within KDM will negatively impact users of Rancher. 

This step typically occurs during the preparation of the final RC for a rancher/rancher release, or immediately after a KDM release for a given version of Rancher (i.e. _not_ OOB releases). A final RC should **never** use KDMs development branches.  

Here's an example PR for this change: 
+ https://github.com/rancher/rancher/pull/47654

## Post Rancher Release Tasks

Once Rancher has been released, make sure to revert the KDM branches to the corresponding dev branches. 
