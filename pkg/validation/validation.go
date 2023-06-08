package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	utiliies "github.com/rancher/kontainer-driver-metadata/pkg"
	"github.com/rancher/rke/types/kdm"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

var (
	releaseDataURL = "https://releases.rancher.com/kontainer-driver-metadata/%s/data.json"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		logrus.Fatal("Usage: go run validation.go <release> [ <release>...]")
	}

	dev, err := utiliies.FromLocalFile()
	if err != nil {
		logrus.Fatalf("failed to get the KDM data from the local file: %v", err)
	}

	for _, release := range args[1:] {
		logrus.Infof("validating [%s]", release)
		released, err := utiliies.FromURL(fmt.Sprintf(releaseDataURL, release))
		if err != nil {
			logrus.Fatalf("failed to get the KDM data for release [%s]: %v", release, err)
		}
		if err = validate(dev, released); err != nil {
			logrus.Fatalf("failed to validte the release [%s]: %v", release, err)
		}
	}
	logrus.Info("validation is passed")
	return
}

// validate checks the versions in the local data.json by comparing with the released data.json,
// Supported releases are RKE, RKE2 and K3s.
func validate(dev, released kdm.Data) error {
	for _, distro := range []string{utiliies.RKE, utiliies.RKE2, utiliies.K3S} {
		if err := validateDistro(distro, dev, released); err != nil {
			return fmt.Errorf("failed to validate the distro [%s]: %v", distro, err)
		}
	}
	return nil
}

func validateDistro(distro string, dev, released kdm.Data) error {
	logrus.Infof("validating the distro [%s]", distro)
	versionsInDev, versionsInRelease, err := getVersions(distro, dev, released)
	if err != nil {
		return fmt.Errorf("failed to get versions for [%s]: %v", distro, err)
	}
	if len(versionsInDev) < len(versionsInRelease) {
		return fmt.Errorf("the number of versions found in the dev is less than in the released")
	}
	dv := make(map[string]bool, len(versionsInDev))
	for _, v := range versionsInDev {
		dv[v] = true
	}
	for _, version := range versionsInRelease {
		// RKE2 and K3s releases may need to be fixed after the fact,
		// so we just make sure we don't remove a released version
		if !dv[version] {
			return fmt.Errorf("a released version [%s] is missing in the dev data", version)
		}
		if distro == utiliies.RKE {
			// RKE release images cannot be changed after the fact
			if !reflect.DeepEqual(released.K8sVersionRKESystemImages[version], dev.K8sVersionRKESystemImages[version]) {
				return fmt.Errorf("image(s) for a released version [%s] is changed in the dev data", version)
			}
		}
	}
	// check charts for RKE2 release
	if distro == utiliies.RKE2 {
		raw, _, err := unstructured.NestedSlice(dev.RKE2, "releases")
		if err != nil {
			return err
		}
		for _, r := range raw {
			release, ok := r.(map[string]interface{})
			if !ok {
				return err
			}
			if err := validateRKE2Charts(release); err != nil {
				logrus.Infof("the release: %v", release)
				return fmt.Errorf("failed to validate RKE2 charts: %v", err)
			}
		}
	}
	return nil
}

func validateRKE2Charts(release map[string]interface{}) error {
	rke2Version, _, err := unstructured.NestedString(release, "version")
	if err != nil {
		return err
	}
	charts, found, err := unstructured.NestedMap(release, "charts")
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	logrus.Infof("checking RKE2 %s chart metadata against rke2-runtime chart manifests", rke2Version)
	dir, err := os.MkdirTemp("", rke2Version)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	image := fmt.Sprintf("docker.io/rancher/rke2-runtime:%s", strings.ReplaceAll(rke2Version, "+", "-"))
	if err := exec.Command("wharfie", image, fmt.Sprintf("/charts:%s/charts", dir)).Run(); err != nil {
		logrus.Warnf("unable to extract rke2 runtime image %s; skipping chart validation. ", image)
		return nil
	}
	for chartName := range charts {
		repo, _, err := unstructured.NestedString(charts, chartName, "repo")
		if err != nil {
			return err
		}
		chartVersion, _, err := unstructured.NestedString(charts, chartName, "version")
		if err != nil {
			return err
		}
		logrus.Infof("checking RKE2 %s %s/%s:%s", rke2Version, repo, chartName, chartVersion)
		var info map[string]interface{}
		bytes, err := os.ReadFile(fmt.Sprintf("%s/charts/%s.yaml", dir, chartName))
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(bytes, &info); err != nil {
			return fmt.Errorf("failed to unmarshal the chart yaml: %v", err)
		}
		chartURL, _, err := unstructured.NestedString(info, "metadata", "annotations", "helm.cattle.io/chart-url")
		if err != nil {
			return err
		}
		expectedChartTarball := fmt.Sprintf("%s-%s.tgz", chartName, chartVersion)
		if !strings.Contains(chartURL, expectedChartTarball) {
			return fmt.Errorf("unexpected chart URL for %s/%s:%s: %s", repo, chartName, chartVersion, chartURL)
		}
	}
	return nil
}

// getVersions returns the versions found from the dev and released data, and an error if anything goes wrong
func getVersions(distro string, dev, released kdm.Data) (devVersions, releasedVersions []string, err error) {
	helper := func(source map[string]interface{}) ([]string, error) {
		var results []string
		versions, _, err := unstructured.NestedSlice(source, "releases")
		if err != nil {
			return nil, err
		}
		for _, version := range versions {
			version, ok := version.(map[string]interface{})
			if ok {
				results = append(results, version["version"].(string))
			}
		}
		return results, nil
	}

	switch distro {
	case utiliies.RKE:
		for version := range dev.K8sVersionRKESystemImages {
			devVersions = append(devVersions, version)
		}
		for version := range released.K8sVersionRKESystemImages {
			releasedVersions = append(releasedVersions, version)
		}
	case utiliies.RKE2:
		devVersions, err = helper(dev.RKE2)
		if err != nil {
			return devVersions, releasedVersions, err
		}
		releasedVersions, err = helper(released.RKE2)
		if err != nil {
			return devVersions, releasedVersions, err
		}
	case utiliies.K3S:
		devVersions, err = helper(dev.K3S)
		if err != nil {
			return devVersions, releasedVersions, err
		}
		releasedVersions, err = helper(released.K3S)
		if err != nil {
			return devVersions, releasedVersions, err
		}
	}
	return devVersions, releasedVersions, err
}
