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
{{- if eq .UseAciCniPriorityClass "true"}}
apiVersion: scheduling.k8s.io/v1beta1
kind: PriorityClass
metadata:
  name: acicni-priority
value: 1000000000
globalDefault: false
description: "This priority class is used for ACI-CNI resources"
---
{{- end }}
{{- if ne .UseAciAnywhereCRD "false"}}
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: epgs.aci.aw
spec:
  group: aci.aw
  names:
    kind: Epg
    listKind: EpgList
    plural: epgs
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: contracts.aci.aw
spec:
  group: aci.aw
  names:
    kind: Contract
    listKind: ContractList
    plural: contracts
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: podifs.aci.aw
spec:
  group: aci.aw
  names:
    kind: PodIF
    listKind: PodIFList
    plural: podifs
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: gbpsstates.aci.aw
spec:
  group: aci.aw
  names:
    kind: GBPSState
    listKind: GBPSStateList
    plural: gbpsstates
  scope: Namespaced
  version: v1
  subresources:
    status: {}
---
{{- end }}
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
  name: snatlocalinfos.aci.snat
spec:
  group: aci.snat
  names:
    kind: SnatLocalInfo
    listKind: SnatLocalInfoList
    plural: snatlocalinfos
    singular: snatlocalinfo
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
  name: qospolicies.aci.qos
spec:
  group: aci.qos
  version: v1
  names:
    kind: QosPolicy
    listKind: QosPolicyList
    plural: qospolicies
    singular: qospolicy
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          type: string
        kind:
          type: string
        spec:
          properties:
            podSelector:
              description: 'Selection of Pods'
              properties:
                matchLabels:
                  type: object
                  description:
            ingress:
              properties:
                policing_rate:
                  type: integer
                  minimum: 0
                policing_burst:
                  type: integer
                  minimum: 0
            egress:
              properties:
                policing_rate:
                  type: integer
                  minimum: 0
                policing_burst:
                  type: integer
                  minimum: 0
            dscpmark:
              properties:
                dscp_marking:
                  type: integer
                  minimum: 0
                  maximum: 56
