package rke

import (
	"github.com/rancher/rke/types/kdm"
)

const (
	// reasons for not applicable checks
	reasonForMalformed          = `The argument --repair-malformed-updates has been removed as of Kubernetes version 1.14.`
	reasonNoConfigFileAPIServer = `Clusters provisioned by RKE doesn't require or maintain a configuration file for kube-apiserver.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileControllerManager = `Clusters provisioned by RKE doesn't require or maintain a configuration file for controller-manager.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileEtcd = `Clusters provisioned by RKE doesn't require or maintain a configuration file for etcd.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileKubelet = `Clusters provisioned by RKE doesn’t require or maintain a configuration file for the kubelet.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileKubeletSvc = `Clusters provisioned by RKE doesn’t require or maintain a configuration file for the kubelet service.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileScheduler = `Clusters provisioned by RKE doesn't require or maintain a configuration file for scheduler.
All configuration is passed in as arguments at container run time.`
	reasonForNoHostnameOverride = `Clusters provisioned by RKE clusters and most cloud providers require hostnames.`
	reasonNoKubeConfigDefault   = `Clusters provisioned by RKE does not store the kubernetes default kubeconfig credentials file on the nodes.`
	reasonForNotRotatingCerts   = `Clusters provisioned by RKE handles certificate rotation directly through RKE.`

	// reasons for skipped checks
	reasonForAlwaysPullImages                    = `Enabling AlwaysPullImages can use significant bandwidth.`
	reasonForDefaultSA                           = `Kubernetes provides default service accounts to be used.`
	reasonForDefaultNS                           = `Kubernetes provides a default namespace.`
	reasonForEtcdDataDir                         = `A system service account is required for etcd data directory ownership. Refer to Rancher's hardening guide for more details on how to configure this ownership.`
	reasonForEncryption                          = `Enabling encryption changes how data can be recovered as data is encrypted.`
	reasonForEventRateLimit                      = `EventRateLimit needs to be tuned depending on the cluster.`
	reasonForKubeletServerCerts                  = `When generating serving certificates, functionality could break in conjunction with hostname overrides which are required for certain cloud providers.`
	reasonForLocalhostListeningControllerManager = `Adding this argument prevents Rancher's monitoring tool to collect metrics on the controller manager.`
	reasonForLocalhostListeningScheduler         = `Adding this argument prevents Rancher's monitoring tool to collect metrics on the scheduler.`
	reasonForNetPol                              = `Enabling Network Policies can prevent certain applications from communicating with each other.`
	reasonForProtectKernelDefaults               = `System level configurations are required prior to provisioning the cluster in order for this argument to be set to true.`
	reasonForPSP                                 = `Enabling Pod Security Policy can cause applications to unexpectedly fail.`
)

var rkeCIS14NotApplicableChecks = map[string]string{
	"1.1.9":  reasonForMalformed,
	"1.3.6":  reasonForNotRotatingCerts,
	"1.4.1":  reasonNoConfigFileAPIServer,
	"1.4.2":  reasonNoConfigFileAPIServer,
	"1.4.3":  reasonNoConfigFileControllerManager,
	"1.4.4":  reasonNoConfigFileControllerManager,
	"1.4.5":  reasonNoConfigFileScheduler,
	"1.4.6":  reasonNoConfigFileScheduler,
	"1.4.7":  reasonNoConfigFileEtcd,
	"1.4.8":  reasonNoConfigFileEtcd,
	"1.4.13": reasonNoKubeConfigDefault,
	"1.4.14": reasonNoKubeConfigDefault,
	"2.1.8":  reasonForNoHostnameOverride,
	"2.1.12": reasonForNotRotatingCerts,
	"2.1.13": reasonForNotRotatingCerts,
	"2.2.3":  reasonNoConfigFileKubeletSvc,
	"2.2.4":  reasonNoConfigFileKubeletSvc,
	"2.2.9":  reasonNoConfigFileKubelet,
	"2.2.10": reasonNoConfigFileKubelet,
}

