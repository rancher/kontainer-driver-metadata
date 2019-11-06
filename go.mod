module github.com/rancher/kontainer-driver-metadata

go 1.12

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/types v0.0.0-20190930165650-6bbedae77a35
	github.com/sirupsen/logrus v1.4.2
)
