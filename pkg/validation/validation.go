package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/exp/slices"

	utiliies "github.com/rancher/kontainer-driver-metadata/pkg"
	"github.com/rancher/kontainer-driver-metadata/pkg/data"
	"github.com/rancher/kontainer-driver-metadata/pkg/images"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const (
	rancherChart    = "https://charts.rancher.io"
	oldRancherChart = "https://github.com/rancher/charts"
	rke2Chart       = "https://rke2-charts.rancher.io"
)

var (
	releaseDataURL    = "https://releases.rancher.com/kontainer-driver-metadata/%s/data.json"
	releaseRegSyncURL = "https://raw.githubusercontent.com/rancher/kontainer-driver-metadata/%s/regsync.yaml"
	versionsToSkip    = map[string]bool{
		"v1.30.12+rke2r1": true,
		"v1.31.8+rke2r1":  true,
		"v1.32.4+rke2r1":  true,
		"v1.30.12+k3s1":   true,
		"v1.31.8+k3s1":    true,
		"v1.32.4+k3s1":    true,
		"v1.33.0+rke2r1":  true,
		"v1.33.0+k3s1":    true,
	}

	validMessageTypes = []string{
		"info",
		"warning",
		"critical",
	}
	validMessageSeverity = []string{
		"low",
		"medium",
		"high",
	}
	validMessageFields = map[string]struct{}{
		"id":       {},
		"type":     {},
		"severity": {},
		"summary":  {},
		"message":  {},
		"url":      {},
	}
)

// imageTags holds images and their tags as nested maps to make the comparison easy
type imageTags map[string]map[string]bool

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
			logrus.Fatalf("failed to validte the KDM data for the release [%s]: %v", release, err)
		}
		if err := validateRegSync(release); err != nil {
			logrus.Fatalf("failed to validte the regsync file for the release [%s]: %v", release, err)
		}
	}
	logrus.Info("validation is passed")
}

func validateRegSync(release string) error {
	raw, err := utiliies.DownloadFromURL(fmt.Sprintf(releaseRegSyncURL, release))
	if err != nil {
		return fmt.Errorf("failed to download the upstream regsync file: %v", err)
	}
	upstream, err := getImageTags([]byte(raw))
	if err != nil {
		return fmt.Errorf("failed to extract images and tags from the upstream: %v", err)
	}
	file, err := os.ReadFile(images.RegSyncFilePath)
	if err != nil {
		return err
	}
	local, err := getImageTags(file)
	if err != nil {
		return fmt.Errorf("failed to extract images and tags from the local: %v", err)
	}
	// RKE2 and K3s releases may need to be fixed after the fact,
	// so we just make sure we don't remove any released image or tag
	for name, tags := range upstream {
		// RKE entries removed; validate-ci fails as it compares against release-v2.11 which still includes RKE.
		// This will be removed once release-v2.12 is branched and validation is updated accordingly.
		excluded := map[string]bool{
			"docker.io/rancher/flannel-cni":                                      true,
			"docker.io/rancher/mirrored-k8s-dns-node-cache":                      true,
			"docker.io/rancher/mirrored-flannelcni-flannel":                      true,
			"docker.io/rancher/mirrored-coredns-coredns":                         true,
			"docker.io/rancher/hyperkube":                                        true,
			"docker.io/rancher/mirrored-calico-pod2daemon-flexvol":               true,
			"docker.io/rancher/mirrored-calico-kube-controllers":                 true,
			"docker.io/rancher/mirrored-k8s-dns-dnsmasq-nanny":                   true,
			"docker.io/rancher/mirrored-coreos-etcd":                             true,
			"docker.io/rancher/mirrored-calico-node":                             true,
			"docker.io/rancher/mirrored-calico-cni":                              true,
			"docker.io/rancher/mirrored-calico-typha":                            true,
			"docker.io/rancher/mirrored-calico-kube-proxy":                       true,
			"docker.io/rancher/mirrored-calico-ctl":                              true,
			"docker.io/rancher/mirrored-k8s-dns-kube-dns":                        true,
			"docker.io/rancher/mirrored-metrics-server":                          true,
			"docker.io/rancher/mirrored-coreos-flannel":                          true,
			"docker.io/rancher/calico-cni":                                       true,
			"docker.io/rancher/mirrored-flannel-flannel":                         true,
			"docker.io/rancher/rke-tools":                                        true,
			"docker.io/rancher/nginx-ingress-controller":                         true,
			"docker.io/rancher/mirrored-cluster-proportional-autoscaler":         true,
			"docker.io/rancher/mirrored-k8s-dns-sidecar":                         true,
			"docker.io/rancher/mirrored-nginx-ingress-controller-defaultbackend": true,
			"docker.io/rancher/mirrored-ingress-nginx-kube-webhook-certgen":      true,
		}

		if excluded[name] {
			continue
		}

		localTags, found := local[name]
		if !found {
			return fmt.Errorf("a released image [%s] is missing in the dev regSync file", name)
		}
		for tag := range tags {
			if !localTags[tag] {
				return fmt.Errorf("a released tag [%s:%s] is missing in the dev regSync file", name, tag)
			}
		}
	}
	return nil
}

