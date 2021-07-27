package rke

import (
	"testing"

	"github.com/blang/semver"
)

func TestLoadRKEDefaultK8sVersions(t *testing.T) {
	var latest semver.Version

	defaultVersion, _ := semver.Parse(loadRKEDefaultK8sVersions()["default"][1:])

	systemImages := loadK8sRKESystemImages()

	for k := range systemImages {
		v, _ := semver.Parse(k[1:])
		if v.GT(latest) {

			latest = v
		}
	}

	if !defaultVersion.EQ(latest) {
		t.Fail()
	}
}
