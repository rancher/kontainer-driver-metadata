module github.com/rancher/kontainer-driver-metadata

go 1.13

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	github.com/rancher/types => github.com/StrongMonkey/types-1 v0.0.0-20200212172741-c1a01949463e
	k8s.io/client-go => k8s.io/client-go v0.17.2
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/types v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.4.2
)
