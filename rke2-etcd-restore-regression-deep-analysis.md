# RKE2 Etcd Restore Regression — Deep Root Cause Analysis

## Executive Summary

The `Custom_EtcdSnapshotOperationsOnNewCombinedNode` test began failing across ALL K8s versions (1.33–1.36) in the May 2026 RKE2 release cycle. This is an **RKE2-side regression** caused by the **Calico v3.31.5 → v3.32.0 upgrade** that shipped in all May patches.

The R2 toleration fix (`rke2-calico` chart v3.32.002) resolves the `Single_Node_All_Roles_Drain` test but does **not** fix the etcd restore test. This document explains why, identifies the root causes, and proposes RKE2-side fixes.

---

## The Regression Boundary

| Release | RKE2 Versions | Calico | tigera-operator | rke2-calico chart | Etcd Restore | Drain |
|---------|--------------|--------|-----------------|-------------------|-------------|-------|
| **April R1** | v1.33.11, v1.34.7, v1.35.4 | **v3.31.5** | **v1.40.8** | **v3.31.500** | ✅ PASS | ✅ PASS |
| **May R1** | v1.33.12, v1.34.8, v1.35.5, v1.36.1 | **v3.32.0** | **v1.42.0** | **v3.32.001** | ❌ FAIL | ❌ FAIL |
| **May R2** | same + rke2r2 | **v3.32.0** | **v1.42.0** | **v3.32.002** | ❌ **STILL FAILS** | ✅ FIXED |

