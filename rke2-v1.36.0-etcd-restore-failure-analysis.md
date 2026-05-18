# RKE2 v1.36.0 — `Custom_EtcdSnapshotOperationsOnNewCombinedNode` Failure Analysis

## Executive Summary

The test fails during the **etcd snapshot restore phase**. The Rancher planner gets permanently stuck at:

> *"Waiting for Cluster control plane to be initialized, waiting for cluster agent to connect"*

The `cattle-cluster-agent` in the downstream cluster never reconnects to Rancher after the restore, creating a deadlock that times out after ~42 minutes.

---

## Timeline of Events

| Time (UTC) | Event |
|------------|-------|
| 20:35:02 | Test starts, initial cluster provisioning begins |
| 20:37:06 | Cluster agent connects — initial provisioning succeeds ✅ |
| 20:39:28 | Etcd snapshot taken (`on-demand-control-etcd-test-node-1778791168`) |
| 20:40:07 | Old control plane node's RKE2 detects server failure |
| 20:40:12 | Rancher **uninitializes rkecontrolplane**, force-deletes old etcd machine |
| 20:40:19 | New machine `custom-1d328b12d63d` registered on node `test-node-qdjnj` |
| 20:40:21 | **Worker node `test-node-qfvzp` killed** (`rke2-killall.sh`) — never restarted |
| 20:40:23 | Restore plan starts on control plane node |
| 20:41:03 | `rke2 server --cluster-reset` completes successfully |
| 20:41:09 | `systemctl --no-block restart rke2-server` |
| 20:41:29 | RKE2 v1.36.0+rke2r1 starts, etcd bootstrap recognized |
| 20:41:32 | kube-apiserver, kube-controller-manager, kube-scheduler, **cloud-controller-manager** all started |
| 20:41:33 | **`rke2-cloud-provider:v1.36.0-rc2...` image import FAILS** — `not found` in airgap import |
| 20:41:35 | Kubelet starts **with** `--cloud-provider=external` (adds uninitialized taint) |
| 20:41:49 | Node `test-node-qdjnj` successfully registered |
| 20:41:52 | **"Waiting for untainted node"** — CCM hasn't removed the cloud-provider taint yet |
| 20:41:55 | `cluster-agent.yaml` manifest applied |
| 20:41:56 | CCM static pod starts (using previously-cached image) ✅ |
| 20:43:07 | `wait_for_ready.sh` completes — node reports Ready |
| 20:43:07–09 | Post-restore cleanup deletes: CoreDNS, calico, metrics-server, ingress, **cattle-cluster-agent** pods |
| 20:43:10 | CoreDNS replacement pods start ✅ |
| 20:43:11 | Calico-node replacement pods start ✅ |
| 20:43:15 | Rancher planner: "running full reconcile during etcd restore to initially restart cluster" |
| 20:43:18 | **STUCK: "Waiting for Cluster control plane to be initialized, waiting for cluster agent to connect"** |
| 20:43:18 → 21:25 | Cluster agent **never reconnects**. Test times out. ❌ |

---

## Key Findings

### 1. `rke2-cloud-provider` Image Import Fails (non-blocking but suspicious)

```
Failed to process image event: failed to import cloud-controller-manager-image.txt:
image "index.docker.io/rancher/rke2-cloud-provider:v1.36.0-rc2.0.20260427154526-d239025e2a23-build20260429": not found
```

