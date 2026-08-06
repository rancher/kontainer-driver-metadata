module github.com/rancher/kontainer-driver-metadata

go 1.25.0

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	k8s.io/api => k8s.io/api v0.35.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.35.3
	k8s.io/apiserver => k8s.io/apiserver v0.35.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.35.3
	k8s.io/client-go => k8s.io/client-go v0.35.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.35.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.35.3
	k8s.io/code-generator => k8s.io/code-generator v0.35.3
	k8s.io/component-base => k8s.io/component-base v0.35.3
	k8s.io/component-helpers => k8s.io/component-helpers v0.35.3
	k8s.io/controller-manager => k8s.io/controller-manager v0.35.3
	k8s.io/cri-api => k8s.io/cri-api v0.35.3
	k8s.io/cri-client => k8s.io/cri-client v0.35.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.35.3
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.35.3
	k8s.io/endpointslice => k8s.io/endpointslice v0.35.3
	k8s.io/externaljwt => k8s.io/externaljwt v0.35.3
	k8s.io/kms => k8s.io/kms v0.35.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.35.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.35.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.35.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.35.3
	k8s.io/kubectl => k8s.io/kubectl v0.35.3
	k8s.io/kubelet => k8s.io/kubelet v0.35.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.32.1
	k8s.io/metrics => k8s.io/metrics v0.35.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.35.3
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.35.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.35.3
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.35.3
	k8s.io/sample-controller => k8s.io/sample-controller v0.35.3
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/sirupsen/logrus v1.9.4
	golang.org/x/mod v0.38.0
	k8s.io/apimachinery v0.35.3
	sigs.k8s.io/yaml v1.6.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250910181357-589584f1c912 // indirect
	k8s.io/utils v0.0.0-20251002143259-bc988d571ff4 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v6 v6.3.0 // indirect
)
