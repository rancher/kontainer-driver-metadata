package rke

import (
	"sort"
	"testing"

	v3 "github.com/rancher/rke/types"
	"github.com/stretchr/testify/assert"
)

func TestLinuxAndWindowsKubernetesVersionsMatch(t *testing.T) {
	// arrange
	a := assert.New(t)

	// act
	windowsVersions := keys(loadK8sVersionWindowsServiceOptions())
	sort.Strings(windowsVersions)
	// need just the versions where naming is the same.
	tw := windowsVersions[2:]
	linuxVersions := keys(loadK8sVersionServiceOptions())
	sort.Strings(linuxVersions)
	// need to trim versions of Linux too to remove the last sort value and to keep length correct.
	in := len(linuxVersions) - len(windowsVersions)
	tl := linuxVersions[in+1 : len(linuxVersions)-1]

	// assert
	a.Equal(len(tl), len(tw))
	a.Equal(tl, tw)
	for _, v := range tl {
		a.Contains(tw, v)
	}
}

func keys(m map[string]v3.KubernetesServicesOptions) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