---
{{- if ne .InstallIstio "false"}}
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
{{- end}}
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
        "aci-policy-tenant": "{{.Tenant}}",
        "require-netpol-annot": false,
{{- if ne .CApic "false"}}
        "lb-type": "None",
{{- end}}
        "install-istio": {{.InstallIstio}},
        "istio-profile": "{{.IstioProfile}}",
{{- if ne .CApic "true"}}
        "aci-podbd-dn": "uni/tn-{{.Tenant}}/BD-aci-containers-{{.SystemIdentifier}}-pod-bd",
        "aci-nodebd-dn": "uni/tn-{{.Tenant}}/BD-aci-containers-{{.SystemIdentifier}}-node-bd",
{{- end}}
        "aci-service-phys-dom": "{{.SystemIdentifier}}-pdom",
        "aci-service-encap": "vlan-{{.ServiceVlan}}",
        "aci-service-monitor-interval": {{.ServiceMonitorInterval}},
        "aci-pbr-tracking-non-snat": {{.PBRTrackingNonSnat}},
        "aci-vrf-tenant": "{{.VRFTenant}}",
        "aci-l3out": "{{.L3Out}}",
        "aci-ext-networks": {{.L3OutExternalNetworks}},
{{- if ne .CApic "true"}}
        "aci-vrf": "{{.VRFName}}",
{{- else}}
        "aci-vrf": "{{.OverlayVRFName}}",
{{- end}}
        "default-endpoint-group": {
            "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
{{- else}}
            "name": "aci-containers-{{.SystemIdentifier}}"
{{- end}}
        },
        "max-nodes-svc-graph": {{.MaxNodesSvcGraph}},
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "istio-operator": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "istio-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "kube-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-prometheus": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-logging": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            }        },
        "service-ip-pool": [
            {
                "end": "{{.ServiceIPEnd}}",
                "start": "{{.ServiceIPStart}}"
            }
        ],
        "snat-contract-scope": "{{.SnatContractScope}}",
        "static-service-ip-pool": [
            {
                "end": "{{.StaticServiceIPEnd}}",
                "start": "{{.StaticServiceIPStart}}"
            }
        ],
        "pod-ip-pool": [
            {
                "end": "{{.PodIPEnd}}",
                "start": "{{.PodIPStart}}"
            }
        ],
        "pod-subnet-chunk-size": {{.PodSubnetChunkSize}},
        "node-service-ip-pool": [
            {
                "end": "{{.NodeServiceIPEnd}}",
                "start": "{{.NodeServiceIPStart}}"
            }
        ],
        "node-service-subnets": [
            "{{.ServiceGraphSubnet}}"
        ],
        "enable_endpointslice": {{.EnableEndpointSlice}}
    }
  host-agent-config: |-
    {
        "app-profile": "aci-containers-{{.SystemIdentifier}}",
{{- if ne .EpRegistry ""}}
        "ep-registry": "{{.EpRegistry}}",
{{- else}}
        "ep-registry": null,
{{- end}}
{{- if ne .OpflexMode ""}}
        "opflex-mode": "{{.OpflexMode}}",
{{- else}}
        "opflex-mode": null,
{{- end}}
        "log-level": "{{.HostAgentLogLevel}}",
        "aci-snat-namespace": "{{.SnatNamespace}}",
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
{{- if ne .CApic "true"}}
        "aci-vrf": "{{.VRFName}}",
{{- else}}
        "aci-vrf": "{{.OverlayVRFName}}",
{{- end}}
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
            "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
{{- else}}
            "name": "aci-containers-default"
{{- end}}
        },
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "istio-operator": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "istio-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "kube-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-prometheus": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-logging": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            }        },
        "enable-drop-log": {{.DropLogEnable}},
        "enable_endpointslice": {{.EnableEndpointSlice}}
    }
  opflex-agent-config: |-
    {
        "log": {
            "level": "{{.OpflexAgentLogLevel}}"
        },
        "opflex": {
{{- if eq .OpflexClientSSL "false"}}
          "ssl": { "mode": "disabled"}
{{- end}}
        }
    }
{{- if eq .RunGbpContainer "true"}}
  gbp-server-config: |-
   {
        "aci-policy-tenant": "{{.Tenant}}",
        "aci-vrf": "{{.OverlayVRFName}}",
{{- if ne .VmmDomain ""}}
        "aci-vmm-domain": "{{.VmmDomain}}",
{{- else}}
        "aci-vmm-domain": "{{.SystemIdentifier}}",
{{- end}}
{{- if ne .CApic "true"}}
        "pod-subnet": "{{.GbpPodSubnet}}"
{{- else}}
        "pod-subnet": "{{.GbpPodSubnet}}",
        "apic": {
            "apic-hosts": {{.ApicHosts}},
            "apic-username": {{.ApicUserName}},
            "apic-private-key-path": "/usr/local/etc/aci-cert/user.key",
            "kafka": {
                "brokers": {{.KafkaBrokers}},
                "client-key-path": "/certs/kafka-client.key",
                "client-cert-path": "/certs/kafka-client.crt",
                "ca-cert-path": "/certs/ca.crt",
                "topic": {{.SystemIdentifier}}
            },
            "cloud-info": {
                "cluster-name": {{.SystemIdentifier}},
                "subnet": {{.SubnetDomainName}},
                "vrf": {{.VRFDomainName}}
            }
        }
{{- end}}
   }
{{- end}}
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
    "start": "{{.SnatPortRangeStart}}"
    "end": "{{.SnatPortRangeEnd}}"
    "ports-per-node": "{{.SnatPortsPerNode}}"
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
{{- if eq .CApic "true"}}
apiVersion: v1
kind: Secret
metadata:
  name: kafka-client-certificates
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
data:
  ca.crt: {{.KafkaClientCrt}}
  kafka-client.crt: {{.KafkaClientCrt}}
  kafka-client.key: {{.KafkaClientKey}}
---
{{- end}}
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
{{- if ne .InstallIstio "false"}}
  - serviceaccounts
{{- end}}
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
{{- if ne .InstallIstio "false"}}
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
{{- end}}
- apiGroups:
  - "networking.k8s.io"
  resources:
  - networkpolicies
  verbs:
  - list
  - watch
  - get
{{- if ne .UseAciAnywhereCRD "false"}}
- apiGroups:
  - "aci.aw"
  resources:
  - epgs
  - contracts
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "aci.aw"
  resources:
  - podifs
  - gbpsstates
  - gbpsstates/status
  verbs:
  - '*'
{{- end}}
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
  - "aci.qos"
  resources:
  - qospolicies
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
  - patch
