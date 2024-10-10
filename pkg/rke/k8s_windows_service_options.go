package rke

import v3 "github.com/rancher/rke/types"

func loadK8sVersionWindowsServiceOptions() map[string]v3.KubernetesServicesOptions {
	// since 1.14, windows has been supported
	return map[string]v3.KubernetesServicesOptions{
		"v1.31": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.30": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.29": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.28": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.27": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.26": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.25": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.24": {
			Kubelet:   getWindowsKubeletOptions124(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.23.4-rancher1-2": {
			Kubelet:   getWindowsKubeletOptions121(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.23": {
			Kubelet:   getWindowsKubeletOptions121(),
			Kubeproxy: getWindowsKubeProxyOptions123(),
		},
		"v1.22.7-rancher1-2": {
			Kubelet:   getWindowsKubeletOptions121(),
			Kubeproxy: getWindowsKubeProxyOptions121(),
		},
		"v1.22": {
			Kubelet:   getWindowsKubeletOptions121(),
			Kubeproxy: getWindowsKubeProxyOptions121(),
		},
		"v1.21": {
			Kubelet:   getWindowsKubeletOptions121(),
			Kubeproxy: getWindowsKubeProxyOptions121(),
		},
		"v1.20": {
			Kubelet:   getWindowsKubeletOptions116(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
		"v1.19": {
			Kubelet:   getWindowsKubeletOptions116(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
		"v1.18": {
			Kubelet:   getWindowsKubeletOptions116(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
		"v1.17": {
			Kubelet:   getWindowsKubeletOptions116(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
		"v1.16": {
			Kubelet:   getWindowsKubeletOptions116(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
		"v1.15": {
			Kubelet:   getWindowsKubeletOptions115(),
			Kubeproxy: getWindowsKubeProxyOptions(),
		},
	}
}

func getWindowsKubeletOptions() map[string]string {
	kubeletOptions := getKubeletOptions()

	// doesn't support cgroups
	kubeletOptions["cgroups-per-qos"] = "false"
	kubeletOptions["enforce-node-allocatable"] = "''"
	// doesn't support dns
	kubeletOptions["resolv-conf"] = "''"
	// add prefix path for directory options
	kubeletOptions["cni-bin-dir"] = "[PREFIX_PATH]/opt/cni/bin"
	kubeletOptions["cni-conf-dir"] = "[PREFIX_PATH]/etc/cni/net.d"
	kubeletOptions["cert-dir"] = "[PREFIX_PATH]/var/lib/kubelet/pki"
	kubeletOptions["volume-plugin-dir"] = "[PREFIX_PATH]/var/lib/kubelet/volumeplugins"
	// add reservation for kubernetes components
	kubeletOptions["kube-reserved"] = "cpu=500m,memory=500Mi,ephemeral-storage=1Gi"
	// add reservation for system
	kubeletOptions["system-reserved"] = "cpu=1000m,memory=2Gi,ephemeral-storage=2Gi"
	// increase image pulling deadline
	kubeletOptions["image-pull-progress-deadline"] = "30m"
	// enable some windows features
	kubeletOptions["feature-gates"] = "HyperVContainer=true,WindowsGMSA=true"

	return kubeletOptions
}

func getWindowsKubeletOptions115() map[string]string {
	kubeletOptions := getWindowsKubeletOptions()

	// doesn't support `allow-privileged`
	delete(kubeletOptions, "allow-privileged")

	return kubeletOptions
}

func getWindowsKubeletOptions116() map[string]string {
	kubeletOptions := getWindowsKubeletOptions()

	// doesn't support `allow-privileged`
	delete(kubeletOptions, "allow-privileged")

	return kubeletOptions
}

func getWindowsKubeletOptions121() map[string]string {
	kubeletOptions := getWindowsKubeletOptions()

	// doesn't support `allow-privileged`
	delete(kubeletOptions, "allow-privileged")
	delete(kubeletOptions, "feature-gates")

	return kubeletOptions
}

func getWindowsKubeletOptions124() map[string]string {
	kubeletOptions := getWindowsKubeletOptions()

	// doesn't support `allow-privileged`
	delete(kubeletOptions, "allow-privileged")
	delete(kubeletOptions, "feature-gates")
	delete(kubeletOptions, "cni-conf-dir")
	delete(kubeletOptions, "cni-bin-dir")
	delete(kubeletOptions, "image-pull-progress-deadline")
	return kubeletOptions
}

func getWindowsKubeProxyOptions123() map[string]string {
	kubeProxyOptions := getKubeProxyOptions()

	// use kernelspace proxy mode
	kubeProxyOptions["proxy-mode"] = "kernelspace"
	// disable Windows DSR support explicitly
	kubeProxyOptions["enable-dsr"] = "false"

	return kubeProxyOptions
}

func getWindowsKubeProxyOptions121() map[string]string {
	kubeProxyOptions := getKubeProxyOptions()

	// use kernelspace proxy mode
	kubeProxyOptions["proxy-mode"] = "kernelspace"
	// disable Windows IPv6DualStack support, WinOverlay already defaults to true
	kubeProxyOptions["feature-gates"] = "IPv6DualStack=false"
	// disable Windows DSR support explicitly
	kubeProxyOptions["enable-dsr"] = "false"

	return kubeProxyOptions
}

func getWindowsKubeProxyOptions() map[string]string {
	kubeProxyOptions := getKubeProxyOptions()

	// use kernelspace proxy mode
	kubeProxyOptions["proxy-mode"] = "kernelspace"
	// enable Windows Overlay support
	kubeProxyOptions["feature-gates"] = "WinOverlay=true"
	// disable Windows DSR support explicitly
	kubeProxyOptions["enable-dsr"] = "false"

	return kubeProxyOptions
}
