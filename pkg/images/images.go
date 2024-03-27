package images

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/blang/semver"
	utiliies "github.com/rancher/kontainer-driver-metadata/pkg"
	v3 "github.com/rancher/rke/types"
	image2 "github.com/rancher/rke/types/image"
	"github.com/rancher/rke/types/kdm"
	"github.com/rancher/rke/util"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	RegSyncFilePath = "./regsync.yaml"

	rke2LinuxImageURL   = "https://github.com/rancher/rke2/releases/download/%s/rke2-images-all.linux-amd64.txt"
	rke2WindowsImageURL = "https://github.com/rancher/rke2/releases/download/%s/rke2-images.windows-amd64.txt"
	k3sLinuxImageURL    = "https://github.com/k3s-io/k3s/releases/download/%s/k3s-images.txt"

	systemAgentInstallerImage = "rancher/system-agent-installer-%s:%s"
	upgradeImage              = "rancher/%s-upgrade:%s"
	releasesKey               = "releases"
	templateFilePath          = "./pkg/images/template.go.tmpl"

	linux  = "linux"
	window = "windows"
)

var (
	v270 = semver.MustParse("2.7.0")
	v290 = semver.MustParse("2.9.0")

	data     = kdm.Data{}
	releases = map[string][]map[string]interface{}{}
	urls     = map[string]map[string]string{}
)

// GenerateRegSyncFile generates the regsync.yaml file which contains all images used by
// the supported RKE, RKE2, and K3s k8s version in Rancher v2.8.x.
func GenerateRegSyncFile() {
	logrus.Info("generating the regsync.yaml file")
	if err := initData(); err != nil {
		logrus.Fatalf("failed to initilize data: %v", err)
	}
	var imageTag = map[string]map[string]bool{}
	for _, distro := range []string{utiliies.RKE, utiliies.RKE2, utiliies.K3S} {
		versions, err := filterVersions(distro)
		if err != nil {
			logrus.Fatalf("failed to filter versions for %s: %v", distro, err)
		}
		images, err := getImages(distro, versions)
		if err != nil {
			logrus.Fatalf("failed to get images for %s: %v", distro, err)
		}
		if err = unique(imageTag, images); err != nil {
			logrus.Fatalf("failed to passe: %v", err)
		}
	}

	temp, err := template.ParseFiles(templateFilePath)
	if err != nil {
		logrus.Fatalf("failed to parse the template file: %v ", err)
	}
	file, err := os.Create(RegSyncFilePath)
	if err != nil {
		logrus.Fatalf("failed to create the regsync file: %v ", err)
	}
	defer file.Close()

	if err = temp.Execute(file, imageTag); err != nil {
		logrus.Fatalf("failed to render the file: %v", err)
	}
	logrus.Info("finished generating regsync.yaml")
}

// unique extracts the image names and tags from the provided images slice into the provided map.
// For each (key,value) pair in the map, the key is the image name, the value is a map which keys are tags.
func unique(imageTag map[string]map[string]bool, images []string) error {
	for _, image := range images {
		if image == "" {
			continue
		}
		parts := strings.Split(image, ":")
		if len(parts) != 2 {
			return fmt.Errorf("failed to get image and tag from %s ", image)
		}
		name, tag := parts[0], parts[1]
		if !strings.HasPrefix(name, "rancher") {
			return fmt.Errorf("image name %s is not prefixed by rancher", image)
		}
		if tag == "" {
			return fmt.Errorf("tag is missing from %s ", image)
		}
		it, found := imageTag[name]
		if !found {
			imageTag[name] = map[string]bool{tag: true}
		} else {
			if !it[tag] {
				it[tag] = true
			}
		}
	}
	return nil
}

