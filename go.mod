module github.com/rancher/kontainer-driver-metadata

go 1.16

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	k8s.io/client-go => k8s.io/client-go v0.20.0
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/rke v1.2.4-0.20201118155135-4c1ad81cd691
	github.com/sirupsen/logrus v1.4.2
	sigs.k8s.io/yaml v1.2.0
)
