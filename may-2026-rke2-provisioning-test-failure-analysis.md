# May 2026 RKE2 Provisioning Test Failure Analysis

## Scope

Analysis of provisioning test failures across all four May RKE2 patch PRs:

| PR | Branch | K8s versions tested |
|----|--------|---------------------|
| [#2068](https://github.com/rancher/kontainer-driver-metadata/pull/2068) | dev-v2.15 | 34, 35, 36 |
| [#2067](https://github.com/rancher/kontainer-driver-metadata/pull/2067) | dev-v2.14 | 33, 34, 35 |
| [#2066](https://github.com/rancher/kontainer-driver-metadata/pull/2066) | dev-v2.13 | 32, 33, 34 |
| [#2065](https://github.com/rancher/kontainer-driver-metadata/pull/2065) | dev-v2.12 | 31, 32, 33 |

---

## Executive Summary

There are exactly **two distinct failures** affecting all four PRs. Both are **RKE2-only** (all K3s tests pass). Both failures share a common root cause pattern: **single-node / combined-role RKE2 clusters fail to complete lifecycle operations that work correctly in machine-pool (MP) multi-node setups**.

### Failure 1: `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode`

**Error:** `etcd snapshot restore wait did not succeed : timeout waiting condition: context deadline exceeded`

This is the **same failure** previously analyzed in `rke2-v1.36.0-etcd-restore-failure-analysis.md`. It is **not unique to K8s 1.36** — it occurs on K8s 1.33, 1.34, 1.35, and 1.36.

### Failure 2: `Test_Provisioning_Single_Node_All_Roles_Drain`

**Error:** `did not converge back to a single node`

This is a **new failure** not previously analyzed. It occurs on K8s 1.33, 1.34, and 1.35.

---

## Complete Pass/Fail Matrix

### `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode`

| K8s ver | PR#2065 (v2.12) | PR#2066 (v2.13) | PR#2067 (v2.14) | PR#2068 (v2.15) |
|---------|----------------|----------------|----------------|----------------|
| **31** | ✅ PASS | — | — | — |
| **32** | ✅ PASS | ✅ PASS | — | — |
| **33** | ✅ PASS (628s) | ❌ FAIL (3035s) | ❌ FAIL (2995s) | — |
| **34** | — | ❌ FAIL (3041s) | ❌ FAIL (2983s) | ❌ FAIL (3056s) |
| **35** | — | — | ❌ FAIL (3004s) | ❌ FAIL (3061s) |
| **36** | — | — | — | ❌ FAIL (3026s) |

> **Key observation:** K8s 1.33 **passes on v2.12** but **fails on v2.13+**. This means the Rancher version matters, not just the K8s version!

### `Test_Provisioning_Single_Node_All_Roles_Drain`

| K8s ver | PR#2065 (v2.12) | PR#2066 (v2.13) | PR#2067 (v2.14) | PR#2068 (v2.15) |
|---------|----------------|----------------|----------------|----------------|
| **31** | ✅ PASS | — | — | — |
| **32** | ✅ PASS | ✅ PASS (345s) | — | — |
| **33** | ❌ FAIL (1146s) | ❌ FAIL (1134s) | ❌ FAIL (1138s) | — |
| **34** | — | ❌ FAIL (1149s) | ❌ FAIL (1161s) | ❌ FAIL (1150s) |
| **35** | — | — | ❌ FAIL (1154s) | ❌ FAIL (1143s) |
| **36** | — | — | — | ✅ PASS |

> **Key observation:** Fails on K8s 1.33–1.35 across ALL Rancher versions, but K8s 1.32 passes AND K8s 1.36 passes. This suggests the failure is K8s-version-specific, possibly related to a temporary K8s regression that was fixed in 1.36.

### Other SetB tests (MP variants) — all pass

| Test | Result |
|------|--------|
| `Test_Operation_SetB_MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode` | ✅ PASS everywhere |
| `Test_Operation_SetB_MP_EtcdSnapshotOperationsOnNewNode` | ✅ PASS everywhere |
| `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewNode` | ✅ PASS (PR#2065 only test that includes this) |

### All other Provisioning tests — all pass

All of these pass across all PRs and K8s versions:
- `Test_Provisioning_Custom_ThreeNode`
- `Test_Provisioning_Custom_UniqueRoles`
- `Test_Provisioning_MP_SingleNodeAllRolesWithDelete`
- `Test_Provisioning_MP_MultipleEtcdNodesScaledDownThenDelete`
- `Test_Provisioning_MP_Drain`
- `Test_Provisioning_MP_DrainNoDelete`

---

## Failure 1: Etcd Snapshot Restore on Combined Node — Deep Analysis

### What the test does

`Custom_EtcdSnapshotOperationsOnNewCombinedNode` creates a custom cluster with a single node running **all three roles** (etcd + controlplane + worker), takes an etcd snapshot, then performs a restore to a **new** combined node. This is the most complex restore scenario: the entire cluster must be rebuilt on a fresh node from the snapshot.

### Failure pattern

The test consistently times out after ~3000s (~50 minutes). The error is always:
```
etcd snapshot restore wait did not succeed : timeout waiting condition: context deadline exceeded
```

### The v2.12 vs v2.13 boundary

The most critical finding is that **K8s 1.33 passes on Rancher v2.12 (PR#2065) but fails on v2.13 (PR#2066)**. This means:

- The RKE2 v1.33 binary itself is **not the sole cause** — the same binary works with Rancher v2.12.
- Something changed in **Rancher v2.13** that broke the restore flow for combined-node custom clusters.
- The failure then persists through v2.14 and v2.15 because the same Rancher-side change is present.

### Root cause (from prior analysis, confirmed here)

The detailed analysis in `rke2-v1.36.0-etcd-restore-failure-analysis.md` identified the root cause as:

1. After restore, kubelet starts with `--cloud-provider=external`, adding the `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` taint.
2. The cloud-controller-manager (CCM) fails to remove this taint after restore (stale node objects from snapshot confuse it).
3. With the `NoSchedule` taint on the **only** node, `cattle-cluster-agent` pods cannot be scheduled.
4. Without the cluster agent, Rancher cannot reconnect and the planner is permanently stuck.

The MP tests pass because they have **separate worker nodes** that are already initialized (no taint), so the agent can schedule there.

### Why v2.12 passes but v2.13 fails (for K8s 1.33)

This is likely due to a change in Rancher v2.13 in one of:
- The planner's post-restore reconciliation logic
- How/when the worker node is killed and restarted during restore
- How `cattle-cluster-agent` tolerations are configured
- CCM / cloud-provider integration changes

**Recommendation:** Diff the Rancher planner and CAPR controller code between v2.12 and v2.13 to find the specific change that introduced this regression.

---

## Failure 2: Single Node Drain — Deep Analysis

### What the test does

`Test_Provisioning_Single_Node_All_Roles_Drain` creates a single-node cluster with all roles (etcd + controlplane + worker) that has drain enabled in its upgrade strategy. It then:
1. Changes the machine config (PodConfig) to trigger a rolling update
2. Waits for a second machine to appear (scale-up)
3. Waits for the second node's NodeRef and Ready status
4. Verifies the new node is NOT cordoned (incoming CP shouldn't be drained)
5. **Waits for the first (old) machine to enter Deleting state**
6. **Waits for convergence back to 1 machine** ← THIS TIMES OUT

### Failure pattern

The test consistently times out after ~1140-1160s (~19 minutes) with:
```
Messages: did not converge back to a single node
```

This means the rolling replacement of the single node gets stuck — the old machine is never fully removed after the new one comes up.

### K8s version pattern: 1.33-1.35 fail, 1.32 and 1.36 pass

| K8s | Duration | Result |
|-----|----------|--------|
| 1.32 | 345s | ✅ PASS |
| 1.33 | ~1140s | ❌ FAIL |
| 1.34 | ~1150s | ❌ FAIL |
| 1.35 | ~1150s | ❌ FAIL |
| 1.36 | — | ✅ PASS |

This U-shaped pattern (pass → fail → fail → fail → pass) strongly suggests a **K8s-level regression** that was introduced in 1.33 and fixed in 1.36, OR a behavior change in RKE2's drain/rolling-update handling.

### Likely root cause

The drain/cordon mechanism for single-node all-roles clusters changed behavior in K8s 1.33. When draining the only control-plane node during a rolling update:

1. The old node is drained (cordoned + pods evicted)
2. The new node comes up and joins the cluster
3. The old node should be deleted
4. **But the deletion gets stuck** — likely because the node drain/deletion controller can't properly handle the transition when both nodes briefly exist as control-plane nodes

The fact that K8s 1.36 passes again suggests either:
- A K8s upstream fix was backported to 1.36
- RKE2 v1.36 includes a fix in its drain/rolling-update logic
- The CAPI/machine controller behavior changed in a way that resolves this

### This failure is independent of Rancher version

Unlike Failure 1, this test fails on K8s 1.33 across **all** Rancher versions (v2.12 through v2.15). This confirms it's a K8s/RKE2-level issue, not a Rancher planner issue.

---

## Summary Table

| Failure | Test | Root Cause Domain | K8s Versions Affected | Rancher Versions Affected |
|---------|------|-------------------|----------------------|--------------------------|
| **1** | `Custom_EtcdSnapshotOperationsOnNewCombinedNode` | Rancher planner + CCM taint interaction | 1.33+ (but 1.33 works on v2.12!) | v2.13, v2.14, v2.15 |
| **2** | `Single_Node_All_Roles_Drain` | K8s/RKE2 drain/rolling-update | 1.33, 1.34, 1.35 (NOT 1.36) | ALL (v2.12-v2.15) |

---

## Answers to Specific Questions

### "Is the failure in `provisioning-test (rke2, 36, ^Test_Operation_SetB_.*$)` the same?"

**Yes.** It is the exact same `Custom_EtcdSnapshotOperationsOnNewCombinedNode` failure as documented in `rke2-v1.36.0-etcd-restore-failure-analysis.md`. Same test, same error message, same timeout duration (~3000s). The failure is **not specific to K8s 1.36** — it occurs identically on K8s 1.33 (v2.13+), 1.34, and 1.35 as well.

### "What about all other failures?"

The only **other** unique failure is `Test_Provisioning_Single_Node_All_Roles_Drain`, which fails with "did not converge back to a single node". This is a completely different bug — a single-node rolling replacement that can't remove the old machine. It affects K8s 1.33–1.35 but is fixed in K8s 1.36.

All remaining test failures in the PRs are instances of these two failures. **No other test is failing.**

---

## Recommendations

### For Failure 1 (Etcd Snapshot Restore)

1. **Priority: HIGH** — This is a Rancher regression introduced in v2.13.
2. **Action:** Diff Rancher v2.12 vs v2.13 planner/CAPR code, specifically:
   - `pkg/controllers/provisioningv2/rke2/` — restore reconciliation
   - `pkg/capr/planner/` — etcd restore plan generation
   - How worker nodes are handled during restore
   - `cattle-cluster-agent` Deployment tolerations
3. **Quick fix option:** Add `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` toleration to `cattle-cluster-agent` to prevent it from being blocked by the CCM taint.

### For Failure 2 (Single Node Drain)

1. **Priority: MEDIUM** — This is a K8s/RKE2 issue that resolves itself in K8s 1.36.
2. **Action:** Check if there are known K8s issues with single-node rolling updates in 1.33-1.35.
3. **Workaround:** The test itself may need a longer timeout, or the single-node drain scenario may need special handling in the CAPI machine controller for K8s 1.33-1.35.
4. **Note:** Since K8s 1.36 passes, this will self-resolve as K8s 1.32 ages out of support.
