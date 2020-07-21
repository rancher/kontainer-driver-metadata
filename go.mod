module github.com/rancher/kontainer-driver-metadata

go 1.13

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	k8s.io/client-go => k8s.io/client-go v0.18.4
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/rke v1.2.0-rc2.0.20200721222223-d1fd91f27af2
	github.com/sirupsen/logrus v1.4.2
	sigs.k8s.io/yaml v1.2.0
)
