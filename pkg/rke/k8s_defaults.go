package rke

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/blang/semver"
	"github.com/rancher/kontainer-driver-metadata/pkg/rke/templates"
	"github.com/rancher/rke/types/image"
	"github.com/rancher/rke/types/kdm"
	"github.com/sirupsen/logrus"
)

const (
	DataFilePath = "./data/data.json"
)

var (
	DriverData     kdm.Data
	MissedTemplate map[string][]string
	m              = image.Mirror
)

func initData() {
	DriverData = kdm.Data{
		K8sVersionRKESystemImages: loadK8sRKESystemImages(),
		K3S:                       map[string]interface{}{},
		RKE2:                      map[string]interface{}{},
	}

	for version, images := range DriverData.K8sVersionRKESystemImages {
		longName := "rancher/hyperkube:" + version
		if !strings.HasPrefix(longName, images.Kubernetes) {
			panic(fmt.Sprintf("For K8s version %s, the Kubernetes image tag should be a substring of %s, currently it is %s", version, version, images.Kubernetes))
		}
	}

	DriverData.RKEDefaultK8sVersions = loadRKEDefaultK8sVersions()
	DriverData.RancherDefaultK8sVersions = loadRancherDefaultK8sVersions()

	validateDefaultPresent(DriverData.RKEDefaultK8sVersions)

	DriverData.K8sVersionedTemplates = templates.LoadK8sVersionedTemplates()

	validateTemplateMatch()

	DriverData.K8sVersionServiceOptions = loadK8sVersionServiceOptions()

	DriverData.K8sVersionInfo = loadK8sVersionInfo()

	validateVersionInfo()

	// init Windows versions
	DriverData.K8sVersionWindowsServiceOptions = loadK8sVersionWindowsServiceOptions()
	DriverData.K8sVersionDockerInfo = loadK8sVersionDockerInfo()

	// CIS
	DriverData.CisConfigParams = loadCisConfigParams()
	DriverData.CisBenchmarkVersionInfo = loadCisBenchmarkVersionInfo()

	if err := readFile("./channels.yaml", DriverData.K3S); err != nil {
		panic(err)
	}
	if err := readFile("./channels-rke2.yaml", DriverData.RKE2); err != nil {
		panic(err)
	}
}

func validateVersionInfo() {
	var errorsFound bool
	var incompleteVersions []string
	versionRangesNeedSpecificVersionInfo := []string{
		// 1.15.12-rancher1-1 comes from 2.2.13, doesn't need version info
		">=1.15.11-rancher1-1 <1.15.12-rancher1-1",
		">=1.15.12-rancher2-2 <1.16.0-alpha",
		">=1.16.8-rancher1-1 <1.17.0-alpha",
		">=1.17.4-rancher1-1 <1.18.0-alpha"}
	for k8sVersion := range DriverData.K8sVersionRKESystemImages {
		toMatch, err := semver.Make(k8sVersion[1:])
		if err != nil {
			panic(fmt.Sprintf("k8sVersion not sem-ver %s %v", k8sVersion, err))
		}
		for _, versionRange := range versionRangesNeedSpecificVersionInfo {
			parsedVersionRange, err := semver.ParseRange(versionRange)
			if err != nil {
				panic(fmt.Sprintf("range not sem-ver %v %v", versionRange, err))
			}
			if parsedVersionRange(toMatch) {
				// check specific version info
				if _, ok := DriverData.K8sVersionInfo[k8sVersion]; !ok {
					incompleteVersions = append(incompleteVersions, k8sVersion)
					errorsFound = true
				}
			}
		}
	}
	if errorsFound {
		panic(fmt.Sprintf("following versions do not have specific version info specified: %v", strings.Join(incompleteVersions, ",")))
	}
}

func validateDefaultPresent(versions map[string]string) {
	for _, defaultK8s := range versions {
		if _, ok := DriverData.K8sVersionRKESystemImages[defaultK8s]; !ok {
			panic(fmt.Sprintf("Default K8s version %v is not found in the K8sVersionToRKESystemImages", defaultK8s))
		}
	}
}

func validateTemplateMatch() {
	MissedTemplate = map[string][]string{}
	for k8sVersion := range DriverData.K8sVersionRKESystemImages {
		toMatch, err := semver.Make(k8sVersion[1:])
		if err != nil {
			panic(fmt.Sprintf("k8sVersion not sem-ver %s %v", k8sVersion, err))
		}
		for plugin, pluginData := range DriverData.K8sVersionedTemplates {
			if plugin == kdm.TemplateKeys {
				continue
			}
			matchedKey := ""
			matchedRange := ""
			for toTestRange, key := range pluginData {
				testRange, err := semver.ParseRange(toTestRange)
				if err != nil {
					panic(fmt.Sprintf("range for %s not sem-ver %v %v", plugin, toTestRange, err))
				}
				if testRange(toMatch) {
					// only one range should be matched
					if matchedKey != "" {
						panic(fmt.Sprintf("k8sVersion %s for plugin %s passing range %s, conflict range matching with %s",
							k8sVersion, plugin, toTestRange, matchedRange))
					}
					matchedKey = key
					matchedRange = toTestRange
				}
			}

			// no template found
			if matchedKey == "" {
				// check if plugin was introduced later
				if templateRanges, ok := templates.TemplateIntroducedRanges[plugin]; ok {
					// as we want to use the logic outside this loop, we check every range and if its matched, we set pluginCheck to true
					// in the end, we check if any of the ranges have matched, if so, we dont skip the logic to add a missing template (because every version matched in the range should have a template)
					var pluginCheck bool
					// plugin has ranges configured
					for _, toTestRange := range templateRanges {
						testRange, err := semver.ParseRange(toTestRange)
						if err != nil {
							panic(fmt.Sprintf("range for %s not sem-ver %v %v", plugin, testRange, err))
						}
						if testRange(toMatch) {
							pluginCheck = true
						}
					}
					if !pluginCheck {
						// logrus.Warnf("skipping %s for %s", k8sVersion, plugin)
						continue
					}

				}

				// if version not already found for that plugin, append it, else create it
				if val, ok := MissedTemplate[plugin]; ok {
					val = append(val, k8sVersion)
					MissedTemplate[plugin] = val
				} else {
					MissedTemplate[plugin] = []string{k8sVersion}
				}
				continue
			}
		}
	}
}

func readFile(input string, data map[string]interface{}) error {
	bytes, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, &data)
}

func GenerateData() {
	initData()

	if len(MissedTemplate) != 0 {
		logrus.Warnf("found k8s versions without a template")
		for plugin, data := range MissedTemplate {
			logrus.Warnf("no %s template for k8sVersions %v \n", plugin, data)
		}
	}

	// todo: zip file
	strData, _ := json.MarshalIndent(DriverData, "", " ")
	jsonFile, err := os.Create(DataFilePath)
	if err != nil {
		panic(fmt.Errorf("err creating data file %v", err))
	}
	defer jsonFile.Close()
	_, err = jsonFile.Write(strData)
	if err != nil {
		panic(fmt.Errorf("err writing jsonFile %v", err))
	}
	fmt.Println("finished generating data.json")
}
