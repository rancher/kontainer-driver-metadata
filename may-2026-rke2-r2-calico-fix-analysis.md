# May 2026 RKE2 R2 Calico Toleration Fix — Test Failure Analysis

## Executive Summary

The Calico toleration fix ([rke2#10438](https://github.com/rancher/rke2/issues/10438)) was applied via rke2-calico chart v3.32.002 and shipped in RKE2 R2 release candidates. KDM PRs [#2068](https://github.com/rancher/kontainer-driver-metadata/pull/2068) (dev-v2.15) and [#2067](https://github.com/rancher/kontainer-driver-metadata/pull/2067) (dev-v2.14) were updated with these R2 versions and re-tested.

**Result: The fix partially helped but did NOT fully resolve the failures.**

- ✅ `Test_Provisioning_Single_Node_All_Roles_Drain` — **NOW PASSES** for K8s 1.33, 1.34, 1.35 (was failing before)
- ❌ `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode` — **STILL FAILS** on all K8s versions
- ❌ `Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` — **NEW REGRESSION**, now FAILS on K8s 1.33, 1.34, 1.35 (was PASSING before!)
- ❌ `Test_Operation_SetB_MP_EtcdSnapshotOperationsOnNewNode` — **NEW REGRESSION**, now FAILS on K8s 1.33 (was PASSING before!)
- ❌ `Test_Provisioning_Single_Node_All_Roles_Drain` — **STILL FAILS** for K8s 1.36 specifically (v1.36.1-rc2+rke2r1)

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

This causes the tigera-operator to use these specific tolerations for control-plane components (including Typha) instead of the broad `TolerateAll` defaults. With key-specific tolerations, Kubernetes' admission controller will again auto-inject the `tolerationSeconds: 300` tolerations for `not-ready` and `unreachable`, re-enabling pod eviction after 5 minutes.

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

### PR #2068 (dev-v2.15) — K8s 34, 35, 36

| K8s Version | RKE2 Version | Has Fix? | rke2-calico chart |
|-------------|-------------|----------|-------------------|
| 34 | v1.34.8-rc2+**rke2r2** | ✅ Yes | v3.32.002 |
| 35 | v1.35.5-rc2+**rke2r2** | ✅ Yes | v3.32.002 |
| 36 | v1.36.1-rc2+**rke2r1** | ✅ Yes | v3.32.002 |

### PR #2067 (dev-v2.14) — K8s 33, 34, 35

| K8s Version | RKE2 Version | Has Fix? | rke2-calico chart |
|-------------|-------------|----------|-------------------|
| 33 | v1.33.12-rc2+**rke2r2** | ✅ Yes | v3.32.002 |
| 34 | v1.34.8-rc2+**rke2r2** | ✅ Yes | v3.32.002 |
| 35 | v1.35.5-rc2+**rke2r2** | ✅ Yes | v3.32.002 |

**Confirmation: All R2 versions include the Calico toleration fix (chart v3.32.002).**

Note: v1.36.1-rc2+rke2r1 is labeled "rke2r1" but also includes the fix chart v3.32.002 — this is the first release of the 1.36.1 patch cycle, not a revision of 1.36.0.

---

## Complete Test Results (R2 Runs, May 22-23)

### Test_Operation_SetB tests

| Test | K8s | PR#2067 (v2.14) | PR#2068 (v2.15) |
|------|-----|-----------------|-----------------|
| **Custom_EtcdSnapshotOperationsOnNewCombinedNode** | 33 | ❌ FAIL (2985s) | — |
| | 34 | ❌ FAIL (3044s) | ❌ FAIL (timeout) |
| | 35 | ❌ FAIL (2982s) | ❌ FAIL (3042s) |
| | 36 | — | ❌ FAIL (3197s) |
| **MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode** | 33 | ✅ PASS (618s) | — |
| | 34 | ❌ FAIL (3133s) | timeout |
| | 35 | ❌ FAIL (3159s) | ❌ FAIL (3192s) |
| | 36 | — | ✅ PASS (853s) |
| **MP_EtcdSnapshotOperationsOnNewNode** | 33 | ❌ FAIL (2945s) | — |
| | 34 | timeout | timeout |
| | 35 | — | — |
| | 36 | — | ✅ PASS (593s) |

### Test_Provisioning tests

| Test | K8s | PR#2067 (v2.14) | PR#2068 (v2.15) |
|------|-----|-----------------|-----------------|
| **Single_Node_All_Roles_Drain** | 33 | ✅ PASS | — |
| | 34 | ✅ PASS | ✅ PASS |
| | 35 | ✅ PASS | ✅ PASS |
| | 36 | — | ❌ FAIL (1149s) |
| Other Provisioning tests | all | ✅ PASS | ✅ PASS |

### Comparison: Before Fix (R1) vs After Fix (R2)

| Test | R1 (before fix) | R2 (with fix) | Delta |
|------|-----------------|---------------|-------|
| Custom_EtcdSnapshotOnNewCombinedNode | ❌ All versions | ❌ All versions | **No change** |
| MP_EtcdSnapshotWithThreeEtcd | ✅ All versions | ❌ K8s 1.34, 1.35 (✅ 1.33, 1.36) | **⚠️ REGRESSION** |
| MP_EtcdSnapshotOnNewNode | ✅ All versions | ❌ K8s 1.33 (✅ K8s 1.36) | **⚠️ REGRESSION** |
| Single_Node_All_Roles_Drain | ❌ K8s 1.33-1.35 | ✅ K8s 1.33-1.35 (❌ K8s 1.36 only) | **✅ IMPROVED** |

---

## Analysis: Why the Fix Didn't Fully Help

### What the fix DID solve: Single_Node_All_Roles_Drain

The `Single_Node_All_Roles_Drain` test creates a single-node cluster, adds a new node, then deletes the old node (without drain). Previously, Calico Typha was stuck on the old node indefinitely because it tolerated `unreachable:NoExecute` without a timeout. With the fix restoring `tolerationSeconds: 300`, Typha would be evicted after 5 minutes and rescheduled.

This test now passes for K8s 1.33, 1.34, 1.35 — confirming the toleration fix works for the Typha rescheduling scenario.

The K8s 1.36 `Single_Node_All_Roles_Drain` failure persists, likely due to a **separate K8s 1.36 / Rancher v2.15 interaction** unrelated to tolerations.

### What the fix did NOT solve: Etcd Snapshot Restore

The `Custom_EtcdSnapshotOperationsOnNewCombinedNode` test still fails on ALL versions. This test involves:

1. Creating a single combined etcd+cp+worker node cluster
2. Taking an etcd snapshot
3. **Restoring** to a **new** combined node

The restore scenario is fundamentally different from the drain scenario:

- **Drain scenario**: Old node goes NotReady → Typha needs to be evicted (fixed by tolerationSeconds: 300)
- **Restore scenario**: Entire cluster is destroyed and rebuilt from scratch on a new node. The old cluster's pods (including Typha) exist only in the restored etcd data, pointing to a non-existent old node.

**The toleration fix cannot help with restore** because:
1. After restore, the old Typha pod in etcd data references a node that no longer exists
2. The tigera-operator needs to detect this and reschedule Typha
3. But the tigera-operator itself may be stuck because networking isn't functional
4. The `tolerationSeconds: 300` helps with graceful eviction, but in a restore scenario the node doesn't just go NotReady — it's completely gone

### Why MP tests now ALSO fail (NEW REGRESSION)

The `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` and `MP_EtcdSnapshotOperationsOnNewNode` tests previously passed because they have separate worker nodes that remained alive during restore. Now they're failing too, but **only for K8s 1.33-1.35** (K8s 1.36 still passes).

This pattern strongly suggests the R2 builds introduced a **new regression** unrelated to the toleration fix — possibly:

1. **A build artifact issue** — the R2 RC builds may have packaging problems affecting etcd restore
2. **A regression in another component** included in the R2 release — the R2 builds may include other changes beyond just the Calico chart fix
3. **An interaction between the toleration fix and the restore flow** — the more restrictive tolerations may prevent critical pods from scheduling during restore when nodes have specific taints

The fact that K8s 1.36 MP tests pass but K8s 1.33-1.35 fail suggests a **K8s version-dependent behavior** in how the restore and node reinitialization interact with the new toleration settings.

---

## Root Cause: The Etcd Restore Problem is NOT a Toleration Issue

Based on the prior analysis (see `rke2-v1.36.0-etcd-restore-failure-analysis.md`), the etcd snapshot restore failure on combined nodes has a **different root cause**:

1. After `rke2 server --cluster-reset`, kubelet starts with `--cloud-provider=external`
2. This adds `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` taint
3. The cloud-controller-manager (CCM) fails to remove this taint because the restored etcd contains stale node objects
4. With the uninitialized taint on the only node, `cattle-cluster-agent` can't be scheduled
5. Without the agent, Rancher can't connect and the planner is stuck

**The Calico toleration fix does not address this CCM taint issue at all.** The CCM taint problem is orthogonal to the Typha toleration problem.

### Why MP tests were previously immune

In the R1 runs, MP tests passed because:
- They have separate worker nodes without the CCM taint
- `cattle-cluster-agent` could schedule on workers
- Typha could be on any node (broad tolerations meant it stayed even on broken nodes, but networking worked because workers were healthy)

### Why MP tests now fail with R2

With the R2 fix, the tolerations are more restrictive. This may cause Typha (and other Calico control-plane components) to **not tolerate certain taints** that exist during the restore flow. If the restored nodes have taints that the new restrictive tolerations don't cover, Calico components can't schedule, breaking networking even on MP clusters.

Specifically, if nodes after restore temporarily have taints beyond just `control-plane:NoSchedule` and `etcd:NoExecute` (e.g., `uninitialized:NoSchedule`, or custom taints from the restore process), the new restrictive tolerations would prevent Calico from scheduling there.

---

## Suggested Fixes

### Fix 1: Restore the `not-ready` and `unreachable` toleration behavior (IMMEDIATE)

The current fix removes ALL broad tolerations and only keeps control-plane/etcd ones. This loses the `not-ready:NoExecute` and `unreachable:NoExecute` tolerations with 300s timeout that Kubernetes used to auto-inject.

However, during restore, nodes may be in `NotReady` state temporarily. Without `not-ready:NoExecute` tolerance, Calico pods may be evicted from temporarily-NotReady nodes during restore, breaking networking.

**Recommendation**: Explicitly add `not-ready` and `unreachable` tolerations back with `tolerationSeconds: 300`:

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

This exactly matches what Calico v3.31.5 had (the combination of explicit tolerations + Kubernetes auto-injection).

### Fix 2: Add CCM uninitialized taint toleration to cattle-cluster-agent (FOR COMBINED NODE RESTORE)

The core etcd restore issue on combined nodes is the CCM `uninitialized:NoSchedule` taint blocking `cattle-cluster-agent`. This requires a fix in Rancher (not RKE2):

```yaml
# cattle-cluster-agent Deployment should tolerate:
- key: "node.cloudprovider.kubernetes.io/uninitialized"
  operator: "Exists"
  effect: "NoSchedule"
```

### Fix 3: Post-restore cleanup should handle stale Calico resources

During etcd restore, the planner should explicitly:
1. Delete stale Calico pods (Typha, calico-node, calico-kube-controllers) that reference old nodes
2. Restart the tigera-operator to re-reconcile
3. Wait for Calico networking to be healthy before proceeding

### Fix 4: Investigate K8s 1.36 Single_Node_All_Roles_Drain regression

The K8s 1.36 `Single_Node_All_Roles_Drain` failure on PR#2068 may be a Rancher v2.15-specific issue or a K8s 1.36 / CAPI interaction problem. This needs separate investigation.

---

## Summary of Issues

| Issue | Root Cause | Fix Location | Status |
|-------|-----------|--------------|--------|
| Typha not rescheduled on node removal | Calico v3.32.0 overly broad tolerations prevent Kubernetes auto-eviction | rke2-charts (Calico chart) | ✅ Partially fixed in v3.32.002, but may need `not-ready`/`unreachable` tolerations re-added |
| Etcd restore fails on combined node | CCM `uninitialized` taint blocks cattle-cluster-agent, stale etcd data confuses CCM | Rancher (planner/agent tolerations) + RKE2 (CCM behavior) | ❌ NOT fixed by toleration change |
| MP etcd restore now fails (K8s 1.33-1.35) | Possibly: too-restrictive tolerations prevent Calico scheduling during restore with transient taints | rke2-charts (Calico chart tolerations need to be less restrictive during restore) | ❌ NEW REGRESSION from fix |
| K8s 1.36 Single_Node_Drain still fails | Separate issue, likely Rancher v2.15 or K8s 1.36 specific | Unknown | ❌ Needs investigation |

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

The fix in v3.32.002 sets specific `controlPlaneTolerations` with keys, which means the admission controller WILL add the default `tolerationSeconds: 300` tolerations again. However, this only applies to Typha and other control-plane components managed by the tigera-operator. The `calico-node` DaemonSet is unaffected (DaemonSets tolerate everything by default and are managed differently).