func getImageTags(source []byte) (imageTags, error) {
	var upstream map[string]interface{}
	if err := yaml.Unmarshal(source, &upstream); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}
	sync, _, err := unstructured.NestedSlice(upstream, "sync")
	if err != nil {
		return nil, err
	}
	upstreamImageTag := imageTags{}
	for _, item := range sync {
		source, _, err := unstructured.NestedString(item.(map[string]interface{}), "source")
		if err != nil {
			return nil, err
		}
		allowTags, _, err := unstructured.NestedSlice(item.(map[string]interface{}), "tags", "allow")
		if err != nil {
			return nil, err
		}
		tags := map[string]bool{}
		for _, tag := range allowTags {
			t, _ := tag.(string)
			tags[t] = true
		}
		upstreamImageTag[source] = tags
	}
	return upstreamImageTag, nil
}

// validate checks the versions in the local data.json by comparing with the released data.json,
// Supported releases are RKE, RKE2 and K3s.
func validate(dev, released data.Data) error {
	for _, distro := range []string{utiliies.RKE2, utiliies.K3S} {
		if err := validateDistro(distro, dev, released); err != nil {
			return fmt.Errorf("failed to validate the distro [%s]: %v", distro, err)
		}
	}
	return nil
}

func validateDistro(distro string, dev, released data.Data) error {
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
			if err := validateEncryptedKeyRotation(release); err != nil {
				return fmt.Errorf("failed to validate rke2 encrypted key rotation: %v", err)
			}
			if err := validateMessages(release); err != nil {
				return fmt.Errorf("failed to validate messages: %v", err)
			}
		}
	}

	if distro == utiliies.K3S {
		raw, _, err := unstructured.NestedSlice(dev.K3S, "releases")
		if err != nil {
			return err
		}
		for _, r := range raw {
			release, ok := r.(map[string]interface{})
			if !ok {
				return errors.New("failed to parse map")
			}
			if err := validateEncryptedKeyRotation(release); err != nil {
				return fmt.Errorf("failed to validate k3s encrypted key rotation: %w", err)
			}
			if err := validateMessages(release); err != nil {
				return fmt.Errorf("failed to validate messages: %v", err)
			}
		}
	}
	return nil
}

func validateEncryptedKeyRotation(release map[string]interface{}) error {
	version, _, err := unstructured.NestedString(release, "version")
	if err != nil {
		return err
	}

	// this is the first version that hasn't reached its end of life that requires
	// the encrypted-key-rotation key to exist when this validation is being written
	const firstVersionToCheckEncryptedKeyRotation = "v1.25.11"
	compareVersions := semver.Compare(firstVersionToCheckEncryptedKeyRotation, version)
	if compareVersions > 0 || versionsToSkip[version] {
		return nil
	}
	logrus.Info("validating encrypted key rotation key on version: " + version)

	featureVersions, foundFeatureVersions, err := unstructured.NestedMap(release, "featureVersions")
	if err != nil {
		return err
	}
	if !foundFeatureVersions {
		return errors.New("missing featureVersions on version: " + version)
	}

	_, foundEncryptionKeyRotation := featureVersions["encryption-key-rotation"]

	if !foundEncryptionKeyRotation {
		return errors.New("missing encryption-key-rotation on version: " + version)
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
		if chartVersion == "0.0.0" {
			continue
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
		var isValidRepo bool
		switch repo {
		case "rancher-charts":
			isValidRepo = strings.HasPrefix(chartURL, rancherChart) || strings.HasPrefix(chartURL, oldRancherChart)
		case "rancher-rke2-charts":
			isValidRepo = strings.HasPrefix(chartURL, rke2Chart)
		default:
			isValidRepo = strings.HasPrefix(chartURL, "https://"+repo)
		}
		expectedChartTarball := fmt.Sprintf("%s-%s.tgz", chartName, chartVersion)
		if !strings.Contains(chartURL, expectedChartTarball) || !isValidRepo {
			return fmt.Errorf("unexpected chart URL for %s/%s:%s: %s", repo, chartName, chartVersion, chartURL)
		}
	}
	return nil
}

// getVersions returns the versions found from the dev and released data, and an error if anything goes wrong
func getVersions(distro string, dev, released data.Data) (devVersions, releasedVersions []string, err error) {
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

// validateMessages validates the messages field in a release if present
func validateMessages(release map[string]interface{}) error {
	version, _, err := unstructured.NestedString(release, "version")
	if err != nil {
		return err
	}

	messages, found, err := unstructured.NestedSlice(release, "messages")
	if err != nil {
		return fmt.Errorf("failed to extract messages from release %s: %v", version, err)
	}

	if !found {
		// messages field is optional
		return nil
	}

	logrus.Debugf("validating messages for version %s", version)

	for i, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid message format at index %d in version %s: expected map", i, version)
		}

		// Check for unexpected fields
		for field := range msgMap {
			if _, ok := validMessageFields[field]; !ok {
				return fmt.Errorf("message at index %d in version %s has unexpected field '%s'", i, version, field)
			}
		}

		if msgType, found, _ := unstructured.NestedString(msgMap, "type"); !found || msgType == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'type'", i, version)
		} else if !slices.Contains(validMessageTypes, msgType) {
			return fmt.Errorf("message at index %d in version %s has invalid type '%s': must be one of %v", i, version, msgType, validMessageTypes)
		}

		if severity, found, _ := unstructured.NestedString(msgMap, "severity"); found && !slices.Contains(validMessageSeverity, severity) {
			return fmt.Errorf("message at index %d in version %s has severity '%s': must be one of %v", i, version, severity, validMessageSeverity)
		}

		if summary, found, _ := unstructured.NestedString(msgMap, "summary"); !found || summary == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'summary'", i, version)
		}

		if message, found, _ := unstructured.NestedString(msgMap, "message"); !found || message == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'message'", i, version)
		}

		logrus.Debugf("message %d in version %s is valid", i, version)
	}

	return nil
}
