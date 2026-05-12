package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	utiliies "github.com/rancher/kontainer-driver-metadata/pkg"
	"github.com/rancher/kontainer-driver-metadata/pkg/data"
	"github.com/rancher/kontainer-driver-metadata/pkg/images"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
	yamlv3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

const (
	rancherChart    = "https://charts.rancher.io"
	oldRancherChart = "https://github.com/rancher/charts"
	rke2Chart       = "https://rke2-charts.rancher.io"
)

type MessageType string

const (
	MessageTypeInfo     MessageType = "info"
	MessageTypeWarning  MessageType = "warning"
	MessageTypeCritical MessageType = "critical"
)

type MessageSeverity string

const (
	MessageSeverityLow    MessageSeverity = "low"
	MessageSeverityMedium MessageSeverity = "medium"
	MessageSeverityHigh   MessageSeverity = "high"
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

	validMessageTypes = map[MessageType]struct{}{
		MessageTypeInfo:     {},
		MessageTypeWarning:  {},
		MessageTypeCritical: {},
	}
	validMessageSeverity = map[MessageSeverity]struct{}{
		MessageSeverityLow:    {},
		MessageSeverityMedium: {},
		MessageSeverityHigh:   {},
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

func loadReleases(source map[string]interface{}) (*RKEReleases, error) {
	raw, err := json.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal release data: %w", err)
	}

	var releases RKEReleases
	decoder := yamlv3.NewDecoder(bytes.NewReader(raw))
	decoder.KnownFields(true)
	if err := decoder.Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to unmarshal release data: %w", err)
	}

	return &releases, nil
}

// validate checks the versions in the local data.json by comparing with the released data.json,
// Supported releases are RKE, RKE2 and K3s.
func validate(dev, released data.Data) error {
	seenMessageIDs := make(map[string]bool)
	devRKE2, err := loadReleases(dev.RKE2)
	if err != nil {
		return fmt.Errorf("failed to load dev RKE2 releases: %v", err)
	}
	releasedRKE2, err := loadReleases(released.RKE2)
	if err != nil {
		return fmt.Errorf("failed to load released RKE2 releases: %v", err)
	}
	devK3S, err := loadReleases(dev.K3S)
	if err != nil {
		return fmt.Errorf("failed to load dev K3S releases: %v", err)
	}
	releasedK3S, err := loadReleases(released.K3S)
	if err != nil {
		return fmt.Errorf("failed to load released K3S releases: %v", err)
	}
	for _, distro := range []string{utiliies.RKE2, utiliies.K3S} {
		if err := validateDistro(distro, devRKE2, releasedRKE2, devK3S, releasedK3S, seenMessageIDs); err != nil {
			return fmt.Errorf("failed to validate the distro [%s]: %v", distro, err)
		}
	}
	return nil
}

func validateDistro(distro string, devRKE2, releasedRKE2, devK3S, releasedK3S *RKEReleases, seenMessageIDs map[string]bool) error {
	logrus.Infof("validating the distro [%s]", distro)
	var versionsInDev, versionsInRelease []string
	var err error
	switch distro {
	case utiliies.RKE2:
		versionsInDev, versionsInRelease, err = getVersions(devRKE2.Releases, releasedRKE2.Releases)
	case utiliies.K3S:
		versionsInDev, versionsInRelease, err = getVersions(devK3S.Releases, releasedK3S.Releases)
	}
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
		for _, release := range devRKE2.Releases {
			if err := validateRKE2Charts(release); err != nil {
				logrus.Infof("the release: %+v", release)
				return fmt.Errorf("failed to validate RKE2 charts: %v", err)
			}
			if err := validateEncryptedKeyRotation(release); err != nil {
				return fmt.Errorf("failed to validate rke2 encrypted key rotation: %v", err)
			}
			if err := validateMessages(release, seenMessageIDs); err != nil {
				return fmt.Errorf("failed to validate messages: %v", err)
			}
		}
	}

	if distro == utiliies.K3S {
		for _, release := range devK3S.Releases {
			if err := validateEncryptedKeyRotation(release); err != nil {
				return fmt.Errorf("failed to validate k3s encrypted key rotation: %w", err)
			}
			if err := validateMessages(release, seenMessageIDs); err != nil {
				return fmt.Errorf("failed to validate messages: %v", err)
			}
		}
	}
	return nil
}

func validateEncryptedKeyRotation(release Releases) error {
	version := release.Version
	// this is the first version that hasn't reached its end of life that requires
	// the encrypted-key-rotation key to exist when this validation is being written
	const firstVersionToCheckEncryptedKeyRotation = "v1.25.11"
	compareVersions := semver.Compare(firstVersionToCheckEncryptedKeyRotation, version)
	if compareVersions > 0 || versionsToSkip[version] {
		return nil
	}
	logrus.Info("validating encrypted key rotation key on version: " + version)

	if release.FeatureVersions == nil {
		return errors.New("missing featureVersions on version: " + version)
	}

	if release.FeatureVersions.EncryptionKeyRotation == nil {
		return errors.New("missing encryption-key-rotation on version: " + version)
	}

	return nil
}

func validateRKE2Charts(release Releases) error {
	if release.Charts == nil {
		return nil
	}
	rke2Version := release.Version
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
	for chartName, chart := range release.Charts {
		if chart == nil || chart.Repo == nil || chart.Version == nil {
			return fmt.Errorf("missing chart metadata for %s in release %s", chartName, rke2Version)
		}
		repo := *chart.Repo
		chartVersion := *chart.Version
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
func getVersions(dev, released []Releases) (devVersions, releasedVersions []string, err error) {
	for _, release := range dev {
		devVersions = append(devVersions, release.Version)
	}
	for _, release := range released {
		releasedVersions = append(releasedVersions, release.Version)
	}
	return devVersions, releasedVersions, nil
}

// validateMessages validates the messages field in a release if present
func validateMessages(release Releases, seenMessageIDs map[string]bool) error {
	version := release.Version
	if len(release.Messages) == 0 {
		// messages field is optional
		return nil
	}

	logrus.Debugf("validating messages for version %s", version)

	for i, msg := range release.Messages {
		if msg.ID == nil || *msg.ID == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'id'", i, version)
		}
		msgID := *msg.ID

		// Check for duplicate IDs globally
		if seenMessageIDs[msgID] {
			return fmt.Errorf("message at index %d in version %s has duplicate id '%s'", i, version, msgID)
		}
		seenMessageIDs[msgID] = true

		if msg.Type == nil || *msg.Type == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'type'", i, version)
		} else if _, validType := validMessageTypes[MessageType(*msg.Type)]; !validType {
			return fmt.Errorf("message at index %d in version %s has invalid type '%s': must be one of info, warning, critical", i, version, *msg.Type)
		}

		if msg.Severity != nil {
			if _, validSeverity := validMessageSeverity[MessageSeverity(*msg.Severity)]; !validSeverity {
				return fmt.Errorf("message at index %d in version %s has invalid severity '%s': must be one of low, medium, high", i, version, *msg.Severity)
			}
		}

		if msg.Summary == nil || *msg.Summary == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'summary'", i, version)
		}

		if msg.Message == nil || *msg.Message == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'message'", i, version)
		}

		logrus.Debugf("message %d in version %s is valid", i, version)
	}

	return nil
}
