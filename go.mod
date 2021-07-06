module github.com/rancher/kontainer-driver-metadata

go 1.16

replace (
	github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
	github.com/matryer/moq => github.com/rancher/moq v0.0.0-20190404221404-ee5226d43009
	k8s.io/client-go => k8s.io/client-go v0.21.0
)

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/rancher/lasso v0.0.0-20210408231703-9ddd9378d08d // indirect
	github.com/rancher/rke v1.3.0-rc8.0.20210702150239-b8e94c98f8cd
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	sigs.k8s.io/yaml v1.2.0
)
