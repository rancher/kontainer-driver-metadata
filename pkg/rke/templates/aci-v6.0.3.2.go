package templates

const AciTemplateV6032 = `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: acicontainersoperators.aci.ctrl
spec:
  group: aci.ctrl
  names:
    kind: AciContainersOperator
    listKind: AciContainersOperatorList
    plural: acicontainersoperators
    singular: acicontainersoperator
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    subresources:
      status: {}
    schema:
      openAPIV3Schema:
        description: acicontainersoperator owns the lifecycle of ACI objects in the cluster
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            description: AciContainersOperatorSpec defines the desired spec for ACI Objects
            properties:
              flavor:
                type: string
              config:
                type: string
            type: object
          status:
            description: AciContainersOperatorStatus defines the successful completion of AciContainersOperator
            properties:
              status:
                type: boolean
            type: object
        required:
        - spec
        type: object
---
apiVersion: v1
kind: Namespace
metadata:
  name: aci-containers-system
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: nodepodifs.aci.aw
spec:
  group: aci.aw
  names:
    kind: NodePodIF
    listKind: NodePodIFList
    plural: nodepodifs
    singular: nodepodif
  scope: Namespaced
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            type: object
            properties:
              podifs:
                type: array
                items:
                  type: object
                  properties:
                    containerID:
                      type: string
                    epg:
                      type: string
                    ifname:
                      type: string
                    ipaddr:
                      type: string
                    macaddr:
                      type: string
                    podname:
                      type: string
                    podns:
                      type: string
                    vtep:
                      type: string
        required:
        - spec
        type: object
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
apiVersion: apiextensions.k8s.io/v1
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
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        description: SnatGlobalInfo is the Schema for the snatglobalinfos API
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              globalInfos:
                additionalProperties:
                  items:
                    properties:
                      macAddress:
                        type: string
                      portRanges:
                        items:
                          properties:
                            end:
                              maximum: 65535
                              minimum: 1
                              type: integer
                            start:
                              maximum: 65535
                              minimum: 1
                              type: integer
                          type: object
                        type: array
                      snatIp:
                        type: string
                      snatIpUid:
                        type: string
                      snatPolicyName:
                        type: string
                    required:
                    - macAddress
                    - portRanges
                    - snatIp
                    - snatIpUid
                    - snatPolicyName
                    type: object
                  type: array
                type: object
            required:
            - globalInfos
            type: object
          status:
            description: SnatGlobalInfoStatus defines the observed state of SnatGlobalInfo
            type: object
        type: object
---
apiVersion: apiextensions.k8s.io/v1
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
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            description: SnatLocalInfoSpec defines the desired state of SnatLocalInfo
            properties:
              localInfos:
                items:
                  properties:
                    podName:
                      type: string
                    podNamespace:
                      type: string
                    podUid:
                      type: string
                    snatPolicies:
                      items:
                        properties:
                          destIp:
                            items:
                              type: string
                            type: array
                          name:
                            type: string
                          snatIp:
                            type: string
                        required:
                        - destIp
                        - name
                        - snatIp
                        type: object
                      type: array
                  required:
                  - podName
                  - podNamespace
                  - podUid
                  - snatPolicies
                  type: object
                type: array
            required:
            - localInfos
            type: object
        type: object
---
apiVersion: apiextensions.k8s.io/v1
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
  versions:
  - name: v1
    served: true
    storage: true
    subresources:
      status: {}
    schema:
      openAPIV3Schema:
        type: object
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            type: object
            properties:
              selector:
                type: object
                properties:
                  labels:
                    type: object
                    description: 'Selection of Pods'
                    properties:
                    additionalProperties:
                      type: string
                  namespace:
                    type: string
                type: object
              snatIp:
                type: array
                items:
                  type: string
              destIp:
                type: array
                items:
                  type: string
            type: object
          status:
            type: object
            properties:
            additionalProperties:
              type: string
---
apiVersion: apiextensions.k8s.io/v1
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
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              macaddress:
                type: string
              snatpolicynames:
                additionalProperties:
                  type: boolean
                type: object
            type: object
          status:
            description: NodeinfoStatus defines the observed state of Nodeinfo
            type: object
        type: object
---
apiVersion: apiextensions.k8s.io/v1
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
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              discoveredsubnets:
                items:
                  type: string
                type: array
              usersubnets:
                items:
                  type: string
                type: array
            type: object
          status:
            description: NodeinfoStatus defines the observed state of Nodeinfo
            type: object
        type: object
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: networkpolicies.aci.netpol
spec:
  group: aci.netpol
  names:
    kind: NetworkPolicy
    listKind: NetworkPolicyList
    plural: networkpolicies
    singular: networkpolicy
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: Network Policy describes traffic flow at IP address or port level
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              appliedTo:
                properties:
                  namespaceSelector:
                    properties:
                      matchExpressions:
                        items:
                          properties:
                            key:
                              type: string
                            operator:
                              type: string
                            values:
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                      matchLabels:
                        additionalProperties:
                          type: string
                        type: object
                    type: object
                  podSelector:
                    description: allow ingress from the same namespace
                    properties:
                      matchExpressions:
                        items:
                          properties:
                            key:
                              type: string
                            operator:
                              type: string
                            values:
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                      matchLabels:
                        additionalProperties:
                          type: string
                        type: object
                    type: object
                type: object
              egress:
                description: Set of egress rules evaluated based on the order in which they are set.
                items:
                  properties:
                    action:
                      description: Action specifies the action to be applied on the rule.
                      type: string
                    enableLogging:
                      description: EnableLogging is used to indicate if agent should generate logs default to false.
                      type: boolean
                    ports:
                      description: Set of port and protocol allowed/denied by the rule. If this field is unset or empty, this rule matches all ports.
                      items:
                        description: NetworkPolicyPort describes the port and protocol to match in a rule.
                        properties:
                          endPort:
                            description: EndPort defines the end of the port range, being the end included within the range. It can only be specified when a numerical port is specified.
                            format: int32
                            type: integer
                          port:
                            anyOf:
                            - type: integer
                            - type: string
                            description: The port on the given protocol. This can be either a numerical or named port on a Pod. If this field is not provided, this matches all port names and numbers.
                            x-kubernetes-int-or-string: true
                          protocol:
                            default: TCP
                            description: The protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this field defaults to TCP.
                            type: string
                        type: object
                      type: array
                    to:
                      description: Rule is matched if traffic is intended for workloads selected by this field. If this field is empty or missing, this rule matches all destinations.
                      items:
                        properties:
                          ipBlock:
                            description: IPBlock describes the IPAddresses/IPBlocks that is matched in to/from. IPBlock cannot be set as part of the AppliedTo field. Cannot be set with any other selector.
                            properties:
                              cidr:
                                description: CIDR is a string representing the IP Block Valid examples are "192.168.1.1/24" or "2001:db9::/64"
                                type: string
                              except:
                                description: Except is a slice of CIDRs that should not be included within an IP Block Valid examples are "192.168.1.1/24" or "2001:db9::/64" Except values will be rejected if they are outside the CIDR range
                                items:
                                  type: string
                                type: array
                            required:
                            - cidr
                            type: object
                          namespaceSelector:
                            description: Select all Pods from Namespaces matched by this selector, as workloads in To/From fields. If set with PodSelector, Pods are matched from Namespaces matched by the NamespaceSelector. Cannot be set with any other selector except PodSelector or ExternalEntitySelector.
                            properties:
                              matchExpressions:
                                items:
                                  properties:
                                    key:
                                      type: string
                                    operator:
                                      description: operator represents a keys relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                                      type: string
                                    values:
                                      description: values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                                      items:
                                        type: string
                                      type: array
                                  required:
                                  - key
                                  - operator
                                  type: object
                                type: array
                              matchLabels:
                                additionalProperties:
                                  type: string
                                type: object
                            type: object
                          podSelector:
                            description: Select Pods from NetworkPolicys Namespace as workloads in AppliedTo/To/From fields. If set with NamespaceSelector, Pods are matched from Namespaces matched by the NamespaceSelector. Cannot be set with any other selector except NamespaceSelector.
                            properties:
                              matchExpressions:
                                items:
                                  properties:
                                    key:
                                      type: string
                                    operator:
                                      description: operator represents a keys relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                                      type: string
                                    values:
                                      description: values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                                      items:
                                        type: string
                                      type: array
                                  required:
                                  - key
                                  - operator
                                  type: object
                                type: array
                              matchLabels:
                                additionalProperties:
                                  type: string
                                type: object
                            type: object
                        type: object
                      type: array
                    toFqDn:
                      properties:
                        matchNames:
                          items:
                            type: string
                          type: array
                      required:
                      - matchNames
                      type: object
                  required:
                  - enableLogging
                  - toFqDn
                  type: object
                type: array
              ingress:
                description: Set of ingress rules evaluated based on the order in which they are set.
                items:
                  properties:
                    action:
                      description: Action specifies the action to be applied on the rule.
                      type: string
                    enableLogging:
                      description: EnableLogging is used to indicate if agent should generate logs when rules are matched. Should be default to false.
                      type: boolean
                    from:
                      description: Rule is matched if traffic originates from workloads selected by this field. If this field is empty, this rule matches all sources.
                      items:
                        properties:
                          ipBlock:
                            description: IPBlock describes the IPAddresses/IPBlocks that is matched in to/from. IPBlock cannot be set as part of the AppliedTo field. Cannot be set with any other selector.
                            properties:
                              cidr:
                                description: CIDR is a string representing the IP Block Valid examples are "192.168.1.1/24" or "2001:db9::/64"
                                type: string
                              except:
                                description: Except is a slice of CIDRs that should not be included within an IP Block Valid examples are "192.168.1.1/24" or "2001:db9::/64" Except values will be rejected if they are outside the CIDR range
                                items:
                                  type: string
                                type: array
                            required:
                            - cidr
                            type: object
                          namespaceSelector:
                            properties:
                              matchExpressions:
                                items:
                                  properties:
                                    key:
                                      type: string
                                    operator:
                                      description: operator represents a keys relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                                      type: string
                                    values:
                                      description: values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                                      items:
                                        type: string
                                      type: array
                                  required:
                                  - key
                                  - operator
                                  type: object
                                type: array
                              matchLabels:
                                additionalProperties:
                                  type: string
                                type: object
                            type: object
                          podSelector:
                            description: Select Pods from NetworkPolicys Namespace as workloads in AppliedTo/To/From fields. If set with NamespaceSelector, Pods are matched from Namespaces matched by the NamespaceSelector. Cannot be set with any other selector except NamespaceSelector.
                            properties:
                              matchExpressions:
                                description: matchExpressions is a list of label selector requirements. The requirements are ANDed.
                                items:
                                  properties:
                                    key:
                                      type: string
                                    operator:
                                      description: operator represents a keys relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                                      type: string
                                    values:
                                      description: values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                                      items:
                                        type: string
                                      type: array
                                  required:
                                  - key
                                  - operator
                                  type: object
                                type: array
                              matchLabels:
                                additionalProperties:
                                  type: string
                                type: object
                            type: object
                        type: object
                      type: array
                    ports:
                      description: Set of port and protocol allowed/denied by the rule. If this field is unset or empty, this rule matches all ports.
                      items:
                        description: NetworkPolicyPort describes the port and protocol to match in a rule.
                        properties:
                          endPort:
                            description: EndPort defines the end of the port range, being the end included within the range. It can only be specified when a numerical port is specified.
                            format: int32
                            type: integer
                          port:
                            anyOf:
                            - type: integer
                            - type: string
                            description: The port on the given protocol. This can be either a numerical or named port on a Pod. If this field is not provided, this matches all port names and numbers.
                            x-kubernetes-int-or-string: true
                          protocol:
                            default: TCP
                            description: The protocol (TCP, UDP, or SCTP) which traffic must match. If not specified, this field defaults to TCP.
                            type: string
                        type: object
                      type: array
                  type: object
                type: array
              policyTypes:
                items:
                  description: Policy Type string describes the NetworkPolicy type This type is beta-level in 1.8
                  type: string
                type: array
              priority:
                description: Priority specfies the order of the NetworkPolicy relative to other NetworkPolicies.
                type: integer
              type:
                description: type of the policy.
                type: string
            required:
            - type
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: dnsnetworkpolicies.aci.dnsnetpol
spec:
  group: aci.dnsnetpol
  names:
    kind: DnsNetworkPolicy
    listKind: DnsNetworkPolicyList
    plural: dnsnetworkpolicies
    singular: dnsnetworkpolicy
  scope: Namespaced
  versions:
  - name: v1beta
    schema:
      openAPIV3Schema:
        description: dns network Policy
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            properties:
              appliedTo:
                properties:
                  namespaceSelector:
                    properties:
                      matchExpressions:
                        items:
                          properties:
                            key:
                              type: string
                            operator:
                              type: string
                            values:
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                      matchLabels:
                        additionalProperties:
                          type: string
                        type: object
                    type: object
                  podSelector:
                    description: allow ingress from the same namespace
                    properties:
                      matchExpressions:
                        items:
                          properties:
                            key:
                              type: string
                            operator:
                              description: operator represents a keys relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                              type: string
                            values:
                              description: values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                              items:
                                type: string
                              type: array
                          required:
                          - key
                          - operator
                          type: object
                        type: array
                      matchLabels:
                        additionalProperties:
                          type: string
                        type: object
                    type: object
                type: object
              egress:
                description: Set of egress rules evaluated based on the order in which they are set.
                properties:
                  toFqdn:
                    properties:
                      matchNames:
                        items:
                          type: string
                        type: array
                    required:
                    - matchNames
                    type: object
                required:
                - toFqdn
                type: object
            type: object
        required:
        - spec
        type: object
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: qospolicies.aci.qos
spec:
  group: aci.qos
  names:
    kind: QosPolicy
    listKind: QosPolicyList
    plural: qospolicies
    singular: qospolicy
  scope: Namespaced
  preserveUnknownFields: false
  versions:
  - name: v1
    served: true
    storage: true
    subresources:
      status: {}
    schema:
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
              podSelector:
                description: 'Selection of Pods'
                type: object
                properties:
                  matchLabels:
                    type: object
                    description:
              ingress:
                type: object
                properties:
                  policing_rate:
                    type: integer
                    minimum: 0
                  policing_burst:
                    type: integer
                    minimum: 0
              egress:
                type: object
                properties:
                  policing_rate:
                    type: integer
                    minimum: 0
                  policing_burst:
                    type: integer
                    minimum: 0
              dscpmark:
                type: integer
                default: 0
                minimum: 0
                maximum: 63
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
  preserveUnknownFields: false
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
                    default: 2055
                  flowType:
                    type: string
                    enum:
                      - netflow
                      - ipfix
                    default: netflow
                  activeFlowTimeOut:
                    type: integer
                    minimum: 0
                    maximum: 3600
                    default: 60
                  idleFlowTimeOut:
                    type: integer
                    minimum: 0
                    maximum: 600
                    default: 15
                  samplingRate:
                    type: integer
                    minimum: 0
                    maximum: 1000
                    default: 0
                required:
                - destIp
                type: object
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: erspanpolicies.aci.erspan
spec:
  group: aci.erspan
  names:
    kind: ErspanPolicy
    listKind: ErspanPolicyList
    plural: erspanpolicies
    singular: erspanpolicy
  scope: Cluster
  preserveUnknownFields: false
  versions:
  - name: v1alpha
    served: true
    storage: true
    schema:
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
              selector:
                type: object
                description: 'Selection of Pods'
                properties:
                  labels:
                    type: object
                    properties:
                    additionalProperties:
                      type: string
                  namespace:
                    type: string
              source:
                type: object
                properties:
                  adminState:
                    description: Administrative state.
                    default: start
                    type: string
                    enum:
                      - start
                      - stop
                  direction:
                    description: Direction of the packets to monitor.
                    default: both
                    type: string
                    enum:
                      - in
                      - out
                      - both
              destination:
                type: object
                properties:
                  destIP:
                    description: Destination IP of the ERSPAN packet.
                    type: string
                  flowID:
                    description: Unique flow ID of the ERSPAN packet.
                    default: 1
                    type: integer
                    minimum: 1
                    maximum: 1023
                required:
                - destIP
                type: object
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: enabledroplogs.aci.droplog
spec:
  group: aci.droplog
  names:
    kind: EnableDropLog
    listKind: EnableDropLogList
    plural: enabledroplogs
    singular: enabledroplog
  scope: Cluster
  versions:
  - name: v1alpha1
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
            description: Defines the desired state of EnableDropLog
            type: object
            properties:
              disableDefaultDropLog:
                description: Disables the default droplog enabled by acc-provision.
                default: false
                type: boolean
              nodeSelector:
                type: object
                description: Drop logging is enabled on nodes selected based on labels
                properties:
                  labels:
                    type: object
                    properties:
                    additionalProperties:
                      type: string
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: prunedroplogs.aci.droplog
spec:
  group: aci.droplog
  names:
    kind: PruneDropLog
    listKind: PruneDropLogList
    plural: prunedroplogs
    singular: prunedroplog
  scope: Cluster
  versions:
  - name: v1alpha1
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
            description: Defines the desired state of PruneDropLog
            type: object
            properties:
              nodeSelector:
                type: object
                description: Drop logging filters are applied to nodes selected based on labels
                properties:
                  labels:
                    type: object
                    properties:
                    additionalProperties:
                      type: string
              dropLogFilters:
                type: object
                properties:
                  srcIP:
                    type: string
                  destIP:
                    type: string
                  srcMAC:
                    type: string
                  destMAC:
                    type: string
                  srcPort:
                    type: integer
                  destPort:
                    type: integer
                  ipProto:
                    type: integer
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: accprovisioninputs.aci.ctrl
spec:
  group: aci.ctrl
  names:
    kind: AccProvisionInput
    listKind: AccProvisionInputList
    plural: accprovisioninputs
    singular: accprovisioninput
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    subresources:
      status: {}
    schema:
      openAPIV3Schema:
        description: accprovisioninput defines the input configuration for ACI CNI
        properties:
          apiVersion:
            type: string
          kind:
            type: string
          metadata:
            type: object
          spec:
            description: AccProvisionInputSpec defines the desired spec for accprovisioninput object
            properties:
              acc_provision_input:
                type: object
                properties:
                  operator_managed_config:
                    type: object
                    properties:
                      enable_updates:
                        type: boolean
                  aci_config:
                    type: object
                    properties:
                      sync_login:
                        type: object
                        properties:
                          certfile:
                            type: string
                          keyfile:
                            type: string
                      client_ssl:
                        type: boolean
                  net_config:
                    type: object
                    properties:
                      interface_mtu:
                        type: integer
                      service_monitor_interval:
                        type: integer
                      pbr_tracking_non_snat:
                        type: boolean
                      pod_subnet_chunk_size:
                        type: integer
                      disable_wait_for_network:
                        type: boolean
                      duration_wait_for_network:
                        type: integer
                  registry:
                    type: object
                    properties:
                      image_prefix:
                        type: string
                      image_pull_secret:
                        type: string
                      aci_containers_operator_version:
                        type: string
                      aci_containers_controller_version:
                        type: string
                      aci_containers_host_version:
                        type: string
                      acc_provision_operator_version:
                        type: string
                      aci_cni_operator_version:
                        type: string
                      cnideploy_version:
                        type: string
                      opflex_agent_version:
                        type: string
                      openvswitch_version:
                        type: string
                      gbp_version:
                        type: string
                  logging:
                    type: object
                    properties:
                      controller_log_level:
                        type: string
                      hostagent_log_level:
                        type: string
                      opflexagent_log_level:
                        type: string
                  istio_config:
                    type: object
                    properties:
                      install_istio:
                        type: boolean
                      install_profile:
                        type: string
                  multus:
                    type: object
                    properties:
                      disable:
                        type: boolean
                  drop_log_config:
                    type: object
                    properties:
                      enable:
                        type: boolean
                  nodepodif_config:
                    type: object
                    properties:
                      enable:
                        type: boolean
                  sriov_config:
                    type: object
                    properties:
                      enable:
                        type: boolean
                  kube_config:
                    type: object
                    properties:
                      ovs_memory_limit:
                        type: string
                      use_privileged_containers:
                        type: boolean
                      image_pull_policy:
                        type: string
                      reboot_opflex_with_ovs:
                        type: string
                      snat_operator:
                        type: object
                        properties:
                          port_range:
                            type: object
                            properties:
                              start:
                                type: integer
                              end:
                                type: integer
                              ports_per_node:
                                type: integer
                          contract_scope:
                            type: string
                          disable_periodic_snat_global_info_sync:
                            type: boolean
            type: object
          status:
            description: AccProvisionInputStatus defines the successful completion of AccProvisionInput
            properties:
              status:
                type: boolean
            type: object
        required:
        - spec
        type: object
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
{{- if ne .AciMultipod "false" }}
        "aci-multipod": {{.AciMultipod}},
{{- end}}
{{- if .OpflexDeviceReconnectWaitTimeout }}
        "opflex-device-reconnect-wait-timeout": {{.OpflexDeviceReconnectWaitTimeout}},
{{- end}}
        "apic-refreshtime": "{{.ApicRefreshTime}}",
        "apic-subscription-delay": {{.ApicSubscriptionDelay}},
        "apic_refreshticker_adjust": "{{.ApicRefreshTickerAdjust}}",
        "apic-username": "{{.ApicUserName}}",
        "apic-private-key-path": "/usr/local/etc/aci-cert/user.key",
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
{{- if ne .CApic "false"}}
        "lb-type": "None",
{{- end}}
{{- if ne .HppOptimization "false"}}
        "hpp-optimization": {{.HppOptimization}},
{{- end}}
{{- if ne .NoWaitForServiceEpReadiness "false"}}
        "no-wait-for-service-ep-readiness": {{.NoWaitForServiceEpReadiness}},
{{- end}}
{{- if ne .ServiceGraphEndpointAddDelay "0"}}
        "service-graph-endpoint-add-delay" : {
            "delay": {{.ServiceGraphEndpointAddDelay}},
            "services": [{{- range $index, $item :=.ServiceGraphEndpointAddServices }}{{- if $index}},{{end}}{ {{- range $k, $v := $item }}"{{ $k }}": "{{ $v }}"{{if eq $k "name"}},{{end}}{{- end}}}{{end}}]
        },
{{- end}}
{{- if ne .AddExternalSubnetsToRdconfig "false"}}
        "add-external-subnets-to-rdconfig": {{.AddExternalSubnetsToRdconfig}},
{{- end}}
{{- if ne .DisablePeriodicSnatGlobalInfoSync "false"}}
        "disable-periodic-snat-global-info-sync": {{.DisablePeriodicSnatGlobalInfoSync}},
{{- end}}
{{- if .NodeSnatRedirectExclude }}
        "node-snat-redirect-exclude": [{{ range $index,$item := .NodeSnatRedirectExclude}}{{- if $index}}, {{end }}{"group": "{{ index $item "group" }}", "labels": {{ index $item "labels" }}}{{ end }}],
{{- end }}
        "opflex-device-delete-timeout": {{.OpflexDeviceDeleteTimeout}},
        "sleep-time-snat-global-info-sync": {{.SleepTimeSnatGlobalInfoSync}},
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
        "app-profile": "aci-containers-{{.SystemIdentifier}}",
{{- if ne .AddExternalContractToDefaultEpg "false"}}
        "add-external-contract-to-default-epg": {{.AddExternalContractToDefaultEpg}},
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
            "service-ip-pool": [{{- range $index, $item := .ServiceIPPool }}{{- if $index}},{{end}}{ "start": "{{ $item.Start }}", "end": "{{ $item.End}}" }{{end}}],
        "extern-static": [{{- range $index, $item := .StaticExternalSubnet }}{{- if $index}},{{end}}{{$item}}{{end}}],
        "extern-dynamic": [{{- range $index, $item := .DynamicExternalSubnet }}{{- if $index}},{{end}}{{$item}}{{end}}],
        "snat-contract-scope": "{{.SnatContractScope}}",
        "static-service-ip-pool": [{{- range $index, $item := .StaticServiceIPPool }}{{- if $index}},{{end}}{ "start": "{{ $item.Start }}", "end": "{{ $item.End }}" }{{end}}],
        "pod-ip-pool": [{{- range $index, $item := .PodIPPool }}{{- if $index}},{{end}}{ "start": "{{ $item.Start }}", "end": "{{ $item.End}}" }{{end}}],
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
{{- if ne .AciMultipod "false" }}
        "aci-multipod": {{.AciMultipod}},
{{- end}}
{{- if ne .DhcpRenewMaxRetryCount "0" }}
        "dhcp-renew-max-retry-count": {{.DhcpRenewMaxRetryCount}},
{{- end}}
{{- if ne .DhcpDelay "0" }}
        "dhcp-delay": {{.DhcpDelay}},
{{- end}}
{{- if ne .EnableOpflexAgentReconnect "false"}}
        "enable-opflex-agent-reconnect": {{.EnableOpflexAgentReconnect}},
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
{{- if ne .HppOptimization "false"}}
        "hpp-optimization": {{.HppOptimization}},
{{- end}}
        "pod-subnet":  [{{- range $index, $item := .PodSubnet }}{{- if $index}},{{end}}{{$item}}{{end}}],
        "node-subnet": [{{- range $index, $item := .NodeSubnet }}{{- if $index}},{{end}}{{$item}}{{end}}],
        "encap-type": "{{.EncapType}}",
        "aci-infra-vlan": {{.InfraVlan}},
{{- if .MTU}}
{{- if ne .MTU 0}}
        "interface-mtu": {{.MTU}},
{{- end}}
{{- end}}
{{- if .MTUHeadRoom}}
{{- if ne .MTUHeadRoom "0"}}
        "interface-mtu-headroom": {{.MTUHeadRoom}},
{{- end}}
{{- end}}
        "cni-netconfig": [{{- range $index, $item := .PodNetwork }}{{- if $index}},{{end}}{ "gateway": "{{ $item.Gateway }}", "subnet": "{{ $item.Subnet }}", "routes": [{ "dst": "0.0.0.0/0", "gw": "{{ $item.Gateway }}" }]}{{end}}],
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
        "enable_endpointslice": {{.EnableEndpointSlice}},
        "enable-nodepodif": {{.NodePodIfEnable}},
        "enable-ovs-hw-offload": {{.SriovEnable}}
    }
  opflex-agent-config: |-
    {
        "log": {
            "level": "{{.OpflexAgentLogLevel}}"
        },
        "opflex": {
{{- if eq .OpflexClientSSL "false"}}
          "ssl": { "mode": "disabled"},
{{- end}}
{{- if eq .OpflexAgentStatistics "false"}}
          "statistics" : { "mode" : "off" },
{{- end}}
{{- if ne .OpflexAgentPolicyRetryDelayTimer "10" }}
          "timers" : { "policy-retry-delay": {{.OpflexAgentPolicyRetryDelayTimer}} },
{{- end}}
          "notif" : { "enabled" : "false" },
          "asyncjson": { "enabled" : {{.OpflexAgentOpflexAsyncjsonEnabled}} }
        },
        "ovs": {
          "asyncjson": { "enabled" : {{.OpflexAgentOvsAsyncjsonEnabled}} }
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
{{- if eq .UseClusterRole "true"}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
  name: aci-containers-controller
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
  - update
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
  - "rbac.authorization.k8s.io"
  resources:
  - clusterroles
  - clusterrolebindings
  verbs:
  - '*'
{{- if ne .InstallIstio "false"}}
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
  - nodepodifs
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
- apiGroups:
  - "aci.netpol"
  resources:
  - networkpolicies
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
- apiGroups:
  - "aci.dnsnetpol"
  resources:
  - dnsnetworkpolicies
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - delete
---
{{- end}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    aci-containers-config-version: "{{.Token}}"
    network-plugin: aci-containers
  name: aci-containers-host-agent
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
  - "aci.droplog"
  resources:
  - enabledroplogs
  - prunedroplogs
  verbs:
  - list
  - watch
  - get
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
  - delete
- apiGroups:
  - discovery.k8s.io
  resources:
  - endpointslices
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - "aci.netpol"
  resources:
  - networkpolicies
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - "aci.aw"
  resources:
  - nodepodifs
  verbs:
  - "*"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: aci-containers-controller
  labels:
    aci-containers-config-version: "{{.Token}}"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: aci-containers-controller
subjects:
- kind: ServiceAccount
  name: aci-containers-controller
  namespace: aci-containers-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: aci-containers-host-agent
  labels:
    aci-containers-config-version: "{{.Token}}"
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: aci-containers-host-agent
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
{{- if ne .UseSystemNodePriorityClass "false"}}
      priorityClassName: system-node-critical
{{- else if .UseAciContainersHostPriorityClass}}        
      priorityClassName: aci-containers-host
{{- else}}        
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-cluster-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
{{- end}}
      containers:
        - name: aci-containers-host
          image: {{.AciHostContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
{{- if or ( .AciContainersHostMemoryLimit ) ( .AciContainersHostMemoryRequest )}}
          resources:
            limits:
{{- if .AciContainersHostMemoryLimit }}
              memory: "{{ .AciContainersHostMemoryLimit }}"
{{- else}}
              memory: "{{ .AciContainersMemoryLimit }}"
{{- end}}
            requests:
{{- if .AciContainersHostMemoryRequest }}
              memory: "{{ .AciContainersHostMemoryRequest }}"
{{- else}}
              memory: "{{ .AciContainersMemoryRequest }}"
{{- end}}
{{- end}}
          securityContext:
{{- if eq .UsePrivilegedContainer "true"}}
            privileged: true
{{- end}}
            capabilities:
              add:
                - SYS_ADMIN
                - NET_ADMIN
                - SYS_PTRACE
                - NET_RAW
          env:
            - name: KUBERNETES_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: TENANT
              value: "{{.Tenant}}"
{{- if ne .MultusDisable "true"}}
            - name: MULTUS
              value: true
{{- end}}
{{- if eq .DisableWaitForNetwork "true"}}
            - name: DISABLE_WAIT_FOR_NETWORK
              value: true
{{- else}}
            - name: DURATION_WAIT_FOR_NETWORK
              value: "{{.DurationWaitForNetwork}}"
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
            - name: varlogpods
              mountPath: /var/log/pods
              readOnly: true
            - name: varlogcontainers
              mountPath: /var/log/containers
              readOnly: true
            - name: varlibdocker
              mountPath: /var/lib/docker
              readOnly: true
{{- if eq .AciMultipod "true" }}
            - name: dhclient
              mountPath: /var/lib/dhclient
{{- end}}
{{- if eq .UseHostNetnsVolume "true"}}
            - mountPath: /run/netns
              name: host-run-netns
              readOnly: true
              mountPropagation: HostToContainer
{{- end}}
{{- if ne .MultusDisable "true"}}
            - name: multus-cni-conf
              mountPath: /mnt/multus-cni-conf
{{- end}}
          livenessProbe:
            failureThreshold: 10
            httpGet:
              path: /status
              port: 8090
              scheme: HTTP
            initialDelaySeconds: 120
            periodSeconds: 60
            successThreshold: 1
            timeoutSeconds: 30
        - name: opflex-agent
          env:
            - name: REBOOT_WITH_OVS
              value: "true"
{{- if ne .OpflexOpensslCompat "false"}}
            - name: OPENSSL_CONF
              value: "/etc/pki/tls/openssl11.cnf" 
{{- end}}
          image: {{.AciOpflexContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
{{- if or ( .OpflexAgentMemoryLimit ) ( .OpflexAgentMemoryRequest )}}
          resources:
            limits:
{{- if .OpflexAgentMemoryLimit }}
              memory: "{{ .OpflexAgentMemoryLimit }}"
{{- else}}
              memory: "{{ .AciContainersMemoryLimit }}"
{{- end}}
            requests:
{{- if .OpflexAgentMemoryRequest }}
              memory: "{{ .OpflexAgentMemoryRequest }}"
{{- else}}
              memory: "{{ .AciContainersMemoryRequest }}"
{{- end}}
{{- end}}                    
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
          image: {{.AciOpflexContainer}}
          command: ["/bin/sh"]
          args: ["/usr/local/bin/launch-opflexserver.sh"]
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
{{- if ne .OpflexMode "overlay"}}
        - name: mcast-daemon
          image: {{.AciMcastContainer}}
          command: ["/bin/sh"]
          args: ["/usr/local/bin/launch-mcastdaemon.sh"]
          imagePullPolicy: {{.ImagePullPolicy}}
{{- if or ( .McastDaemonMemoryLimit ) ( .McastDaemonMemoryRequest )}}
          resources:
            limits:
{{- if .McastDaemonMemoryLimit }}
              memory: "{{ .McastDaemonMemoryLimit }}"
{{- else}}
              memory: "{{ .AciContainersMemoryLimit }}"
{{- end}}
            requests:
{{- if .McastDaemonMemoryRequest }}
              memory: "{{ .McastDaemonMemoryRequest }}"
{{- else}}
              memory: "{{ .AciContainersMemoryRequest }}"
{{- end}}
{{- end}}
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
{{- end}}
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
        - name: varlogpods
          hostPath:
            path: /var/log/pods
        - name: varlogcontainers
          hostPath:
            path: /var/log/containers
        - name: varlibdocker
          hostPath:
            path: /var/lib/docker
{{- if eq .AciMultipod "true" }}
{{- if eq .AciMultipodUbuntu "true" }}
        - name: dhclient
          hostPath:
            path: /var/lib/dhcp
{{- else}}
        - name: dhclient
          hostPath:
            path: /var/lib/dhclient
{{- end}}
{{- end}}
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
{{- if ne .MultusDisable  "true" }}
        - name: multus-cni-conf
          hostPath:
            path: /var/run/multus/
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
{{- if ne .UseSystemNodePriorityClass "false"}}
      priorityClassName: system-node-critical
{{- else if .UseAciContainersOpenvswitchPriorityClass}}        
      priorityClassName: aci-containers-openvswitch
{{- else}}        
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-cluster-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
{{- end}}
      containers:
        - name: aci-containers-openvswitch
          image: {{.AciOpenvSwitchContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
          resources:
            limits:
              memory: "{{.OVSMemoryLimit}}"
            requests:
              memory: "{{.OVSMemoryRequest}}"
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
            - name: varlogpods
              mountPath: /var/log/pods
              readOnly: true
            - name: varlogcontainers
              mountPath: /var/log/containers
              readOnly: true
            - name: varlibdocker
              mountPath: /var/lib/docker
              readOnly: true
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
        - name: varlogpods
          hostPath:
            path: /var/log/pods
        - name: varlogcontainers
          hostPath:
            path: /var/log/containers
        - name: varlibdocker
          hostPath:
            path: /var/lib/docker
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
    spec:
      hostNetwork: true
      serviceAccountName: aci-containers-controller
{{- if ne .ImagePullSecret ""}}
      imagePullSecrets:
        - name: {{.ImagePullSecret}}
{{- end}}
{{- if .Tolerations }}
      tolerations:
{{ toYaml .Tolerations | indent 6}}
{{- else }}
      tolerations:
      - effect: NoExecute
        operator: Exists
        tolerationSeconds: {{ .TolerationSeconds }}
      - effect: NoSchedule
        key: node.kubernetes.io/not-ready
        operator: Exists
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
{{- end }}
{{- if ne .UseSystemNodePriorityClass "false"}}
      priorityClassName: system-node-critical
{{- else if .UseAciContainersControllerPriorityClass}}        
      priorityClassName: aci-containers-controller
{{- else}}        
{{- if ne .NoPriorityClass "true"}}
      priorityClassName: system-node-critical
{{- end}}
{{- if eq .UseAciCniPriorityClass "true"}}
      priorityClassName: acicni-priority
{{- end}}
{{- end}}
      containers:
        - name: aci-containers-controller
          image: {{.AciControllerContainer}}
          imagePullPolicy: {{.ImagePullPolicy}}
{{- if or ( .AciContainersControllerMemoryLimit ) ( .AciContainersControllerMemoryRequest )}}
          resources:
            limits:
{{- if .AciContainersControllerMemoryLimit }}
              memory: "{{ .AciContainersControllerMemoryLimit }}"
{{- else}}
              memory: "{{ .AciContainersMemoryLimit }}"
{{- end}}
            requests:
{{- if .AciContainersControllerMemoryRequest }}
              memory: "{{ .AciContainersControllerMemoryRequest }}"
{{- else}}
              memory: "{{ .AciContainersMemoryRequest }}"
{{- end}}
{{- end}}
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
            - name: varlogpods
              mountPath: /var/log/pods
              readOnly: true
            - name: varlogcontainers
              mountPath: /var/log/containers
              readOnly: true
            - name: varlibdocker
              mountPath: /var/lib/docker
              readOnly: true
            - name: aci-user-cert-volume
              mountPath: /usr/local/etc/aci-cert/
          livenessProbe:
            failureThreshold: 10
            httpGet:
              path: /status
              port: 8091
              scheme: HTTP
            initialDelaySeconds: 120
            periodSeconds: 60
            successThreshold: 1
            timeoutSeconds: 30
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
        - name: varlogpods
          hostPath:
            path: /var/log/pods
        - name: varlogcontainers
          hostPath:
            path: /var/log/containers
        - name: varlibdocker
          hostPath:
            path: /var/lib/docker
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
---
apiVersion: v1
kind: LimitRange
metadata:
  name: memory-limit-range
  namespace: aci-containers-system
spec:
  limits:
  - default:
      memory: {{ .AciContainersMemoryLimit }}
    defaultRequest:
      memory: {{ .AciContainersMemoryRequest }}
    type: Container
`
