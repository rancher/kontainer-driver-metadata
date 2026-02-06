package templates

/*
FlannelTemplateV0_28_1 is based on upstream flannel v0.28.1
Rancher Specifics Preserved:
- Namespace: kube-system (Upstream uses kube-flannel)
- Container: install-cni sidecar (Upstream uses initContainers; RKE uses sidecar)
- CNI Config: Retained forceAddress: true for RKE node detection
- Templates: {{.ClusterCIDR}}, {{.Image}}, etc.
- Fix: Added BlackholeRoute conditional logic
*/

const FlannelTemplateV0_28_1 = `
{{- if eq .RBACConfig "rbac"}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  labels:
    k8s-app: flannel
  name: flannel
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  labels:
    k8s-app: flannel
  name: flannel
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: flannel
subjects:
- kind: ServiceAccount
  name: flannel
  namespace: kube-system
---
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    k8s-app: flannel
  name: flannel
  namespace: kube-system
{{- end}}
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    tier: node
    k8s-app: flannel
    app: flannel
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "flannel",
          "delegate": {
            "forceAddress": true,
            "hairpinMode": true,
            "isDefaultGateway": true
          }
        },
        {
          "type": "portmap",
          "capabilities": {
            "portMappings": true
          }
        }
      ]
    }
  net-conf.json: |
    {
      "Network": "{{.ClusterCIDR}}",
      "EnableNFTables": false,
      "Backend": {
        "Type": "{{.FlannelBackend.Type}}",
        "VNI": {{.FlannelBackend.VNI}},
        "Port": {{.FlannelBackend.Port}}
      }
    }
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: kube-flannel
  namespace: kube-system
  labels:
    tier: node
    k8s-app: flannel
    app: flannel
spec:
  selector:
    matchLabels:
      k8s-app: flannel
  template:
    metadata:
      labels:
        tier: node
        k8s-app: flannel
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/os
                operator: In
                values:
                - linux
{{if .NodeSelector}}
      nodeSelector:
      {{ range $k, $v := .NodeSelector }}
        {{ $k }}: "{{ $v }}"
      {{ end }}
{{end}}
      hostNetwork: true
      priorityClassName: {{ .KubeFlannelPriorityClassName | default "system-node-critical" }}
      tolerations:
      {{- if ge .ClusterVersion "v1.12" }}
      - operator: Exists
        effect: NoSchedule
      - operator: Exists
        effect: NoExecute
      {{- else }}
      - key: node-role.kubernetes.io/controlplane
        operator: Exists
        effect: NoSchedule
      - key: node-role.kubernetes.io/etcd
        operator: Exists
        effect: NoExecute
      {{- end }}
      serviceAccountName: flannel
      containers:
      # RKE1 Sidecar Pattern (Maintains compatibility vs Upstream InitContainer)
      - name: install-cni
        image: {{.CNIImage}}
        command: ["/install-cni.sh"]
        env:
        # The CNI network config to install on each node.
        - name: CNI_NETWORK_CONFIG
          valueFrom:
            configMapKeyRef:
              name: kube-flannel-cfg
              key: cni-conf.json
        - name: CNI_CONF_NAME
          value: "10-flannel.conflist"
        volumeMounts:
        - name: cni
          mountPath: /host/etc/cni/net.d
        - name: host-cni-bin
          mountPath: /host/opt/cni/bin/
      - name: kube-flannel
        image: {{.Image}}
        command:
        - /opt/bin/flanneld
        args:
        - --ip-masq
        - --kube-subnet-mgr
        {{- if .FlannelInterface}}
        - --iface={{.FlannelInterface}}
        {{end}}
        # --- BLACKHOLE ROUTE FIX ---
        {{- if .BlackholeRoute }}
        - --ip-blackhole-route=true
        - --iptables-forward-rules=true
        {{- end }}
        # ---------------------------
        resources:
          requests:
            cpu: "100m"
            memory: "50Mi"
        securityContext:
          seLinuxOptions:
            type: rke_network_t
          privileged: false
          capabilities:
            add: ["NET_ADMIN", "NET_RAW"]
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: EVENT_QUEUE_DEPTH
          value: "5000"
        - name: CONT_WHEN_CACHE_NOT_READY
          value: "false"
        volumeMounts:
        - name: run
          mountPath: /run/flannel
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
        - name: xtables-lock
          mountPath: /run/xtables.lock
      volumes:
      - name: run
        hostPath:
          path: /run/flannel
      - name: host-cni-bin
        hostPath:
          path: /opt/cni/bin
      - name: cni
        hostPath:
          path: /etc/cni/net.d
      - name: flannel-cfg
        configMap:
          name: kube-flannel-cfg
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
  updateStrategy:
{{if .UpdateStrategy}}
{{ toYaml .UpdateStrategy | indent 4}}
{{else}}
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 20%
{{end}}
`