// getImages returns all images, for Linux and Windows, used by the provided distro and k8s versions.
// It returns an empty slice and an error if something goes wrong.
func getImages(distro string, versions []interface{}) (all []string, err error) {
	for _, version := range versions {
		switch distro {
		case utiliies.RKE:
			v, _ := version.(string)
			images, ok := data.K8sVersionRKESystemImages[v]
			if !ok {
				return nil, fmt.Errorf("failed to find the RKE system images for the version %s", v)
			}
			obj := map[string]interface{}{}

			bytes, err := json.Marshal(images)
			if err != nil {
				return nil, err
			}
			if err = json.Unmarshal(bytes, &obj); err != nil {
				return nil, err
			}
			for _, o := range obj {
				image, ok := o.(string)
				if !ok || image == "" {
					continue
				}
				converted := image2.Mirror(image)
				// skip images that we do not mirror and maintain
				if strings.HasPrefix(image, "weaveworks") || strings.HasPrefix(image, "noiro") {
					continue
				}
				// all images should be prefixed by "rancher"
				if !strings.HasPrefix(converted, "rancher") {
					return nil, fmt.Errorf("RKE system image %s does not start with rancher", converted)
				}
				logrus.Tracef("distro %s version %s adds %s ", distro, version, converted)
				all = append(all, converted)
			}
		case utiliies.RKE2, utiliies.K3S:
			temp, _ := version.(map[string]interface{})
			v, _ := temp["version"].(string)
			for _, os := range []string{linux, window} {
				url, _ := urls[distro][os]
				if url != "" {
					images, err := getImagesFromURL(fmt.Sprintf(url, v))
					if err != nil {
						return nil, fmt.Errorf("failed to get images: %o", err)
					}
					if len(images) > 0 {
						logrus.Tracef("distro %s version %s adds %s ", distro, version, images)
						all = append(all, images...)
					}
				}
			}
			// add the following images which are not in the upstream images.txt file
			// - rancher/rancher-<distro>-upgrade
			// - rancher/system-agent-installer-<distro>
			safeVersion := strings.ReplaceAll(v, "+", "-")
			upgradeImage := fmt.Sprintf(upgradeImage, strings.ToLower(distro), safeVersion)
			systemAgentInstallerImage := fmt.Sprintf(systemAgentInstallerImage, strings.ToLower(distro), safeVersion)
			all = append(all, upgradeImage, systemAgentInstallerImage)
		}
	}
	return all, nil
}

// getImagesFromURL returns the images that are downloaded and read from the provided url.
// It removes the prefix "docker.io/" from images' names.
// It returns an empty slice and an error when something goes wrong.
func getImagesFromURL(url string) ([]string, error) {
	if url == "" {
		return nil, nil
	}
	raw, err := utiliies.DownloadFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to dowlaod from the url %s: %v", url, err)
	}
	var results []string
	images := strings.Split(raw, "\n")
	for _, item := range images {
		if item != "" {
			item = strings.TrimPrefix(item, "docker.io/")
			results = append(results, item)
		}
	}
	return results, nil
}

// initData loads and generate relevant objects from the local data.json for further operations.
// It returns an error when something goes wrong.
func initData() (err error) {
	data, err = utiliies.FromLocalFile()
	if err != nil {
		logrus.Fatalf("failed to get the KDM data from the local file: %v", err)
	}

	getReleases := func(obj map[string]interface{}, fields ...string) ([]map[string]interface{}, error) {
		var releases []map[string]interface{}
		raw, _, err := unstructured.NestedSlice(obj, fields...)
		if err != nil {
			return nil, err
		}
		for _, r := range raw {
			release, ok := r.(map[string]interface{})
			if !ok {
				return nil, err
			}
			releases = append(releases, release)
		}
		return releases, nil
	}

	rke2, err := getReleases(data.RKE2, releasesKey)
	if err != nil {
		return fmt.Errorf("failed to initData RKE2 releases: %s", err)
	}
	releases[utiliies.RKE2] = rke2
	k3s, err := getReleases(data.K3S, releasesKey)
	if err != nil {
		return fmt.Errorf("failed to initData K3s releases: %s", err)
	}
	releases[utiliies.K3S] = k3s

	urls[utiliies.K3S] = map[string]string{linux: k3sLinuxImageURL}
	urls[utiliies.RKE2] = map[string]string{linux: rke2LinuxImageURL, window: rke2WindowsImageURL}
	return nil
}

