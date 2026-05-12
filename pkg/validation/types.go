package main

type RKEReleases struct {
	Channels    []Channels    `yaml:"channels" json:"channels"`
	AppDefaults []AppDefaults `yaml:"appDefaults" json:"appDefaults"`
	Releases    []Releases    `yaml:"releases" json:"releases"`
}

// Channels represents a release channel
type Channels struct {
	Name   string `yaml:"name" json:"name"`
	Latest string `yaml:"latest" json:"latest"`
}

// Defaults represents default version information
type Defaults struct {
	AppVersion     string `yaml:"appVersion" json:"appVersion"`
	DefaultVersion string `yaml:"defaultVersion" json:"defaultVersion"`
}

// AppDefaults represents application default versions
type AppDefaults struct {
	AppName  string     `yaml:"appName" json:"appName"`
	Defaults []Defaults `yaml:"defaults" json:"defaults"`
}

// Common argument types for reuse

// GenericArgument represents a simple string argument
type GenericArgument struct {
	Type string `yaml:"type" json:"type"`
}

// BooleanArgument represents a boolean argument with a default value
type BooleanArgument struct {
	Type    string `yaml:"type" json:"type"`
	Default bool   `yaml:"default" json:"default"`
}

// ArrayArgument represents an array argument with optional options
type ArrayArgument struct {
	Type    string   `yaml:"type" json:"type"`
	Options []string `yaml:"options,omitempty" json:"options,omitempty"`
}

// EnumArgument represents an enum argument with options and optional default

type EnumArgument struct {
	Type     string      `yaml:"type" json:"type"`
	Default  interface{} `yaml:"default,omitempty" json:"default,omitempty"`
	Options  []string    `yaml:"options" json:"options"`
	Nullable bool        `yaml:"nullable,omitempty" json:"nullable,omitempty"`
}

// CNIArgument represents the CNIArgument configuration with options
type CNIArgument struct {
	Default string   `yaml:"default" json:"default"`
	Type    string   `yaml:"type" json:"type"`
	Options []string `yaml:"options" json:"options"`
}