- The image is **not in the airgap tarball** on disk.
- However, the CCM static pod **does start** at 20:41:56 using a **previously-cached image** in containerd (confirmed by kubelet's `pod_startup_latency_tracker` showing 0ms image pull time).
- This error is therefore **not the direct cause**, but it indicates this image version is not being pre-staged correctly for airgap/offline use cases — worth investigating separately.

### 2. `"Waiting for untainted node"` Is Never Resolved

The K3s/RKE2 server code (`pkg/daemons/control/server.go`) starts a goroutine before the kube-scheduler that watches all nodes, waiting for any node to **not** have the `node.cloudprovider.kubernetes.io/uninitialized` taint:

```go
func waitForUntaintedNode(ctx context.Context, kubeConfig string) error {
    // Watches nodes until one is found without TaintExternalCloudProvider
    ...
}
```

This message is logged once at **20:41:52** and is **never followed by a resolution message**. This strongly suggests the cloud-controller-manager is **not removing the `TaintExternalCloudProvider` taint** from node `test-node-qdjnj`.

### 3. Worker Node Is Killed and Never Restarted

- Worker node `test-node-qfvzp` is killed at 20:40:21 via `rke2-killall.sh`.
- Its last log entry is at 20:40:23 — it is **never restarted**.
- The Rancher planner will not restart workers until the control plane is re-initialized — **creating a chicken-and-egg problem**.
- With the worker dead, the **only available node** for pod scheduling is the control-plane node `test-node-qdjnj`.

### 4. Cattle-Cluster-Agent Pods Are Deleted But Never Replaced

- The cleanup at 20:43:09 deletes the old cattle-cluster-agent pods.
- Replacement pods for CoreDNS (20:43:10) and calico-node (20:43:11) are confirmed started.
- **No cattle-cluster-agent pod startup is recorded** in the remaining 8 seconds of captured node logs.
- Rancher management logs confirm the agent **never reconnects** through the entire 42-minute timeout.

### 5. K3s 1.36 Passes — RKE2's CCM Is the Differentiator

- K3s does **not** use `--cloud-provider=external` and does **not** run a cloud-controller-manager.
- The K3s v1.36.0 etcd restore tests **pass**.
- RKE2's CCM and the external cloud provider taint interaction is the unique failure point.

### 6. Evidence That the CI Registry Cache Is NOT the Root Cause

- The images `rancher/hardened-kubernetes:v1.36.0-rke2r1-build20260429` and `rancher/hardened-etcd:v3.6.7-k3s1-build20260415` are **confirmed already pulled** in containerd (logs say "has already been pulled").
- Initial cluster provisioning **succeeded** — which would have failed if critical images were missing from the registry.
- The `registry-cache` is reachable and returns `MANIFEST_UNKNOWN` (not a connection error) for missing images, indicating fall-through to Docker Hub works.
- Therefore, the images being available on DockerHub is consistent with passing initial provisioning, and **missing images are not the cause of the restore failure**.

---

## Root Cause Hypotheses

### Hypothesis 1 (Most Likely): CCM Fails to Remove the Cloud-Provider Taint After Restore

**Flow:**

1. Kubelet restarts with `--cloud-provider=external`, adding `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` to the node.
2. The CCM static pod starts at 20:41:56.
3. The CCM **fails to initialize node `test-node-qdjnj`** because:
   - The restored etcd data contains stale node objects from the **old** node (`control-etcd-test-node`).
   - The CCM's cloud-provider state is inconsistent — it may be confused about which nodes it should be managing.
   - The RKE2 cloud-provider may skip nodes it believes are already initialized (based on stale etcd data).
4. With `NoSchedule` taint on the only available node, the `cattle-cluster-agent` Deployment pods **cannot be scheduled** (unlike static pods and DaemonSets which tolerate all taints).
5. Without the cluster agent, Rancher cannot connect, and the planner is permanently stuck.

**Evidence:**

- `"Waiting for untainted node"` is logged and never resolved.
- The cattle-cluster-agent pods are never re-created after deletion.
- CoreDNS and calico-node pods that restart successfully are either DaemonSets (tolerate all taints) or were scheduled in a pre-taint window.

### Hypothesis 2: `rancher/rancher-agent:head` Image Is Not in Containerd After Restore

After the restore and containerd restart:

- The `rancher/rancher-agent:head` image (used by cattle-cluster-agent) is a **CI-built image** that may not persist in containerd's content store across a full service restart.
- If the pull fails (image not in registry-cache, `:head` tag not on Docker Hub), the pods would be stuck in `ImagePullBackOff`.
- This cannot be conclusively verified because the node logs end at 20:43:17, only 8 seconds after the pod deletion.

### Hypothesis 3: Stale Node Objects From Snapshot Confuse CCM Initialization

After the cluster-reset, etcd is restored from a snapshot taken from the **old** node. The new kubelet registers `test-node-qdjnj`, but the etcd data still contains node objects for `control-etcd-test-node`. The CCM may:

- Attempt to initialize the old (non-existent) node and succeed, satisfying its internal tracking.
- Skip initializing `test-node-qdjnj` because the restored etcd indicates it was already handled.
- Result: `test-node-qdjnj` retains the `uninitialized` taint permanently.

### Hypothesis 4: K8s 1.36 Cloud-Provider Taint Behavior Change

In Kubernetes 1.31+, the `DisableCloudProviders` and `DisableKubeletCloudCredentialProviders` feature gates went GA (locked to `true`). This means:

- The kubelet's `--cloud-provider` flag may be **silently ignored** in 1.36.
- The taint may now be added/removed **only** by the CCM itself, not the kubelet.
- If the CCM in RKE2 v1.36 has changed its taint-management behavior, it might add the taint but fail to remove it under certain conditions (e.g., after a cluster-reset restore with a new node hostname).

---

## Comparison with Passing Machine Pool (MP) Test

The `MP_EtcdSnapshotOperationsOnNewNode` test **also** goes through "waiting for cluster agent to connect" (at 20:48:31) but **eventually succeeds** at 20:59:28.

**Key difference:** In the MP test:

- The cluster has **separate node pools** for etcd, control-plane, and worker roles.
- The **worker node is NOT killed** during the restore — it remains up and running.
- The `cattle-cluster-agent` Deployment can be scheduled on the **live worker node**, which does not have the cloud-provider uninitialized taint (since it was already initialized before the restore).
- The agent connects from the worker, the control plane gets re-initialized, and the planner proceeds.

In the custom combined-node test:

- The single combined etcd+control-plane+worker node is the **only** node after the restore.
- The worker node is killed and not restarted until after the control plane is initialized.
- If that single node has the `uninitialized` taint blocking the cattle-cluster-agent, **there is no alternative node** to schedule it on.

---

## Recommendations for Further Investigation

1. **Capture node taints post-restore**: Add `kubectl get nodes -o jsonpath='{.items[*].spec.taints}'` to the post-restore diagnostics to confirm whether the `node.cloudprovider.kubernetes.io/uninitialized` taint is present and whether it's ever removed.

2. **Capture CCM pod logs**: The `cloud-controller-manager` static pod logs should show why it's failing (or succeeding at) initializing the node. These logs are not currently captured in the test data bundle.

3. **Capture post-restore pod status**: Run `kubectl get pods -A -o wide` from the downstream cluster after the restore to see if cattle-cluster-agent pods are `Pending` with `SchedulingDisabled` or `ImagePullBackOff`.

4. **Test with explicit toleration**: Add `node.cloudprovider.kubernetes.io/uninitialized:NoSchedule` toleration to the cattle-cluster-agent Deployment in the test and verify if this resolves the issue.

5. **Compare RKE2 1.35 vs 1.36 CCM behavior**: If RKE2 1.35 restores work (to be confirmed), the delta in `rke2-cloud-provider` behavior between the two versions (especially around node taint management after a cluster-reset) is the most promising investigation path.

6. **Check if `rancher/rancher-agent:head` survives containerd restart**: Verify that the cattle-cluster-agent image persists in containerd's content store after the `rke2-killall.sh` + restart sequence.
