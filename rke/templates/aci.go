package templates

const AciTemplateV500 = `
apiVersion: v1
kind: Namespace
metadata:
  name: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: snatglobalinfos.aci.snat
spec:
  group: aci.snat
  names:
    kind: SnatGlobalInfo
    listKind: SnatGlobalInfoList
    plural: snatglobalinfos
    singular: snatglobalinfo
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: snatpolicies.aci.snat
spec:
  group: aci.snat
  names:
    kind: SnatPolicy
    listKind: SnatPolicyList
    plural: snatpolicies
    singular: snatpolicy
  scope: Cluster
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          properties:
            selector:
              properties:
                labels:
                  type: object
                  properties:
                    additionalProperties:
                      type: string
                namespace:
                  type: string
              type: object
            snatIp:
              type: array
            destIp:
              type: array
          type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: nodeinfos.aci.snat
spec:
  group: aci.snat
  names:
    kind: NodeInfo
    listKind: NodeInfoList
    plural: nodeinfos
    singular: nodeinfo
  scope: Namespaced
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: rdconfigs.aci.snat
spec:
  group: aci.snat
  names:
    kind: RdConfig
    listKind: RdConfigList
    plural: rdconfigs
    singular: rdconfig
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: aciistiooperators.aci.istio
spec:
  group: aci.istio
  names:
    kind: AciIstioOperator
    listKind: AciIstioOperatorList
    plural: aciistiooperators
    singular: aciistiooperator
  scope: Namespaced
  version: v1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: aci-containers-config
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
data:
  controller-config: |-
    {
        "log-level": "{{.ControllerLogLevel}}",
        "apic-hosts": {{.ApicHosts}},
        "apic-refreshtime": "{{.ApicRefreshTime}}",
        "apic-username": "{{.ApicUserName}}",
        "apic-private-key-path": "/usr/local/etc/aci-cert/user.key",
        "apic-use-inst-tag": true,
        "aci-prefix": "{{.SystemIdentifier}}",
        "aci-vmm-type": "Kubernetes",
{{- if ne .VmmDomain ""}}
        "aci-vmm-domain": "{{.VmmDomain}}",
{{- else}}
        "aci-vmm-domain": "{{.SystemIdentifier}}",
{{- end}}
{{- if ne .VmmController ""}}
        "aci-vmm-controller": "{{.VmmController}}",
{{- else}}
        "aci-vmm-controller": "{{.SystemIdentifier}}",
{{- end}}
        "aci-policy-tenant": "{{.VRFTenant}}",
        "require-netpol-annot": false,
        "install-istio": {{.InstallIstio}},
        "istio-profile": "{{.IstioProfile}}",
        "aci-podbd-dn": "uni/tn-{{.SystemIdentifier}}/BD-aci-containers-{{.SystemIdentifier}}-pod-bd",
        "aci-nodebd-dn": "uni/tn-{{.SystemIdentifier}}/BD-aci-containers-{{.SystemIdentifier}}-node-bd",
        "aci-service-phys-dom": "{{.SystemIdentifier}}-pdom",
        "aci-service-encap": "vlan-{{.ServiceVlan}}",
        "aci-service-monitor-interval": {{.ServiceMonitorInterval}},
        "aci-pbr-tracking-non-snat": {{.PBRTrackingNonSnat}},
        "aci-vrf-tenant": "{{.VRFTenant}}",
        "aci-l3out": "{{.L3Out}}",
        "aci-ext-networks": {{.L3OutExternalNetworks}},
        "aci-vrf": "{{.VRFName}}",
        "default-endpoint-group": {
            "policy-space": "{{.SystemIdentifier}}",
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
        },
        "max-nodes-svc-graph": 32,
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
            },
            "istio-operator": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
            },
            "istio-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
            },
            "kube-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
            }        },
        "service-ip-pool": [
            {
                "end": "{{.ServiceIpEnd}}",
                "start": "{{.ServiceIpStart}}"
            }
        ],
        "snat-contract-scope": "global",
        "static-service-ip-pool": [
            {
                "end": "{{.StaticServiceIpEnd}}",
                "start": "{{.StaticServiceIpStart}}"
            }
        ],
        "pod-ip-pool": [
            {
                "end": "{{.PodIpEnd}}",
                "start": "{{.PodIpStart}}"
            }
        ],
        "pod-subnet-chunk-size": 32,
        "node-service-ip-pool": [
            {
                "end": "{{.NodeServiceIpEnd}}",
                "start": "{{.NodeServiceIpStart}}"
            }
        ],
        "node-service-subnets": [
            "{{.ServiceGraphSubnet}}"
        ]
    }
  host-agent-config: |-
    {
        "ep-registry": null,
        "opflex-mode": null,
        "log-level": "{{.HostAgentLogLevel}}",
        "aci-snat-namespace": "aci-containers-system",
        "aci-vmm-type": "Kubernetes",
{{- if ne .VmmDomain ""}}
        "aci-vmm-domain": "{{.VmmDomain}}",
{{- else}}
        "aci-vmm-domain": "{{.SystemIdentifier}}",
{{- end}}
{{- if ne .VmmController ""}}
        "aci-vmm-controller": "{{.VmmController}}",
{{- else}}
        "aci-vmm-controller": "{{.SystemIdentifier}}",
{{- end}}
        "aci-prefix": "{{.SystemIdentifier}}",
        "aci-vrf": "{{.VRFName}}",
        "aci-vrf-tenant": "{{.VRFTenant}}",
        "service-vlan": {{.ServiceVlan}},
        "kubeapi-vlan": {{.KubeAPIVlan}},
        "pod-subnet": "{{.ClusterCIDR}}",
        "node-subnet": "{{.NodeSubnet}}",
        "encap-type": "{{.EncapType}}",
        "aci-infra-vlan": {{.InfraVlan}},
{{- if .MTU}}
{{- if ne .MTU 0}}
        "interface-mtu": {{.MTU}},
{{- end}}
{{- end}}
        "cni-netconfig": [
            {
                "gateway": "{{.PodGateway}}",
                "routes": [
                    {
                        "dst": "0.0.0.0/0",
                        "gw": "{{.PodGateway}}"
                    }
                ],
                "subnet": "{{.ClusterCIDR}}"
            }
        ],
        "default-endpoint-group": {
            "policy-space": "{{.SystemIdentifier}}",
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
        },
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
            },
            "istio-operator": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
            },
            "istio-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
            },
            "kube-system": {
                "policy-space": "{{.SystemIdentifier}}",
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
            }        },
        "enable-drop-log": {{.DropLogEnable}}
    }
  opflex-agent-config: |-
    {
        "log": {
            "level": "{{.OpflexAgentLogLevel}}"
        },
        "opflex": {
        }
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: snat-operator-config
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
data:
    "start": "5000"
    "end": "65000"
    "ports-per-node": "3000"
---
apiVersion: v1
kind: Secret
metadata:
  name: aci-user-cert
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
data:
  user.key: {{.ApicUserKey}}
  user.crt: {{.ApicUserCrt}}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aci-containers-controller
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: aci-containers-host-agent
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
  name: aci-containers:controller
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  - namespaces
  - pods
  - endpoints
  - services
  - events
  - replicationcontrollers
  - serviceaccounts
  verbs:
  - list
  - watch
  - get
  - patch
  - create
  - delete
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
- apiGroups:
  - "rbac.authorization.k8s.io"
  resources:
  - clusterroles
  - clusterrolebindings
  verbs:
  - '*'
- apiGroups:
  - "apiextensions.k8s.io"
  resources:
  - customresourcedefinitions
  verbs:
  - '*'
- apiGroups:
  - "install.istio.io"
  resources:
  - istiocontrolplanes
  - istiooperators
  verbs:
  - '*'
- apiGroups:
  - "aci.istio"
  resources:
  - aciistiooperators
  - aciistiooperator
  verbs:
  - '*'
- apiGroups:
  - "networking.k8s.io"
  resources:
  - networkpolicies
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "apps"
  resources:
  - deployments
  - replicasets
  - daemonsets
  - statefulsets
  verbs:
  - '*'
- apiGroups:
  - ""
  resources:
  - nodes
  - services/status
  verbs:
  - update
- apiGroups:
  - "monitoring.coreos.com"
  resources:
  - servicemonitors
  verbs:
  - get
  - create
- apiGroups:
  - "aci.snat"
  resources:
  - snatpolicies/finalizers
  - snatpolicies/status
  - nodeinfos
  verbs:
  - update
  - create
  - list
  - watch
  - get
  - delete
- apiGroups:
  - "aci.snat"
  resources:
  - snatglobalinfos
  - snatpolicies
  - nodeinfos
  - rdconfigs
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
- apiGroups:
  - apps.openshift.io
  resources:
  - deploymentconfigs
  verbs:
  - list
  - watch
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
  name: aci-containers:host-agent
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  - namespaces
  - pods
  - endpoints
  - services
  - replicationcontrollers
  verbs:
  - list
  - watch
  - get
  - update
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - "networking.k8s.io"
  resources:
  - networkpolicies
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "apps"
  resources:
  - deployments
  - replicasets
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "aci.snat"
  resources:
  - snatpolicies
  - snatglobalinfos
  - rdconfigs
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "aci.snat"
  resources:
  - nodeinfos
  verbs:
  - create
  - update
  - list
  - watch
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: aci-containers:controller
  labels:
    aci-containers-config-version: "{{.Token}}"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: aci-containers:controller
subjects:
- kind: ServiceAccount
  name: aci-containers-controller
  namespace: aci-containers-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: aci-containers:host-agent
  labels:
    aci-containers-config-version: "{{.Token}}"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: aci-containers:host-agent
subjects:
- kind: ServiceAccount
  name: aci-containers-host-agent
  namespace: aci-containers-system
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: aci-containers-host
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
spec:
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      name: aci-containers-host
      network-plugin: aci-containers
  template:
    metadata:
      labels:
        name: aci-containers-host
        network-plugin: aci-containers
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      hostNetwork: true
      hostPID: true
      hostIPC: true
      serviceAccountName: aci-containers-host-agent
{{- if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{- end}}
      tolerations:
        - operator: Exists
      priorityClassName: system-cluster-critical
      containers:
        - name: aci-containers-host
          image: {{.AciHostContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
            capabilities:
              add:
                - SYS_ADMIN
                - NET_ADMIN
                - SYS_PTRACE
          env:
            - name: KUBERNETES_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: TENANT
{{- if ne .Tenant ""}}
              value: "{{.Tenant}}"
{{- else}}
              value: "{{.SystemIdentifier}}"
{{- end}}
            - name: NODE_EPG
              value: "aci-containers-{{.SystemIdentifier}}|aci-containers-nodes"
          volumeMounts:
            - name: cni-bin
              mountPath: /mnt/cni-bin
            - name: cni-conf
              mountPath: /mnt/cni-conf
            - name: hostvar
              mountPath: /usr/local/var
            - name: hostrun
              mountPath: /run
            - name: hostrun
              mountPath: /usr/local/run
            - name: opflex-hostconfig-volume
              mountPath: /usr/local/etc/opflex-agent-ovs/base-conf.d
            - name: host-config-volume
              mountPath: /usr/local/etc/aci-containers/
          livenessProbe:
            httpGet:
              path: /status
              port: 8090
        - name: opflex-agent
          env:
            - name: REBOOT_WITH_OVS
              value: "true"
          image: {{.AciOpflexContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
          volumeMounts:
            - name: hostvar
              mountPath: /usr/local/var
            - name: hostrun
              mountPath: /run
            - name: hostrun
              mountPath: /usr/local/run
            - name: opflex-hostconfig-volume
              mountPath: /usr/local/etc/opflex-agent-ovs/base-conf.d
            - name: opflex-config-volume
              mountPath: /usr/local/etc/opflex-agent-ovs/conf.d
        - name: mcast-daemon
          image: {{.AciMcastContainer}}
          command: ["/bin/sh"]
          args: ["/usr/local/bin/launch-mcastdaemon.sh"]
          imagePullPolicy: {{.ImagePullPolicy}}
          volumeMounts:
            - name: hostvar
              mountPath: /usr/local/var
            - name: hostrun
              mountPath: /run
            - name: hostrun
              mountPath: /usr/local/run
      restartPolicy: Always
      volumes:
        - name: cni-bin
          hostPath:
            path: /opt
        - name: cni-conf
          hostPath:
            path: /etc
        - name: hostvar
          hostPath:
            path: /var
        - name: hostrun
          hostPath:
            path: /run
        - name: host-config-volume
          configMap:
            name: aci-containers-config
            items:
              - key: host-agent-config
                path: host-agent.conf
        - name: opflex-hostconfig-volume
          emptyDir:
            medium: Memory
        - name: opflex-config-volume
          configMap:
            name: aci-containers-config
            items:
              - key: opflex-agent-config
                path: local.conf
---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: aci-containers-openvswitch
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
spec:
  updateStrategy:
    type: RollingUpdate
  selector:
    matchLabels:
      name: aci-containers-openvswitch
      network-plugin: aci-containers
  template:
    metadata:
      labels:
        name: aci-containers-openvswitch
        network-plugin: aci-containers
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      hostNetwork: true
      hostPID: true
      hostIPC: true
      serviceAccountName: aci-containers-host-agent
{{ if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{end}}
      tolerations:
        - operator: Exists
      priorityClassName: system-cluster-critical
      containers:
        - name: aci-containers-openvswitch
          image: {{.AciOpenvSwitchContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          resources:
            limits:
              memory: "{{.OVSMemoryLimit}}"
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
                - SYS_MODULE
                - SYS_NICE
                - IPC_LOCK
          env:
            - name: OVS_RUNDIR
              value: /usr/local/var/run/openvswitch
          volumeMounts:
            - name: hostvar
              mountPath: /usr/local/var
            - name: hostrun
              mountPath: /run
            - name: hostrun
              mountPath: /usr/local/run
            - name: hostetc
              mountPath: /usr/local/etc
            - name: hostmodules
              mountPath: /lib/modules
          livenessProbe:
            exec:
              command:
                - /usr/local/bin/liveness-ovs.sh
      restartPolicy: Always
      volumes:
        - name: hostetc
          hostPath:
            path: /etc
        - name: hostvar
          hostPath:
            path: /var
        - name: hostrun
          hostPath:
            path: /run
        - name: hostmodules
          hostPath:
            path: /lib/modules
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aci-containers-controller
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
    name: aci-containers-controller
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      name: aci-containers-controller
      network-plugin: aci-containers
  template:
    metadata:
      name: aci-containers-controller
      namespace: aci-containers-system
      labels:
        name: aci-containers-controller
        network-plugin: aci-containers
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      hostNetwork: true
      serviceAccountName: aci-containers-controller
{{ if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{end}}
      tolerations:
        - operator: Exists
          effect: NoSchedule
      priorityClassName: system-node-critical
      containers:
        - name: aci-containers-controller
          image: {{.AciControllerContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          env:
            - name: WATCH_NAMESPACE
              value: ""
            - name: ACI_SNAT_NAMESPACE
              value: "aci-containers-system"
            - name: ACI_SNAGLOBALINFO_NAME
              value: "snatglobalinfo"
            - name: ACI_RDCONFIG_NAME
              value: "routingdomain-config"
            - name: SYSTEM_NAMESPACE
              value: "aci-containers-system"
          volumeMounts:
            - name: controller-config-volume
              mountPath: /usr/local/etc/aci-containers/
            - name: aci-user-cert-volume
              mountPath: /usr/local/etc/aci-cert/
          livenessProbe:
            httpGet:
              path: /status
              port: 8091
      volumes:
        - name: aci-user-cert-volume
          secret:
            secretName: aci-user-cert
        - name: controller-config-volume
          configMap:
            name: aci-containers-config
            items:
              - key: controller-config
                path: controller.conf
`
