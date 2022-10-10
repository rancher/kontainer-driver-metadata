package templates

/*
FlannelTemplateV0_19_2 is based on upstream flannel v0.19.2
Source: https://raw.githubusercontent.com/flannel-io/flannel/master/Documentation/kube-flannel.yml
Upstream Changelog:
- Remove PodSecurityPolicy, and use the PodSecurity Admission Controller instead
- Use init containers to install CNI
- Add /run/xtables.lock mount to prevent iptables contention with kube-proxy and the host OS
Rancher Changelog:
- Remove duplicated sections for NodeSelector and priorityClassName
*/

const FlannelTemplateV0_19_2 = `
{{- if eq .RBACConfig "rbac"}}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
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
# Rancher specific change
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
      initContainers:
      - name: install-cni-plugin
        image: {{.CNIImage}}
        command:
        - cp
        args:
        - -f
        - /flannel
        - /opt/cni/bin/flannel
        volumeMounts:
        - name: cni-plugin
          mountPath: /opt/cni/bin
      - name: install-cni
        image: {{.Image}}
        command:
        - cp
        args:
        - -f
        - /etc/kube-flannel/cni-conf.json
        - /etc/cni/net.d/10-flannel.conflist
        volumeMounts:
        - name: cni
          mountPath: /etc/cni/net.d
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      containers:
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
        resources:
          requests:
            cpu: "100m"
            memory: "50Mi"
          limits:
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
      - name: cni-plugin
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
