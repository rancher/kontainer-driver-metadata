package rke

import (
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
)

const (
	// reasons for not applicable checks
	reasonNoConfigFileApiServer = `Cluster provisioned by RKE doesn't require or maintain a configuration file for kube-apiserver.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileEtcd = `Cluster provisioned by RKE doesn't require or maintain a configuration file for etcd.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileScheduler = `Cluster provisioned by RKE doesn't require or maintain a configuration file for scheduler.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileControllerManager = `Cluster provisioned by RKE doesn't require or maintain a configuration file for controller-manager.
All configuration is passed in as arguments at container run time.`
	reasonNoKubeConfigDefault    = `RKE does not store the kubernetes default kubeconfig credentials file on the nodes.`
	reasonNoConfigFileKubeletSvc = `RKE doesn’t require or maintain a configuration file for the kubelet service.
All configuration is passed in as arguments at container run time.`
	reasonNoConfigFileKubelet = `RKE doesn’t require or maintain a configuration file for the kubelet.
All configuration is passed in as arguments at container run time.`

	// reasons for skipped checks
	reasonForEtcdDataDir           = `TODO`
	reasonForKubeletCA             = `TODO`
	reasonForPSP                   = `Enabling Pod Security Policy can cause issues with many helm chart installations`
	reasonForAuditLog              = `TODO`
	reasonForEncryption            = `TODO`
	reasonForKubeletCertRotation   = `TODO`
	reasonForProtectKernelDefaults = `TODO`
	reasonForKubeletTLS            = `TODO`
	reasonForDefaultSA             = `TODO`
	reasonForNetPol                = `Enabling Network Policies can cause lot of unintended network traffic disruptions`
	reasonForDefaultNS             = `A default namespace provides a flexible workspace to try out various deployments`
	reasonForRotateCerts           = `TODO`
	reasonForHostnameOverride      = `TODO`
	reasonForAlwaysPullImages      = `TODO`
	reasonForEventRateLimit        = `TODO`
)

var rkeCIS14NotApplicableChecks = map[string]string{
	"1.4.1":  reasonNoConfigFileApiServer,
	"1.4.2":  reasonNoConfigFileApiServer,
	"1.4.3":  reasonNoConfigFileControllerManager,
	"1.4.4":  reasonNoConfigFileControllerManager,
	"1.4.5":  reasonNoConfigFileScheduler,
	"1.4.6":  reasonNoConfigFileScheduler,
	"1.4.7":  reasonNoConfigFileEtcd,
	"1.4.8":  reasonNoConfigFileEtcd,
	"1.4.13": reasonNoKubeConfigDefault,
	"1.4.14": reasonNoKubeConfigDefault,
	"2.2.3":  reasonNoConfigFileKubeletSvc,
	"2.2.4":  reasonNoConfigFileKubeletSvc,
	"2.2.9":  reasonNoConfigFileKubelet,
	"2.2.10": reasonNoConfigFileKubelet,
}

var rkeCIS14SkippedChecks = map[string]string{
	"1.1.11":  reasonForAlwaysPullImages,
	"1.1.15":  reasonForAuditLog,
	"1.1.16":  reasonForAuditLog,
	"1.1.17":  reasonForAuditLog,
	"1.1.18":  reasonForAuditLog,
	"1.1.24":  reasonForPSP,
	"1.1.34":  reasonForEncryption,
	"1.1.35":  reasonForEncryption,
	"1.1.36":  reasonForEventRateLimit,
	"1.1.37a": reasonForAuditLog,
	"1.1.37b": reasonForAuditLog,
	"1.3.6":   reasonForKubeletCertRotation,
	"1.4.12":  reasonForEtcdDataDir,
	"1.7.2":   reasonForPSP,
	"1.7.3":   reasonForPSP,
	"1.7.4":   reasonForPSP,
	"1.7.5":   reasonForPSP,
	"2.1.6":   reasonForProtectKernelDefaults,
	"2.1.8":   reasonForHostnameOverride,
	"2.1.10":  reasonForKubeletTLS,
	"2.1.12":  reasonForRotateCerts,
}

var rkeCIS15NotApplicableChecks = map[string]string{
	"1.1.1":  reasonNoConfigFileApiServer,
	"1.1.2":  reasonNoConfigFileApiServer,
	"1.1.3":  reasonNoConfigFileApiServer,
	"1.1.4":  reasonNoConfigFileApiServer,
	"1.1.5":  reasonNoConfigFileApiServer,
	"1.1.6":  reasonNoConfigFileApiServer,
	"1.1.7":  reasonNoConfigFileApiServer,
	"1.1.8":  reasonNoConfigFileApiServer,
	"1.1.13": reasonNoKubeConfigDefault,
	"1.1.14": reasonNoKubeConfigDefault,
	"1.1.15": reasonNoConfigFileScheduler,
	"1.1.16": reasonNoConfigFileScheduler,
	"1.1.17": reasonNoConfigFileControllerManager,
	"1.1.18": reasonNoConfigFileControllerManager,
	"4.1.1":  reasonNoConfigFileKubeletSvc,
	"4.1.2":  reasonNoConfigFileKubeletSvc,
	"4.1.9":  reasonNoConfigFileKubelet,
	"4.1.10": reasonNoConfigFileKubelet,
}

var rkeCIS15SkippedChecks = map[string]string{
	"1.1.12": reasonForEtcdDataDir,
	"1.2.6":  reasonForKubeletCA,
	"1.2.16": reasonForPSP,
	"1.2.22": reasonForAuditLog,
	"1.2.23": reasonForAuditLog,
	"1.2.24": reasonForAuditLog,
	"1.2.25": reasonForAuditLog,
	"1.2.33": reasonForEncryption,
	"1.2.34": reasonForEncryption,
	"1.3.6":  reasonForKubeletCertRotation,
	"3.2.1":  reasonForAuditLog,
	"4.2.6":  reasonForProtectKernelDefaults,
	"4.2.10": reasonForKubeletTLS,
	"4.2.12": reasonForKubeletCertRotation,
	"5.1.5":  reasonForDefaultSA,
	"5.2.2":  reasonForPSP,
	"5.2.3":  reasonForPSP,
	"5.2.4":  reasonForPSP,
	"5.2.5":  reasonForPSP,
	"5.3.2":  reasonForNetPol,
	"5.6.4":  reasonForDefaultNS,
}

func loadCisConfigParams() map[string]v3.CisConfigParams {
	return map[string]v3.CisConfigParams{
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
	}
}

func loadCisBenchmarkVersionInfo() map[string]v3.CisBenchmarkVersionInfo {
	return map[string]v3.CisBenchmarkVersionInfo{
		"rke-cis-1.4": {
			Managed:              true,
			MinKubernetesVersion: "1.13",
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