// ServerArgs represents server configuration arguments
type ServerArgs struct {
	KubeApiserverArg              GenericArgument `yaml:"kube-apiserver-arg" json:"kube-apiserver-arg"`
	KubeSchedulerArg              GenericArgument `yaml:"kube-scheduler-arg" json:"kube-scheduler-arg"`
	ClusterDNS                    GenericArgument `yaml:"cluster-dns" json:"cluster-dns"`
	ServiceCidr                   GenericArgument `yaml:"service-cidr" json:"service-cidr"`
	KubeControllerManagerArg      GenericArgument `yaml:"kube-controller-manager-arg" json:"kube-controller-manager-arg"`
	KubeProxyArg                  GenericArgument `yaml:"kube-proxy-arg" json:"kube-proxy-arg"`
	TLSSan                        ArrayArgument   `yaml:"tls-san" json:"tls-san"`
	EtcdExposeMetrics             BooleanArgument `yaml:"etcd-expose-metrics" json:"etcd-expose-metrics"`
	Disable                       ArrayArgument   `yaml:"disable" json:"disable"`
	DisableKubeProxy              BooleanArgument `yaml:"disable-kube-proxy" json:"disable-kube-proxy"`
	ClusterCidr                   GenericArgument `yaml:"cluster-cidr" json:"cluster-cidr"`
	AuditPolicyFile               GenericArgument `yaml:"audit-policy-file" json:"audit-policy-file"`
	ServiceNodePortRange          GenericArgument `yaml:"service-node-port-range" json:"service-node-port-range"`
	ClusterDomain                 GenericArgument `yaml:"cluster-domain" json:"cluster-domain"`
	Cni                           CNIArgument     `yaml:"cni" json:"cni"`
	ContainerRuntimeEndpoint      GenericArgument `yaml:"container-runtime-endpoint" json:"container-runtime-endpoint"`
	Snapshotter                   GenericArgument `yaml:"snapshotter" json:"snapshotter"`
	KubeApiserverImage            GenericArgument `yaml:"kube-apiserver-image" json:"kube-apiserver-image"`
	KubeControllerManagerImage    GenericArgument `yaml:"kube-controller-manager-image" json:"kube-controller-manager-image"`
	KubeSchedulerImage            GenericArgument `yaml:"kube-scheduler-image" json:"kube-scheduler-image"`
	PauseImage                    GenericArgument `yaml:"pause-image" json:"pause-image"`
	RuntimeImage                  GenericArgument `yaml:"runtime-image" json:"runtime-image"`
	EtcdImage                     GenericArgument `yaml:"etcd-image" json:"etcd-image"`
	DisableScheduler              GenericArgument `yaml:"disable-scheduler" json:"disable-scheduler"`
	DisableCloudController        GenericArgument `yaml:"disable-cloud-controller" json:"disable-cloud-controller"`
	DisableApiserver              BooleanArgument `yaml:"disable-apiserver,omitempty" json:"disable-apiserver,omitempty"`
	DisableControllerManager      BooleanArgument `yaml:"disable-controller-manager,omitempty" json:"disable-controller-manager,omitempty"`
	DisableEtcd                   BooleanArgument `yaml:"disable-etcd,omitempty" json:"disable-etcd,omitempty"`
	DisableNetworkPolicy          BooleanArgument `yaml:"disable-network-policy,omitempty" json:"disable-network-policy,omitempty"`
	KubeletPath                   GenericArgument `yaml:"kubelet-path" json:"kubelet-path"`
	EtcdArg                       GenericArgument `yaml:"etcd-arg,omitempty" json:"etcd-arg,omitempty"`
	EtcdS3BucketLookupType        EnumArgument    `yaml:"etcd-s3-bucket-lookup-type,omitempty" json:"etcd-s3-bucket-lookup-type,omitempty"`
	EgressSelectorMode            GenericArgument `yaml:"egress-selector-mode,omitempty" json:"egress-selector-mode,omitempty"`
	DefaultLocalStoragePath       GenericArgument `yaml:"default-local-storage-path,omitempty" json:"default-local-storage-path,omitempty"`
	FlannelBackend                EnumArgument    `yaml:"flannel-backend,omitempty" json:"flannel-backend,omitempty"`
	FlannelIPv6Masq               BooleanArgument `yaml:"flannel-ipv6-masq,omitempty" json:"flannel-ipv6-masq,omitempty"`
	KineTLS                       BooleanArgument `yaml:"kine-tls,omitempty" json:"kine-tls,omitempty"`
	SecretsEncryption             BooleanArgument `yaml:"secrets-encryption,omitempty" json:"secrets-encryption,omitempty"`
	SecretsEncryptionProvider     EnumArgument    `yaml:"secrets-encryption-provider,omitempty" json:"secrets-encryption-provider,omitempty"`
	TLSSanSecurity                GenericArgument `yaml:"tls-san-security,omitempty" json:"tls-san-security,omitempty"`
	KubeCloudControllerManagerArg GenericArgument `yaml:"kube-cloud-controller-manager-arg,omitempty" json:"kube-cloud-controller-manager-arg,omitempty"`
	DatastoreEndpoint             GenericArgument `yaml:"datastore-endpoint,omitempty" json:"datastore-endpoint,omitempty"`
	DatastoreCafile               GenericArgument `yaml:"datastore-cafile,omitempty" json:"datastore-cafile,omitempty"`
	DatastoreCertfile             GenericArgument `yaml:"datastore-certfile,omitempty" json:"datastore-certfile,omitempty"`
	DatastoreKeyfile              GenericArgument `yaml:"datastore-keyfile,omitempty" json:"datastore-keyfile,omitempty"`
	SupervisorMetrics             GenericArgument `yaml:"supervisor-metrics,omitempty" json:"supervisor-metrics,omitempty"`
	WriteKubeconfigGroup          GenericArgument `yaml:"write-kubeconfig-group,omitempty" json:"write-kubeconfig-group,omitempty"`
	IngressController             EnumArgument    `yaml:"ingress-controller,omitempty" json:"ingress-controller,omitempty"`
	HelmJobImage                  GenericArgument `yaml:"helm-job-image,omitempty" json:"helm-job-image,omitempty"`
	ServicelbNamespace            GenericArgument `yaml:"servicelb-namespace,omitempty" json:"servicelb-namespace,omitempty"`
	EnableServicelb               GenericArgument `yaml:"enable-servicelb,omitempty" json:"enable-servicelb,omitempty"`
	EmbeddedRegistry              GenericArgument `yaml:"embedded-registry,omitempty" json:"embedded-registry,omitempty"`
}

