package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	utilities "github.com/rancher/kontainer-driver-metadata/pkg"
	"github.com/rancher/kontainer-driver-metadata/pkg/data"
	"github.com/sirupsen/logrus"
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
	releaseDataURL = "https://releases.rancher.com/kontainer-driver-metadata/%s/data.json"

	validMessageTypes      = []MessageType{MessageTypeInfo, MessageTypeWarning, MessageTypeCritical}
	validMessageSeverities = []MessageSeverity{MessageSeverityLow, MessageSeverityMedium, MessageSeverityHigh}
)

func Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: go run ./cmd/validation <release> [ <release>...]")
	}

	dev, err := utilities.FromLocalFile()
	if err != nil {
		return fmt.Errorf("failed to get the KDM data from the local file: %w", err)
	}

	for _, release := range args[1:] {
		logrus.Infof("validating [%s]", release)
		released, err := utilities.FromURL(fmt.Sprintf(releaseDataURL, release))
		if err != nil {
			return fmt.Errorf("failed to get the KDM data for release [%s]: %w", release, err)
		}
		if err = validate(dev, released); err != nil {
			return fmt.Errorf("failed to validate the KDM data for the release [%s]: %w", release, err)
		}
	}
	logrus.Info("validation is passed")
	return nil
}

func loadReleases(source map[string]any) (*RKEReleases, error) {
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

type releaseValidator func(Releases) error

type distroValidationTarget struct {
	name       string
	dev        []Releases
	released   []Releases
	validators []releaseValidator
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
	targets := []distroValidationTarget{
		{
			name:       utilities.RKE2,
			dev:        devRKE2.Releases,
			released:   releasedRKE2.Releases,
			validators: []releaseValidator{validateRKE2Charts},
		},
		{
			name:     utilities.K3S,
			dev:      devK3S.Releases,
			released: releasedK3S.Releases,
		},
	}
	for _, target := range targets {
		if err := validateDistro(target, seenMessageIDs); err != nil {
			return fmt.Errorf("failed to validate the distro [%s]: %v", target.name, err)
		}
	}
	return nil
}

func validateDistro(target distroValidationTarget, seenMessageIDs map[string]bool) error {
	logrus.Infof("validating the distro [%s]", target.name)
	versionsInDev, versionsInRelease, err := getVersions(target.dev, target.released)
	if err != nil {
		return fmt.Errorf("failed to get versions for [%s]: %v", target.name, err)
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

	for _, release := range target.dev {
		for _, validator := range target.validators {
			if err := validator(release); err != nil {
				return fmt.Errorf("failed to run distro validator for release [%s]: %v", release.Version, err)
			}
		}
		if err := validateMessages(release, seenMessageIDs); err != nil {
			return fmt.Errorf("failed to validate messages: %v", err)
		}
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
		var info map[string]any
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
		if msg.ID == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'id'", i, version)
		}

		// Check for duplicate IDs globally
		if seenMessageIDs[msg.ID] {
			return fmt.Errorf("message at index %d in version %s has duplicate id '%s'", i, version, msg.ID)
		}
		seenMessageIDs[msg.ID] = true

		if msg.Type == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'type'", i, version)
		} else if !slices.Contains(validMessageTypes, msg.Type) {
			return fmt.Errorf("message at index %d in version %s has invalid type '%s': must be one of info, warning, critical", i, version, msg.Type)
		}

		if msg.Severity != "" && !slices.Contains(validMessageSeverities, msg.Severity) {
			return fmt.Errorf("message at index %d in version %s has invalid severity '%s': must be one of low, medium, high", i, version, msg.Severity)
		}

		if msg.Summary == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'summary'", i, version)
		}

		if msg.Message == "" {
			return fmt.Errorf("message at index %d in version %s is missing required field 'message'", i, version)
		}

		logrus.Debugf("message %d in version %s is valid", i, version)
	}

	return nil
}