The ONLY significant change between April and May releases is the Calico/tigera-operator upgrade, delivered via the May CNI update PRs:
- [rancher/rke2#10350](https://github.com/rancher/rke2/pull/10350) (release-1.33)
- [rancher/rke2#10349](https://github.com/rancher/rke2/pull/10349) (release-1.34)
- [rancher/rke2#10348](https://github.com/rancher/rke2/pull/10348) (release-1.35)
- [rancher/rke2#10384](https://github.com/rancher/rke2/pull/10384) (main/1.36)

These PRs bumped:
- `rancher/hardened-calico`: v3.31.5-build20260415 → v3.32.0-build20260507
- `rancher/mirrored-calico-operator`: v1.40.8 → v1.42.0
- `rancher/mirrored-calico-typha`: v3.31.5 → v3.32.0
- `rancher/mirrored-calico-node`: v3.31.5 → v3.32.0
- All other Calico component images: v3.31.5 → v3.32.0

---

## Finding 1: YAML Patch Bug in rke2-calico v3.32.001 (Root Cause of Drain Failure)

### Discovery

The rke2-calico chart uses a `values.yaml.patch` to overlay RKE2-specific settings onto the upstream tigera-operator chart. A **YAML duplicate key bug** in v3.32.001 caused the `controlPlaneTolerations` override to be silently ignored.

### The Bug

**Upstream Calico v3.31.5** `values.yaml` — NO `controlPlaneTolerations` field exists:
```yaml
installation:
  enabled: true
  kubernetesProvider: ""
  # ... (no controlPlaneTolerations)
```

**Upstream Calico v3.32.0** `values.yaml` — `controlPlaneTolerations: []` ADDED at line ~67:
```yaml
installation:
  # ... (many fields)
  controlPlaneReplicas: 2
  controlPlaneNodeSelector: {}
  controlPlaneTolerations: []     # <-- NEW in v3.32.0
  nonPrivileged: "Disabled"
```

The rke2-calico chart patch for **both** v3.31.500 and v3.32.001 attempted to set `controlPlaneTolerations` at position 9 (near the top of `installation:`):

```yaml
# Patch adds at position 9:
installation:
  controlPlaneTolerations:
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node-role.kubernetes.io/etcd"
    operator: "Exists"
    effect: "NoExecute"
  enabled: true
```

**In v3.31.500**: No conflict — the upstream doesn't have `controlPlaneTolerations`. Patch works correctly. ✅

**In v3.32.001**: YAML now has TWO `controlPlaneTolerations` entries under the same `installation:` mapping:
- Position 9: `controlPlaneTolerations: [{cp:NoSchedule}, {etcd:NoExecute}]` (from rke2 patch)
- Position 67: `controlPlaneTolerations: []` (from upstream v3.32.0, NOT modified by patch)

**YAML duplicate key rule: last value wins.** The empty list at position 67 silently overrides the patch at position 9.

### Effect on Typha Tolerations

The tigera-operator's Typha deployment code ([typha.go](https://github.com/tigera/operator/blob/v1.42.0/pkg/render/typha.go)):
```go
tolerations := rmeta.TolerateAll    // default: tolerate everything
if len(c.cfg.Installation.ControlPlaneTolerations) != 0 {
    tolerations = c.cfg.Installation.ControlPlaneTolerations  // override with chart values
}
```

**This code is IDENTICAL in operator v1.40.8 and v1.42.0.**

| Chart Version | Effective `ControlPlaneTolerations` | `len()` | Typha gets... |
|--------------|-------------------------------------|---------|---------------|
| v3.31.500 | `[{cp:NoSchedule}, {etcd:NoExecute}]` | 2 | Chart values (key-specific) |
| v3.32.001 | `[]` (empty, due to YAML override) | 0 | `TolerateAll` (blanket) |
| v3.32.002 | `[{cp:NoSchedule}, {etcd:NoExecute}]` | 2 | Chart values (key-specific) |

And `TolerateAll` is ([meta.go](https://github.com/tigera/operator/blob/v1.42.0/pkg/render/common/meta/meta.go)):
```go
TolerateAll = []corev1.Toleration{
    {Key: "CriticalAddonsOnly", Operator: Exists},
    {Effect: NoSchedule, Operator: Exists},   // blanket: tolerates ALL NoSchedule
    {Effect: NoExecute, Operator: Exists},     // blanket: tolerates ALL NoExecute
}
```

The blanket `NoExecute:Exists` toleration (without a key) **matches all NoExecute taints**, causing the Kubernetes `DefaultTolerationSeconds` admission controller to **skip adding** the default `tolerationSeconds: 300` for `not-ready:NoExecute` and `unreachable:NoExecute`. This means Typha tolerates dead/unreachable nodes **indefinitely** → never evicted → networking never recovers on the new node.

### How v3.32.002 Fixed It

The [rke2-charts#1002](https://github.com/rancher/rke2-charts/pull/1002) fix:
1. **Removed** the patch at position 9 (which was being overridden anyway)
2. **Modified** the upstream `controlPlaneTolerations: []` at position 67 to contain the actual values

This eliminates the YAML duplicate key issue. Now `ControlPlaneTolerations` is correctly set, the operator uses the chart values, and the admission controller adds `tolerationSeconds: 300`. Typha is evicted from dead nodes after 5 minutes. **The Drain test passes.**

---

## Finding 2: Etcd Restore Failure Has a SEPARATE Root Cause Within Calico v3.32.0

### The paradox

After the v3.32.002 fix, v3.31.500 (April) and v3.32.002 (May R2) have **identical Typha toleration behavior**:
- Both set `controlPlaneTolerations: [{cp:NoSchedule}, {etcd:NoExecute}]`
- Both get admission-controller-injected `tolerationSeconds: 300`
- Both evict Typha from dead nodes after 5 minutes

Yet:
- **v3.31.500 PASSES** the etcd restore test (in ~628s ≈ 10.5 min)
- **v3.32.002 FAILS** the etcd restore test (after ~3000s ≈ 50 min timeout)

**This means the etcd restore failure is NOT a toleration issue.** Something else in the Calico v3.32.0 upgrade breaks the post-restore recovery flow.

### What happens during etcd restore (from CI artifacts)

```
1. etcd snapshot taken on original combined-role node
2. Original node killed (rke2-killall.sh)
3. New node starts → cluster-reset restores etcd from snapshot
4. RKE2 restarts → kube-apiserver, etcd, kubelet start
5. Post-restore cleanup deletes stale pods (calico, CoreDNS, cattle-cluster-agent)
6. New calico-node DaemonSet pod starts on new node
7. Felix needs to connect to Typha to initialize dataplane
8. New Typha Deployment pod needs to be scheduled on new node
9. Rancher planner waits for "calico" probe and cluster agent
```

The test fails at step 9 — the planner is permanently stuck:
```
[planner] configuring bootstrap node(s):
  Waiting for Cluster control plane to be initialized, waiting for cluster agent to connect
```

### Candidate root causes within Calico v3.32.0

Since the toleration behavior is identical to v3.31.500, the regression must be in one of these Calico v3.32.0 changes:

#### Candidate A: Felix/calico-node startup or readiness changes

Calico v3.32.0 introduced major changes to Felix:
- **Route programming handoff from BIRD to Felix** (`ProgramClusterRoutes` option in `BGPConfiguration`)
- **Live migration support** with new per-workload state machine
- **ClusterNetworkPolicy support** with new static tiers (`kube-admin`, `kube-baseline`)
- **conntrack handling changes** for eBPF mode

Any of these could add new initialization requirements that make Felix/calico-node take longer to become Ready, or require state from Typha that isn't available during the post-restore bootstrap window.

**Investigation needed**: Compare `calico-node` pod readiness timing between v3.31.5 and v3.32.0 in a post-restore scenario.

#### Candidate B: tigera-operator v1.42.0 reconciliation changes

The operator v1.42.0 may handle post-restore reconciliation differently. After cluster-reset:
1. The restored etcd has old Installation CR, old Deployments, old DaemonSets
2. The operator needs to reconcile all of these with the current state
3. If v1.42.0 has stricter initialization requirements or different ordering, it may fail to bring up Calico after restore

**Investigation needed**: Check operator logs after etcd restore to see if reconciliation fails or hangs.

#### Candidate C: CRD chart packaging change

In v3.32.0, the `rke2-calico-crd` chart was separated with a new `package: rke2-calico-crd` field in `chart_versions.yaml`. During bootstrap after restore, if the CRD chart fails to apply or applies after the operator chart, the operator may not be able to create/reconcile the Installation CR.

```yaml
# v3.31.500 (April):
  - version: v3.31.500
    filename: /charts/rke2-calico-crd.yaml
    bootstrap: true
    # No explicit package field

# v3.32.001 (May):
  - version: v3.32.001
    filename: /charts/rke2-calico-crd.yaml
    bootstrap: true
    package: rke2-calico-crd    # <-- NEW
```

**Investigation needed**: Check if CRDs are properly available after etcd restore with v3.32.0.

#### Candidate D: calico-kube-controllers v3.32.0 startup regression

The v3.32.0 release fixed "calico-kube-controllers IPAM GC controller getting stuck during rapid scale-down" ([calico#11906](https://github.com/projectcalico/calico/pull/11906)). This fix or other changes may have introduced new startup dependencies that cause the controller to crash-loop during the post-restore phase when networking isn't yet available.

**Investigation needed**: Check calico-kube-controllers pod status and logs after etcd restore.

#### Candidate E: Typha v3.32.0 wire protocol change

Calico v3.32.0 added a Typha fix that "rejects oversized inbound client gob frames before reading them" ([calico#12590](https://github.com/projectcalico/calico/pull/12590)). This changes the wire protocol between Felix and Typha. If the restored cluster has a mix of old and new Typha/Felix configurations in etcd, there could be a version mismatch during the bootstrap phase.

---

## Finding 3: The `TolerateBootstrap` Change is NOT in v1.42.0

An important clarification: the [tigera/operator#4820](https://github.com/tigera/operator/pull/4820) PR ("Respect node cordoning for Typha and host-networked Deployments") introduced `TolerateBootstrap` and changed Typha's default from `TolerateAll` to `TolerateBootstrap`. However, this was merged on **May 19, 2026** for milestone **v1.43.0**, NOT v1.42.0.

Both operator v1.40.8 (Calico 3.31.5) and v1.42.0 (Calico 3.32.0) use `TolerateAll` as the Typha default. The `TolerateBootstrap` change is relevant for future releases but is **not the cause of the current regression**.

`TolerateBootstrap` notably includes `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` toleration, which would help with post-restore scenarios. This change should be considered for backport alongside any Calico fix.

---

## CI Test Results (All Attempts)

### `Custom_EtcdSnapshotOperationsOnNewCombinedNode` (etcd restore to new combined node)

| K8s | RKE2 Version | Calico | PR | Run | Result | Duration |
|-----|-------------|--------|-----|-----|--------|----------|
| **April R1 (baseline)** | | | | | | |
| 1.31 | v1.31.x+rke2r1 | v3.31.5 | #2065 (v2.12) | — | ✅ PASS | — |
| 1.32 | v1.32.x+rke2r1 | v3.31.5 | #2065/2066 | — | ✅ PASS | — |
| 1.33 | v1.33.11+rke2r1 | v3.31.5 | #2065 (v2.12) | — | ✅ PASS | 628s |
| **May R1** | | | | | | |
| 1.33 | v1.33.12+rke2r1 | v3.32.0 | #2067 | A1 | ❌ FAIL | 2985s |
| 1.33 | v1.33.12+rke2r1 | v3.32.0 | #2067 | A2 | ❌ FAIL | 3051s |
| 1.34 | v1.34.8+rke2r1 | v3.32.0 | #2067 | A1 | ❌ FAIL | 3044s |
| 1.34 | v1.34.8+rke2r1 | v3.32.0 | #2068 | A1 | ❌ FAIL | — |
| 1.34 | v1.34.8+rke2r1 | v3.32.0 | #2067 | A2 | ❌ FAIL | 2975s |
| 1.34 | v1.34.8+rke2r1 | v3.32.0 | #2068 | A2 | ❌ FAIL | — |
| 1.35 | v1.35.5+rke2r1 | v3.32.0 | #2067 | A1 | ❌ FAIL | 2982s |
| 1.35 | v1.35.5+rke2r1 | v3.32.0 | #2068 | A1 | ❌ FAIL | 3042s |
| 1.35 | v1.35.5+rke2r1 | v3.32.0 | #2067 | A2 | ❌ FAIL | — |
| 1.35 | v1.35.5+rke2r1 | v3.32.0 | #2068 | A2 | ❌ FAIL | — |
| 1.36 | v1.36.1+rke2r1 | v3.32.0 | #2068 | A1 | ❌ FAIL | 3197s |
| 1.36 | v1.36.1+rke2r1 | v3.32.0 | #2068 | A2 | ❌ FAIL | 3037s |
| **May R2 (toleration fix)** | | | | | | |
| 1.34 | v1.34.8+rke2r2 | v3.32.002 | #2072 | A2 | ❌ FAIL | 3056s |
| 1.34 | v1.34.8+rke2r2 | v3.32.002 | #2072 | A3 | ❌ FAIL | 2977s |
| 1.35 | v1.35.5+rke2r2 | v3.32.002 | #2072 | A2 | ❌ FAIL | 3033s |
| 1.35 | v1.35.5+rke2r2 | v3.32.002 | #2072 | A3 | ❌ FAIL | — |
| 1.36 | v1.36.1+rke2r2 | v3.32.002 | #2072 | A2 | ❌ FAIL | — |
| 1.36 | v1.36.1+rke2r2 | v3.32.002 | #2072 | A3 | ❌ FAIL | — |

**100% failure rate with Calico v3.32.0 (both R1 and R2), 100% pass rate with Calico v3.31.5.**

### `Single_Node_All_Roles_Drain`

| K8s | RKE2 Version | Calico | PR | Run | Result |
|-----|-------------|--------|-----|-----|--------|
| **April R1** | | | | | |
| 1.33–1.35 | April+rke2r1 | v3.31.5 | #2065–2067 | — | ✅ PASS |
| **May R1** | | | | | |
| 1.33–1.36 | May+rke2r1 | v3.32.001 | #2067/2068 | A1,A2 | ❌ FAIL (100%) |
| **May R2** | | | | | |
| 1.33–1.36 | May+rke2r2 | v3.32.002 | #2067/2072 | All | ✅ PASS (100%) |

**Toleration fix definitively resolves the Drain test.** ✅

### `MP_EtcdSnapshotOperationsWithThreeEtcdNodesOnNewNode`

| K8s | Calico | PR | Run | Result | Duration |
|-----|--------|-----|-----|--------|----------|
| **April R1** | v3.31.5 | All | — | ✅ PASS | — |
| **May R1** | v3.32.001 | #2068 K8s 1.36 | A1 | ✅ PASS | 853s |
| **May R1** | v3.32.001 | #2068 K8s 1.36 | A2 | ✅ PASS | 747s |
| **May R1** | v3.32.001 | Others | — | ❌ FAIL / ⏱️ timeout | ~3100–3600s |
| **May R2** | v3.32.002 | #2067 K8s 1.33 | A1 | ✅ PASS | 618s |
| **May R2** | v3.32.002 | Others | — | ❌ FAIL / ⏱️ timeout | ~3100–3600s |

**Flaky regression** — passes occasionally but fails ~70-80% of the time.

---

## Suggested Fixes (All RKE2-Side)

### Fix 1: Investigate and fix the Calico v3.32.0 post-restore recovery regression (**CRITICAL**)

The etcd restore failure is caused by something in Calico v3.32.0 beyond tolerations. The RKE2 team should:

1. **Reproduce locally** by running the `Custom_EtcdSnapshotOperationsOnNewCombinedNode` test with Calico v3.31.5 (passes) and v3.32.0 (fails)

2. **Capture diagnostic data** during the failing restore:
   - `kubectl get pods -A -o wide` immediately after restore
   - `kubectl describe pod <typha-pod>` and `kubectl describe pod <calico-node-pod>`
   - `kubectl get nodes -o jsonpath='{.items[*].spec.taints}'`
   - tigera-operator logs
   - calico-node (Felix) logs
   - calico-kube-controllers logs

3. **Bisect the Calico changes** — the most productive approach:
   - Build an rke2-calico chart with Calico v3.32.0 images but operator v1.40.8 (isolates operator changes)
   - Build with operator v1.42.0 but Calico v3.31.5 node/typha images (isolates component changes)
   - Test each combination against the etcd restore test

4. **Check CRD bootstrap ordering** — ensure the `rke2-calico-crd` chart is properly installed before the `rke2-calico` chart during post-restore bootstrap

### Fix 2: Add CCM uninitialized taint toleration to `controlPlaneTolerations` (**RECOMMENDED**)

The current v3.32.002 fix sets `controlPlaneTolerations` with only two tolerations. The upcoming operator v1.43.0 `TolerateBootstrap` includes `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule`, which is essential during bootstrap and post-restore scenarios (before the CCM initializes the node).

Update the rke2-calico chart overlay:

```yaml
controlPlaneTolerations:
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node-role.kubernetes.io/etcd"
    operator: "Exists"
    effect: "NoExecute"
  - key: "node.cloudprovider.kubernetes.io/uninitialized"   # <-- ADD
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node.kubernetes.io/not-ready"                      # <-- ADD
    operator: "Exists"
    effect: "NoSchedule"
  - key: "node.kubernetes.io/network-unavailable"            # <-- ADD
    operator: "Exists"
    effect: "NoSchedule"
```

This aligns with what the operator v1.43.0 will use as the default (`TolerateBootstrap`), and ensures Typha and other control-plane Calico components can schedule on nodes during bootstrap, before the CCM has initialized them and before the CNI is fully up.

**Important**: This is an RKE2 chart change, not a Rancher change. It's delivered via the rke2-calico chart overlay and takes effect with any Rancher version via KDM.

### Fix 3: Consider reverting to Calico v3.31.5 for the May release cycle (**FALLBACK**)

If the root cause of the etcd restore regression can't be identified quickly, reverting the Calico v3.31.5 → v3.32.0 upgrade would restore all tests to passing. This would be a temporary measure until the v3.32.0 issue is resolved.

This revert would be an RKE2-side change (reverting the CNI update PRs) and would be delivered to all Rancher versions via KDM.

---

## Why Fixes Must Be RKE2-Side

New RKE2 versions are delivered to existing Rancher installations via KDM (Kontainer Driver Metadata) updates. This means:

- Rancher v2.12, v2.13, v2.14, and v2.15 all receive new RKE2 versions without upgrading Rancher
- Any fix that requires a Rancher code change would NOT protect existing Rancher installations
- The regression was introduced by an RKE2-side change (Calico upgrade), so it MUST be fixed on the RKE2 side
- The rke2-calico chart overlay is the correct fix location — it's part of the RKE2 binary, not Rancher

---

## Appendix A: The Calico v3.31.5 → v3.32.0 Change Detail

### Chart values.yaml structural change

Calico v3.32.0 significantly expanded the `installation:` section in the tigera-operator Helm chart values.yaml:

**v3.31.5** (21 lines under `installation:`):
```yaml
installation:
  enabled: true
  kubernetesProvider: ""
  imagePullSecrets: []
  kubeletVolumePluginPath: "None"
```

**v3.32.0** (77 lines under `installation:`):
```yaml
installation:
  enabled: true
  # ... (many new fields)
  controlPlaneReplicas: 2           # NEW
  controlPlaneNodeSelector: {}       # NEW
  controlPlaneTolerations: []        # NEW <-- this caused the duplicate key bug
  nonPrivileged: "Disabled"          # NEW
  # ... (more new fields)
```

The addition of `controlPlaneTolerations: []` as an explicit upstream field is what caused the YAML duplicate key conflict with the rke2-calico patch.

### tigera-operator: v1.40.8 vs v1.42.0 Typha toleration logic

**IDENTICAL** in both versions:
```go
// pkg/render/typha.go
tolerations := rmeta.TolerateAll
if len(c.cfg.Installation.ControlPlaneTolerations) != 0 {
    tolerations = c.cfg.Installation.ControlPlaneTolerations
}
```

The `TolerateBootstrap` default (which includes CCM taint toleration) is only available from operator v1.43.0 onward ([tigera/operator#4820](https://github.com/tigera/operator/pull/4820), merged May 19, 2026).

### Key Calico v3.32.0 features that could affect restore

| Feature | Description | Risk to Restore |
|---------|-------------|-----------------|
| Route programming handoff to Felix | New `ProgramClusterRoutes` option delegates cluster route programming from BIRD to Felix | May add startup dependencies |
| Live migration support | New per-workload state machine in Felix | Additional Typha state requirements |
| ClusterNetworkPolicy | New static tiers and removed old ANP support | Different policy initialization |
| CRD chart separation | CRDs now in separate `crd.projectcalico.org.v1` chart | Bootstrap ordering risk |
| calico-kube-controllers IPAM fix | Fix for GC controller stuck during scale-down | May affect startup during restore |
| Typha frame size limit | Rejects oversized gob frames | Wire protocol change |

## Appendix B: Kubernetes DefaultTolerationSeconds Admission Controller

The admission controller ([source](https://github.com/kubernetes/kubernetes/blob/master/plugin/pkg/admission/defaulttolerationseconds/admission.go#L68-L73)) automatically adds:

```yaml
- key: node.kubernetes.io/not-ready
  effect: NoExecute
  tolerationSeconds: 300
- key: node.kubernetes.io/unreachable
  effect: NoExecute
  tolerationSeconds: 300
```

**But only if** the pod doesn't already have a toleration matching `NoExecute` without a specific key. The blanket `{Effect: NoExecute, Operator: Exists}` in `TolerateAll` matches all NoExecute taints, so the admission controller skips adding the 300s tolerations. With key-specific tolerations (like in v3.31.500 and v3.32.002), the admission controller correctly adds them.

## Appendix C: Calico v3.32.0 Image Changes in RKE2

From the CNI update PRs (e.g., [#10350](https://github.com/rancher/rke2/pull/10350)):

```diff
# Standalone Calico images:
- rancher/mirrored-calico-operator:v1.40.8
+ rancher/mirrored-calico-operator:v1.42.0
- rancher/mirrored-calico-typha:v3.31.5
+ rancher/mirrored-calico-typha:v3.32.0
- rancher/mirrored-calico-node:v3.31.5
+ rancher/mirrored-calico-node:v3.32.0
- rancher/mirrored-calico-kube-controllers:v3.31.5
+ rancher/mirrored-calico-kube-controllers:v3.32.0
- rancher/mirrored-calico-cni:v3.31.5
+ rancher/mirrored-calico-cni:v3.32.0

# Canal images:
- rancher/hardened-calico:v3.31.5-build20260415
+ rancher/hardened-calico:v3.32.0-build20260507

# rke2-calico chart:
- v3.31.500
+ v3.32.001  (buggy, see Finding 1)
# Then fixed to:
+ v3.32.002  (toleration fix only)
```