- apiGroups:
  - apps.openshift.io
  resources:
  - deploymentconfigs
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
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
{{- if ne .DropLogEnable "false"}}
  - update
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
{{- end}}
{{- if ne .UseAciAnywhereCRD "false"}}
- apiGroups:
  - "aci.aw"
  resources:
  - podifs
  - podifs/status
  verbs:
  - "*"
{{- end}}
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
  - "aci.qos"
  resources:
  - qospolicies
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
  - patch
- apiGroups:
  - "aci.snat"
  resources:
  - nodeinfos
  - snatlocalinfos
  verbs:
  - create
  - update
  - list
  - watch
  - get
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
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
        prometheus.io/scrape: "true"
        prometheus.io/port: "9612"
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
      initContainers:
        - name: cnideploy
          image: {{.AciCniDeployContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
{{- if eq .UsePrivilegedContainer "true"}}
            privileged: true
{{- end}}
            capabilities:
              add:
                - SYS_ADMIN
          volumeMounts:
            - name: cni-bin
              mountPath: /mnt/cni-bin
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-cluster-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
      containers:
        - name: aci-containers-host
          image: {{.AciHostContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
{{- if eq .UsePrivilegedContainer "true"}}
            privileged: true
{{- end}}
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
              value: "{{.Tenant}}"
{{- if eq .RunGbpContainer "true"}}
            - name: NODE_EPG
              value: aci-containers-nodes"
            - name: OPFLEX_MODE
              value: overlay
{{- else}}
            - name: NODE_EPG
              value: "aci-containers-{{.SystemIdentifier}}|aci-containers-nodes"
{{- end}}
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
{{- if eq .UseHostNetnsVolume "true"}}
            - mountPath: /run/netns
              name: host-run-netns
              readOnly: true
              mountPropagation: HostToContainer
{{- end}}
          livenessProbe:
            httpGet:
              path: /status
              port: 8090
        - name: opflex-agent
          env:
            - name: REBOOT_WITH_OVS
              value: "true"
{{- if eq .RunGbpContainer "true"}}
            - name: SSL_MODE
              value: disabled
{{- end}}
          image: {{.AciOpflexContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
{{- if eq .UsePrivilegedContainer "true"}}
            privileged: true
{{- end}}
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
{{- if eq .RunOpflexServerContainer "true"}}
        - name: opflex-server
          image: {{.AciOpflexServerContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
          ports:
            - containerPort: {{.OpflexServerPort}}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - name: opflex-server-config-volume
              mountPath: /usr/local/etc/opflex-server
            - name: hostvar
              mountPath: /usr/local/var
{{- end}}
        - name: mcast-daemon
          image: {{.AciMcastContainer}}
          command: ["/bin/sh"]
          args: ["/usr/local/bin/launch-mcastdaemon.sh"]
          imagePullPolicy: {{.ImagePullPolicy}}
{{- if eq .UsePrivilegedContainer "true"}}
          securityContext:
            privileged: true
{{- end}}
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
{{- if eq .UseOpflexServerVolume "true"}}
        - name: opflex-server-config-volume
{{- end}}
{{- if eq .UseHostNetnsVolume "true"}}
        - name: host-run-netns
          hostPath:
            path: /run/netns
{{- end}}
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
{{- if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{end}}
      tolerations:
        - operator: Exists
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-cluster-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
      containers:
        - name: aci-containers-openvswitch
          image: {{.AciOpenvSwitchContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          resources:
            limits:
              memory: "{{.OVSMemoryLimit}}"
          securityContext:
{{- if eq .UsePrivilegedContainer "true"}}
            privileged: true
{{- end}}
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
{{- if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{end}}
{{- if .Tolerations }}
      tolerations:
{{ toYaml .Tolerations | indent 6}}
{{- else }}
      tolerations:
        - operator: Exists
          effect: NoSchedule
{{- end }}
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-node-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
      containers:
{{- if eq .RunGbpContainer "true"}}
        - name: aci-gbpserver
          image: {{.AciGbpServerContainer}}
          imagePullPolicy: {{ .ImagePullPolicy }}
          volumeMounts:
            - name: controller-config-volume
              mountPath: /usr/local/etc/aci-containers/
{{- if eq .CApic "true"}}
            - name: kafka-certs
              mountPath: /certs
            - name: aci-user-cert-volume
              mountPath: /usr/local/etc/aci-cert/
{{- end}}
          env:
            - name: GBP_SERVER_CONF
              value: /usr/local/etc/aci-containers/gbp-server.conf
{{- end}}
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
{{- if eq .CApic "true"}}
        - name: kafka-certs
          secret:
            secretName: kafka-client-certificates
{{- end}}
        - name: aci-user-cert-volume
          secret:
            secretName: aci-user-cert
        - name: controller-config-volume
          configMap:
            name: aci-containers-config
            items:
              - key: controller-config
                path: controller.conf
{{- if eq .RunGbpContainer "true"}}
              - key: gbp-server-config
                path: gbp-server.conf
{{- end}}
{{- if eq .CApic "true"}}
---
apiVersion: aci.aw/v1
kind: PodIF
metadata:
  name: inet-route
  namespace: kube-system
status:
  epg: aci-containers-inet-out
  ipaddr: 0.0.0.0/0
{{- end}}
`

const AciTemplateV513 = `
apiVersion: v1
kind: Namespace
metadata:
  name: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
---
{{- if ne .UseAciAnywhereCRD "false"}}
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: epgs.aci.aw
spec:
  group: aci.aw
  names:
    kind: Epg
    listKind: EpgList
    plural: epgs
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: contracts.aci.aw
spec:
  group: aci.aw
  names:
    kind: Contract
    listKind: ContractList
    plural: contracts
  scope: Namespaced
  version: v1
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: gbpsstates.aci.aw
spec:
  group: aci.aw
  names:
    kind: GBPSState
    listKind: GBPSStateList
    plural: gbpsstates
  scope: Namespaced
  version: v1
  subresources:
    status: {}
---
{{- end }}
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: podifs.aci.aw
spec:
  group: aci.aw
  names:
    kind: PodIF
    listKind: PodIFList
    plural: podifs
  scope: Namespaced
  version: v1
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
  name: snatlocalinfos.aci.snat
spec:
  group: aci.snat
  names:
    kind: SnatLocalInfo
    listKind: SnatLocalInfoList
    plural: snatlocalinfos
    singular: snatlocalinfo
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
  name: qospolicies.aci.qos
spec:
  group: aci.qos
  version: v1
  names:
    kind: QosPolicy
    listKind: QosPolicyList
    plural: qospolicies
    singular: qospolicy
  scope: Namespaced
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          type: string
        kind:
          type: string
        spec:
          properties:
            podSelector:
              description: 'Selection of Pods'
              properties:
                matchLabels:
                  type: object
                  description:
            ingress:
              properties:
                policing_rate:
                  type: integer
                  minimum: 0
                policing_burst:
                  type: integer
                  minimum: 0
            egress:
              properties:
                policing_rate:
                  type: integer
                  minimum: 0
                policing_burst:
                  type: integer
                  minimum: 0
            dscpmark:
              properties:
                dscp_marking:
                  type: integer
                  minimum: 0
                  maximum: 56
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: netflowpolicies.aci.netflow
spec:
  group: aci.netflow
  names:
    kind: NetflowPolicy
    listKind: NetflowPolicyList
    plural: netflowpolicies
    singular: netflowpolicy
  scope: Cluster
  versions:
  - name: v1alpha
    served: true
    storage: true
    schema:
   # openAPIV3Schema is the schema for validating custom objects.
      openAPIV3Schema:
        type: object
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          spec:
            type: object
            properties:
              flowSamplingPolicy:
                type: object
                properties:
                  destIp:
                    type: string
                  destPort:
                    type: integer
                    minimum: 0
                    maximum: 65535
                  flowType:
                    type: string
                    enum:
                      - netflow
                      - ipfix
                  activeFlowTimeOut:
                    type: integer
                    minimum: 0
                    maximum: 3600
                  idleFlowTimeOut:
                    type: integer
                    minimum: 0
                    maximum: 600
                  samplingRate:
                    type: integer
                    minimum: 0
                    maximum: 1000
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: erspanpolicies.aci.erspan
spec:
  group: aci.erspan
  version: v1alpha
  names:
    kind: ErspanPolicy
    listKind: ErspanPolicyList
    plural: erspanpolicies
    singular: erspanpolicy
  scope: Cluster
  validation:
    openAPIV3Schema:
    # openAPIV3Schema is the schema for validating custom objects.
      type: object
      properties:
        apiVersion:
          type: string
        kind:
          type: string
        spec:
          type: object
          properties:
            podSelector:
              type: object
              description: 'Selection of Pods'
              properties:
                matchLabels:
                  type: object
                  description:
            source:
              type: object
              properties:
                admin_state:
                  type: string
                  enum:
                    - start
                    - stop
                direction:
                  type: string
                  enum:
                    - in
                    - out
                    - both
                tag:
                  type: string
            destination:
              type: object
              properties:
                destIp:
                  type: string
                  minimum: 0
                flowId:
                  type: integer
                  minimum: 1
                  maximum: 1023
  versions:
  - name: v1alpha
    served: true
    storage: true
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
        "log-level": "info",
        "apic-hosts": {{.ApicHosts}},
        "apic-refreshtime": "{{.ApicRefreshTime}}",
        "apic-username": "{{.ApicUserName}}",
        "apic-private-key-path": "/usr/local/etc/aci-cert/user.key",
        "apic-use-inst-tag": true,
        "aci-prefix": "{{.SystemIdentifier}}",
        "aci-vmm-type": "Kubernetes",
        "aci-vmm-domain": "{{.SystemIdentifier}}",
        "aci-vmm-controller": "{{.SystemIdentifier}}",
        "aci-policy-tenant": "{{.Tenant}}",
        "require-netpol-annot": false,
{{- if ne .CApic "false"}}
        "lb-type": "None",
{{- end}}
        "install-istio": false,
        "istio-profile": "demo",
{{- if ne .CApic "true"}}
        "aci-podbd-dn": "uni/tn-{{.Tenant}}/BD-aci-containers-{{.SystemIdentifier}}-pod-bd",
        "aci-nodebd-dn": "uni/tn-{{.Tenant}}/BD-aci-containers-{{.SystemIdentifier}}-node-bd",
{{- end}}
        "aci-service-phys-dom": "{{.SystemIdentifier}}-pdom",
        "aci-service-encap": "vlan-{{.ServiceVlan}}",
        "aci-service-monitor-interval": {{.ServiceMonitorInterval}},
        "aci-pbr-tracking-non-snat": {{.PBRTrackingNonSnat}},
        "aci-vrf-tenant": "{{.VRFTenant}}",
        "aci-l3out": "{{.L3Out}}",
        "aci-ext-networks": {{.L3OutExternalNetworks}},
{{- if ne .CApic "true"}}
        "aci-vrf": "{{.VRFName}}",
{{- else}}
        "aci-vrf": "{{.OverlayVRFName}}",
{{- end}}
        "default-endpoint-group": {
            "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
{{- else}}
            "name": "aci-containers-{{.SystemIdentifier}}"
{{- end}}
        },
        "max-nodes-svc-graph": {{.MaxNodesSvcGraph}},
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "istio-operator": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "istio-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "kube-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-prometheus": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-logging": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            }        },
        "service-ip-pool": [
            {
                "end": "{{.ServiceIPEnd}}",
                "start": "{{.ServiceIPStart}}"
            }
        ],
        "snat-contract-scope": "{{.SnatContractScope}}",
        "static-service-ip-pool": [
            {
                "end": "{{.StaticServiceIPEnd}}",
                "start": "{{.StaticServiceIPStart}}"
            }
        ],
        "pod-ip-pool": [
            {
                "end": "{{.PodIPEnd}}",
                "start": "{{.PodIPStart}}"
            }
        ],
        "pod-subnet-chunk-size": {{.PodSubnetChunkSize}},
        "node-service-ip-pool": [
            {
                "end": "{{.NodeServiceIPEnd}}",
                "start": "{{.NodeServiceIPStart}}"
            }
        ],
        "node-service-subnets": [
            "{{.ServiceGraphSubnet}}"
        ],
        "enable_endpointslice": {{.EnableEndpointSlice}}
    }
  host-agent-config: |-
    {
        "app-profile": "aci-containers-{{.SystemIdentifier}}",
{{- if ne .EpRegistry ""}}
        "ep-registry": "{{.EpRegistry}}",
{{- else}}
        "ep-registry": null,
{{- end}}
{{- if ne .OpflexMode ""}}
        "opflex-mode": "{{.OpflexMode}}",
{{- else}}
        "opflex-mode": null,
{{- end}}
        "log-level": "info",
        "aci-snat-namespace": "{{.SnatNamespace}}",
        "aci-vmm-type": "Kubernetes",
        "aci-vmm-domain": "{{.SystemIdentifier}}",
        "aci-vmm-controller": "{{.SystemIdentifier}}",
        "aci-prefix": "{{.SystemIdentifier}}",
{{- if ne .CApic "true"}}
        "aci-vrf": "{{.VRFName}}",
{{- else}}
        "aci-vrf": "{{.OverlayVRFName}}",
{{- end}}
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
            "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
            "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-default"
{{- else}}
            "name": "aci-containers-default"
{{- end}}
        },
        "namespace-default-endpoint-group": {
            "aci-containers-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "istio-operator": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "istio-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-istio"
{{- else}}
                "name": "aci-containers-istio"
{{- end}}
            },
            "kube-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-system": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-prometheus": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            },
            "cattle-logging": {
                "policy-space": "{{.Tenant}}",
{{- if ne .CApic "true"}}
                "name": "aci-containers-{{.SystemIdentifier}}|aci-containers-system"
{{- else}}
                "name": "aci-containers-system"
{{- end}}
            }        },
        "enable-drop-log": {{.DropLogEnable}},
        "enable_endpointslice": {{.EnableEndpointSlice}}
    }
  opflex-agent-config: |-
    {
        "log": {
            "level": "info"
        },
        "opflex": {
          "notif" : { "enabled" : "false"}
{{- if eq .OpflexClientSSL "false"}}
         ,"ssl": { "mode": "disabled"}
{{- end}}
{{- if eq .RunGbpContainer "true"}}
         ,"statistics" : { "mode" : "off" }
{{- end}}
        }
    }
{{- if eq .RunGbpContainer "true"}}
  gbp-server-config: |-
   {
        "aci-policy-tenant": "{{.Tenant}}",
        "aci-vrf": "{{.OverlayVRFName}}",
        "aci-vmm-domain": "{{.SystemIdentifier}}",
{{- if ne .CApic "true"}}
        "pod-subnet": "{{.GbpPodSubnet}}"
{{- else}}
        "pod-subnet": "{{.GbpPodSubnet}}",
        "apic": {
            "apic-hosts": {{.ApicHosts}},
            "apic-username": {{.ApicUserName}},
            "apic-private-key-path": "/usr/local/etc/aci-cert/user.key",
            "kafka": {
                "brokers": {{.KafkaBrokers}},
                "client-key-path": "/certs/kafka-client.key",
                "client-cert-path": "/certs/kafka-client.crt",
                "ca-cert-path": "/certs/ca.crt",
                "topic": {{.SystemIdentifier}}
            },
            "cloud-info": {
                "cluster-name": {{.SystemIdentifier}},
                "subnet": {{.SubnetDomainName}},
                "vrf": {{.VRFDomainName}}
            }
        }
{{- end}}
   }
{{- end}}
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
    "start": "{{.SnatPortRangeStart}}"
    "end": "{{.SnatPortRangeEnd}}"
    "ports-per-node": "{{.SnatPortsPerNode}}"
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
{{- if eq .CApic "true"}}
apiVersion: v1
kind: Secret
metadata:
  name: kafka-client-certificates
  namespace: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
data:
  ca.crt: {{.KafkaClientCrt}}
  kafka-client.crt: {{.KafkaClientCrt}}
  kafka-client.key: {{.KafkaClientKey}}
---
{{- end}}
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
  - "apiextensions.k8s.io"
  resources:
  - customresourcedefinitions
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
{{- if ne .UseAciAnywhereCRD "false"}}
- apiGroups:
  - "aci.aw"
  resources:
  - epgs
  - contracts
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - "aci.aw"
  resources:
  - gbpsstates
  - gbpsstates/status
  verbs:
  - '*'
{{- end}}
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
  - "aci.qos"
  resources:
  - qospolicies
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
  - patch
- apiGroups:
  - "aci.netflow"
  resources:
  - netflowpolicies
  verbs:
  - list
  - watch
  - get
  - update
- apiGroups:
  - "aci.erspan"
  resources:
  - erspanpolicies
  verbs:
  - list
  - watch
  - get
  - update
- apiGroups:
  - "aci.aw"
  resources:
  - podifs
  verbs:
  - '*'
- apiGroups:
  - apps.openshift.io
  resources:
  - deploymentconfigs
  verbs:
  - list
  - watch
  - get
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
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
{{- if ne .DropLogEnable "false"}}
  - update
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
{{- end}}
- apiGroups:
  - "apiextensions.k8s.io"
  resources:
  - customresourcedefinitions
  verbs:
  - list
  - watch
  - get
{{- if ne .UseAciAnywhereCRD "false"}}
- apiGroups:
  - "aci.aw"
  resources:
  - podifs
  - podifs/status
  verbs:
  - "*"
{{- end}}
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
  - "aci.qos"
  resources:
  - qospolicies
  verbs:
  - list
  - watch
  - get
  - create
  - update
  - delete
  - patch
- apiGroups:
  - "aci.netflow"
  resources:
  - netflowpolicies
  verbs:
  - list
  - watch
  - get
  - update
- apiGroups:
  - "aci.snat"
  resources:
  - nodeinfos
  - snatlocalinfos
  verbs:
  - create
  - update
  - list
  - watch
  - get
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
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
        prometheus.io/scrape: "true"
        prometheus.io/port: "9612"
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
      initContainers:
        - name: cnideploy
          image: {{.AciCniDeployContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
            capabilities:
              add:
                - SYS_ADMIN
          volumeMounts:
            - name: cni-bin
              mountPath: /mnt/cni-bin
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
              value: "{{.Tenant}}"
{{- if eq .RunGbpContainer "true"}}
            - name: NODE_EPG
              value: aci-containers-nodes"
            - name: OPFLEX_MODE
              value: overlay
{{- else}}
            - name: NODE_EPG
              value: "aci-containers-{{.SystemIdentifier}}|aci-containers-nodes"
{{- end}}
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
{{- if eq .RunGbpContainer "true"}}
            - name: SSL_MODE
              value: disabled
{{- end}}
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
{{- if eq .RunOpflexServerContainer "true"}}
        - name: opflex-server
          image: {{.AciOpflexServerContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          securityContext:
            capabilities:
              add:
                - NET_ADMIN
          ports:
            - containerPort: {{.OpflexServerPort}}
            - name: metrics
              containerPort: 9632
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
          volumeMounts:
            - name: opflex-server-config-volume
              mountPath: /usr/local/etc/opflex-server
            - name: hostvar
              mountPath: /usr/local/var
{{- end}}
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
{{- if ne .ImagePullSecret ""}}
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
{{- if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{end}}
{{- if .Tolerations }}
      tolerations:
{{ toYaml .Tolerations | indent 6}}
{{- else }}
      tolerations:
        - operator: Exists
          effect: NoSchedule
{{- end }}
      priorityClassName: system-node-critical
      containers:
{{- if eq .RunGbpContainer "true"}}
        - name: aci-gbpserver
          image: {{.AciGbpServerContainer}}
          imagePullPolicy: {{ .ImagePullPolicy }}
          volumeMounts:
            - name: controller-config-volume
              mountPath: /usr/local/etc/aci-containers/
{{- if eq .CApic "true"}}
            - name: kafka-certs
              mountPath: /certs
            - name: aci-user-cert-volume
              mountPath: /usr/local/etc/aci-cert/
{{- end}}
          env:
            - name: GBP_SERVER_CONF
              value: /usr/local/etc/aci-containers/gbp-server.conf
{{- end}}
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
{{- if eq .CApic "true"}}
        - name: kafka-certs
          secret:
            secretName: kafka-client-certificates
{{- end}}
        - name: aci-user-cert-volume
          secret:
            secretName: aci-user-cert
        - name: controller-config-volume
          configMap:
            name: aci-containers-config
            items:
              - key: controller-config
                path: controller.conf
{{- if eq .RunGbpContainer "true"}}
              - key: gbp-server-config
                path: gbp-server.conf
{{- end}}
{{- if eq .CApic "true"}}
---
apiVersion: aci.aw/v1
kind: PodIF
metadata:
  name: inet-route
  namespace: kube-system
status:
  epg: aci-containers-inet-out
  ipaddr: 0.0.0.0/0
{{- end}}
`