// releaseToKeep returns a boolean to indicate if the provided RKE2/K3s release should be available in Rancher v2.8.x.
// it returns false and an error when something goes wrong.
func releaseToKeep(release map[string]interface{}) (bool, error) {
	version, ok := release["version"].(string)
	if !ok || version == "" {
		return false, fmt.Errorf("the value of the version is missing for the release %v", release)
	}
	minChannelServerVersion, ok := release["minChannelServerVersion"].(string)
	if ok && minChannelServerVersion != "" {
		if min, err := semver.ParseTolerant(minChannelServerVersion); err != nil || min.GE(v290) {
			return false, err
		}
	}
	maxChannelServerVersion, ok := release["maxChannelServerVersion"].(string)
	if ok && maxChannelServerVersion != "" {
		if max, err := semver.ParseTolerant(maxChannelServerVersion); err != nil || max.LT(v270) {
			return false, err
		}
	}
	return true, nil
}

// toKeep returns a boolean to indicate if the provided RKE k8s version should be available in Rancher v2.8.x.
// it returns false and an error when something goes wrong.
func toKeep(info v3.K8sVersionInfo) (bool, error) {
	if info.DeprecateRancherVersion != "" {
		drv, err := semver.ParseTolerant(info.DeprecateRancherVersion)
		if err != nil {
			return false, fmt.Errorf("failed to parse %s: %v", info.DeprecateRancherVersion, err)
		}
		if drv.LE(v270) {
			return false, nil
		}
	}
	if info.MinRancherVersion != "" {
		min, err := semver.ParseTolerant(info.MinRancherVersion)
		if err != nil {
			return false, fmt.Errorf("failed to parse %s: %v", info.MinRancherVersion, err)
		}
		if min.GE(v290) {
			return false, nil
		}
	}

	if info.MaxRancherVersion != "" {
		max, err := semver.ParseTolerant(info.MaxRancherVersion)
		if err != nil {
			return false, fmt.Errorf("failed to parse %s: %v", info.MaxRancherVersion, err)
		}
		if max.LT(v270) {
			return false, nil
		}
	}
	return true, nil
}

// filterVersions returns suitable versions of the provided distro that are available for Rancher 2.8.x.
// It returns an empty slice and an error when something goes wrong.
// In the case of RKE, the type of the returned object is slice of Strings;
// in the case of RKE2 or K3s, it is a slice of Maps (map[string]interface).
func filterVersions(distro string) (filteredVersions []interface{}, err error) {
	switch distro {
	case utiliies.RKE:
		for k8sVersion := range data.K8sVersionRKESystemImages {
			if info, ok := data.K8sVersionInfo[k8sVersion]; ok {
				keep, err := toKeep(info)
				if err != nil {
					return nil, err
				}
				if !keep {
					continue
				}
			}
			majorVersion := util.GetTagMajorVersion(k8sVersion)
			if info, ok := data.K8sVersionInfo[majorVersion]; ok {
				keep, err := toKeep(info)
				if err != nil {
					return nil, err
				}
				if !keep {
					continue
				}
			}
			logrus.Debugf("adding %v", k8sVersion)
			filteredVersions = append(filteredVersions, k8sVersion)
		}
	case utiliies.RKE2, utiliies.K3S:
		for _, r := range releases[distro] {
			if r == nil {
				continue
			}
			keep, err := releaseToKeep(r)
			if err != nil {
				return nil, err
			}
			if !keep {
				continue
			}
			logrus.Debugf("adding %v", r["version"].(string))
			filteredVersions = append(filteredVersions, r)
		}
	}
	return filteredVersions, err
}