// AgentArgs represents agent configuration arguments
type AgentArgs struct {
	Profile                           EnumArgument    `yaml:"profile" json:"profile"`
	CloudProviderConfig               GenericArgument `yaml:"cloud-provider-config" json:"cloud-provider-config"`
	CloudProviderName                 EnumArgument    `yaml:"cloud-provider-name" json:"cloud-provider-name"`
	Selinux                           GenericArgument `yaml:"selinux" json:"selinux"`
	AuditPolicyFile                   GenericArgument `yaml:"audit-policy-file" json:"audit-policy-file"`
	SystemDefaultRegistry             GenericArgument `yaml:"system-default-registry" json:"system-default-registry"`
	ProtectKernelDefaults             BooleanArgument `yaml:"protect-kernel-defaults" json:"protect-kernel-defaults"`
	KubeletArg                        ArrayArgument   `yaml:"kubelet-arg" json:"kubelet-arg"`
	KubeProxyArg                      ArrayArgument   `yaml:"kube-proxy-arg" json:"kube-proxy-arg"`
	ResolvConf                        GenericArgument `yaml:"resolv-conf" json:"resolv-conf"`
	ControlPlaneResourceRequests      GenericArgument `yaml:"control-plane-resource-requests" json:"control-plane-resource-requests"`
	ControlPlaneResourceLimits        GenericArgument `yaml:"control-plane-resource-limits" json:"control-plane-resource-limits"`
	KubeApiserverExtraMount           ArrayArgument   `yaml:"kube-apiserver-extra-mount" json:"kube-apiserver-extra-mount"`
	KubeSchedulerExtraMount           ArrayArgument   `yaml:"kube-scheduler-extra-mount" json:"kube-scheduler-extra-mount"`
	KubeControllerManagerExtraMount   ArrayArgument   `yaml:"kube-controller-manager-extra-mount" json:"kube-controller-manager-extra-mount"`
	KubeProxyExtraMount               ArrayArgument   `yaml:"kube-proxy-extra-mount" json:"kube-proxy-extra-mount"`
	EtcdExtraMount                    ArrayArgument   `yaml:"etcd-extra-mount" json:"etcd-extra-mount"`
	CloudControllerManagerExtraMount  ArrayArgument   `yaml:"cloud-controller-manager-extra-mount" json:"cloud-controller-manager-extra-mount"`
	KubeApiserverExtraEnv             ArrayArgument   `yaml:"kube-apiserver-extra-env" json:"kube-apiserver-extra-env"`
	KubeSchedulerExtraEnv             ArrayArgument   `yaml:"kube-scheduler-extra-env" json:"kube-scheduler-extra-env"`
	KubeControllerManagerExtraEnv     ArrayArgument   `yaml:"kube-controller-manager-extra-env" json:"kube-controller-manager-extra-env"`
	KubeProxyExtraEnv                 ArrayArgument   `yaml:"kube-proxy-extra-env" json:"kube-proxy-extra-env"`
	EtcdExtraEnv                      ArrayArgument   `yaml:"etcd-extra-env" json:"etcd-extra-env"`
	CloudControllerManagerExtraEnv    ArrayArgument   `yaml:"cloud-controller-manager-extra-env" json:"cloud-controller-manager-extra-env"`
	Debug                             GenericArgument `yaml:"debug" json:"debug"`
	Docker                            BooleanArgument `yaml:"docker,omitempty" json:"docker,omitempty"`
	DisableApiserverLb                BooleanArgument `yaml:"disable-apiserver-lb,omitempty" json:"disable-apiserver-lb,omitempty"`
	FlannelCniConf                    GenericArgument `yaml:"flannel-cni-conf,omitempty" json:"flannel-cni-conf,omitempty"`
	FlannelConf                       GenericArgument `yaml:"flannel-conf,omitempty" json:"flannel-conf,omitempty"`
	FlannelIface                      GenericArgument `yaml:"flannel-iface,omitempty" json:"flannel-iface,omitempty"`
	ImageServiceEndpoint              GenericArgument `yaml:"image-service-endpoint,omitempty" json:"image-service-endpoint,omitempty"`
	NodeExternalDNS                   ArrayArgument   `yaml:"node-external-dns,omitempty" json:"node-external-dns,omitempty"`
	NodeInternalDNS                   ArrayArgument   `yaml:"node-internal-dns,omitempty" json:"node-internal-dns,omitempty"`
	PauseImage                        GenericArgument `yaml:"pause-image,omitempty" json:"pause-image,omitempty"`
	PreferBundledBin                  BooleanArgument `yaml:"prefer-bundled-bin,omitempty" json:"prefer-bundled-bin,omitempty"`
	DefaultRuntime                    GenericArgument `yaml:"default-runtime,omitempty" json:"default-runtime,omitempty"`
	DisableDefaultRegistryEndpoint    GenericArgument `yaml:"disable-default-registry-endpoint,omitempty" json:"disable-default-registry-endpoint,omitempty"`
	EnablePprof                       GenericArgument `yaml:"enable-pprof,omitempty" json:"enable-pprof,omitempty"`
	BindAddress                       GenericArgument `yaml:"bind-address,omitempty" json:"bind-address,omitempty"`
	ContainerRuntimeEndpoint          GenericArgument `yaml:"container-runtime-endpoint,omitempty" json:"container-runtime-endpoint,omitempty"`
	KubeletPath                       GenericArgument `yaml:"kubelet-path,omitempty" json:"kubelet-path,omitempty"`
	CloudControllerManagerImage       GenericArgument `yaml:"cloud-controller-manager-image,omitempty" json:"cloud-controller-manager-image,omitempty"`
	LbServerPort                      GenericArgument `yaml:"lb-server-port,omitempty" json:"lb-server-port,omitempty"`
	KubeProxyImage                    GenericArgument `yaml:"kube-proxy-image,omitempty" json:"kube-proxy-image,omitempty"`
	PodSecurityAdmissionConfigFile    GenericArgument `yaml:"pod-security-admission-config-file,omitempty" json:"pod-security-admission-config-file,omitempty"`
	ControlPlaneProbeConfiguration    GenericArgument `yaml:"control-plane-probe-configuration,omitempty" json:"control-plane-probe-configuration,omitempty"`
	NonrootDevices                    GenericArgument `yaml:"nonroot-devices,omitempty" json:"nonroot-devices,omitempty"`
	NodeNameFromCloudProviderMetadata GenericArgument `yaml:"node-name-from-cloud-provider-metadata,omitempty" json:"node-name-from-cloud-provider-metadata,omitempty"`
	Snapshotter                       GenericArgument `yaml:"snapshotter,omitempty" json:"snapshotter,omitempty"`
	VPNAuth                           GenericArgument `yaml:"vpn-auth,omitempty" json:"vpn-auth,omitempty"`
	VPNAuthFile                       GenericArgument `yaml:"vpn-auth-file,omitempty" json:"vpn-auth-file,omitempty"`
}

