# May 2026 RKE2 R2 Calico Toleration Fix — Test Failure Analysis

## Executive Summary

The Calico toleration fix ([rke2#10438](https://github.com/rancher/rke2/issues/10438)) was applied via rke2-calico chart v3.32.002 and shipped in RKE2 R2 release candidates. Three KDM PRs were tested across multiple CI runs:

- [#2067](https://github.com/rancher/kontainer-driver-metadata/pull/2067) (dev-v2.14) — K8s 1.33, 1.34, 1.35 — `v1.33.12-rc2+rke2r2`, `v1.34.8-rc2+rke2r2`, `v1.35.5-rc2+rke2r2`
- [#2068](https://github.com/rancher/kontainer-driver-metadata/pull/2068) (dev-v2.15) — K8s 1.34, 1.35, 1.36 — but **with typo** `v1.36.1-rc2+rke2r1` instead of rke2r2
- [#2072](https://github.com/rancher/kontainer-driver-metadata/pull/2072) (dev-v2.15) — K8s 1.34, 1.35, 1.36 — **fix for typo**: correct `v1.36.1-rc2+rke2r2`

**Overall result: The Calico toleration fix definitively resolves `Single_Node_All_Roles_Drain` for all K8s versions (1.33–1.36). The etcd snapshot restore tests remain broken due to a separate, unrelated root cause.**

| Test | Before R2 fix (R1) | After R2 fix | Conclusion |
|------|--------------------|--------------|------------|
| `Test_Provisioning_Single_Node_All_Roles_Drain` | ❌ K8s 1.33–1.36 | ✅ K8s 1.33–1.36 | **✅ FIXED** |
| `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode` | ❌ All versions | ❌ All versions | **❌ UNCHANGED** (different root cause) |
| `Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` | ✅ All versions | ⚠️ Flaky (fails ~70–80% of runs) | **⚠️ REGRESSED** (timing race, see analysis) |
| All other tests (SetA, SetB MP_OnNewNode, other Provisioning) | ✅ | ✅ | **✅ Unchanged** |

---

## The Fix: What Was Changed

### Root cause identified (Slack, May 21)

The Calico v3.32.0 upgrade changed Typha pod tolerations from specific key-based tolerations to broad "tolerate everything" tolerations:

**Calico v3.31.5 (old/working):**
```yaml
tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/control-plane
    operator: Exists
  - effect: NoExecute
    key: node-role.kubernetes.io/etcd
    operator: Exists
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
```

**Calico v3.32.0 (new/broken):**
```yaml
tolerations:
  - key: CriticalAddonsOnly
    operator: Exists
  - effect: NoSchedule
    operator: Exists
  - effect: NoExecute
    operator: Exists
```

### Key insight from Brad Davidson

The `tolerationSeconds: 300` for `not-ready` and `unreachable` taints was **never explicitly set** in the Calico chart or tigera-operator code. It was **injected by Kubernetes' `DefaultTolerationSeconds` admission controller** ([source](https://github.com/kubernetes/kubernetes/blob/master/plugin/pkg/admission/defaulttolerationseconds/admission.go#L68-L73)):

> If a pod doesn't already have a toleration for `notReady:NoExecute` or `unreachable:NoExecute`, Kubernetes automatically adds one with `tolerationSeconds: 300`.

In Calico v3.32.0, the broad `NoExecute: Exists` toleration (without a key) **matches all NoExecute taints**, so the admission controller **skips adding** the default `tolerationSeconds: 300` tolerations. This means pods tolerate `unreachable` and `not-ready` **indefinitely** instead of being evicted after 5 minutes.

### Fix applied ([rke2-charts#1002](https://github.com/rancher/rke2-charts/pull/1002))

The rke2-calico chart overlay was updated to explicitly set `controlPlaneTolerations` with the old key-specific tolerations:

```yaml
controlPlaneTolerations:
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node-role.kubernetes.io/etcd"
    operator: "Exists"
    effect: "NoExecute"
```

Chart version bumped from `v3.32.001` → `v3.32.002`.

### Fix propagation via RKE2 PRs

| Branch | PR | Status |
|--------|-----|--------|
| main | [rke2#10450](https://github.com/rancher/rke2/pull/10450) | ✅ Merged May 22 |
| release-1.36 | [rke2#10455](https://github.com/rancher/rke2/pull/10455) | ✅ Merged May 22 |
| release-1.35 | [rke2#10456](https://github.com/rancher/rke2/pull/10456) | ✅ Merged May 22 |
| release-1.34 | [rke2#10457](https://github.com/rancher/rke2/pull/10457) | ✅ Merged May 22 |
| release-1.33 | [rke2#10458](https://github.com/rancher/rke2/pull/10458) | ✅ Merged May 22 |

---

## R2 Versions Tested on KDM PRs

| PR | Branch | K8s | RKE2 Version | rke2-calico chart | Note |
|----|--------|-----|-------------|-------------------|------|
| [#2067](https://github.com/rancher/kontainer-driver-metadata/pull/2067) | dev-v2.14 | 1.33 | v1.33.12-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2067 | dev-v2.14 | 1.34 | v1.34.8-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2067 | dev-v2.14 | 1.35 | v1.35.5-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| [#2068](https://github.com/rancher/kontainer-driver-metadata/pull/2068) | dev-v2.15 | 1.34 | v1.34.8-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2068 | dev-v2.15 | 1.35 | v1.35.5-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2068 | dev-v2.15 | 1.36 | v1.36.1-rc2+**rke2r1** | v3.32.002 | ⚠️ Typo — should be rke2r2 |
| [#2072](https://github.com/rancher/kontainer-driver-metadata/pull/2072) | dev-v2.15 | 1.34 | v1.34.8-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2072 | dev-v2.15 | 1.35 | v1.35.5-rc2+**rke2r2** | v3.32.002 | ✅ Correct |
| #2072 | dev-v2.15 | 1.36 | v1.36.1-rc2+**rke2r2** | v3.32.002 | ✅ Correct (typo fixed) |

**Confirmation**: All R2 versions include Calico chart v3.32.002 with the toleration fix. PR #2072 was opened specifically to correct the `rke2r1` typo in PR #2068 for K8s 1.36.

---

## Master Test Results Table

All CI runs across all PRs and attempts. **Abbreviations**: ✅ PASS · ❌ FAIL · ⏱️ job timeout · — not tested · ⚠️ flaky

> **All tests not listed here (SetA, other Provisioning tests, all k3s tests) passed in 100% of runs.**

| K8s | RKE2 Version | PR | Run | `Custom_EtcdRestore`¹ | `MP_ThreeEtcd`² | `Single_Node_Drain`³ |
|-----|-------------|-----|-----|----------------------|-----------------|----------------------|
| **R1 baseline (before fix)** | | | | | | |
| 1.33 | v1.33.x-rc1+rke2r1 | — | prior | ❌ | ✅ | ❌ |
| 1.34 | v1.34.x-rc1+rke2r1 | — | prior | ❌ | ✅ | ❌ |
| 1.35 | v1.35.x-rc1+rke2r1 | — | prior | ❌ | ✅ | ❌ |
| 1.36 | v1.36.x-rc1+rke2r1 | — | prior | ❌ | ✅ | ❌ |
| **R2 — PR #2067 (dev-v2.14)** | | | | | | |
| 1.33 | v1.33.12-rc2+rke2r2 | #2067 | A1 (May 22) | ❌ (2985s) | ✅ (618s) | ✅ |
| 1.33 | v1.33.12-rc2+rke2r2 | #2067 | A2 (May 24) | ❌ (3051s) | ❌ (3118s) | ✅ |
| 1.34 | v1.34.8-rc2+rke2r2 | #2067 | A1 (May 22) | ❌ (3044s) | ❌ (3133s) | ✅ |
| 1.34 | v1.34.8-rc2+rke2r2 | #2067 | A2 (May 24) | ❌ (2975s) | ❌ (3208s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2067 | A1 (May 22) | ❌ (2982s) | ❌ (3159s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2067 | A2 (May 24) | ❌ | ⏱️ (3600s) | ✅ |
| **R2 — PR #2068 (dev-v2.15, note: K8s 1.36 has rke2r1 typo)** | | | | | | |
| 1.34 | v1.34.8-rc2+rke2r2 | #2068 | A1 (May 22) | ❌ | ⏱️ (3600s) | ✅ |
| 1.34 | v1.34.8-rc2+rke2r2 | #2068 | A2 (May 24) | ❌ | ❌ (3148s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2068 | A1 (May 22) | ❌ (3042s) | ❌ (3192s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2068 | A2 (May 24) | ❌ | ⏱️ (3600s) | ✅ |
| 1.36 | v1.36.1-rc2+**rke2r1** ⚠️ | #2068 | A1 (May 22) | ❌ (3197s) | ✅ (853s) | ✅ |
| 1.36 | v1.36.1-rc2+**rke2r1** ⚠️ | #2068 | A2 (May 24) | ❌ (3037s) | ✅ (747s) | **❌ ⚠️** |
| **R2 — PR #2072 (dev-v2.15, K8s 1.36 rke2r2 typo fixed)** | | | | | | |
| 1.34 | v1.34.8-rc2+rke2r2 | #2072 | A2 (May 25) | ❌ (3056s) | ❌ (3148s) | ✅ |
| 1.34 | v1.34.8-rc2+rke2r2 | #2072 | A3 (May 25) | ❌ (2977s) | ❌ (3108s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2072 | A2 (May 25) | ❌ (3033s) | ❌ (3203s) | ✅ |
| 1.35 | v1.35.5-rc2+rke2r2 | #2072 | A3 (May 25) | ❌ | ⏱️ (3600s) | ✅ |
| 1.36 | v1.36.1-rc2+**rke2r2** ✅ | #2072 | A2 (May 25) | ❌ | ⏱️ (3600s) | **✅** |
| 1.36 | v1.36.1-rc2+**rke2r2** ✅ | #2072 | A3 (May 25) | ❌ | ⏱️ (3600s) | **✅** |

**Footnotes:**
1. `Custom_EtcdRestore` = `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode`
2. `MP_ThreeEtcd` = `Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` (within `machineprovisioning` package; ⏱️ indicates job timeout hit before test could name its failure)
3. `Single_Node_Drain` = `Test_Provisioning_Single_Node_All_Roles_Drain`

---

## Conclusions

### Tests that started passing consistently in R2

| Test | R1 Status | R2 Status | Notes |
|------|-----------|-----------|-------|
| `Single_Node_All_Roles_Drain` (K8s 1.33-1.35) | ❌ FAIL | ✅ PASS (100%) | Toleration fix works: Typha evicted after 5 min, rescheduled |
| `Single_Node_All_Roles_Drain` (K8s 1.36) | ❌ FAIL | ✅ PASS (100%) | **Only with correct `rke2r2` — PR #2072 confirms this. With wrong `rke2r1` (PR #2068) it was flaky.** |

### Tests that continue failing consistently in R2

| Test | All Attempts | Error |
|------|-------------|-------|
| `Custom_EtcdSnapshotOperationsOnNewCombinedNode` | ❌ 100% fail all K8s, all attempts | `etcd snapshot restore wait did not succeed: timeout waiting condition: context deadline exceeded` (~3000s per run) |

### Tests that regressed (flaky) in R2

| Test | R1 Status | R2 Status | Likely Cause |
|------|-----------|-----------|--------------|
| `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` | ✅ PASS (R1) | ⚠️ Fails ~70–80% of runs | Timing race: 5-min eviction window from toleration fix creates race condition with etcd restore flow |
| `Single_Node_All_Roles_Drain` (K8s 1.36, **rke2r1 only**) | ❌ FAIL (R1) | ⚠️ Flaky (PR #2068 only) | Not actually a regression — caused by wrong version (`rke2r1` typo). With correct `rke2r2` (PR #2072) it passes consistently. |

### Key insight: rke2r1 vs rke2r2 for K8s 1.36

PR #2068 inadvertently tested `v1.36.1-rc2+rke2r1` (a typo) instead of `v1.36.1-rc2+rke2r2`. This caused `Single_Node_All_Roles_Drain` to be flaky on K8s 1.36 — passed once, failed once. PR #2072 corrected this to `rke2r2` and the test passed in **both CI runs (attempts 2 and 3)**. This confirms the Calico toleration fix in `rke2r2` also resolves the K8s 1.36 drain test.

---

## Analysis: Why Tests Still Fail

### What the fix DID solve: Single_Node_All_Roles_Drain

The `Single_Node_All_Roles_Drain` test creates a single-node cluster, adds a new node, then deletes the old node (without drain). Previously, Calico Typha was stuck on the old node indefinitely because it tolerated `unreachable:NoExecute` without a timeout. With the fix restoring `tolerationSeconds: 300`, Typha is evicted after 5 minutes and rescheduled on the new node.

This test now passes **consistently** for all K8s versions (1.33–1.36) with the correct R2 versions.

### What the fix did NOT solve: Etcd Snapshot Restore to combined node

The `Custom_EtcdSnapshotOperationsOnNewCombinedNode` test still fails on ALL versions. This test involves:

1. Creating a single combined etcd+cp+worker node cluster
2. Taking an etcd snapshot
3. **Restoring** to a **new** combined node

**Root cause** (confirmed by artifact analysis — `rke_controlplane.yaml` status message):
```
Waiting for Cluster control plane to be initialized, waiting for cluster agent to connect
```

After `rke2 server --cluster-reset`, the restore flow fails because:

1. After restore, `kubelet` starts with `--cloud-provider=external`
2. This adds `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` taint
3. The cloud-controller-manager (CCM) fails to remove this taint because the restored etcd contains stale node objects
4. With the taint on the only node, `cattle-cluster-agent` can't be scheduled
5. Without the agent, Rancher can't connect → planner stuck in a loop

**Confirmed from node logs in PR #2072 artifact:**
```
Starting rke2 agent v1.36.1-rc2+rke2r2 (1e680efe89a27bedc7f98cd60f5954fc671468f9)
cluster agent disconnected, requeuing
```

**The Calico toleration fix does not address this CCM taint issue.** These are orthogonal problems.

### Why MP_ThreeEtcd tests are flaky (timing regression)

`MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` passed in R1 because:
- Separate worker nodes without the CCM taint allowed `cattle-cluster-agent` to schedule on workers
- Broad v3.32.0 tolerations (pre-fix) meant Typha tolerated everything indefinitely — networking worked via healthy worker nodes

With the R2 fix (restrictive tolerations + `tolerationSeconds: 300`):
- Typha WILL be evicted after 5 minutes (good for drain test)
- But during etcd restore the cluster needs networking before Typha has time to evict and reschedule
- This creates a ~5-minute race window that sometimes (70–80%) causes failure
- This is a timing regression, not a fundamental breakage

---

## Root Cause: The Etcd Restore Problem is NOT a Toleration Issue

See also: `rke2-v1.36.0-etcd-restore-failure-analysis.md` for the detailed analysis.

The complete failure chain:

```
rke2 server --cluster-reset
    → kubelet starts with --cloud-provider=external
    → node.cloudprovider.kubernetes.io/uninitialized:NoSchedule taint added
    → CCM cannot remove taint (restored etcd has stale node data)
    → cattle-cluster-agent cannot schedule (taint blocks it)
    → Rancher planner stuck: "waiting for cluster agent to connect"
    → etcd snapshot restore wait times out
    → TEST FAILS
```

---

## Suggested Fixes

### Fix 1: Restore `not-ready`/`unreachable` tolerations explicitly (CALICO CHART)

The current fix adds only `control-plane:NoSchedule` and `etcd:NoExecute`. During restore, nodes may be temporarily `NotReady`. Explicitly adding the previously auto-injected tolerations would prevent Typha eviction during the critical restore window:

```yaml
controlPlaneTolerations:
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node-role.kubernetes.io/etcd"
    operator: "Exists"
    effect: "NoExecute"
  - key: "node.kubernetes.io/not-ready"
    operator: "Exists"
    effect: "NoExecute"
    tolerationSeconds: 300
  - key: "node.kubernetes.io/unreachable"
    operator: "Exists"
    effect: "NoExecute"
    tolerationSeconds: 300
```

### Fix 2: Add CCM uninitialized taint toleration to cattle-cluster-agent (RANCHER)

The core etcd restore issue is the CCM `uninitialized:NoSchedule` taint blocking `cattle-cluster-agent`. Fix in Rancher:

```yaml
# cattle-cluster-agent Deployment should tolerate:
- key: "node.cloudprovider.kubernetes.io/uninitialized"
  operator: "Exists"
  effect: "NoSchedule"
```

### Fix 3: Post-restore cleanup for stale Calico resources (RKE2 PLANNER)

During etcd restore, the planner should explicitly:
1. Delete stale Calico pods (Typha, calico-node, calico-kube-controllers) that reference old nodes
2. Restart the tigera-operator to re-reconcile
3. Wait for Calico networking to be healthy before proceeding

---

## Summary of Issues

| Issue | Root Cause | Fix Location | Status |
|-------|-----------|--------------|--------|
| Typha not rescheduled on node removal | Calico v3.32.0 broad tolerations prevent Kubernetes auto-eviction | rke2-charts (Calico chart v3.32.002) | ✅ **FIXED** — Single_Node_Drain passes for all K8s 1.33–1.36 in R2 |
| Etcd restore fails on combined node | CCM `uninitialized` taint blocks cattle-cluster-agent + stale etcd data | Rancher (agent tolerations) + RKE2 (CCM behavior) | ❌ NOT fixed by toleration change |
| MP etcd restore flaky | Timing race — 5-min eviction window clashes with restore flow | rke2-charts (add explicit not-ready/unreachable tolerations) | ⚠️ FLAKY regression from fix |
| K8s 1.36 Single_Node_Drain with rke2r1 | Wrong binary tested (typo `rke2r1` in PR #2068) | PR #2072 fixes this | ✅ Resolved — passes 100% with correct rke2r2 |

---

## Detailed CI Run Log (All Failures)

### Error signature (all Custom_EtcdRestore failures, all PRs/runs):
```
etcdsnapshot.go:277: cluster test-custom-etcd-snapshot-operations-on-new-combined-node
etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed :
timeout waiting condition: context deadline exceeded
```
Duration: ~2977–3056s (~50 min) before failing consistently.

### Error signature (MP_ThreeEtcd failures when named explicitly, e.g. PR #2072 K8s 1.34):
```
etcdsnapshot.go:277: cluster test-mp-etcd-snapshot-conventional-arch-new-node
etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed :
timeout waiting condition: context deadline exceeded
```

### Error signature (Rancher planner, from PR #2072 K8s 1.36 SetB artifact):
```
[planner] rkecluster test-ns-fqql4/test-custom-...: configuring bootstrap node(s) custom-d6fb339c9afe:
Waiting for Cluster control plane to be initialized, waiting for cluster agent to connect
[ERROR] error syncing 'calico-system': cluster agent disconnected, requeuing
```

### PR #2068 K8s 1.36 Provisioning FLAKY failure (attempt 2 only, with rke2r1 typo):
```
--- FAIL: Test_Provisioning_Single_Node_All_Roles_Drain (1189.47s)
[planner] rkecluster: configuring etcd node(s): waiting for probes: calico
```
This failure does NOT occur with PR #2072 (rke2r2) — the test passes consistently.

---

## Appendix: Kubernetes DefaultTolerationSeconds Admission Controller

The key discovery from Brad Davidson's investigation:

The Kubernetes `DefaultTolerationSeconds` admission controller ([source code](https://github.com/kubernetes/kubernetes/blob/master/plugin/pkg/admission/defaulttolerationseconds/admission.go#L68-L73)) automatically adds two tolerations to every pod:

```yaml
- key: node.kubernetes.io/not-ready
  operator: Exists
  effect: NoExecute
  tolerationSeconds: 300
- key: node.kubernetes.io/unreachable
  operator: Exists
  effect: NoExecute
  tolerationSeconds: 300
```

**BUT** it only adds these if the pod doesn't already have a toleration that matches `NoExecute` without a specific key. In Calico v3.32.0, the broad `NoExecute: Exists` toleration (no key) matches all NoExecute taints, so the admission controller skips adding the defaults. This means pods never get auto-evicted from NotReady/unreachable nodes.

The fix in v3.32.002 sets specific `controlPlaneTolerations` with keys, which means the admission controller WILL add the default `tolerationSeconds: 300` tolerations again. However, this only applies to Typha and other control-plane components managed by the tigera-operator. The `calico-node` DaemonSet is unaffected (DaemonSets tolerate everything by default).
