package rke

import (
	"fmt"

	v3 "github.com/rancher/rke/types"
)

const (
	rkeIssuer              = "rke"
	serviceAccountKeyPath  = "/etc/kubernetes/ssl/kube-service-account-token-key.pem"
	tlsCipherSuites        = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305"
	enableAdmissionPlugins = "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,NodeRestriction"
)

func loadK8sVersionServiceOptions() map[string]v3.KubernetesServicesOptions {
	return map[string]v3.KubernetesServicesOptions{
		"v1.29": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions127(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.28": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions127(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.27": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions127(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.26": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions124(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.25": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions124(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.24": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions124(),
			Kubelet:        getKubeletOptions124(),
			KubeController: getKubeControllerOptions124(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions124(),
		},
		"v1.23.4-rancher1-2": {
			Etcd:           getETCDOptions122v352(),
			KubeAPI:        getKubeAPIOptions121(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.22.7-rancher1-2": {
			Etcd:           getETCDOptions122v352(),
			KubeAPI:        getKubeAPIOptions121(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.23": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions121(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.22": {
			Etcd:           getETCDOptions122(),
			KubeAPI:        getKubeAPIOptions121(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.21": {
			Etcd:           getETCDOptions117(),
			KubeAPI:        getKubeAPIOptions121(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.20": {
			Etcd:           getETCDOptions117(),
			KubeAPI:        getKubeAPIOptions120(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.19": {
			Etcd:           getETCDOptions117(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.18": {
			Etcd:           getETCDOptions117(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.17": {
			Etcd:           getETCDOptions117(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.3-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.4-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.6-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.6-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.7-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.8-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.8-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.10-rancher2-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.13-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.13-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.14-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.15-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.15-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16.15-rancher1-4": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.16": {
			KubeAPI:        getKubeAPIOptions116(),
			Kubelet:        getKubeletOptions116(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.6-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.7-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.9-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.10-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.11-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.11-rancher1-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.12-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.12-rancher2-2": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.12-rancher2-3": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.12-rancher2-5": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15.12-rancher2-8": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.15": {
			KubeAPI:        getKubeAPIOptions115(),
			Kubelet:        getKubeletOptions115(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.14.9-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions114(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.14.10-rancher1-1": {
			Etcd:           getETCDOptions(),
			KubeAPI:        getKubeAPIOptions114(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.14": {
			KubeAPI:        getKubeAPIOptions114(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.13": {
			KubeAPI:        getKubeAPIOptions(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.12": {
			KubeAPI:        getKubeAPIOptions(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.11": {
			KubeAPI:        getKubeAPIOptions(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.10": {
			KubeAPI:        getKubeAPIOptions(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
		"v1.9": {
			KubeAPI:        getKubeAPIOptions19(),
			Kubelet:        getKubeletOptions(),
			KubeController: getKubeControllerOptions(),
			Kubeproxy:      getKubeProxyOptions(),
			Scheduler:      getSchedulerOptions(),
		},
	}
}

func getKubeAPIOptions() map[string]string {
	data := map[string]string{
		"tls-cipher-suites":                  tlsCipherSuites,
		"enable-admission-plugins":           enableAdmissionPlugins, // order doesn't matter >= 1.10
		"allow-privileged":                   "true",
		"anonymous-auth":                     "false",
		"bind-address":                       "0.0.0.0",
		"insecure-port":                      "0",
		"kubelet-preferred-address-types":    "InternalIP,ExternalIP,Hostname",
		"profiling":                          "false",
		"requestheader-extra-headers-prefix": "X-Remote-Extra-",
		"requestheader-group-headers":        "X-Remote-Group",
		"requestheader-username-headers":     "X-Remote-User",
		"secure-port":                        "6443",
		"service-account-lookup":             "true",
		"storage-backend":                    "etcd3",
		"runtime-config":                     "authorization.k8s.io/v1beta1=true",
	}
	return data
}

func getKubeAPIOptions19() map[string]string {
	kubeAPIOptions := getKubeAPIOptions()
	kubeAPIOptions["admission-control"] = "ServiceAccount,NamespaceLifecycle,LimitRanger,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds,NodeRestriction"
	return kubeAPIOptions
}

func getKubeAPIOptions114() map[string]string {
	kubeAPIOptions := getKubeAPIOptions()
	kubeAPIOptions["enable-admission-plugins"] = fmt.Sprintf("%s,%s", enableAdmissionPlugins, "Priority")
	kubeAPIOptions["runtime-config"] = "authorization.k8s.io/v1beta1=true"
	return kubeAPIOptions
}

func getKubeAPIOptions115() map[string]string {
	kubeAPIOptions := getKubeAPIOptions114()
	kubeAPIOptions["enable-admission-plugins"] = fmt.Sprintf("%s,%s", kubeAPIOptions["enable-admission-plugins"], "TaintNodesByCondition,PersistentVolumeClaimResize")
	kubeAPIOptions["runtime-config"] = "authorization.k8s.io/v1beta1=true"
	return kubeAPIOptions
}

func getKubeAPIOptions116() map[string]string {
	kubeAPIOptions := getKubeAPIOptions114()
	kubeAPIOptions["enable-admission-plugins"] = fmt.Sprintf("%s,%s", kubeAPIOptions["enable-admission-plugins"], "TaintNodesByCondition,PersistentVolumeClaimResize")
	kubeAPIOptions["runtime-config"] = "authorization.k8s.io/v1beta1=true"
	return kubeAPIOptions
}

func getKubeAPIOptions120() map[string]string {
	kubeAPIOptions := getKubeAPIOptions116()
	// need to turn off this feature in k8s 1.20 since we are not using it https://github.com/kubernetes/kubernetes/pull/91921
	kubeAPIOptions["feature-gates"] = "ServiceAccountIssuerDiscovery=false"
	kubeAPIOptions["service-account-issuer"] = rkeIssuer
	kubeAPIOptions["service-account-signing-key-file"] = serviceAccountKeyPath
	kubeAPIOptions["api-audiences"] = "unknown"
	return kubeAPIOptions
}

func getKubeAPIOptions121() map[string]string {
	kubeAPIOptions := getKubeAPIOptions116()
	kubeAPIOptions["service-account-issuer"] = rkeIssuer
	kubeAPIOptions["service-account-signing-key-file"] = serviceAccountKeyPath
	kubeAPIOptions["api-audiences"] = "unknown"
	return kubeAPIOptions
}

func getKubeAPIOptions124() map[string]string {
	kubeAPIOptions := getKubeAPIOptions116()
	kubeAPIOptions["service-account-issuer"] = rkeIssuer
	kubeAPIOptions["service-account-signing-key-file"] = serviceAccountKeyPath
	kubeAPIOptions["api-audiences"] = "unknown"
	delete(kubeAPIOptions, "insecure-port")
	return kubeAPIOptions
}

// getKubeletOptions provides the root options for windows
// note: please double-check on windows side if changing the following options
func getKubeletOptions() map[string]string {
	return map[string]string{
		"tls-cipher-suites":                 tlsCipherSuites,
		"address":                           "0.0.0.0",
		"allow-privileged":                  "true",
		"anonymous-auth":                    "false",
		"authentication-token-webhook":      "true",
		"cgroups-per-qos":                   "True",
		"cni-bin-dir":                       "/opt/cni/bin",
		"cni-conf-dir":                      "/etc/cni/net.d",
		"enforce-node-allocatable":          "",
		"event-qps":                         "0",
		"make-iptables-util-chains":         "true",
		"network-plugin":                    "cni",
		"read-only-port":                    "0",
		"resolv-conf":                       "/etc/resolv.conf",
		"streaming-connection-idle-timeout": "30m",
		"volume-plugin-dir":                 "/var/lib/kubelet/volumeplugins",
		"v":                                 "2",
		"authorization-mode":                "Webhook",
	}
}

func getKubeletOptions115() map[string]string {
	kubeletOptions := getKubeletOptions()
	kubeletOptions["authorization-mode"] = "Webhook"
	delete(kubeletOptions, "allow-privileged")
	return kubeletOptions
}

func getKubeletOptions116() map[string]string {
	kubeletOptions := getKubeletOptions()
	kubeletOptions["authorization-mode"] = "Webhook"
	delete(kubeletOptions, "allow-privileged")
	return kubeletOptions
}

func getKubeletOptions124() map[string]string {
	kubeletOptions := getKubeletOptions()
	kubeletOptions["authorization-mode"] = "Webhook"
	delete(kubeletOptions, "allow-privileged")
	delete(kubeletOptions, "cni-conf-dir")
	delete(kubeletOptions, "cni-bin-dir")
	delete(kubeletOptions, "network-plugin")
	return kubeletOptions
}

func getKubeControllerOptions() map[string]string {
	return map[string]string{
		"address":                     "0.0.0.0",
		"allow-untagged-cloud":        "true",
		"allocate-node-cidrs":         "true",
		"configure-cloud-routes":      "false",
		"enable-hostpath-provisioner": "false",
		"leader-elect":                "true",
		"node-monitor-grace-period":   "40s",
		"pod-eviction-timeout":        "5m0s",
		"profiling":                   "false",
		"terminated-pod-gc-threshold": "1000",
		"v":                           "2",
	}
}

func getKubeControllerOptions124() map[string]string {
	kubeControllerOptions := getKubeControllerOptions()
	delete(kubeControllerOptions, "address")
	return kubeControllerOptions
}

func getKubeControllerOptions127() map[string]string {
	kubeControllerOptions := getKubeControllerOptions()
	delete(kubeControllerOptions, "address")
	delete(kubeControllerOptions, "pod-eviction-timeout")
	return kubeControllerOptions
}

// getKubeProxyOptions provides the root options for windows
// note: please double-check on windows side if changing the following options
func getKubeProxyOptions() map[string]string {
	return map[string]string{
		"v":                    "2",
		"healthz-bind-address": "127.0.0.1",
	}
}

func getSchedulerOptions() map[string]string {
	return map[string]string{
		"leader-elect": "true",
		"v":            "2",
		"address":      "0.0.0.0",
		"profiling":    "false",
	}
}

func getSchedulerOptions124() map[string]string {
	kubeSchedulerOptions := getSchedulerOptions()
	delete(kubeSchedulerOptions, "address")
	return kubeSchedulerOptions
}

func getETCDOptions() map[string]string {
	return map[string]string{
		"client-cert-auth":      "true",
		"peer-client-cert-auth": "true",
	}
}

func getETCDOptions117() map[string]string {
	return map[string]string{
		"client-cert-auth":      "true",
		"peer-client-cert-auth": "true",
		"enable-v2":             "true",
	}
}

func getETCDOptions122() map[string]string {
	return map[string]string{
		"client-cert-auth":      "true",
		"peer-client-cert-auth": "true",
	}
}

func getETCDOptions122v352() map[string]string {
	return map[string]string{
		"client-cert-auth":                   "true",
		"peer-client-cert-auth":              "true",
		"experimental-initial-corrupt-check": "true",
	}
}
