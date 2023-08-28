package templates

/*
TigeraOperatorTemplateV3_26_1 is based on upstream calico v3.26.1
Source: https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml
Notes: the upstream manifest is too large to be stored in a ConfigMap
so we get a timeout when trying to apply it through rke.
To mitigate the issue, the description fields have been removed from the template with this command:
yq  'del(.. | .description?)' tigera-operator.yaml > tigera-operator-nodesc.yaml
Upstream Changelog:
- Initial use of this manifest
Rancher Changelog:
- Initial use of this manifest
*/

const TigeraOperatorTemplateV3_26_1 = `
{{- $cidrs := splitList "," .ClusterCIDR }}
apiVersion: v1
kind: Namespace
metadata:
  name: tigera-operator
  labels:
    name: tigera-operator
---
# Source: crds/calico/crd.projectcalico.org_bgpconfigurations.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: bgpconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BGPConfiguration
    listKind: BGPConfigurationList
    plural: bgpconfigurations
    singular: bgpconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                asNumber:
                  format: int32
                  type: integer
                bindMode:
                  type: string
                communities:
                  items:
                    properties:
                      name:
                        type: string
                      value:
                        pattern: ^(\d+):(\d+)$|^(\d+):(\d+):(\d+)$
                        type: string
                    type: object
                  type: array
                ignoredInterfaces:
                  items:
                    type: string
                  type: array
                listenPort:
                  maximum: 65535
                  minimum: 1
                  type: integer
                logSeverityScreen:
                  type: string
                nodeMeshMaxRestartTime:
                  type: string
                nodeMeshPassword:
                  properties:
                    secretKeyRef:
                      properties:
                        key:
                          type: string
                        name:
                          type: string
                        optional:
                          type: boolean
                      required:
                        - key
                      type: object
                  type: object
                nodeToNodeMeshEnabled:
                  type: boolean
                prefixAdvertisements:
                  items:
                    properties:
                      cidr:
                        type: string
                      communities:
                        items:
                          type: string
                        type: array
                    type: object
                  type: array
                serviceClusterIPs:
                  items:
                    properties:
                      cidr:
                        type: string
                    type: object
                  type: array
                serviceExternalIPs:
                  items:
                    properties:
                      cidr:
                        type: string
                    type: object
                  type: array
                serviceLoadBalancerIPs:
                  items:
                    properties:
                      cidr:
                        type: string
                    type: object
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_bgpfilters.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: (devel)
  creationTimestamp: null
  name: bgpfilters.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BGPFilter
    listKind: BGPFilterList
    plural: bgpfilters
    singular: bgpfilter
  scope: Cluster
  versions:
    - name: v1
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
                exportV4:
                  items:
                    properties:
                      action:
                        type: string
                      cidr:
                        type: string
                      matchOperator:
                        type: string
                    required:
                      - action
                      - cidr
                      - matchOperator
                    type: object
                  type: array
                exportV6:
                  items:
                    properties:
                      action:
                        type: string
                      cidr:
                        type: string
                      matchOperator:
                        type: string
                    required:
                      - action
                      - cidr
                      - matchOperator
                    type: object
                  type: array
                importV4:
                  items:
                    properties:
                      action:
                        type: string
                      cidr:
                        type: string
                      matchOperator:
                        type: string
                    required:
                      - action
                      - cidr
                      - matchOperator
                    type: object
                  type: array
                importV6:
                  items:
                    properties:
                      action:
                        type: string
                      cidr:
                        type: string
                      matchOperator:
                        type: string
                    required:
                      - action
                      - cidr
                      - matchOperator
                    type: object
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_bgppeers.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: bgppeers.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BGPPeer
    listKind: BGPPeerList
    plural: bgppeers
    singular: bgppeer
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                asNumber:
                  format: int32
                  type: integer
                filters:
                  items:
                    type: string
                  type: array
                keepOriginalNextHop:
                  type: boolean
                maxRestartTime:
                  type: string
                node:
                  type: string
                nodeSelector:
                  type: string
                numAllowedLocalASNumbers:
                  format: int32
                  type: integer
                password:
                  properties:
                    secretKeyRef:
                      properties:
                        key:
                          type: string
                        name:
                          type: string
                        optional:
                          type: boolean
                      required:
                        - key
                      type: object
                  type: object
                peerIP:
                  type: string
                peerSelector:
                  type: string
                reachableBy:
                  type: string
                sourceAddress:
                  type: string
                ttlSecurity:
                  type: integer
              type: object
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
# Source: crds/calico/crd.projectcalico.org_blockaffinities.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: blockaffinities.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: BlockAffinity
    listKind: BlockAffinityList
    plural: blockaffinities
    singular: blockaffinity
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                cidr:
                  type: string
                deleted:
                  type: string
                node:
                  type: string
                state:
                  type: string
              required:
                - cidr
                - deleted
                - node
                - state
              type: object
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
# Source: crds/calico/crd.projectcalico.org_caliconodestatuses.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: (devel)
  creationTimestamp: null
  name: caliconodestatuses.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: CalicoNodeStatus
    listKind: CalicoNodeStatusList
    plural: caliconodestatuses
    singular: caliconodestatus
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                classes:
                  items:
                    type: string
                  type: array
                node:
                  type: string
                updatePeriodSeconds:
                  format: int32
                  type: integer
              type: object
            status:
              properties:
                agent:
                  properties:
                    birdV4:
                      properties:
                        lastBootTime:
                          type: string
                        lastReconfigurationTime:
                          type: string
                        routerID:
                          type: string
                        state:
                          type: string
                        version:
                          type: string
                      type: object
                    birdV6:
                      properties:
                        lastBootTime:
                          type: string
                        lastReconfigurationTime:
                          type: string
                        routerID:
                          type: string
                        state:
                          type: string
                        version:
                          type: string
                      type: object
                  type: object
                bgp:
                  properties:
                    numberEstablishedV4:
                      type: integer
                    numberEstablishedV6:
                      type: integer
                    numberNotEstablishedV4:
                      type: integer
                    numberNotEstablishedV6:
                      type: integer
                    peersV4:
                      items:
                        properties:
                          peerIP:
                            type: string
                          since:
                            type: string
                          state:
                            type: string
                          type:
                            type: string
                        type: object
                      type: array
                    peersV6:
                      items:
                        properties:
                          peerIP:
                            type: string
                          since:
                            type: string
                          state:
                            type: string
                          type:
                            type: string
                        type: object
                      type: array
                  required:
                    - numberEstablishedV4
                    - numberEstablishedV6
                    - numberNotEstablishedV4
                    - numberNotEstablishedV6
                  type: object
                lastUpdated:
                  format: date-time
                  nullable: true
                  type: string
                routes:
                  properties:
                    routesV4:
                      items:
                        properties:
                          destination:
                            type: string
                          gateway:
                            type: string
                          interface:
                            type: string
                          learnedFrom:
                            properties:
                              peerIP:
                                type: string
                              sourceType:
                                type: string
                            type: object
                          type:
                            type: string
                        type: object
                      type: array
                    routesV6:
                      items:
                        properties:
                          destination:
                            type: string
                          gateway:
                            type: string
                          interface:
                            type: string
                          learnedFrom:
                            properties:
                              peerIP:
                                type: string
                              sourceType:
                                type: string
                            type: object
                          type:
                            type: string
                        type: object
                      type: array
                  type: object
              type: object
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
# Source: crds/calico/crd.projectcalico.org_clusterinformations.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: clusterinformations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: ClusterInformation
    listKind: ClusterInformationList
    plural: clusterinformations
    singular: clusterinformation
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                calicoVersion:
                  type: string
                clusterGUID:
                  type: string
                clusterType:
                  type: string
                datastoreReady:
                  type: boolean
                variant:
                  type: string
              type: object
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
# Source: crds/calico/crd.projectcalico.org_felixconfigurations.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: felixconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: FelixConfiguration
    listKind: FelixConfigurationList
    plural: felixconfigurations
    singular: felixconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                allowIPIPPacketsFromWorkloads:
                  type: boolean
                allowVXLANPacketsFromWorkloads:
                  type: boolean
                awsSrcDstCheck:
                  enum:
                    - DoNothing
                    - Enable
                    - Disable
                  type: string
                bpfConnectTimeLoadBalancingEnabled:
                  type: boolean
                bpfDSROptoutCIDRs:
                  items:
                    type: string
                  type: array
                bpfDataIfacePattern:
                  type: string
                bpfDisableUnprivileged:
                  type: boolean
                bpfEnabled:
                  type: boolean
                bpfEnforceRPF:
                  type: string
                bpfExtToServiceConnmark:
                  type: integer
                bpfExternalServiceMode:
                  type: string
                bpfHostConntrackBypass:
                  type: boolean
                bpfKubeProxyEndpointSlicesEnabled:
                  type: boolean
                bpfKubeProxyIptablesCleanupEnabled:
                  type: boolean
                bpfKubeProxyMinSyncPeriod:
                  type: string
                bpfL3IfacePattern:
                  type: string
                bpfLogLevel:
                  type: string
                bpfMapSizeConntrack:
                  type: integer
                bpfMapSizeIPSets:
                  type: integer
                bpfMapSizeIfState:
                  type: integer
                bpfMapSizeNATAffinity:
                  type: integer
                bpfMapSizeNATBackend:
                  type: integer
                bpfMapSizeNATFrontend:
                  type: integer
                bpfMapSizeRoute:
                  type: integer
                bpfPSNATPorts:
                  anyOf:
                    - type: integer
                    - type: string
                  pattern: ^.*
                  x-kubernetes-int-or-string: true
                bpfPolicyDebugEnabled:
                  type: boolean
                chainInsertMode:
                  type: string
                dataplaneDriver:
                  type: string
                dataplaneWatchdogTimeout:
                  type: string
                debugDisableLogDropping:
                  type: boolean
                debugMemoryProfilePath:
                  type: string
                debugSimulateCalcGraphHangAfter:
                  type: string
                debugSimulateDataplaneHangAfter:
                  type: string
                defaultEndpointToHostAction:
                  type: string
                deviceRouteProtocol:
                  type: integer
                deviceRouteSourceAddress:
                  type: string
                deviceRouteSourceAddressIPv6:
                  type: string
                disableConntrackInvalidCheck:
                  type: boolean
                endpointReportingDelay:
                  type: string
                endpointReportingEnabled:
                  type: boolean
                externalNodesList:
                  items:
                    type: string
                  type: array
                failsafeInboundHostPorts:
                  items:
                    properties:
                      net:
                        type: string
                      port:
                        type: integer
                      protocol:
                        type: string
                    required:
                      - port
                      - protocol
                    type: object
                  type: array
                failsafeOutboundHostPorts:
                  items:
                    properties:
                      net:
                        type: string
                      port:
                        type: integer
                      protocol:
                        type: string
                    required:
                      - port
                      - protocol
                    type: object
                  type: array
                featureDetectOverride:
                  type: string
                featureGates:
                  type: string
                floatingIPs:
                  enum:
                    - Enabled
                    - Disabled
                  type: string
                genericXDPEnabled:
                  type: boolean
                healthEnabled:
                  type: boolean
                healthHost:
                  type: string
                healthPort:
                  type: integer
                healthTimeoutOverrides:
                  items:
                    properties:
                      name:
                        type: string
                      timeout:
                        type: string
                    required:
                      - name
                      - timeout
                    type: object
                  type: array
                interfaceExclude:
                  type: string
                interfacePrefix:
                  type: string
                interfaceRefreshInterval:
                  type: string
                ipipEnabled:
                  type: boolean
                ipipMTU:
                  type: integer
                ipsetsRefreshInterval:
                  type: string
                iptablesBackend:
                  type: string
                iptablesFilterAllowAction:
                  type: string
                iptablesFilterDenyAction:
                  type: string
                iptablesLockFilePath:
                  type: string
                iptablesLockProbeInterval:
                  type: string
                iptablesLockTimeout:
                  type: string
                iptablesMangleAllowAction:
                  type: string
                iptablesMarkMask:
                  format: int32
                  type: integer
                iptablesNATOutgoingInterfaceFilter:
                  type: string
                iptablesPostWriteCheckInterval:
                  type: string
                iptablesRefreshInterval:
                  type: string
                ipv6Support:
                  type: boolean
                kubeNodePortRanges:
                  items:
                    anyOf:
                      - type: integer
                      - type: string
                    pattern: ^.*
                    x-kubernetes-int-or-string: true
                  type: array
                logDebugFilenameRegex:
                  type: string
                logFilePath:
                  type: string
                logPrefix:
                  type: string
                logSeverityFile:
                  type: string
                logSeverityScreen:
                  type: string
                logSeveritySys:
                  type: string
                maxIpsetSize:
                  type: integer
                metadataAddr:
                  type: string
                metadataPort:
                  type: integer
                mtuIfacePattern:
                  type: string
                natOutgoingAddress:
                  type: string
                natPortRange:
                  anyOf:
                    - type: integer
                    - type: string
                  pattern: ^.*
                  x-kubernetes-int-or-string: true
                netlinkTimeout:
                  type: string
                openstackRegion:
                  type: string
                policySyncPathPrefix:
                  type: string
                prometheusGoMetricsEnabled:
                  type: boolean
                prometheusMetricsEnabled:
                  type: boolean
                prometheusMetricsHost:
                  type: string
                prometheusMetricsPort:
                  type: integer
                prometheusProcessMetricsEnabled:
                  type: boolean
                prometheusWireGuardMetricsEnabled:
                  type: boolean
                removeExternalRoutes:
                  type: boolean
                reportingInterval:
                  type: string
                reportingTTL:
                  type: string
                routeRefreshInterval:
                  type: string
                routeSource:
                  type: string
                routeSyncDisabled:
                  type: boolean
                routeTableRange:
                  properties:
                    max:
                      type: integer
                    min:
                      type: integer
                  required:
                    - max
                    - min
                  type: object
                routeTableRanges:
                  items:
                    properties:
                      max:
                        type: integer
                      min:
                        type: integer
                    required:
                      - max
                      - min
                    type: object
                  type: array
                serviceLoopPrevention:
                  type: string
                sidecarAccelerationEnabled:
                  type: boolean
                usageReportingEnabled:
                  type: boolean
                usageReportingInitialDelay:
                  type: string
                usageReportingInterval:
                  type: string
                useInternalDataplaneDriver:
                  type: boolean
                vxlanEnabled:
                  type: boolean
                vxlanMTU:
                  type: integer
                vxlanMTUV6:
                  type: integer
                vxlanPort:
                  type: integer
                vxlanVNI:
                  type: integer
                wireguardEnabled:
                  type: boolean
                wireguardEnabledV6:
                  type: boolean
                wireguardHostEncryptionEnabled:
                  type: boolean
                wireguardInterfaceName:
                  type: string
                wireguardInterfaceNameV6:
                  type: string
                wireguardKeepAlive:
                  type: string
                wireguardListeningPort:
                  type: integer
                wireguardListeningPortV6:
                  type: integer
                wireguardMTU:
                  type: integer
                wireguardMTUV6:
                  type: integer
                wireguardRoutingRulePriority:
                  type: integer
                workloadSourceSpoofing:
                  type: string
                xdpEnabled:
                  type: boolean
                xdpRefreshInterval:
                  type: string
              type: object
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
# Source: crds/calico/crd.projectcalico.org_globalnetworkpolicies.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalnetworkpolicies.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: GlobalNetworkPolicy
    listKind: GlobalNetworkPolicyList
    plural: globalnetworkpolicies
    singular: globalnetworkpolicy
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                applyOnForward:
                  type: boolean
                doNotTrack:
                  type: boolean
                egress:
                  items:
                    properties:
                      action:
                        type: string
                      destination:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                      http:
                        properties:
                          methods:
                            items:
                              type: string
                            type: array
                          paths:
                            items:
                              properties:
                                exact:
                                  type: string
                                prefix:
                                  type: string
                              type: object
                            type: array
                        type: object
                      icmp:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      ipVersion:
                        type: integer
                      metadata:
                        properties:
                          annotations:
                            additionalProperties:
                              type: string
                            type: object
                        type: object
                      notICMP:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      notProtocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      protocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      source:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                    required:
                      - action
                    type: object
                  type: array
                ingress:
                  items:
                    properties:
                      action:
                        type: string
                      destination:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                      http:
                        properties:
                          methods:
                            items:
                              type: string
                            type: array
                          paths:
                            items:
                              properties:
                                exact:
                                  type: string
                                prefix:
                                  type: string
                              type: object
                            type: array
                        type: object
                      icmp:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      ipVersion:
                        type: integer
                      metadata:
                        properties:
                          annotations:
                            additionalProperties:
                              type: string
                            type: object
                        type: object
                      notICMP:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      notProtocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      protocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      source:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                    required:
                      - action
                    type: object
                  type: array
                namespaceSelector:
                  type: string
                order:
                  type: number
                preDNAT:
                  type: boolean
                selector:
                  type: string
                serviceAccountSelector:
                  type: string
                types:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_globalnetworksets.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: globalnetworksets.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: GlobalNetworkSet
    listKind: GlobalNetworkSetList
    plural: globalnetworksets
    singular: globalnetworkset
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                nets:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_hostendpoints.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: hostendpoints.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: HostEndpoint
    listKind: HostEndpointList
    plural: hostendpoints
    singular: hostendpoint
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                expectedIPs:
                  items:
                    type: string
                  type: array
                interfaceName:
                  type: string
                node:
                  type: string
                ports:
                  items:
                    properties:
                      name:
                        type: string
                      port:
                        type: integer
                      protocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                    required:
                      - name
                      - port
                      - protocol
                    type: object
                  type: array
                profiles:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_ipamblocks.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamblocks.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMBlock
    listKind: IPAMBlockList
    plural: ipamblocks
    singular: ipamblock
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                affinity:
                  type: string
                allocations:
                  items:
                    type: integer
                    # TODO: This nullable is manually added in. We should update controller-gen
                    # to handle []*int properly itself.
                    nullable: true
                  type: array
                attributes:
                  items:
                    properties:
                      handle_id:
                        type: string
                      secondary:
                        additionalProperties:
                          type: string
                        type: object
                    type: object
                  type: array
                cidr:
                  type: string
                deleted:
                  type: boolean
                sequenceNumber:
                  default: 0
                  format: int64
                  type: integer
                sequenceNumberForAllocation:
                  additionalProperties:
                    format: int64
                    type: integer
                  type: object
                strictAffinity:
                  type: boolean
                unallocated:
                  items:
                    type: integer
                  type: array
              required:
                - allocations
                - attributes
                - cidr
                - strictAffinity
                - unallocated
              type: object
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
# Source: crds/calico/crd.projectcalico.org_ipamconfigs.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamconfigs.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMConfig
    listKind: IPAMConfigList
    plural: ipamconfigs
    singular: ipamconfig
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                autoAllocateBlocks:
                  type: boolean
                maxBlocksPerHost:
                  maximum: 2147483647
                  minimum: 0
                  type: integer
                strictAffinity:
                  type: boolean
              required:
                - autoAllocateBlocks
                - strictAffinity
              type: object
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
# Source: crds/calico/crd.projectcalico.org_ipamhandles.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ipamhandles.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPAMHandle
    listKind: IPAMHandleList
    plural: ipamhandles
    singular: ipamhandle
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                block:
                  additionalProperties:
                    type: integer
                  type: object
                deleted:
                  type: boolean
                handleID:
                  type: string
              required:
                - block
                - handleID
              type: object
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
# Source: crds/calico/crd.projectcalico.org_ippools.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: ippools.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPPool
    listKind: IPPoolList
    plural: ippools
    singular: ippool
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                allowedUses:
                  items:
                    type: string
                  type: array
                blockSize:
                  type: integer
                cidr:
                  type: string
                disableBGPExport:
                  type: boolean
                disabled:
                  type: boolean
                ipip:
                  properties:
                    enabled:
                      type: boolean
                    mode:
                      type: string
                  type: object
                ipipMode:
                  type: string
                nat-outgoing:
                  type: boolean
                natOutgoing:
                  type: boolean
                nodeSelector:
                  type: string
                vxlanMode:
                  type: string
              required:
                - cidr
              type: object
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
# Source: crds/calico/crd.projectcalico.org_ipreservations.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: (devel)
  creationTimestamp: null
  name: ipreservations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: IPReservation
    listKind: IPReservationList
    plural: ipreservations
    singular: ipreservation
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                reservedCIDRs:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_kubecontrollersconfigurations.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: kubecontrollersconfigurations.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: KubeControllersConfiguration
    listKind: KubeControllersConfigurationList
    plural: kubecontrollersconfigurations
    singular: kubecontrollersconfiguration
  preserveUnknownFields: false
  scope: Cluster
  versions:
    - name: v1
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
                controllers:
                  properties:
                    namespace:
                      properties:
                        reconcilerPeriod:
                          type: string
                      type: object
                    node:
                      properties:
                        hostEndpoint:
                          properties:
                            autoCreate:
                              type: string
                          type: object
                        leakGracePeriod:
                          type: string
                        reconcilerPeriod:
                          type: string
                        syncLabels:
                          type: string
                      type: object
                    policy:
                      properties:
                        reconcilerPeriod:
                          type: string
                      type: object
                    serviceAccount:
                      properties:
                        reconcilerPeriod:
                          type: string
                      type: object
                    workloadEndpoint:
                      properties:
                        reconcilerPeriod:
                          type: string
                      type: object
                  type: object
                debugProfilePort:
                  format: int32
                  type: integer
                etcdV3CompactionPeriod:
                  type: string
                healthChecks:
                  type: string
                logSeverityScreen:
                  type: string
                prometheusMetricsPort:
                  type: integer
              required:
                - controllers
              type: object
            status:
              properties:
                environmentVars:
                  additionalProperties:
                    type: string
                  type: object
                runningConfig:
                  properties:
                    controllers:
                      properties:
                        namespace:
                          properties:
                            reconcilerPeriod:
                              type: string
                          type: object
                        node:
                          properties:
                            hostEndpoint:
                              properties:
                                autoCreate:
                                  type: string
                              type: object
                            leakGracePeriod:
                              type: string
                            reconcilerPeriod:
                              type: string
                            syncLabels:
                              type: string
                          type: object
                        policy:
                          properties:
                            reconcilerPeriod:
                              type: string
                          type: object
                        serviceAccount:
                          properties:
                            reconcilerPeriod:
                              type: string
                          type: object
                        workloadEndpoint:
                          properties:
                            reconcilerPeriod:
                              type: string
                          type: object
                      type: object
                    debugProfilePort:
                      format: int32
                      type: integer
                    etcdV3CompactionPeriod:
                      type: string
                    healthChecks:
                      type: string
                    logSeverityScreen:
                      type: string
                    prometheusMetricsPort:
                      type: integer
                  required:
                    - controllers
                  type: object
              type: object
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
# Source: crds/calico/crd.projectcalico.org_networkpolicies.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: networkpolicies.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: NetworkPolicy
    listKind: NetworkPolicyList
    plural: networkpolicies
    singular: networkpolicy
  preserveUnknownFields: false
  scope: Namespaced
  versions:
    - name: v1
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
                egress:
                  items:
                    properties:
                      action:
                        type: string
                      destination:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                      http:
                        properties:
                          methods:
                            items:
                              type: string
                            type: array
                          paths:
                            items:
                              properties:
                                exact:
                                  type: string
                                prefix:
                                  type: string
                              type: object
                            type: array
                        type: object
                      icmp:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      ipVersion:
                        type: integer
                      metadata:
                        properties:
                          annotations:
                            additionalProperties:
                              type: string
                            type: object
                        type: object
                      notICMP:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      notProtocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      protocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      source:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                    required:
                      - action
                    type: object
                  type: array
                ingress:
                  items:
                    properties:
                      action:
                        type: string
                      destination:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                      http:
                        properties:
                          methods:
                            items:
                              type: string
                            type: array
                          paths:
                            items:
                              properties:
                                exact:
                                  type: string
                                prefix:
                                  type: string
                              type: object
                            type: array
                        type: object
                      icmp:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      ipVersion:
                        type: integer
                      metadata:
                        properties:
                          annotations:
                            additionalProperties:
                              type: string
                            type: object
                        type: object
                      notICMP:
                        properties:
                          code:
                            type: integer
                          type:
                            type: integer
                        type: object
                      notProtocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      protocol:
                        anyOf:
                          - type: integer
                          - type: string
                        pattern: ^.*
                        x-kubernetes-int-or-string: true
                      source:
                        properties:
                          namespaceSelector:
                            type: string
                          nets:
                            items:
                              type: string
                            type: array
                          notNets:
                            items:
                              type: string
                            type: array
                          notPorts:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          notSelector:
                            type: string
                          ports:
                            items:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^.*
                              x-kubernetes-int-or-string: true
                            type: array
                          selector:
                            type: string
                          serviceAccounts:
                            properties:
                              names:
                                items:
                                  type: string
                                type: array
                              selector:
                                type: string
                            type: object
                          services:
                            properties:
                              name:
                                type: string
                              namespace:
                                type: string
                            type: object
                        type: object
                    required:
                      - action
                    type: object
                  type: array
                order:
                  type: number
                selector:
                  type: string
                serviceAccountSelector:
                  type: string
                types:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/calico/crd.projectcalico.org_networksets.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: networksets.crd.projectcalico.org
spec:
  group: crd.projectcalico.org
  names:
    kind: NetworkSet
    listKind: NetworkSetList
    plural: networksets
    singular: networkset
  preserveUnknownFields: false
  scope: Namespaced
  versions:
    - name: v1
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
                nets:
                  items:
                    type: string
                  type: array
              type: object
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
# Source: crds/operator.tigera.io_apiservers_crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
  name: apiservers.operator.tigera.io
spec:
  group: operator.tigera.io
  names:
    kind: APIServer
    listKind: APIServerList
    plural: apiservers
    singular: apiserver
  scope: Cluster
  versions:
    - name: v1
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
                apiServerDeployment:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-apiserver
                                          - tigera-queryserver
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                initContainers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-apiserver-certs-key-cert-provisioner
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
              type: object
            status:
              properties:
                state:
                  type: string
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
# Source: crds/operator.tigera.io_imagesets_crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
  name: imagesets.operator.tigera.io
spec:
  group: operator.tigera.io
  names:
    kind: ImageSet
    listKind: ImageSetList
    plural: imagesets
    singular: imageset
  scope: Cluster
  versions:
    - name: v1
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
                images:
                  items:
                    properties:
                      digest:
                        type: string
                      image:
                        type: string
                    required:
                      - digest
                      - image
                    type: object
                  type: array
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
# Source: crds/operator.tigera.io_installations_crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.11.3
  name: installations.operator.tigera.io
spec:
  group: operator.tigera.io
  names:
    kind: Installation
    listKind: InstallationList
    plural: installations
    singular: installation
  scope: Cluster
  versions:
    - name: v1
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
                calicoKubeControllersDeployment:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                          x-kubernetes-map-type: atomic
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-kube-controllers
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
                calicoNetwork:
                  properties:
                    bgp:
                      enum:
                        - Enabled
                        - Disabled
                      type: string
                    containerIPForwarding:
                      enum:
                        - Enabled
                        - Disabled
                      type: string
                    hostPorts:
                      enum:
                        - Enabled
                        - Disabled
                      type: string
                    ipPools:
                      items:
                        properties:
                          blockSize:
                            format: int32
                            type: integer
                          cidr:
                            type: string
                          disableBGPExport:
                            default: false
                            type: boolean
                          encapsulation:
                            enum:
                              - IPIPCrossSubnet
                              - IPIP
                              - VXLAN
                              - VXLANCrossSubnet
                              - None
                            type: string
                          natOutgoing:
                            enum:
                              - Enabled
                              - Disabled
                            type: string
                          nodeSelector:
                            type: string
                        required:
                          - cidr
                        type: object
                      type: array
                    linuxDataplane:
                      enum:
                        - Iptables
                        - BPF
                        - VPP
                      type: string
                    mtu:
                      format: int32
                      type: integer
                    multiInterfaceMode:
                      enum:
                        - None
                        - Multus
                      type: string
                    nodeAddressAutodetectionV4:
                      properties:
                        canReach:
                          type: string
                        cidrs:
                          items:
                            type: string
                          type: array
                        firstFound:
                          type: boolean
                        interface:
                          type: string
                        kubernetes:
                          enum:
                            - NodeInternalIP
                          type: string
                        skipInterface:
                          type: string
                      type: object
                    nodeAddressAutodetectionV6:
                      properties:
                        canReach:
                          type: string
                        cidrs:
                          items:
                            type: string
                          type: array
                        firstFound:
                          type: boolean
                        interface:
                          type: string
                        kubernetes:
                          enum:
                            - NodeInternalIP
                          type: string
                        skipInterface:
                          type: string
                      type: object
                  type: object
                calicoNodeDaemonSet:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                          x-kubernetes-map-type: atomic
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-node
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                initContainers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - install-cni
                                          - hostpath-init
                                          - flexvol-driver
                                          - mount-bpffs
                                          - node-certs-key-cert-provisioner
                                          - calico-node-prometheus-server-tls-key-cert-provisioner
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
                calicoWindowsUpgradeDaemonSet:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                          x-kubernetes-map-type: atomic
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-windows-upgrade
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
                certificateManagement:
                  properties:
                    caCert:
                      format: byte
                      type: string
                    keyAlgorithm:
                      enum:
                        - ""
                        - RSAWithSize2048
                        - RSAWithSize4096
                        - RSAWithSize8192
                        - ECDSAWithCurve256
                        - ECDSAWithCurve384
                        - ECDSAWithCurve521
                      type: string
                    signatureAlgorithm:
                      enum:
                        - ""
                        - SHA256WithRSA
                        - SHA384WithRSA
                        - SHA512WithRSA
                        - ECDSAWithSHA256
                        - ECDSAWithSHA384
                        - ECDSAWithSHA512
                      type: string
                    signerName:
                      type: string
                  required:
                    - caCert
                    - signerName
                  type: object
                cni:
                  properties:
                    ipam:
                      properties:
                        type:
                          enum:
                            - Calico
                            - HostLocal
                            - AmazonVPC
                            - AzureVNET
                          type: string
                      required:
                        - type
                      type: object
                    type:
                      enum:
                        - Calico
                        - GKE
                        - AmazonVPC
                        - AzureVNET
                      type: string
                  required:
                    - type
                  type: object
                componentResources:
                  items:
                    properties:
                      componentName:
                        enum:
                          - Node
                          - Typha
                          - KubeControllers
                        type: string
                      resourceRequirements:
                        properties:
                          limits:
                            additionalProperties:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            type: object
                          requests:
                            additionalProperties:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            type: object
                        type: object
                    required:
                      - componentName
                      - resourceRequirements
                    type: object
                  type: array
                controlPlaneNodeSelector:
                  additionalProperties:
                    type: string
                  type: object
                controlPlaneReplicas:
                  format: int32
                  type: integer
                controlPlaneTolerations:
                  items:
                    properties:
                      effect:
                        type: string
                      key:
                        type: string
                      operator:
                        type: string
                      tolerationSeconds:
                        format: int64
                        type: integer
                      value:
                        type: string
                    type: object
                  type: array
                csiNodeDriverDaemonSet:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                          x-kubernetes-map-type: atomic
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - csi-node-driver
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
                fipsMode:
                  enum:
                    - Enabled
                    - Disabled
                  type: string
                flexVolumePath:
                  type: string
                imagePath:
                  type: string
                imagePrefix:
                  type: string
                imagePullSecrets:
                  items:
                    properties:
                      name:
                        type: string
                    type: object
                    x-kubernetes-map-type: atomic
                  type: array
                kubeletVolumePluginPath:
                  type: string
                kubernetesProvider:
                  enum:
                    - ""
                    - EKS
                    - GKE
                    - AKS
                    - OpenShift
                    - DockerEnterprise
                    - RKE2
                  type: string
                logging:
                  properties:
                    cni:
                      properties:
                        logFileMaxAgeDays:
                          format: int32
                          type: integer
                        logFileMaxCount:
                          format: int32
                          type: integer
                        logFileMaxSize:
                          anyOf:
                            - type: integer
                            - type: string
                          pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                          x-kubernetes-int-or-string: true
                        logSeverity:
                          enum:
                            - Error
                            - Warning
                            - Debug
                            - Info
                          type: string
                      type: object
                  type: object
                nodeMetricsPort:
                  format: int32
                  type: integer
                nodeUpdateStrategy:
                  properties:
                    rollingUpdate:
                      properties:
                        maxSurge:
                          anyOf:
                            - type: integer
                            - type: string
                          x-kubernetes-int-or-string: true
                        maxUnavailable:
                          anyOf:
                            - type: integer
                            - type: string
                          x-kubernetes-int-or-string: true
                      type: object
                    type:
                      type: string
                  type: object
                nonPrivileged:
                  type: string
                registry:
                  type: string
                typhaAffinity:
                  properties:
                    nodeAffinity:
                      properties:
                        preferredDuringSchedulingIgnoredDuringExecution:
                          items:
                            properties:
                              preference:
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
                                  matchFields:
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
                                type: object
                                x-kubernetes-map-type: atomic
                              weight:
                                format: int32
                                type: integer
                            required:
                              - preference
                              - weight
                            type: object
                          type: array
                        requiredDuringSchedulingIgnoredDuringExecution:
                          properties:
                            nodeSelectorTerms:
                              items:
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
                                  matchFields:
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
                                type: object
                                x-kubernetes-map-type: atomic
                              type: array
                          required:
                            - nodeSelectorTerms
                          type: object
                          x-kubernetes-map-type: atomic
                      type: object
                  type: object
                typhaDeployment:
                  properties:
                    metadata:
                      properties:
                        annotations:
                          additionalProperties:
                            type: string
                          type: object
                        labels:
                          additionalProperties:
                            type: string
                          type: object
                      type: object
                    spec:
                      properties:
                        minReadySeconds:
                          format: int32
                          maximum: 2147483647
                          minimum: 0
                          type: integer
                        strategy:
                          properties:
                            rollingUpdate:
                              properties:
                                maxSurge:
                                  anyOf:
                                    - type: integer
                                    - type: string
                                  x-kubernetes-int-or-string: true
                                maxUnavailable:
                                  anyOf:
                                    - type: integer
                                    - type: string
                                  x-kubernetes-int-or-string: true
                              type: object
                          type: object
                        template:
                          properties:
                            metadata:
                              properties:
                                annotations:
                                  additionalProperties:
                                    type: string
                                  type: object
                                labels:
                                  additionalProperties:
                                    type: string
                                  type: object
                              type: object
                            spec:
                              properties:
                                affinity:
                                  properties:
                                    nodeAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              preference:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - preference
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          properties:
                                            nodeSelectorTerms:
                                              items:
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
                                                  matchFields:
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
                                                type: object
                                                x-kubernetes-map-type: atomic
                                              type: array
                                          required:
                                            - nodeSelectorTerms
                                          type: object
                                          x-kubernetes-map-type: atomic
                                      type: object
                                    podAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                    podAntiAffinity:
                                      properties:
                                        preferredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              podAffinityTerm:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              weight:
                                                format: int32
                                                type: integer
                                            required:
                                              - podAffinityTerm
                                              - weight
                                            type: object
                                          type: array
                                        requiredDuringSchedulingIgnoredDuringExecution:
                                          items:
                                            properties:
                                              labelSelector:
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
                                                x-kubernetes-map-type: atomic
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
                                                x-kubernetes-map-type: atomic
                                              namespaces:
                                                items:
                                                  type: string
                                                type: array
                                              topologyKey:
                                                type: string
                                            required:
                                              - topologyKey
                                            type: object
                                          type: array
                                      type: object
                                  type: object
                                containers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - calico-typha
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                initContainers:
                                  items:
                                    properties:
                                      name:
                                        enum:
                                          - typha-certs-key-cert-provisioner
                                        type: string
                                      resources:
                                        properties:
                                          limits:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                          requests:
                                            additionalProperties:
                                              anyOf:
                                                - type: integer
                                                - type: string
                                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                              x-kubernetes-int-or-string: true
                                            type: object
                                        type: object
                                    required:
                                      - name
                                    type: object
                                  type: array
                                nodeSelector:
                                  additionalProperties:
                                    type: string
                                  type: object
                                terminationGracePeriodSeconds:
                                  format: int64
                                  type: integer
                                tolerations:
                                  items:
                                    properties:
                                      effect:
                                        type: string
                                      key:
                                        type: string
                                      operator:
                                        type: string
                                      tolerationSeconds:
                                        format: int64
                                        type: integer
                                      value:
                                        type: string
                                    type: object
                                  type: array
                                topologySpreadConstraints:
                                  items:
                                    properties:
                                      labelSelector:
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
                                        x-kubernetes-map-type: atomic
                                      matchLabelKeys:
                                        items:
                                          type: string
                                        type: array
                                        x-kubernetes-list-type: atomic
                                      maxSkew:
                                        format: int32
                                        type: integer
                                      minDomains:
                                        format: int32
                                        type: integer
                                      nodeAffinityPolicy:
                                        type: string
                                      nodeTaintsPolicy:
                                        type: string
                                      topologyKey:
                                        type: string
                                      whenUnsatisfiable:
                                        type: string
                                    required:
                                      - maxSkew
                                      - topologyKey
                                      - whenUnsatisfiable
                                    type: object
                                  type: array
                              type: object
                          type: object
                      type: object
                  type: object
                typhaMetricsPort:
                  format: int32
                  type: integer
                variant:
                  enum:
                    - Calico
                    - TigeraSecureEnterprise
                  type: string
              type: object
            status:
              properties:
                calicoVersion:
                  type: string
                computed:
                  properties:
                    calicoKubeControllersDeployment:
                      properties:
                        metadata:
                          properties:
                            annotations:
                              additionalProperties:
                                type: string
                              type: object
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                          type: object
                        spec:
                          properties:
                            minReadySeconds:
                              format: int32
                              maximum: 2147483647
                              minimum: 0
                              type: integer
                            template:
                              properties:
                                metadata:
                                  properties:
                                    annotations:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    labels:
                                      additionalProperties:
                                        type: string
                                      type: object
                                  type: object
                                spec:
                                  properties:
                                    affinity:
                                      properties:
                                        nodeAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  preference:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - preference
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              properties:
                                                nodeSelectorTerms:
                                                  items:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  type: array
                                              required:
                                                - nodeSelectorTerms
                                              type: object
                                              x-kubernetes-map-type: atomic
                                          type: object
                                        podAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                        podAntiAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                      type: object
                                    containers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - calico-kube-controllers
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    nodeSelector:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    tolerations:
                                      items:
                                        properties:
                                          effect:
                                            type: string
                                          key:
                                            type: string
                                          operator:
                                            type: string
                                          tolerationSeconds:
                                            format: int64
                                            type: integer
                                          value:
                                            type: string
                                        type: object
                                      type: array
                                  type: object
                              type: object
                          type: object
                      type: object
                    calicoNetwork:
                      properties:
                        bgp:
                          enum:
                            - Enabled
                            - Disabled
                          type: string
                        containerIPForwarding:
                          enum:
                            - Enabled
                            - Disabled
                          type: string
                        hostPorts:
                          enum:
                            - Enabled
                            - Disabled
                          type: string
                        ipPools:
                          items:
                            properties:
                              blockSize:
                                format: int32
                                type: integer
                              cidr:
                                type: string
                              disableBGPExport:
                                default: false
                                type: boolean
                              encapsulation:
                                enum:
                                  - IPIPCrossSubnet
                                  - IPIP
                                  - VXLAN
                                  - VXLANCrossSubnet
                                  - None
                                type: string
                              natOutgoing:
                                enum:
                                  - Enabled
                                  - Disabled
                                type: string
                              nodeSelector:
                                type: string
                            required:
                              - cidr
                            type: object
                          type: array
                        linuxDataplane:
                          enum:
                            - Iptables
                            - BPF
                            - VPP
                          type: string
                        mtu:
                          format: int32
                          type: integer
                        multiInterfaceMode:
                          enum:
                            - None
                            - Multus
                          type: string
                        nodeAddressAutodetectionV4:
                          properties:
                            canReach:
                              type: string
                            cidrs:
                              items:
                                type: string
                              type: array
                            firstFound:
                              type: boolean
                            interface:
                              type: string
                            kubernetes:
                              enum:
                                - NodeInternalIP
                              type: string
                            skipInterface:
                              type: string
                          type: object
                        nodeAddressAutodetectionV6:
                          properties:
                            canReach:
                              type: string
                            cidrs:
                              items:
                                type: string
                              type: array
                            firstFound:
                              type: boolean
                            interface:
                              type: string
                            kubernetes:
                              enum:
                                - NodeInternalIP
                              type: string
                            skipInterface:
                              type: string
                          type: object
                      type: object
                    calicoNodeDaemonSet:
                      properties:
                        metadata:
                          properties:
                            annotations:
                              additionalProperties:
                                type: string
                              type: object
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                          type: object
                        spec:
                          properties:
                            minReadySeconds:
                              format: int32
                              maximum: 2147483647
                              minimum: 0
                              type: integer
                            template:
                              properties:
                                metadata:
                                  properties:
                                    annotations:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    labels:
                                      additionalProperties:
                                        type: string
                                      type: object
                                  type: object
                                spec:
                                  properties:
                                    affinity:
                                      properties:
                                        nodeAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  preference:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - preference
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              properties:
                                                nodeSelectorTerms:
                                                  items:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  type: array
                                              required:
                                                - nodeSelectorTerms
                                              type: object
                                              x-kubernetes-map-type: atomic
                                          type: object
                                        podAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                        podAntiAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                      type: object
                                    containers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - calico-node
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    initContainers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - install-cni
                                              - hostpath-init
                                              - flexvol-driver
                                              - mount-bpffs
                                              - node-certs-key-cert-provisioner
                                              - calico-node-prometheus-server-tls-key-cert-provisioner
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    nodeSelector:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    tolerations:
                                      items:
                                        properties:
                                          effect:
                                            type: string
                                          key:
                                            type: string
                                          operator:
                                            type: string
                                          tolerationSeconds:
                                            format: int64
                                            type: integer
                                          value:
                                            type: string
                                        type: object
                                      type: array
                                  type: object
                              type: object
                          type: object
                      type: object
                    calicoWindowsUpgradeDaemonSet:
                      properties:
                        metadata:
                          properties:
                            annotations:
                              additionalProperties:
                                type: string
                              type: object
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                          type: object
                        spec:
                          properties:
                            minReadySeconds:
                              format: int32
                              maximum: 2147483647
                              minimum: 0
                              type: integer
                            template:
                              properties:
                                metadata:
                                  properties:
                                    annotations:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    labels:
                                      additionalProperties:
                                        type: string
                                      type: object
                                  type: object
                                spec:
                                  properties:
                                    affinity:
                                      properties:
                                        nodeAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  preference:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - preference
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              properties:
                                                nodeSelectorTerms:
                                                  items:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  type: array
                                              required:
                                                - nodeSelectorTerms
                                              type: object
                                              x-kubernetes-map-type: atomic
                                          type: object
                                        podAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                        podAntiAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                      type: object
                                    containers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - calico-windows-upgrade
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    nodeSelector:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    tolerations:
                                      items:
                                        properties:
                                          effect:
                                            type: string
                                          key:
                                            type: string
                                          operator:
                                            type: string
                                          tolerationSeconds:
                                            format: int64
                                            type: integer
                                          value:
                                            type: string
                                        type: object
                                      type: array
                                  type: object
                              type: object
                          type: object
                      type: object
                    certificateManagement:
                      properties:
                        caCert:
                          format: byte
                          type: string
                        keyAlgorithm:
                          enum:
                            - ""
                            - RSAWithSize2048
                            - RSAWithSize4096
                            - RSAWithSize8192
                            - ECDSAWithCurve256
                            - ECDSAWithCurve384
                            - ECDSAWithCurve521
                          type: string
                        signatureAlgorithm:
                          enum:
                            - ""
                            - SHA256WithRSA
                            - SHA384WithRSA
                            - SHA512WithRSA
                            - ECDSAWithSHA256
                            - ECDSAWithSHA384
                            - ECDSAWithSHA512
                          type: string
                        signerName:
                          type: string
                      required:
                        - caCert
                        - signerName
                      type: object
                    cni:
                      properties:
                        ipam:
                          properties:
                            type:
                              enum:
                                - Calico
                                - HostLocal
                                - AmazonVPC
                                - AzureVNET
                              type: string
                          required:
                            - type
                          type: object
                        type:
                          enum:
                            - Calico
                            - GKE
                            - AmazonVPC
                            - AzureVNET
                          type: string
                      required:
                        - type
                      type: object
                    componentResources:
                      items:
                        properties:
                          componentName:
                            enum:
                              - Node
                              - Typha
                              - KubeControllers
                            type: string
                          resourceRequirements:
                            properties:
                              limits:
                                additionalProperties:
                                  anyOf:
                                    - type: integer
                                    - type: string
                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                  x-kubernetes-int-or-string: true
                                type: object
                              requests:
                                additionalProperties:
                                  anyOf:
                                    - type: integer
                                    - type: string
                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                  x-kubernetes-int-or-string: true
                                type: object
                            type: object
                        required:
                          - componentName
                          - resourceRequirements
                        type: object
                      type: array
                    controlPlaneNodeSelector:
                      additionalProperties:
                        type: string
                      type: object
                    controlPlaneReplicas:
                      format: int32
                      type: integer
                    controlPlaneTolerations:
                      items:
                        properties:
                          effect:
                            type: string
                          key:
                            type: string
                          operator:
                            type: string
                          tolerationSeconds:
                            format: int64
                            type: integer
                          value:
                            type: string
                        type: object
                      type: array
                    csiNodeDriverDaemonSet:
                      properties:
                        metadata:
                          properties:
                            annotations:
                              additionalProperties:
                                type: string
                              type: object
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                          type: object
                        spec:
                          properties:
                            minReadySeconds:
                              format: int32
                              maximum: 2147483647
                              minimum: 0
                              type: integer
                            template:
                              properties:
                                metadata:
                                  properties:
                                    annotations:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    labels:
                                      additionalProperties:
                                        type: string
                                      type: object
                                  type: object
                                spec:
                                  properties:
                                    affinity:
                                      properties:
                                        nodeAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  preference:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - preference
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              properties:
                                                nodeSelectorTerms:
                                                  items:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  type: array
                                              required:
                                                - nodeSelectorTerms
                                              type: object
                                              x-kubernetes-map-type: atomic
                                          type: object
                                        podAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                        podAntiAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                      type: object
                                    containers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - csi-node-driver
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    nodeSelector:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    tolerations:
                                      items:
                                        properties:
                                          effect:
                                            type: string
                                          key:
                                            type: string
                                          operator:
                                            type: string
                                          tolerationSeconds:
                                            format: int64
                                            type: integer
                                          value:
                                            type: string
                                        type: object
                                      type: array
                                  type: object
                              type: object
                          type: object
                      type: object
                    fipsMode:
                      enum:
                        - Enabled
                        - Disabled
                      type: string
                    flexVolumePath:
                      type: string
                    imagePath:
                      type: string
                    imagePrefix:
                      type: string
                    imagePullSecrets:
                      items:
                        properties:
                          name:
                            type: string
                        type: object
                        x-kubernetes-map-type: atomic
                      type: array
                    kubeletVolumePluginPath:
                      type: string
                    kubernetesProvider:
                      enum:
                        - ""
                        - EKS
                        - GKE
                        - AKS
                        - OpenShift
                        - DockerEnterprise
                        - RKE2
                      type: string
                    logging:
                      properties:
                        cni:
                          properties:
                            logFileMaxAgeDays:
                              format: int32
                              type: integer
                            logFileMaxCount:
                              format: int32
                              type: integer
                            logFileMaxSize:
                              anyOf:
                                - type: integer
                                - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            logSeverity:
                              enum:
                                - Error
                                - Warning
                                - Debug
                                - Info
                              type: string
                          type: object
                      type: object
                    nodeMetricsPort:
                      format: int32
                      type: integer
                    nodeUpdateStrategy:
                      properties:
                        rollingUpdate:
                          properties:
                            maxSurge:
                              anyOf:
                                - type: integer
                                - type: string
                              x-kubernetes-int-or-string: true
                            maxUnavailable:
                              anyOf:
                                - type: integer
                                - type: string
                              x-kubernetes-int-or-string: true
                          type: object
                        type:
                          type: string
                      type: object
                    nonPrivileged:
                      type: string
                    registry:
                      type: string
                    typhaAffinity:
                      properties:
                        nodeAffinity:
                          properties:
                            preferredDuringSchedulingIgnoredDuringExecution:
                              items:
                                properties:
                                  preference:
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
                                      matchFields:
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
                                    type: object
                                    x-kubernetes-map-type: atomic
                                  weight:
                                    format: int32
                                    type: integer
                                required:
                                  - preference
                                  - weight
                                type: object
                              type: array
                            requiredDuringSchedulingIgnoredDuringExecution:
                              properties:
                                nodeSelectorTerms:
                                  items:
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
                                      matchFields:
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
                                    type: object
                                    x-kubernetes-map-type: atomic
                                  type: array
                              required:
                                - nodeSelectorTerms
                              type: object
                              x-kubernetes-map-type: atomic
                          type: object
                      type: object
                    typhaDeployment:
                      properties:
                        metadata:
                          properties:
                            annotations:
                              additionalProperties:
                                type: string
                              type: object
                            labels:
                              additionalProperties:
                                type: string
                              type: object
                          type: object
                        spec:
                          properties:
                            minReadySeconds:
                              format: int32
                              maximum: 2147483647
                              minimum: 0
                              type: integer
                            strategy:
                              properties:
                                rollingUpdate:
                                  properties:
                                    maxSurge:
                                      anyOf:
                                        - type: integer
                                        - type: string
                                      x-kubernetes-int-or-string: true
                                    maxUnavailable:
                                      anyOf:
                                        - type: integer
                                        - type: string
                                      x-kubernetes-int-or-string: true
                                  type: object
                              type: object
                            template:
                              properties:
                                metadata:
                                  properties:
                                    annotations:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    labels:
                                      additionalProperties:
                                        type: string
                                      type: object
                                  type: object
                                spec:
                                  properties:
                                    affinity:
                                      properties:
                                        nodeAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  preference:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - preference
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              properties:
                                                nodeSelectorTerms:
                                                  items:
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
                                                      matchFields:
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
                                                    type: object
                                                    x-kubernetes-map-type: atomic
                                                  type: array
                                              required:
                                                - nodeSelectorTerms
                                              type: object
                                              x-kubernetes-map-type: atomic
                                          type: object
                                        podAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                        podAntiAffinity:
                                          properties:
                                            preferredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  podAffinityTerm:
                                                    properties:
                                                      labelSelector:
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
                                                        x-kubernetes-map-type: atomic
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
                                                        x-kubernetes-map-type: atomic
                                                      namespaces:
                                                        items:
                                                          type: string
                                                        type: array
                                                      topologyKey:
                                                        type: string
                                                    required:
                                                      - topologyKey
                                                    type: object
                                                  weight:
                                                    format: int32
                                                    type: integer
                                                required:
                                                  - podAffinityTerm
                                                  - weight
                                                type: object
                                              type: array
                                            requiredDuringSchedulingIgnoredDuringExecution:
                                              items:
                                                properties:
                                                  labelSelector:
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
                                                    x-kubernetes-map-type: atomic
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
                                                    x-kubernetes-map-type: atomic
                                                  namespaces:
                                                    items:
                                                      type: string
                                                    type: array
                                                  topologyKey:
                                                    type: string
                                                required:
                                                  - topologyKey
                                                type: object
                                              type: array
                                          type: object
                                      type: object
                                    containers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - calico-typha
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    initContainers:
                                      items:
                                        properties:
                                          name:
                                            enum:
                                              - typha-certs-key-cert-provisioner
                                            type: string
                                          resources:
                                            properties:
                                              limits:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                              requests:
                                                additionalProperties:
                                                  anyOf:
                                                    - type: integer
                                                    - type: string
                                                  pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                                                  x-kubernetes-int-or-string: true
                                                type: object
                                            type: object
                                        required:
                                          - name
                                        type: object
                                      type: array
                                    nodeSelector:
                                      additionalProperties:
                                        type: string
                                      type: object
                                    terminationGracePeriodSeconds:
                                      format: int64
                                      type: integer
                                    tolerations:
                                      items:
                                        properties:
                                          effect:
                                            type: string
                                          key:
                                            type: string
                                          operator:
                                            type: string
                                          tolerationSeconds:
                                            format: int64
                                            type: integer
                                          value:
                                            type: string
                                        type: object
                                      type: array
                                    topologySpreadConstraints:
                                      items:
                                        properties:
                                          labelSelector:
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
                                            x-kubernetes-map-type: atomic
                                          matchLabelKeys:
                                            items:
                                              type: string
                                            type: array
                                            x-kubernetes-list-type: atomic
                                          maxSkew:
                                            format: int32
                                            type: integer
                                          minDomains:
                                            format: int32
                                            type: integer
                                          nodeAffinityPolicy:
                                            type: string
                                          nodeTaintsPolicy:
                                            type: string
                                          topologyKey:
                                            type: string
                                          whenUnsatisfiable:
                                            type: string
                                        required:
                                          - maxSkew
                                          - topologyKey
                                          - whenUnsatisfiable
                                        type: object
                                      type: array
                                  type: object
                              type: object
                          type: object
                      type: object
                    typhaMetricsPort:
                      format: int32
                      type: integer
                    variant:
                      enum:
                        - Calico
                        - TigeraSecureEnterprise
                      type: string
                  type: object
                conditions:
                  items:
                    properties:
                      lastTransitionTime:
                        format: date-time
                        type: string
                      message:
                        maxLength: 32768
                        type: string
                      observedGeneration:
                        format: int64
                        minimum: 0
                        type: integer
                      reason:
                        maxLength: 1024
                        minLength: 1
                        pattern: ^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$
                        type: string
                      status:
                        enum:
                          - "True"
                          - "False"
                          - Unknown
                        type: string
                      type:
                        maxLength: 316
                        pattern: ^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$
                        type: string
                    required:
                      - lastTransitionTime
                      - message
                      - reason
                      - status
                      - type
                    type: object
                  type: array
                imageSet:
                  type: string
                mtu:
                  format: int32
                  type: integer
                variant:
                  enum:
                    - Calico
                    - TigeraSecureEnterprise
                  type: string
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}
---
# Source: crds/operator.tigera.io_tigerastatuses_crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
  name: tigerastatuses.operator.tigera.io
spec:
  group: operator.tigera.io
  names:
    kind: TigeraStatus
    listKind: TigeraStatusList
    plural: tigerastatuses
    singular: tigerastatus
  scope: Cluster
  versions:
    - additionalPrinterColumns:
        - jsonPath: .status.conditions[?(@.type=='Available')].status
          name: Available
          type: string
        - jsonPath: .status.conditions[?(@.type=='Progressing')].status
          name: Progressing
          type: string
        - jsonPath: .status.conditions[?(@.type=='Degraded')].status
          name: Degraded
          type: string
        - jsonPath: .status.conditions[?(@.type=='Available')].lastTransitionTime
          name: Since
          type: date
      name: v1
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
            status:
              properties:
                conditions:
                  items:
                    properties:
                      lastTransitionTime:
                        format: date-time
                        type: string
                      message:
                        type: string
                      observedGeneration:
                        format: int64
                        type: integer
                      reason:
                        type: string
                      status:
                        type: string
                      type:
                        type: string
                    required:
                      - lastTransitionTime
                      - status
                      - type
                    type: object
                  type: array
              required:
                - conditions
              type: object
          type: object
      served: true
      storage: true
      subresources:
        status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
---
# Source: tigera-operator/templates/tigera-operator/02-serviceaccount-tigera-operator.yaml
{{if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tigera-operator
  namespace: tigera-operator
imagePullSecrets: []
{{end}}
---
# Source: tigera-operator/templates/tigera-operator/02-role-tigera-operator.yaml
# Permissions required when running the operator for a Calico cluster.
{{if eq .RBACConfig "rbac"}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: tigera-operator
rules:
  - apiGroups:
      - ""
    resources:
      - namespaces
      - pods
      - podtemplates
      - services
      - endpoints
      - events
      - configmaps
      - secrets
      - serviceaccounts
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
  - apiGroups:
      - ""
    resources:
      - resourcequotas
    verbs:
      - list
      - get
      - watch
  - apiGroups:
      - ""
    resources:
      - resourcequotas
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
    resourceNames:
      - calico-critical-pods
      - tigera-critical-pods
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      # Need to update node labels when migrating nodes.
      - get
      - patch
      - list
      # We need this for Typha autoscaling
      - watch
  - apiGroups:
      - rbac.authorization.k8s.io
    resources:
      - clusterroles
      - clusterrolebindings
      - rolebindings
      - roles
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
      - bind
      - escalate
  - apiGroups:
      - apps
    resources:
      - deployments
      - daemonsets
      - statefulsets
    verbs:
      - create
      - get
      - list
      - patch
      - update
      - delete
      - watch
  - apiGroups:
      - apps
    resourceNames:
      - tigera-operator
    resources:
      - deployments/finalizers
    verbs:
      - update
  - apiGroups:
      - operator.tigera.io
    resources:
      - '*'
    verbs:
      - create
      - get
      - list
      - update
      - patch
      - delete
      - watch
  - apiGroups:
      - networking.k8s.io
    resources:
      - networkpolicies
    verbs:
      - create
      - update
      - delete
      - get
      - list
      - watch
  - apiGroups:
      - crd.projectcalico.org
    resources:
      - felixconfigurations
    verbs:
      - create
      - patch
      - list
      - get
      - watch
  - apiGroups:
      - crd.projectcalico.org
    resources:
      - ippools
      - kubecontrollersconfigurations
      - bgpconfigurations
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - scheduling.k8s.io
    resources:
      - priorityclasses
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
  - apiGroups:
      - apiregistration.k8s.io
    resources:
      - apiservices
    verbs:
      - list
      - watch
      - create
      - update
  # Needed for operator lock
  - apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    verbs:
      - create
      - get
      - list
      - update
      - delete
      - watch
  - apiGroups:
      - storage.k8s.io
    resources:
      - csidrivers
    verbs:
      - list
      - watch
      - update
      - get
      - create
      - delete
  # Add the appropriate pod security policy permissions
  - apiGroups:
      - policy
    resources:
      - podsecuritypolicies
    resourceNames:
      - tigera-operator
    verbs:
      - use
  - apiGroups:
      - policy
    resources:
      - podsecuritypolicies
    verbs:
      - get
      - list
      - watch
      - create
      - update
      - delete
      # Add the permissions to monitor the status of certificatesigningrequests when certificate management is enabled.
  - apiGroups:
      - certificates.k8s.io
    resources:
      - certificatesigningrequests
    verbs:
      - list
      - watch
{{end}}
---
# Source: tigera-operator/templates/tigera-operator/02-rolebinding-tigera-operator.yaml
{{if eq .RBACConfig "rbac"}}
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: tigera-operator
subjects:
  - kind: ServiceAccount
    name: tigera-operator
    namespace: tigera-operator
roleRef:
  kind: ClusterRole
  name: tigera-operator
  apiGroup: rbac.authorization.k8s.io
{{end}}
---
# Source: tigera-operator/templates/tigera-operator/02-tigera-operator.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tigera-operator
  namespace: tigera-operator
  labels:
    k8s-app: tigera-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: tigera-operator
  template:
    metadata:
      labels:
        name: tigera-operator
        k8s-app: tigera-operator
    spec:
      nodeSelector:
        kubernetes.io/os: linux
      tolerations:
        - effect: NoExecute
          operator: Exists
        - effect: NoSchedule
          operator: Exists
      {{if eq .RBACConfig "rbac"}}
      serviceAccountName: tigera-operator
      {{end}}
      hostNetwork: true
      # This must be set when hostNetwork is true or else the cluster services won't resolve
      dnsPolicy: ClusterFirstWithHostNet
      containers:
        - name: tigera-operator
          image: {{.OperatorImage}}
          imagePullPolicy: IfNotPresent
          command:
            - operator
          volumeMounts:
            - name: var-lib-calico
              readOnly: true
              mountPath: /var/lib/calico
          env:
            - name: WATCH_NAMESPACE
              value: ""
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "tigera-operator"
            - name: TIGERA_OPERATOR_INIT_IMAGE_VERSION
              value: v1.30.4
          envFrom:
            - configMapRef:
                name: kubernetes-services-endpoint
                optional: true
      volumes:
        - name: var-lib-calico
          hostPath:
            path: /var/lib/calico
---
# This section includes base Calico installation configuration.
# For more information, see: https://projectcalico.docs.tigera.io/master/reference/installation/api#operator.tigera.io/v1.Installation
{{if .CalicoUseOperator }}
apiVersion: operator.tigera.io/v1
kind: Installation
metadata:
  name: default
spec:
  # Configures Calico networking.
  calicoNetwork:
{{- if .MTU }}
{{- if ne .MTU 0 }}
    mtu: {{.MTU}}
{{- end}}
{{- else }}
    mtu: 1440
{{- end}}
    # Note: The ipPools section cannot be modified post-install.
    ipPools:
{{range $cidrs }}
    - cidr: "{{ . }}"
#      encapsulation: IPIP VXLANCrossSubnet
      natOutgoing: Enabled
      nodeSelector: all()
{{end}}
{{end}}
---

# This section configures the Calico API server.
# For more information, see: https://projectcalico.docs.tigera.io/master/reference/installation/api#operator.tigera.io/v1.APIServer
{{if .CalicoUseOperator }}
apiVersion: operator.tigera.io/v1
kind: APIServer
metadata:
  name: default
spec: {}
{{end}}
---
# Source: calico/templates/calico-kube-controllers.yaml
# This manifest creates a Pod Disruption Budget for Controller to allow K8s Cluster Autoscaler to evict
{{if not .CalicoUseOperator}}
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: calico-kube-controllers
  namespace: kube-system
  labels:
    k8s-app: calico-kube-controllers
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      k8s-app: calico-kube-controllers
{{end}}
---
# Source: calico/templates/calico-kube-controllers.yaml
{{if not .CalicoUseOperator}}
{{if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calico-kube-controllers
  namespace: kube-system
{{end}}
{{end}}
---
# Source: calico/templates/calico-node.yaml
{{if not .CalicoUseOperator}}
{{if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calico-node
  namespace: kube-system
{{end}}
{{end}}
---
# Source: calico/templates/calico-node.yaml
{{if not .CalicoUseOperator}}
{{if eq .RBACConfig "rbac"}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: calico-cni-plugin
  namespace: kube-system
{{end}}
{{end}}
---
# Source: calico/templates/calico-config.yaml
# This ConfigMap is used to configure a self-hosted Calico installation.
{{if not .CalicoUseOperator}}
kind: ConfigMap
apiVersion: v1
metadata:
  name: calico-config
  namespace: kube-system
data:
  # Typha is disabled.
  typha_service_name: "none"
  # Configure the backend to use.
  calico_backend: "bird"
  # Configure the MTU to use for workload interfaces and tunnels.
  # By default, MTU is auto-detected, and explicitly setting this field should not be required.
  # You can override auto-detection by providing a non-zero value.
{{- if .MTU }}
{{- if ne .MTU 0 }}
  veth_mtu: "{{.MTU}}"
{{- end}}
{{- else }}
  veth_mtu: "1440"
{{- end}}
  # The CNI network configuration to install on each node. The special
  # values in this config will be automatically populated.
  # Rancher specific change: "assign_ipv6": "true" if dualstack configuration is found
  cni_network_config: |-
    {
      "name": "k8s-pod-network",
      "cniVersion": "0.3.1",
      "plugins": [
        {
          "type": "calico",
          "log_level": "info",
          "log_file_path": "/var/log/calico/cni/cni.log",
          "datastore_type": "kubernetes",
          "nodename": "__KUBERNETES_NODE_NAME__",
          "mtu": __CNI_MTU__,
          "ipam": {
{{- if eq (len $cidrs) 2 }}
              "type": "calico-ipam",
              "assign_ipv4": "true",
              "assign_ipv6": "true"
{{- else }}
              "type": "calico-ipam"
{{- end}}
          },
          "policy": {
              "type": "k8s"
          },
          "kubernetes": {
              "kubeconfig": "__KUBECONFIG_FILEPATH__"
          }
        },
        {
          "type": "portmap",
          "snat": true,
          "capabilities": {"portMappings": true}
        },
        {
          "type": "bandwidth",
          "capabilities": {"bandwidth": true}
        }
      ]
    }
{{end}}
---
# Source: calico/templates/calico-kube-controllers-rbac.yaml
# Include a clusterrole for the kube-controllers component,
# and bind it to the calico-kube-controllers serviceaccount.
{{if not .CalicoUseOperator}}
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: calico-kube-controllers
rules:
  # Nodes are watched to monitor for deletions.
  - apiGroups: [""]
    resources:
      - nodes
    verbs:
      - watch
      - list
      - get
  # Pods are watched to check for existence as part of IPAM controller.
  - apiGroups: [""]
    resources:
      - pods
    verbs:
      - get
      - list
      - watch
  # IPAM resources are manipulated in response to node and block updates, as well as periodic triggers.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - ipreservations
    verbs:
      - list
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - blockaffinities
      - ipamblocks
      - ipamhandles
    verbs:
      - get
      - list
      - create
      - update
      - delete
      - watch
  # Pools are watched to maintain a mapping of blocks to IP pools.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - ippools
    verbs:
      - list
      - watch
  # kube-controllers manages hostendpoints.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - hostendpoints
    verbs:
      - get
      - list
      - create
      - update
      - delete
  # Needs access to update clusterinformations.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - clusterinformations
    verbs:
      - get
      - list
      - create
      - update
      - watch
  # KubeControllersConfiguration is where it gets its config
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - kubecontrollersconfigurations
    verbs:
      # read its own config
      - get
      # create a default if none exists
      - create
      # update status
      - update
      # watch for changes
      - watch
{{end}}
---
# Source: calico/templates/calico-node-rbac.yaml
# Include a clusterrole for the calico-node DaemonSet,
# and bind it to the calico-node serviceaccount.
{{if not .CalicoUseOperator}}
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: calico-node
rules:
  # Used for creating service account tokens to be used by the CNI plugin
  - apiGroups: [""]
    resources:
      - serviceaccounts/token
    resourceNames:
      - calico-cni-plugin
    verbs:
      - create
  # The CNI plugin needs to get pods, nodes, and namespaces.
  - apiGroups: [""]
    resources:
      - pods
      - nodes
      - namespaces
    verbs:
      - get
  # EndpointSlices are used for Service-based network policy rule
  # enforcement.
  - apiGroups: ["discovery.k8s.io"]
    resources:
      - endpointslices
    verbs:
      - watch
      - list
  - apiGroups: [""]
    resources:
      - endpoints
      - services
    verbs:
      # Used to discover service IPs for advertisement.
      - watch
      - list
      # Used to discover Typhas.
      - get
  # Pod CIDR auto-detection on kubeadm needs access to config maps.
  - apiGroups: [""]
    resources:
      - configmaps
    verbs:
      - get
  - apiGroups: [""]
    resources:
      - nodes/status
    verbs:
      # Needed for clearing NodeNetworkUnavailable flag.
      - patch
      # Calico stores some configuration information in node annotations.
      - update
  # Watch for changes to Kubernetes NetworkPolicies.
  - apiGroups: ["networking.k8s.io"]
    resources:
      - networkpolicies
    verbs:
      - watch
      - list
  # Used by Calico for policy information.
  - apiGroups: [""]
    resources:
      - pods
      - namespaces
      - serviceaccounts
    verbs:
      - list
      - watch
  # The CNI plugin patches pods/status.
  - apiGroups: [""]
    resources:
      - pods/status
    verbs:
      - patch
  # Calico monitors various CRDs for config.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - globalfelixconfigs
      - felixconfigurations
      - bgppeers
      - bgpfilters
      - globalbgpconfigs
      - bgpconfigurations
      - ippools
      - ipreservations
      - ipamblocks
      - globalnetworkpolicies
      - globalnetworksets
      - networkpolicies
      - networksets
      - clusterinformations
      - hostendpoints
      - blockaffinities
      - caliconodestatuses
    verbs:
      - get
      - list
      - watch
  # Calico must create and update some CRDs on startup.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - ippools
      - felixconfigurations
      - clusterinformations
    verbs:
      - create
      - update
  # Calico must update some CRDs.
  - apiGroups: [ "crd.projectcalico.org" ]
    resources:
      - caliconodestatuses
    verbs:
      - update
  # Calico stores some configuration information on the node.
  - apiGroups: [""]
    resources:
      - nodes
    verbs:
      - get
      - list
      - watch
  # These permissions are only required for upgrade from v2.6, and can
  # be removed after upgrade or on fresh installations.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - bgpconfigurations
      - bgppeers
    verbs:
      - create
      - update
  # These permissions are required for Calico CNI to perform IPAM allocations.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - blockaffinities
      - ipamblocks
      - ipamhandles
    verbs:
      - get
      - list
      - create
      - update
      - delete
  # The CNI plugin and calico/node need to be able to create a default
  # IPAMConfiguration
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - ipamconfigs
    verbs:
      - get
      - create
  # Block affinities must also be watchable by confd for route aggregation.
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - blockaffinities
    verbs:
      - watch
  # The Calico IPAM migration needs to get daemonsets. These permissions can be
  # removed if not upgrading from an installation using host-local IPAM.
  - apiGroups: ["apps"]
    resources:
      - daemonsets
    verbs:
      - get
{{end}}
---
# Source: calico/templates/calico-node-rbac.yaml
# CNI cluster role
{{if not .CalicoUseOperator}}
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: calico-cni-plugin
rules:
  - apiGroups: [""]
    resources:
      - pods
      - nodes
      - namespaces
    verbs:
      - get
  - apiGroups: [""]
    resources:
      - pods/status
    verbs:
      - patch
  - apiGroups: ["crd.projectcalico.org"]
    resources:
      - blockaffinities
      - ipamblocks
      - ipamhandles
      - clusterinformations
      - ippools
      - ipreservations
      - ipamconfigs
    verbs:
      - get
      - list
      - create
      - update
      - delete
{{end}}
---
# Source: calico/templates/calico-kube-controllers-rbac.yaml
{{if not .CalicoUseOperator}}
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: calico-kube-controllers
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: calico-kube-controllers
subjects:
- kind: ServiceAccount
  name: calico-kube-controllers
  namespace: kube-system
{{end}}
---
# Source: calico/templates/calico-node-rbac.yaml
{{if not .CalicoUseOperator}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: calico-node
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: calico-node
subjects:
- kind: ServiceAccount
  name: calico-node
  namespace: kube-system
{{end}}
---
# Source: calico/templates/calico-node-rbac.yaml
{{if not .CalicoUseOperator}}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: calico-cni-plugin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: calico-cni-plugin
subjects:
- kind: ServiceAccount
  name: calico-cni-plugin
  namespace: kube-system
{{end}}
---
# Source: calico/templates/calico-node.yaml
# This manifest installs the calico-node container, as well
# as the CNI plugins and network config on
# each master and worker node in a Kubernetes cluster.
{{if not .CalicoUseOperator}}
kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: calico-node
  namespace: kube-system
  labels:
    k8s-app: calico-node
spec:
  selector:
    matchLabels:
      k8s-app: calico-node
  updateStrategy:
{{if .UpdateStrategy}}
{{ toYaml .UpdateStrategy | indent 4}}
{{else}}
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
{{end}}
  template:
    metadata:
      labels:
        k8s-app: calico-node
      # Rancher-specific: The annotation for scheduler.alpha.kubernetes.io/critical-pod originated from the v3.13.4 base
      annotations:
        # This, along with the CriticalAddonsOnly toleration below,
        # marks the pod as a critical add-on, ensuring it gets
        # priority scheduling and that its resources are reserved
        # if it ever gets evicted.
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      nodeSelector:
        kubernetes.io/os: linux
      {{ range $k, $v := .NodeSelector }}
        {{ $k }}: "{{ $v }}"
      {{ end }}
      hostNetwork: true
      tolerations:
        # Make sure calico-node gets scheduled on all nodes.
        - effect: NoSchedule
          operator: Exists
        # Mark the pod as a critical add-on for rescheduling.
        - key: CriticalAddonsOnly
          operator: Exists
        - effect: NoExecute
          operator: Exists
      {{if eq .RBACConfig "rbac"}}
      serviceAccountName: calico-node
      {{end}}
      # Minimize downtime during a rolling upgrade or deletion; tell Kubernetes to do a "force
      # deletion": https://kubernetes.io/docs/concepts/workloads/pods/pod/#termination-of-pods.
      terminationGracePeriodSeconds: 0
      # Rancher specific change
      priorityClassName: {{ .CalicoNodePriorityClassName | default "system-node-critical" }}
      initContainers:
        # This container performs upgrade from host-local IPAM to calico-ipam.
        # It can be deleted if this is a fresh installation, or if you have already
        # upgraded to use calico-ipam.
        - name: upgrade-ipam
          image: {{.CNIImage}}
          imagePullPolicy: IfNotPresent
          command: ["/opt/cni/bin/calico-ipam", "-upgrade"]
          envFrom:
          - configMapRef:
              # Allow KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT to be overridden for eBPF mode.
              name: kubernetes-services-endpoint
              optional: true
          env:
            - name: KUBERNETES_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: CALICO_NETWORKING_BACKEND
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: calico_backend
          volumeMounts:
            - mountPath: /var/lib/cni/networks
              name: host-local-net-dir
            - mountPath: /host/opt/cni/bin
              name: cni-bin-dir
          securityContext:
            privileged: true
        # This container installs the CNI binaries
        # and CNI network config file on each node.
        - name: install-cni
          image: docker.io/calico/cni:v3.26.1
          imagePullPolicy: IfNotPresent
          command: ["/opt/cni/bin/install"]
          envFrom:
          - configMapRef:
              # Allow KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT to be overridden for eBPF mode.
              name: kubernetes-services-endpoint
              optional: true
          env:
            # Name of the CNI config file to create.
            - name: CNI_CONF_NAME
              value: "10-calico.conflist"
            # The CNI network config to install on each node.
            - name: CNI_NETWORK_CONFIG
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: cni_network_config
            # Set the hostname based on the k8s node name.
            - name: KUBERNETES_NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            # CNI MTU Config variable
            - name: CNI_MTU
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: veth_mtu
            # Prevents the container from sleeping forever.
            - name: SLEEP
              value: "false"
          volumeMounts:
            - mountPath: /host/opt/cni/bin
              name: cni-bin-dir
            - mountPath: /host/etc/cni/net.d
              name: cni-net-dir
          securityContext:
            privileged: true
        # This init container mounts the necessary filesystems needed by the BPF data plane
        # i.e. bpf at /sys/fs/bpf and cgroup2 at /run/calico/cgroup. Calico-node initialisation is executed
        # in best effort fashion, i.e. no failure for errors, to not disrupt pod creation in iptable mode.
        - name: "mount-bpffs"
          image: {{.NodeImage}}
          imagePullPolicy: IfNotPresent
          command: ["calico-node", "-init", "-best-effort"]
          volumeMounts:
            - mountPath: /sys/fs
              name: sys-fs
              # Bidirectional is required to ensure that the new mount we make at /sys/fs/bpf propagates to the host
              # so that it outlives the init container.
              mountPropagation: Bidirectional
            - mountPath: /var/run/calico
              name: var-run-calico
              # Bidirectional is required to ensure that the new mount we make at /run/calico/cgroup propagates to the host
              # so that it outlives the init container.
              mountPropagation: Bidirectional
            # Mount /proc/ from host which usually is an init program at /nodeproc. It's needed by mountns binary,
            # executed by calico-node, to mount root cgroup2 fs at /run/calico/cgroup to attach CTLB programs correctly.
            - mountPath: /nodeproc
              name: nodeproc
              readOnly: true
          securityContext:
            privileged: true
      containers:
        # Runs calico-node container on each Kubernetes node. This
        # container programs network policy and routes on each
        # host.
        - name: calico-node
          image: {{.NodeImage}}
          imagePullPolicy: IfNotPresent
          envFrom:
          - configMapRef:
              # Allow KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT to be overridden for eBPF mode.
              name: kubernetes-services-endpoint
              optional: true
          env:
            # Use Kubernetes API as the backing datastore.
            - name: DATASTORE_TYPE
              value: "kubernetes"
            # Wait for the datastore.
            - name: WAIT_FOR_DATASTORE
              value: "true"
            # Set based on the k8s node name.
            - name: NODENAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            # Choose the backend to use.
            - name: CALICO_NETWORKING_BACKEND
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: calico_backend
            # Cluster type to identify the deployment type
            - name: CLUSTER_TYPE
              value: "k8s,bgp"
            # Auto-detect the BGP IP address.
            - name: IP
              value: "autodetect"
            # Enable IPIP
            - name: CALICO_IPV4POOL_IPIP
              value: "Always"
            # Enable or Disable VXLAN on the default IP pool.
            - name: CALICO_IPV4POOL_VXLAN
              value: "Never"
            # Enable or Disable VXLAN on the default IPv6 IP pool.
            - name: CALICO_IPV6POOL_VXLAN
              value: "Never"
            # Set MTU for tunnel device used if ipip is enabled
            - name: FELIX_IPINIPMTU
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: veth_mtu
            # Set MTU for the VXLAN tunnel device.
            - name: FELIX_VXLANMTU
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: veth_mtu
            # Set MTU for the Wireguard tunnel device.
            - name: FELIX_WIREGUARDMTU
              valueFrom:
                configMapKeyRef:
                  name: calico-config
                  key: veth_mtu
            # The default IPv4 pool to create on startup if none exists. Pod IPs will be
            # chosen from this range. Changing this value after installation will have
            # no effect. This should fall within --cluster-cidr.
            # Rancher-specific: Explicitly set CALICO_IPV4POOL_CIDR/CALICO_IPV6POOL_CIDR
{{- if eq (len $cidrs) 2 }}
            - name: CALICO_IPV4POOL_CIDR
              value: "{{ first $cidrs }}"
            - name: CALICO_IPV6POOL_CIDR
              value: "{{ last $cidrs }}"
{{- else }}
            - name: CALICO_IPV4POOL_CIDR
              value: "{{.ClusterCIDR}}"
{{- end}}
            # Disable file logging so kubectl logs works.
            - name: CALICO_DISABLE_FILE_LOGGING
              value: "true"
            # Set Felix endpoint to host default action to ACCEPT.
            - name: FELIX_DEFAULTENDPOINTTOHOSTACTION
              value: "ACCEPT"
            # Rancher-specific: Support dualstack
{{- if eq (len $cidrs) 2 }}
            - name: FELIX_IPV6SUPPORT
              value: "true"
            - name: IP6
              value: autodetect
            - name: CALICO_IPV6POOL_NAT_OUTGOING
              value: "true"
{{- else }}
            # Disable IPv6 on Kubernetes.
            - name: FELIX_IPV6SUPPORT
              value: "false"
{{- end}}
            # Rancher-specific: Define and set FELIX_LOGFILEPATH to none to disable felix logging to file
            - name: FELIX_LOGFILEPATH
              value: "none"
            # Rancher-specific: Define and set FELIX_LOGSEVERITYSYS to no value from default info to disable felix logging to syslog
            - name: FELIX_LOGSEVERITYSYS
              value: ""
            # Set Felix logging to "warning"
            # Rancher-specific: Set FELIX_LOGSEVERITYSCREEN to Warning from default info
            - name: FELIX_LOGSEVERITYSCREEN
              value: "Warning"
            - name: FELIX_HEALTHENABLED
              value: "true"
            # Rancher-specific: Set FELIX_IPTABLESBACKEND to auto for autodetection of nftables
            - name: FELIX_IPTABLESBACKEND
              value: "auto"
          securityContext:
            privileged: true
          resources:
            requests:
              cpu: 250m
          lifecycle:
            preStop:
              exec:
                command:
                - /bin/calico-node
                - -shutdown
          livenessProbe:
            exec:
              command:
              - /bin/calico-node
              - -felix-live
              - -bird-live
            periodSeconds: 10
            initialDelaySeconds: 10
            failureThreshold: 6
            timeoutSeconds: 10
          readinessProbe:
            exec:
              command:
              - /bin/calico-node
              - -felix-ready
              - -bird-ready
            periodSeconds: 10
            timeoutSeconds: 10
          volumeMounts:
            # For maintaining CNI plugin API credentials.
            - mountPath: /host/etc/cni/net.d
              name: cni-net-dir
              readOnly: false
            - mountPath: /lib/modules
              name: lib-modules
              readOnly: true
            - mountPath: /run/xtables.lock
              name: xtables-lock
              readOnly: false
            - mountPath: /var/run/calico
              name: var-run-calico
              readOnly: false
            - mountPath: /var/lib/calico
              name: var-lib-calico
              readOnly: false
            - name: policysync
              mountPath: /var/run/nodeagent
            # For eBPF mode, we need to be able to mount the BPF filesystem at /sys/fs/bpf so we mount in the
            # parent directory.
            - name: bpffs
              mountPath: /sys/fs/bpf
            - name: cni-log-dir
              mountPath: /var/log/calico/cni
              readOnly: true
      volumes:
        # Used by calico-node.
        - name: lib-modules
          hostPath:
            path: /lib/modules
        - name: var-run-calico
          hostPath:
            path: /var/run/calico
        - name: var-lib-calico
          hostPath:
            path: /var/lib/calico
        - name: xtables-lock
          hostPath:
            path: /run/xtables.lock
            type: FileOrCreate
        - name: sys-fs
          hostPath:
            path: /sys/fs/
            type: DirectoryOrCreate
        - name: bpffs
          hostPath:
            path: /sys/fs/bpf
            type: Directory
        # mount /proc at /nodeproc to be used by mount-bpffs initContainer to mount root cgroup2 fs.
        - name: nodeproc
          hostPath:
            path: /proc
        # Used to install CNI.
        - name: cni-bin-dir
          hostPath:
            path: /opt/cni/bin
        - name: cni-net-dir
          hostPath:
            path: /etc/cni/net.d
        # Used to access CNI logs.
        - name: cni-log-dir
          hostPath:
            path: /var/log/calico/cni
        # Mount in the directory for host-local IPAM allocations. This is
        # used when upgrading from host-local to calico-ipam, and can be removed
        # if not using the upgrade-ipam init container.
        - name: host-local-net-dir
          hostPath:
            path: /var/lib/cni/networks
        # Used to create per-pod Unix Domain Sockets
        - name: policysync
          hostPath:
            type: DirectoryOrCreate
            path: /var/run/nodeagent
{{end}}
---
# Source: calico/templates/calico-kube-controllers.yaml
# See https://github.com/projectcalico/kube-controllers
{{if not .CalicoUseOperator}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: calico-kube-controllers
  namespace: kube-system
  labels:
    k8s-app: calico-kube-controllers
spec:
  # The controllers can only have a single active instance.
  replicas: 1
  selector:
    matchLabels:
      k8s-app: calico-kube-controllers
  strategy:
    type: Recreate
  template:
    metadata:
      name: calico-kube-controllers
      namespace: kube-system
      labels:
        k8s-app: calico-kube-controllers
      # Added by Rancher to mark as a critical pod.
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      nodeSelector:
        kubernetes.io/os: linux
{{- if .Tolerations }}
      tolerations:
{{ toYaml .Tolerations | indent 6}}
{{- else }}
      tolerations:
        # Rancher-specific: Set tolerations on the calico-kube-controllers so as to let it run on all nodes.
        # Make sure calico-node gets scheduled on all nodes.
        - effect: NoSchedule
          operator: Exists
        # Mark the pod as a critical add-on for rescheduling.
        - key: CriticalAddonsOnly
          operator: Exists
        - key: node-role.kubernetes.io/master
          effect: NoSchedule
        - key: node-role.kubernetes.io/control-plane
          effect: NoSchedule
{{- end }}
      {{if eq .RBACConfig "rbac"}}
      serviceAccountName: calico-kube-controllers
      {{end}}
      # Rancher specific change
      priorityClassName: {{ .CalicoKubeControllersPriorityClassName | default "system-cluster-critical" }}
      containers:
        - name: calico-kube-controllers
          image: {{.ControllersImage}}
          imagePullPolicy: IfNotPresent
          env:
            # Choose which controllers to run.
            - name: ENABLED_CONTROLLERS
              value: node
            - name: DATASTORE_TYPE
              value: kubernetes
          livenessProbe:
            exec:
              command:
              - /usr/bin/check-status
              - -l
            periodSeconds: 10
            initialDelaySeconds: 10
            failureThreshold: 6
            timeoutSeconds: 10
          readinessProbe:
            exec:
              command:
              - /usr/bin/check-status
              - -r
            periodSeconds: 10
{{end}}
`
