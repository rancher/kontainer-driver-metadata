package rke

import v3 "github.com/rancher/rke/types"

func loadRancherDefaultK8sVersions() map[string]string {
	/*
		Just mention the major version, the latest minor version will be
		automatically picked based on Rancher's max/min version information.
	*/
	return map[string]string{
		"2.3.0": "v1.15.x",
		"2.3.1": "v1.15.x",
		"2.3.2": "v1.15.x",
		"2.3.3": "v1.16.x",
		"2.3":   "v1.17.x",
		"2.4.0": "v1.17.x",
		"2.4.1": "v1.17.x",
		"2.4.2": "v1.17.x",
		"2.4.3": "v1.17.x",
		"2.4.4": "v1.17.x",
		"2.4.5": "v1.18.x",
		"2.4":   "v1.17.x",
		// rancher will use default if its version is absent
		"default": "v1.18.x",
	}
}

func loadRKEDefaultK8sVersions() map[string]string {
	return map[string]string{
		"0.3": "v1.16.3-rancher1-1",
		// rke will use default if its version is absent
		"default": "v1.18.3-rancher2-3",
	}
}

/*
MaxRancherVersion: Last Rancher version having this k8s in k8sVersionsCurrent
DeprecateRancherVersion: No create/update allowed for RKE >= DeprecateRancherVersion
MaxRKEVersion: Last RKE version having this k8s in k8sVersionsCurrent
DeprecateRKEVersion: No create/update allowed for RKE >= DeprecateRKEVersion
*/

func loadK8sVersionInfo() map[string]v3.K8sVersionInfo {
	return map[string]v3.K8sVersionInfo{
		"v1.8": {
			MaxRancherVersion: "2.2",
			MaxRKEVersion:     "0.2.2",
		},
		"v1.9": {
			MaxRancherVersion: "2.2",
			MaxRKEVersion:     "0.2.2",
		},
		"v1.10": {
			MaxRancherVersion: "2.2",
			MaxRKEVersion:     "0.2.2",
		},
		"v1.11": {
			MaxRancherVersion: "2.2",
			MaxRKEVersion:     "0.2.2",
		},
		"v1.12": {
			MaxRancherVersion: "2.2",
			MaxRKEVersion:     "0.2.2",
		},
		"v1.13": {
			MaxRancherVersion: "2.3.1",
			MaxRKEVersion:     "0.3.1",
		},
		"v1.14": {
			MaxRancherVersion: "2.3.3",
			MaxRKEVersion:     "1.0.0",
		},
		"v1.15.5-rancher1-1": {
			MaxRancherVersion: "2.2.9",
			MaxRKEVersion:     "0.2.8",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.15.11-rancher1-1": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.15.11-rancher1-2": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.15.11-rancher1-3": {
			MinRancherVersion: "2.3.7",
			MinRKEVersion:     "1.0.7",
		},
		"v1.15.12-rancher1-1": {
			MaxRancherVersion: "2.2.99",
			MaxRKEVersion:     "0.2.99",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		// It also includes ingress-nginx 0.32.0
		"v1.15.12-rancher2-2": {
			MinRancherVersion: "2.4.4-rc0",
			MinRKEVersion:     "1.1.2-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		// It also includes ingress-nginx 0.32.0
		"v1.15.12-rancher2-3": {
			MinRancherVersion: "2.4.5-rc0",
			MinRKEVersion:     "1.1.3-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		// It also includes ingress-nginx 0.32.0
		"v1.15.12-rancher2-4": {
			MinRancherVersion: "2.4.6-rc0",
			MinRKEVersion:     "1.1.4-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.16.8-rancher1-1": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.16.8-rancher1-2": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.16.9-rancher1-1": {
			MinRancherVersion: "2.3.7-rc0",
			MinRKEVersion:     "1.0.7",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		// It also includes ingress-nginx 0.32.0
		"v1.16.10-rancher2-1": {
			MinRancherVersion: "2.4.4-rc0",
			MinRKEVersion:     "1.1.2-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.16.10-rancher2-2": {
			MinRancherVersion: "2.4.5-rc0",
			MinRKEVersion:     "1.1.3-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.16.10-rancher2-3": {
			MinRancherVersion: "2.4.6-rc0",
			MinRKEVersion:     "1.1.4-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.17.4-rancher1-1": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		"v1.17.4-rancher1-2": {
			MinRancherVersion: "2.3.3",
			MinRKEVersion:     "1.0.0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.17.5-rancher1-1": {
			MinRancherVersion: "2.3.7-rc0",
			MinRKEVersion:     "1.0.7",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		// It also includes ingress-nginx 0.32.0
		"v1.17.6-rancher2-1": {
			MinRancherVersion: "2.4.4-rc0",
			MinRKEVersion:     "1.1.2-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.17.6-rancher2-2": {
			MinRancherVersion: "2.4.5-rc0",
			MinRKEVersion:     "1.1.3-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.17.6-rancher2-3": {
			MinRancherVersion: "2.4.6-rc0",
			MinRKEVersion:     "1.1.4-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.18.3-rancher2-1": {
			MinRancherVersion: "2.4.5-rc0",
			MinRKEVersion:     "1.1.3-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.18.3-rancher2-2": {
			MinRancherVersion: "2.4.5-rc0",
			MinRKEVersion:     "1.1.3-rc0",
		},
		// The Calico/Canal template in this version use functions that are only available in RKE v1.0.0 and up
		// This version includes nodelocal dns only available in RKE v1.0.7 and up
		"v1.18.3-rancher2-3": {
			MinRancherVersion: "2.4.6-rc0",
			MinRKEVersion:     "1.1.4-rc0",
		},
		"v1.8.10-rancher1-1": {
			DeprecateRKEVersion:     "0.2.2",
			DeprecateRancherVersion: "2.2",
		},
		"v1.8.11-rancher1": {
			DeprecateRKEVersion:     "0.2.2",
			DeprecateRancherVersion: "2.2",
		},
		"v1.9.7-rancher1": {
			DeprecateRKEVersion:     "0.2.2",
			DeprecateRancherVersion: "2.2",
		},
		"v1.10.1-rancher1": {
			DeprecateRKEVersion:     "0.2.2",
			DeprecateRancherVersion: "2.2",
		},
	}
}
