# Root Cause Analysis: Calico v3.31.5 → v3.32.0 Upgrade in May RKE2 Patches

## Summary

The `Test_Operation_SetB_Custom_EtcdSnapshotOperationsOnNewCombinedNode` failure that was first observed in RKE2 v1.36.0 now affects **all May RKE2 patches** (v1.33.12, v1.34.8, v1.35.5, v1.36.1). Investigation confirms that the common denominator is the **Calico v3.31.5 → v3.32.0 upgrade** and the accompanying **tigera-operator v1.40.8 → v1.42.0 upgrade**, which shipped in all May patches simultaneously.

---

## The Calico Version Change: April vs May Patches

| RKE2 Version | Release Month | Calico Version | Canal Chart | tigera-operator | Result |
|-------------|---------------|----------------|-------------|-----------------|--------|
| v1.33.11+rke2r1 | April 2026 | **v3.31.5** | v3.31.5-build2026041500 | v1.40.8 | ✅ Tests pass (on v2.12) |
| v1.34.7+rke2r1 | April 2026 | **v3.31.5** | v3.31.5-build2026041500 | v1.40.8 | ✅ Tests pass |
| v1.35.4+rke2r1 | April 2026 | **v3.31.5** | v3.31.5-build2026041500 | v1.40.8 | ✅ Tests pass |
| v1.33.12+rke2r1 | May 2026 | **v3.32.0** | v3.32.0-build2026051100 | v1.42.0 | ❌ Fails (v2.13+) |
| v1.34.8+rke2r1 | May 2026 | **v3.32.0** | v3.32.0-build2026051100 | v1.42.0 | ❌ Fails |
| v1.35.5+rke2r1 | May 2026 | **v3.32.0** | v3.32.0-build2026051100 | v1.42.0 | ❌ Fails |
| v1.36.0+rke2r1 | May 2026 | **v3.32.0** | v3.32.0-build2026051100 | v1.42.0 | ❌ Fails |
| v1.36.1+rke2r1 | May 2026 | **v3.32.0** | v3.32.0-build2026051100 | v1.42.0 | ❌ Fails |

### How the Calico bump propagated

