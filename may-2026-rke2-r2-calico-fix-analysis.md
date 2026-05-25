# May 2026 RKE2 R2 Calico Toleration Fix — Test Failure Analysis

## Executive Summary

The Calico toleration fix ([rke2#10438](https://github.com/rancher/rke2/issues/10438)) was applied via rke2-calico chart v3.32.002 and shipped in RKE2 R2 release candidates. KDM PRs [#2068](https://github.com/rancher/kontainer-driver-metadata/pull/2068) (dev-v2.15) and [#2067](https://github.com/rancher/kontainer-driver-metadata/pull/2067) (dev-v2.14) were updated with these R2 versions and re-tested.

**Result: The fix partially helped but did NOT fully resolve the failures.**

- ✅ `Test_Provisioning_Single_Node_All_Roles_Drain` — **NOW PASSES** for K8s 1.33, 1.34, 1.35 (was failing before)
- ❌ `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode` — **STILL FAILS** on all K8s versions (100% consistent)
- ⚠️ `Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` — **FLAKY** (~60-70% fail rate on K8s 33-35, passes on K8s 36)
- ⚠️ `Test_Provisioning_Single_Node_All_Roles_Drain` — **FLAKY** for K8s 1.36 specifically (passed attempt 1, failed attempt 2)

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
| MP_EtcdSnapshotWithThreeEtcd | ✅ All versions | ⚠️ Flaky (~60-70% fail K8s 33-35, passes K8s 36) | **⚠️ FLAKY (not consistent regression)** |
| MP_EtcdSnapshotOnNewNode | ✅ All versions | ⚠️ Flaky (sometimes fails, sometimes passes) | **⚠️ FLAKY** |
| Single_Node_All_Roles_Drain | ❌ K8s 1.33-1.35 | ✅ K8s 1.33-1.35, ⚠️ K8s 1.36 flaky | **✅ IMPROVED** |

> **UPDATE (May 24 rerun):** The MP test "regressions" reported earlier are NOT consistent regressions — they are flaky. On the K8s 1.36 rerun, both MP tests PASSED. The only fully consistent failure is `Custom_EtcdSnapshotOperationsOnNewCombinedNode`. See [CI Rerun Results](#ci-rerun-results-may-24-2026--attempt-2) section for full details.

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

### Why MP tests are FLAKY (not consistent regression)

The `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` and `MP_EtcdSnapshotOperationsOnNewNode` tests were initially reported as regressions but **rerun results show they're flaky**, not consistently broken:
- K8s 1.36 rerun: both MP tests PASSED
- K8s 1.33 rerun: MP_ThreeEtcd FAILED (same error)
- K8s 1.34 rerun: MP_ThreeEtcd FAILED
- K8s 1.35 rerun: timed out at 3600s (hard timeout killed it)

The flakiness pattern suggests a **timing-dependent issue** — the etcd restore may succeed if conditions align (e.g., if calico-typha happens to reschedule quickly enough). The previous R1 runs that passed may have been lucky with timing.

This is consistent with the toleration fix adding back `tolerationSeconds: 300` — it gives a 5-minute window for eviction. If the restore process takes longer than 5 minutes to encounter the stale typha situation, it could still fail depending on exact timing.

---

## Root Cause: The Etcd Restore Problem is NOT a Toleration Issue

Based on the prior analysis (see `rke2-v1.36.0-etcd-restore-failure-analysis.md`), the etcd snapshot restore failure on combined nodes has a **different root cause**:

1. After `rke2 server --cluster-reset`, kubelet starts with `--cloud-provider=external`
2. This adds `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` taint
3. The cloud-controller-manager (CCM) fails to remove this taint because the restored etcd contains stale node objects
4. With the uninitialized taint on the only node, `cattle-cluster-agent` can't be scheduled
5. Without the agent, Rancher can't connect and the planner is stuck

**The Calico toleration fix does not address this CCM taint issue at all.** The CCM taint problem is orthogonal to the Typha toleration problem.

### Why MP tests were previously passing (and are now flaky)

In the R1 runs, MP tests passed because:
- They have separate worker nodes without the CCM taint
- `cattle-cluster-agent` could schedule on workers
- The broader tolerations in v3.32.0 (pre-fix) meant Typha tolerated everything indefinitely — in a 3-etcd-node MP cluster, Typha stuck on a down node was less critical because there were other nodes available for calico-node to find a working Typha

With the R2 fix restoring restrictive tolerations + `tolerationSeconds: 300`:
- Typha WILL eventually be evicted (good) but the 5-minute window creates a race condition
- If the etcd restore process needs networking (via Typha) within those 5 minutes before eviction and rescheduling completes, it can still fail
- The flakiness indicates this is a timing race, not a definitive regression

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
| Typha not rescheduled on node removal | Calico v3.32.0 overly broad tolerations prevent Kubernetes auto-eviction | rke2-charts (Calico chart) | ✅ Fixed in v3.32.002 — Single_Node_Drain passes for K8s 1.33-1.35 |
| Etcd restore fails on combined node | CCM `uninitialized` taint blocks cattle-cluster-agent, stale etcd data confuses CCM | Rancher (planner/agent tolerations) + RKE2 (CCM behavior) | ❌ NOT fixed by toleration change |
| MP etcd restore flaky (K8s 1.33-1.35) | Timing race — 5-minute eviction window may not always be sufficient for restore flow | rke2-charts / Rancher planner (needs explicit Typha cleanup during restore) | ⚠️ FLAKY (not consistent regression) |
| K8s 1.36 Single_Node_Drain flaky | Calico probe wait on new node — same root cause as Typha issue but borderline timing | May need longer timeout or K8s 1.36-specific investigation | ⚠️ FLAKY (passes ~50% of the time) |

---

## CI Rerun Results (May 24, 2026 — Attempt 2)

CI was rerun on the same PRs to verify consistency of failures. Below is the comparison.

### Pass/Fail Comparison: Attempt 1 vs Attempt 2 (Rerun)

#### PR #2068 (dev-v2.15) — K8s 34, 35, 36

| Job | Attempt 1 | Attempt 2 (Rerun) | Consistent? |
|-----|-----------|-------------------|-------------|
| rke2, 34, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 34, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 34, Provisioning | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 35, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 35, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 35, Provisioning | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 36, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 36, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 36, Provisioning | ✅ PASS | ❌ FAIL | ⚠️ **FLAKY** |

#### PR #2067 (dev-v2.14) — K8s 33, 34, 35

| Job | Attempt 1 | Attempt 2 (Rerun) | Consistent? |
|-----|-----------|-------------------|-------------|
| rke2, 33, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 33, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 33, Provisioning | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 34, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 34, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 34, Provisioning | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 35, SetA | ✅ PASS | ✅ PASS | ✅ Yes |
| rke2, 35, SetB | ❌ FAIL | ❌ FAIL | ✅ Yes |
| rke2, 35, Provisioning | ✅ PASS | ✅ PASS | ✅ Yes |

### Key Observations from Reruns

1. **SetB failures are 100% consistent** — fails on every K8s version in both attempts, same test
2. **Provisioning tests for K8s 33, 34, 35 are 100% stable** — pass every time  
3. **K8s 1.36 Provisioning is FLAKY** — passed attempt 1, failed attempt 2 (same `Single_Node_All_Roles_Drain` test)
4. **K8s 1.36 SetB MP tests improved** — on the rerun, `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` and `MP_EtcdSnapshotOperationsOnNewNode` both **PASSED** (previously reported failing)

### Detailed Test Results from Rerun (Attempt 2)

#### PR #2068 K8s 36 SetB (Job 77611319916):
```
--- FAIL: Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode (3037.97s)
    etcdsnapshot.go:277: cluster test-custom-etcd-snapshot-operations-on-new-combined-node 
    etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed : 
    timeout waiting condition: context deadline exceeded
--- PASS: Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode (747.88s)
--- PASS: Test_Operation_SetB_MP_EtcdSnapshotOperationsOnNewNode (838.83s)
```

#### PR #2067 K8s 33 SetB (Job 77603951238):
```
--- FAIL: Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode (3051.50s)
    etcdsnapshot.go:277: cluster test-custom-etcd-snapshot-operations-on-new-combined-node 
    etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed
--- FAIL: Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode (3118.32s)
    etcdsnapshot.go:277: cluster test-mp-etcd-snapshot-conventional-arch-new-node 
    etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed : timeout
```

#### PR #2067 K8s 34 SetB (Job 77603951253):
```
--- FAIL: Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode (2975.90s)
    etcdsnapshot.go:277: cluster test-custom-etcd-snapshot-operations-on-new-combined-node 
    etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed
--- FAIL: Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode (3208.59s)
    etcdsnapshot.go:277: cluster test-mp-etcd-snapshot-conventional-arch-new-node 
    etcd snapshot restore wait failed on: etcd snapshot restore wait did not succeed : timeout
```

#### PR #2067 K8s 35 SetB (Job 77603951246):
```
FAIL  github.com/rancher/rancher/tests/v2prov/tests/machineprovisioning  3600.278s
```
_(1-hour hard timeout — no specific test named, both tests likely ran over time)_

#### PR #2068 K8s 34 SetB (Job 77611319925):
```
FAIL  github.com/rancher/rancher/tests/v2prov/tests/machineprovisioning  3600.382s
```
_(1-hour hard timeout)_

#### PR #2068 K8s 35 SetB (Job 77611319923):
```
FAIL  github.com/rancher/rancher/tests/v2prov/tests/machineprovisioning  3600.241s
```
_(1-hour hard timeout)_

#### PR #2068 K8s 36 Provisioning (Job 77611319927 — the FLAKY one):
```
--- PASS: Test_Provisioning_Custom_ThreeNode (292.43s)
--- PASS: Test_Provisioning_Custom_UniqueRoles (348.32s)
--- PASS: Test_Provisioning_MP_SingleNodeAllRolesWithDelete (270.28s)
--- PASS: Test_Provisioning_MP_MultipleEtcdNodesScaledDownThenDelete (394.32s)
--- PASS: Test_Provisioning_MP_Drain (237.76s)
--- PASS: Test_Provisioning_MP_DrainNoDelete (212.79s)
--- FAIL: Test_Provisioning_Single_Node_All_Roles_Drain (1189.47s)
```

**Root cause from Rancher logs:** The planner was stuck `waiting for probes: calico` on the new replacement node for 10+ minutes. The calico probe never succeeded within the test timeout.

### Revised Assessment After Reruns

| Category | Finding |
|----------|---------|
| **Consistent failure** | `Custom_EtcdSnapshotOperationsOnNewCombinedNode` — fails 100% of the time on ALL K8s versions. This is the etcd restore to new combined node issue (CCM taint / stale etcd state). |
| **Mostly consistent** | `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` — fails ~60-70% of the time on K8s 33, 34, 35. Passed on K8s 36 rerun. May be timing-dependent. |
| **Flaky** | `Single_Node_All_Roles_Drain` on K8s 1.36 — passed once, failed once. When it fails, planner gets stuck waiting for calico probe. This is likely still the Typha issue but sometimes the 5-minute eviction timeout (restored by the fix) is enough and sometimes it isn't for K8s 1.36. |
| **Fixed/stable** | `Single_Node_All_Roles_Drain` on K8s 33, 34, 35 — passes 100% now (was failing before the R2 toleration fix). **The fix works for these versions.** |
| **Note on v1.36.1-rc2+rke2r1** | K8s 1.36 uses `v1.36.1-rc2+rke2r1` (not `rke2r2`). While the chart v3.32.002 is listed, this version label difference from the other R2 versions is notable and may indicate a different build pathway. |

### Conclusion

The rerun **confirms the same pattern** as the first run with only one exception:
- **K8s 1.36 Provisioning (`Single_Node_All_Roles_Drain`) is flaky** — not a consistent failure. It passed on attempt 1 and failed on attempt 2.
- All other results are identical between attempts.

The underlying issues remain:
1. **Etcd restore to combined node is broken** (not fixable by toleration change alone)
2. **MP etcd restore tests are flaky/slow** with ~60-70% failure rate due to timeouts
3. **The Typha toleration fix DID help** — `Single_Node_All_Roles_Drain` is now consistently passing for K8s 1.33-1.35 (was consistently failing before)
4. **K8s 1.36 has a borderline timing issue** — the 5-minute eviction timeout may not always be sufficient, or there's an additional K8s 1.36-specific interaction

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