// Chart represents a Helm chart with repo and version
type Chart struct {
	Repo    *string `yaml:"repo" json:"repo"`
	Version *string `yaml:"version" json:"version"`
}

// Charts represents all available Helm charts keyed by chart name.
type Charts map[string]*Chart

// FeatureVersions represents feature version information
type FeatureVersions struct {
	EncryptionKeyRotation *string `yaml:"encryption-key-rotation" json:"encryption-key-rotation"`
}

// Message represents a release message.
type Message struct {
	ID       *string `yaml:"id" json:"id"`
	Type     *string `yaml:"type" json:"type"`
	Severity *string `yaml:"severity,omitempty" json:"severity,omitempty"`
	Summary  *string `yaml:"summary" json:"summary"`
	Message  *string `yaml:"message" json:"message"`
}

// Releases represents a specific RKE2 release with its configuration
type Releases struct {
	Version                 string           `yaml:"version" json:"version"`
	MinChannelServerVersion string           `yaml:"minChannelServerVersion" json:"minChannelServerVersion"`
	MaxChannelServerVersion string           `yaml:"maxChannelServerVersion" json:"maxChannelServerVersion"`
	ServerArgs              *ServerArgs      `yaml:"serverArgs,omitempty" json:"serverArgs,omitempty"`
	AgentArgs               *AgentArgs       `yaml:"agentArgs,omitempty" json:"agentArgs,omitempty"`
	Charts                  Charts           `yaml:"charts,omitempty" json:"charts,omitempty"`
	FeatureVersions         *FeatureVersions `yaml:"featureVersions,omitempty" json:"featureVersions,omitempty"`
	Messages                []Message        `yaml:"messages,omitempty" json:"messages,omitempty"`
}