1. **RKE2 v1.36.0+rke2r1** was the first release to include Calico v3.32.0 (via [rancher/rke2#10384](https://github.com/rancher/rke2/pull/10384) — "Update CNIs for 2026-05 Release Cycle"). Released ~May 14, 2026.
2. The May patch cycle then backported the same Calico v3.32.0 to all active branches via parallel PRs:
   - v1.33.12: [rancher/rke2#10350](https://github.com/rancher/rke2/pull/10350) + [#10387](https://github.com/rancher/rke2/pull/10387)
   - v1.34.8: [rancher/rke2#10349](https://github.com/rancher/rke2/pull/10349) + [#10386](https://github.com/rancher/rke2/pull/10386)
   - v1.35.5: [rancher/rke2#10348](https://github.com/rancher/rke2/pull/10348) + [#10385](https://github.com/rancher/rke2/pull/10385)
3. All April patches shipped **Calico v3.31.5** (via [rancher/rke2#10227](https://github.com/rancher/rke2/pull/10227) for 1.33, etc.).
4. The previous April patches (v1.33.11, v1.34.7, v1.35.4) all **passed** the `Custom_EtcdSnapshotOperationsOnNewCombinedNode` test.

### Specific image changes in the CNI update PRs

From [rancher/rke2#10350](https://github.com/rancher/rke2/pull/10350/files) (representative of all branches):

```diff
# Canal images (default CNI):
- rancher/hardened-calico:v3.31.5-build20260415
+ rancher/hardened-calico:v3.32.0-build20260507

# Standalone Calico images:
- rancher/mirrored-calico-operator:v1.40.8
+ rancher/mirrored-calico-operator:v1.42.0
- rancher/mirrored-calico-typha:v3.31.5
+ rancher/mirrored-calico-typha:v3.32.0
- rancher/mirrored-calico-node:v3.31.5
+ rancher/mirrored-calico-node:v3.32.0
- rancher/mirrored-calico-kube-controllers:v3.31.5
+ rancher/mirrored-calico-kube-controllers:v3.32.0
# (+ all other calico component images bumped similarly)
```

---

## Calico v3.32.0 Release: What Changed

Calico v3.32.0 ([release notes](https://github.com/projectcalico/calico/blob/release-v3.32/release-notes/v3.32.0-release-notes.md)) is a **major feature release**, not a patch. Key changes relevant to the failure:

### Changes most likely to affect single-node restore/drain scenarios

1. **"Hand off cluster route programming from BIRD to Felix"** — New `ProgramClusterRoutes` option in `BGPConfiguration` changes how cluster node routes are programmed. When set to `Disabled`, Felix programs the routes instead of confd/BIRD. This fundamentally changes the networking initialization path.

2. **"Seamless live migration support for KubeVirt VMs"** — Introduces a new `LiveMigration` API and per-workload state machine in Felix that adjusts route priorities. This adds new state tracking in Felix that could interact with node replacement scenarios.

3. **"Reset TCP connections when a backend pod is deleted (eBPF)"** — Changes conntrack handling behavior when pods are deleted, which happens extensively during etcd restore cleanup.

4. **"ClusterNetworkPolicy support"** — New static tiers (`kube-admin`, `kube-baseline`) and removal of old `AdminNetworkPolicy`/`BaselineAdminNetworkPolicy` support. This changes policy initialization.

5. **Breaking change: CRD chart separation** — `tigera-operator` Helm chart no longer includes CRDs; they're in a separate `crd.projectcalico.org.v1` chart. This could affect chart installation ordering.

### Typha-specific changes

The only Typha-specific change in v3.32.0:

> **Bug fix**: Typha now rejects oversized inbound client gob frames before reading them, preventing a potential denial-of-service caused by excessive memory allocation. ([calico#12590](https://github.com/projectcalico/calico/pull/12590))

This is a security fix in Typha's wire protocol. While it shouldn't cause the failure, it changes the connection handling behavior between Felix (calico-node) and Typha.

### Other relevant bug fixes

- Fix calico-kube-controllers IPAM GC controller getting stuck during rapid scale-down ([calico#11906](https://github.com/projectcalico/calico/pull/11906)) — directly relevant to node replacement
- Fix memory leak in routing table logic when interfaces are removed ([calico#12138](https://github.com/projectcalico/calico/pull/12138)) — relevant to node teardown
- Fix goroutine leak in Felix's interface monitor on netlink reconnect ([calico#12139](https://github.com/projectcalico/calico/pull/12139)) — relevant to node restart

---

## Typha Deployment Configuration: Unchanged Between Versions

Detailed investigation of the tigera-operator code confirms that **Typha deployment configuration did NOT change** between operator v1.40.8 and v1.42.0:

| Setting | v1.40.8 (Calico 3.31.5) | v1.42.0 (Calico 3.32.0) |
|---------|-------------------------|-------------------------|
| Deployment strategy | RollingUpdate (MaxUnavailable=1, MaxSurge=100%) | **Same** |
| Replica count | Autoscaler-managed: 1-2 nodes→1, 3-4→2, 5+→3 | **Same** |
| Tolerations | `rmeta.TolerateAll` unless `controlPlaneTolerations` set | **Same** |
| `controlPlaneReplicas` | Defaults to 2 at runtime (does NOT control Typha) | Made explicit in chart YAML, same runtime default |

The `controlPlaneReplicas: 2` field that was newly added to the upstream v3.32.0 chart's `values.yaml` controls **APIServer, Dex, Webhooks, and Linseed** replicas — NOT Typha. Typha replicas are exclusively managed by the autoscaler goroutine in `pkg/common/autoscale.go`.

### Canal vs Standalone Calico: Typha behavior

**Important distinction:** The default RKE2 CNI is **Canal**, which explicitly **disables Typha**:

```yaml
# rke2-canal chart values.yaml
calico:
  # Typha is disabled.
  typhaServiceName: none
```

The user's additional research references `calico-typha` pods, which means the failing test clusters are using **standalone Calico** (rke2-calico), not Canal. In standalone Calico mode, the tigera-operator deploys Typha with the autoscaler determining replicas (1 replica for ≤2 node clusters).

---

## Failure Mechanism: Calico Typha as Single Point of Failure

Based on the user's additional research, the failure chain in the single-node combined-role etcd restore scenario is:

```
1. Etcd restore starts → old node is drained/killed
   ↓
2. Calico Typha pod (1 replica) is stuck Terminating on old (NotReady) node
   ↓
3. New node comes up → calico-node DaemonSet pod starts
   ↓
4. calico-node (Felix) cannot connect to Typha:
   "Didn't find any ready Typha instances"
   "Typha discovery enabled but discovery failed"
   ↓
5. Felix cannot bootstrap → calico-node stays 0/1 Running (not ready)
   ↓
6. No Calico dataplane programming on new node → pod networking broken
   ↓
7. CoreDNS pods can't serve DNS (network broken)
   ↓
8. cattle-cluster-agent gets "Could not resolve host" → crash-loops
   ↓
9. Rancher: "Cluster agent is not connected" → planner stuck
   ↓
10. Test times out after ~50 minutes ❌
```

### Why this is a chicken-and-egg problem

- **Typha** needs to be rescheduled to the new node, but it's stuck Terminating on the old (unreachable) node.
- **Felix** (calico-node) on the new node needs Typha to bootstrap its dataplane programming.
- Without Felix, **pod networking** doesn't work on the new node.
- Without networking, **cattle-cluster-agent** can't reach Rancher DNS/endpoint.
- Without the agent, **Rancher planner** can't proceed with cluster initialization.
- Without cluster initialization, the **old node** stays NotReady and its pods stay Terminating.

### Why MP (Machine Pool) tests pass

In MP tests with separate node pools, the **worker node stays alive** during the etcd restore. Its calico-node already has a working dataplane, its Typha connection is already established, and cattle-cluster-agent can be scheduled there. The single-node combined-role test has no such fallback.

---

## The Critical Question: Why Did Calico 3.31.5 Work But 3.32.0 Fails?

Typha has been a single point of failure for single-node clusters in **both** versions. The autoscaler logic is identical. So what changed?

### Hypothesis 1: Felix's Typha Dependency Tightened in v3.32.0 (MOST LIKELY)

Calico v3.32.0's major new features (route programming handoff to Felix, live migration state machine, ClusterNetworkPolicy tiers) likely **increased Felix's dependency on the Typha connection for initialization**. In v3.31.5, Felix may have been able to:
- Bootstrap basic networking from cached state or local configuration
- Program routes independently via BIRD/confd without needing Typha
- Proceed with partial functionality even when Typha was unreachable

In v3.32.0, Felix may now **require** the Typha connection to:
- Receive the new `ProgramClusterRoutes` configuration
- Initialize the new ClusterNetworkPolicy tiers (`kube-admin`, `kube-baseline`)
- Sync the new LiveMigration state
- Get route programming instructions (since route programming can now be delegated to Felix)

If Felix in v3.32.0 blocks on Typha connection during startup instead of degrading gracefully, this would explain why the same Typha SPOF scenario that was tolerable in v3.31.5 becomes fatal in v3.32.0.

### Hypothesis 2: Calico-kube-controllers Startup Regression

The v3.32.0 release fixed "calico-kube-controllers IPAM GC controller getting stuck during rapid scale-down" ([calico#11906](https://github.com/projectcalico/calico/pull/11906)). The fix may have introduced a new dependency or startup sequencing requirement that causes the controller to be more sensitive to networking availability during node replacement.

The user's research notes that `calico-kube-controllers` replacement pods on the new node are in **CrashLoopBackOff** — potentially caused by either the broken networking or a new startup dependency in v3.32.0.

### Hypothesis 3: CRD Chart Separation Timing Issue

Calico v3.32.0's breaking change — moving CRDs out of the tigera-operator Helm chart into a separate `crd.projectcalico.org.v1` chart — could cause CRD availability issues during the restore sequence if the chart installation order isn't properly managed. If the CRDs aren't available when the operator starts after restore, Typha and other components might fail to initialize.

### The v2.12 vs v2.13 Anomaly

The pass/fail matrix shows K8s 1.33 (v1.33.12, Calico 3.32.0) **passes on Rancher v2.12** but **fails on v2.13+**. This means the Calico v3.32.0 change is **necessary but not sufficient** to trigger the failure — there's also a Rancher-side factor.

Possible explanations:
1. **Rancher v2.13 changed the restore/drain procedure** in a way that is more disruptive to Calico's recovery (e.g., more aggressive pod cleanup, different node drain ordering, changed toleration handling).
2. **The provisioning test code differs between v2.12 and v2.13** — the test might configure the cluster differently, use different CNI settings, or have different timeouts.
3. **Rancher v2.13 uses a different chart installation flow** that exposes the CRD separation issue (Hypothesis 3).

---

## Recommended Investigation Steps

### Immediate (to confirm root cause)

1. **Compare Felix startup behavior between Calico v3.31.5 and v3.32.0** when Typha is unreachable:
   ```bash
   # In a test cluster with standalone Calico, kill Typha and watch Felix logs:
   kubectl delete pod -n calico-system calico-typha-xxxxx
   # Check if Felix in v3.32.0 blocks harder on Typha reconnection than v3.31.5
   ```

2. **Check if the test configures standalone Calico or Canal** — look at the test's cluster config to confirm CNI selection. If it uses Canal, Typha isn't involved and the failure mechanism is different.

3. **Run the failing test with Calico v3.31.5 on May K8s patches** — if the test passes with Calico 3.31.5 on the same K8s version, the Calico upgrade is confirmed as the root cause.

### Short-term fixes

4. **Pin Calico to v3.31.5 in May patches** — if feasible, downgrade Calico back to v3.31.5 to unblock the May patch release while the root cause is investigated.

5. **Add a PodDisruptionBudget or affinity rule for Typha** — ensure Typha is rescheduled before the old node is fully drained. This would require a chart configuration change.

6. **Increase Typha replicas for single-node clusters** — override the autoscaler minimum from 1 to 2 for clusters where control-plane and worker roles are combined, so there's redundancy during node replacement.

### Medium-term fixes

7. **Ensure Felix can degrade gracefully without Typha** — if this is a v3.32.0 regression, report it upstream to projectcalico/calico. Felix should be able to bootstrap basic networking from local/cached state.

8. **Add Typha pod cleanup to the restore procedure** — during etcd restore cleanup, explicitly delete Typha pods so they're rescheduled on the new node instead of being stuck Terminating on the old node.

9. **Investigate the v2.12/v2.13 Rancher-side boundary** — diff the Rancher planner/CAPR code between v2.12 and v2.13 to identify what restore procedure change exposes the Calico 3.32.0 issue on v2.13 but not v2.12.

---

## Appendix: Full Component Version Comparison

### April Patches (Working)

| Component | v1.33.11 | v1.34.7 | v1.35.4 |
|-----------|----------|---------|---------|
| Calico | v3.31.5 | v3.31.5 | v3.31.5 |
| Canal | v3.31.5-build2026041500 | v3.31.5-build2026041500 | v3.31.5-build2026041500 |
| tigera-operator | v1.40.8 | v1.40.8 | v1.40.8 |
| Cilium | v1.19.3 | v1.19.3 | v1.19.3 |
| CoreDNS | 1.45.209 | 1.45.209 | 1.45.209 |
| Flannel | v0.28.4 | v0.28.4 | v0.28.4 |
| Ingress-Nginx | 4.14.504 | 4.14.504 | 4.14.504 |
| Metrics-Server | 3.13.008 | 3.13.008 | 3.13.008 |
| Traefik | 39.0.701 | 39.0.701 | 39.0.701 |
| snapshot-controller | 4.2.003 | 4.2.003 | 4.2.003 |

### May Patches (Failing)

| Component | v1.33.12 | v1.34.8 | v1.35.5 | v1.36.0 |
|-----------|----------|---------|---------|---------|
| **Calico** | **v3.32.0** | **v3.32.0** | **v3.32.0** | **v3.32.0** |
| **Canal** | **v3.32.0-build2026051100** | **v3.32.0-build2026051100** | **v3.32.0-build2026051100** | **v3.32.0-build2026051100** |
| **tigera-operator** | **v1.42.0** | **v1.42.0** | **v1.42.0** | **v1.42.0** |
| Cilium | v1.19.303 | v1.19.303 | v1.19.303 | v1.19.303 |
| CoreDNS | **1.45.212** | **1.45.212** | **1.45.212** | **1.45.212** |
| Flannel | v0.28.4 | v0.28.4 | v0.28.4 | v0.28.4 |
| Ingress-Nginx | **4.14.508** | **4.14.508** | **4.14.508** | 4.14.506 |
| Metrics-Server | **3.13.011** | **3.13.011** | **3.13.011** | 3.13.008 |
| Traefik | **39.0.703** | **39.0.703** | **39.0.703** | — |
| snapshot-controller | **4.2.005** | **4.2.005** | **4.2.005** | 4.2.003 |

The **Calico/Canal/tigera-operator** version change is the most significant delta, being a major feature release (v3.31→v3.32) rather than a patch. Other changes (CoreDNS, Traefik, Ingress-Nginx, Metrics-Server, snapshot-controller) are minor version bumps unlikely to affect networking or etcd restore behavior.

---

## Timeline

| Date | Event |
|------|-------|
| ~April 23, 2026 | April patches released (v1.33.11, v1.34.7, v1.35.4) with Calico v3.31.5 |
| April 30, 2026 | Calico v3.32.0 released upstream by Project Calico |
| May 8, 2026 | "CNI update May release" PRs open in rancher/rke2, bumping Calico to v3.32.0 |
| May 12, 2026 | "Update CNIs for 2026-05 Release Cycle" PRs merged (Go bump + final CNI updates) |
| ~May 14, 2026 | RKE2 v1.36.0+rke2r1 released — first release with Calico v3.32.0 |
| ~May 18, 2026 | May patches released (v1.33.12, v1.34.8, v1.35.5) — all with Calico v3.32.0 |
| May 18-20, 2026 | KDM PRs #2065-#2068 opened; provisioning tests run and fail |