var rkeCIS14SkippedChecks = map[string]string{
	"1.1.11": reasonForAlwaysPullImages,
	"1.1.21": reasonForKubeletServerCerts,
	"1.1.24": reasonForPSP,
	"1.1.34": reasonForEncryption,
	"1.1.35": reasonForEncryption,
	"1.1.36": reasonForEventRateLimit,
	"1.2.2":  reasonForLocalhostListeningScheduler,
	"1.3.7":  reasonForLocalhostListeningControllerManager,
	"1.4.12": reasonForEtcdDataDir,
	"1.7.2":  reasonForPSP,
	"1.7.3":  reasonForPSP,
	"1.7.4":  reasonForPSP,
	"1.7.5":  reasonForPSP,
	"2.1.6":  reasonForProtectKernelDefaults,
	"2.1.10": reasonForKubeletServerCerts,
}

var rkeCIS15NotApplicableChecks = map[string]string{
	"1.1.1":  reasonNoConfigFileAPIServer,
	"1.1.2":  reasonNoConfigFileAPIServer,
	"1.1.3":  reasonNoConfigFileControllerManager,
	"1.1.4":  reasonNoConfigFileControllerManager,
	"1.1.5":  reasonNoConfigFileScheduler,
	"1.1.6":  reasonNoConfigFileScheduler,
	"1.1.7":  reasonNoConfigFileEtcd,
	"1.1.8":  reasonNoConfigFileEtcd,
	"1.1.13": reasonNoKubeConfigDefault,
	"1.1.14": reasonNoKubeConfigDefault,
	"1.1.15": reasonNoConfigFileScheduler,
	"1.1.16": reasonNoConfigFileScheduler,
	"1.1.17": reasonNoConfigFileControllerManager,
	"1.1.18": reasonNoConfigFileControllerManager,
	"1.3.6":  reasonForNotRotatingCerts,
	"4.1.1":  reasonNoConfigFileKubeletSvc,
	"4.1.2":  reasonNoConfigFileKubeletSvc,
	"4.1.9":  reasonNoConfigFileKubelet,
	"4.1.10": reasonNoConfigFileKubelet,
	"4.2.12": reasonForNotRotatingCerts,
}

var rkeCIS15SkippedChecks = map[string]string{
	"1.1.12": reasonForEtcdDataDir,
	"1.2.6":  reasonForKubeletServerCerts,
	"1.2.16": reasonForPSP,
	"1.2.33": reasonForEncryption,
	"1.2.34": reasonForEncryption,
	"4.2.6":  reasonForProtectKernelDefaults,
	"4.2.10": reasonForKubeletServerCerts,
	"5.1.5":  reasonForDefaultSA,
	"5.2.2":  reasonForPSP,
	"5.2.3":  reasonForPSP,
	"5.2.4":  reasonForPSP,
	"5.2.5":  reasonForPSP,
	"5.3.2":  reasonForNetPol,
	"5.6.4":  reasonForDefaultNS,
}

func loadCisConfigParams() map[string]kdm.CisConfigParams {
	return map[string]kdm.CisConfigParams{
		"default": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.15": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.16": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.17": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.18": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.19": {
			BenchmarkVersion: "rke-cis-1.5",
		},
		"v1.20": {
			BenchmarkVersion: "rke-cis-1.5",
		},
	}
}

func loadCisBenchmarkVersionInfo() map[string]kdm.CisBenchmarkVersionInfo {
	return map[string]kdm.CisBenchmarkVersionInfo{
		"rke-cis-1.4": {
			Managed:              true,
			MinKubernetesVersion: "1.15",
			SkippedChecks:        rkeCIS14SkippedChecks,
			NotApplicableChecks:  rkeCIS14NotApplicableChecks,
		},
		"cis-1.4": {
			Managed:              false,
			MinKubernetesVersion: "1.13",
		},
		"rke-cis-1.5": {
			Managed:              true,
			MinKubernetesVersion: "1.15",
			SkippedChecks:        rkeCIS15SkippedChecks,
			NotApplicableChecks:  rkeCIS15NotApplicableChecks,
		},
		"cis-1.5": {
			Managed:              false,
			MinKubernetesVersion: "1.15",
		},
	}
}
