module github.com/rancher/kontainer-driver-metadata

go 1.13

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	k8s.io/client-go => k8s.io/client-go v0.18.4
	github.com/rancher/types => github.com/noironetworks/types v0.0.0-20200714015041-9795069dca9b
	github.com/rancher/rke => github.com/noironetworks/rke v0.3.1-0.20200715091719-953bab9a9116
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/rke v1.2.0-rc2.0.20200712062933-4c1d3db2b0c1
	github.com/sirupsen/logrus v1.4.2
	sigs.k8s.io/yaml v1.2.0
)
