package templates

/*
Not including Vsphere(cloudProvider) and Authz templates
Will they change and require Rancher to pass them to RKE
depending on Kubernetes version?
*/

import (
	"github.com/rancher/rke/types/kdm"
)

const (
	calicov18            = "calico-v1.8"
	calicov113           = "calico-v1.13"
	calicov115           = "calico-v1.15"
	calicov11512         = "calico-v1.15.12"
	calicov115Privileged = "calico-v1.15-privileged"
	calicov116           = "calico-v1.16"
	calicov117           = "calico-v1.17"
	calicov117Privileged = "calico-v1.17-privileged"
	calicov3160          = "calico-v3.16.0"
	calicov3165          = "calico-v3.16.5"
	calicov3171          = "calico-v3.17.1"

	canalv18                      = "canal-v1.8"
	canalv113                     = "canal-v1.13"
	canalv115                     = "canal-v1.15"
	canalv11512                   = "canal-v1.15.12"
	canalv115Privileged           = "canal-v1.15-privileged"
	canalv115PrivilegedCalico3134 = "canal-v1.15-privileged-calico3134"
	canalv116                     = "canal-v1.16"
	canalv117                     = "canal-v1.17"
	canalv117Privileged           = "canal-v1.17-privileged"
	canalv117PrivilegedCalico3134 = "canal-v1.17-privileged-calico3134"
	canalv3160                    = "canal-v3.16.0"
	canalv3165                    = "canal-v3.16.5"
	canalv3171                    = "canal-v3.17.1"

	flannelv18  = "flannel-v1.8"
	flannelv115 = "flannel-v1.15"
	flannelv116 = "flannel-v1.16"

	coreDnsv18  = "coredns-v1.8"
	coreDnsv116 = "coredns-v1.16"
	coreDnsv117 = "coredns-v1.17"

	kubeDnsv18  = "kubedns-v1.8"
	kubeDnsv116 = "kubedns-v1.16"

	metricsServerv18  = "metricsserver-v1.8"
	metricsServerv120 = "metricsserver-v1.20"

	weavev18  = "weave-v1.8"
	weavev116 = "weave-v1.16"
	weavev120 = "weave-v1.20"
	/* aciv500 Supports k8s versions 1.17 and 1.18 */
	/* Versioning: va.b.c-<special-attr/base-if-none>-x.y.z where a.b.c is ACI version and x.y.z is the k8s version,
	if required
	*/
	aciv500 = "aci-v5.0.0"
	aciv513 = "aci-v5.1.3"

	nginxIngressv18    = "nginxingress-v1.8"
	nginxIngressV115   = "nginxingress-v1.15"
	nginxIngressV11512 = "nginxingress-v1.15.12"

	nodelocalv115 = "nodelocal-v1.15"
)

var TemplateIntroducedRanges = map[string][]string{
	kdm.Nodelocal: {">=1.17.4-rancher1-1", ">=1.16.8-rancher1-1 <1.17.0-alpha", ">=1.15.11-rancher1-1 <1.16.0-alpha"},
}

func LoadK8sVersionedTemplates() map[string]map[string]string {
	return map[string]map[string]string{
		kdm.Calico: {
			">=1.20.4-rancher1-1":                    calicov3171,
			">=1.19.4-rancher1-2 <1.20.4-rancher1-1": calicov3165,
			">=1.19.0-rancher0 <1.19.4-rancher1-2":   calicov3160,
			">=1.17.4-rancher0 <1.19.0-rancher0":     calicov117Privileged,
			">=1.17.0-rancher0 <1.17.4-rancher0":     calicov117,
			">=1.16.8-rancher0 <1.17.0-rancher0":     calicov117Privileged,
			">=1.16.4-rancher1 <1.16.8-rancher0":     calicov117,
			">=1.16.0-alpha <1.16.4-rancher1":        calicov116,

			">=1.15.11-rancher1-1 <1.15.12-rancher1-1": calicov115Privileged,
			// 1.15.12-rancher1-1 comes from 2.2.13, uses calicov115 template with new key calicov11512
			// new key is to enable rancher passing template to rke for an already vendored template
			">=1.15.12-rancher1-1 <1.15.12-rancher2-2": calicov11512,
			">=1.15.12-rancher2-2 <1.16.0-alpha":       calicov115Privileged,

			">=1.15.0-rancher0 <1.15.11-rancher1-1": calicov115,
			">=1.13.0-rancher0 <1.15.0-rancher0":    calicov113,
			">=1.8.0-rancher0 <1.13.0-rancher0":     calicov18,
		},
		kdm.Canal: {
			">=1.20.4-rancher1-1":                      canalv3171,
			">=1.19.4-rancher1-2 <1.20.4-rancher1-1":   canalv3165,
			">=1.19.0-rancher0 <1.19.4-rancher1-2":     canalv3160,
			">=1.17.6-rancher2-1 <1.19.0-rancher0":     canalv117PrivilegedCalico3134,
			">=1.17.4-rancher0 <1.17.6-rancher2-1":     canalv117Privileged,
			">=1.17.0-rancher0 <1.17.4-rancher0":       canalv117,
			">=1.16.10-rancher2-1 <1.17.0-rancher0":    canalv117PrivilegedCalico3134,
			">=1.16.8-rancher0 <1.16.10-rancher2-1":    canalv117Privileged,
			">=1.16.4-rancher1 <1.16.8-rancher0":       canalv117,
			">=1.16.0-alpha <1.16.4-rancher1":          canalv116,
			">=1.15.12-rancher2-2 <1.16.0-alpha":       canalv115PrivilegedCalico3134,
			">=1.15.11-rancher1-1 <1.15.12-rancher1-1": canalv115Privileged,
			// 1.15.12-rancher1-1 comes from 2.2.13, uses old canalv115 template with new key canalv11512
			// new key is to enable rancher passing template to rke for an already vendored template
			">=1.15.12-rancher1-1 <1.15.12-rancher2-2": canalv11512,
			">=1.15.0-rancher0 <1.15.11-rancher1-1":    canalv115,
			">=1.13.0-rancher0 <1.15.0-rancher0":       canalv113,
			">=1.8.0-rancher0 <1.13.0-rancher0":        canalv18,
		},
		kdm.Flannel: {
			">=1.16.0-alpha":                    flannelv116,
			">=1.15.0-rancher0 <1.16.0-alpha":   flannelv115,
			">=1.8.0-rancher0 <1.15.0-rancher0": flannelv18,
		},
		kdm.CoreDNS: {
			">=1.17.0-alpha":                 coreDnsv117,
			">=1.16.0-alpha <1.17.0-alpha":   coreDnsv116,
			">=1.8.0-rancher0 <1.16.0-alpha": coreDnsv18,
		},
		kdm.KubeDNS: {
			">=1.16.0-alpha":                 kubeDnsv116,
			">=1.8.0-rancher0 <1.16.0-alpha": kubeDnsv18,
		},
		kdm.MetricsServer: {
			">=1.20.4-rancher1-1":                 metricsServerv120,
			">=1.8.0-rancher0 <1.20.4-rancher1-1": metricsServerv18,
		},
		kdm.Weave: {
			">=1.20.4-rancher1-1":               weavev120,
			">=1.16.0-alpha <1.20.4-rancher1-1": weavev116,
			">=1.8.0-rancher0 <1.16.0-alpha":    weavev18,
		},
		kdm.Aci: {
			">=1.17.0-alpha <=1.18.18-rancher1-2":    aciv500,
			">=1.19.0-rancher0 <=1.19.10-rancher1-1": aciv500,
			">=1.20.0-rancher0 <=1.20.6-rancher1-1":  aciv500,
			">=1.18.18-rancher1-3 <1.19.0-rancher0":  aciv513,
			">=1.19.10-rancher1-2 <1.20.0-rancher0":  aciv513,
			">=1.20.6-rancher1-2":                    aciv513,
		},
		kdm.NginxIngress: {
			">=1.8.0-rancher0 <1.13.10-rancher1-3":  nginxIngressv18,
			">=1.14.0-rancher0 <=1.14.6-rancher1-1": nginxIngressv18,
			">=1.15.0-rancher0 <=1.15.3-rancher1-1": nginxIngressv18,
			// New ingress template introduced for 1.13.10-rancher1-3, 1.14.6-rancher2, 1.15.3-rancher2
			">=1.13.10-rancher1-3 <1.14.0-rancher0":   nginxIngressV115,
			">=1.14.6-rancher2 <1.15.0-rancher0":      nginxIngressV115,
			">=1.15.3-rancher2 <1.15.12-rancher1-1":   nginxIngressV115,
			">=1.16.1-rancher1-1 <1.16.10-rancher1-1": nginxIngressV115,
			">=1.17.0-rancher1-1 <1.17.6-rancher1-1":  nginxIngressV115,
			// New ingress template introduced for 1.15.12-rancher1-1, 1.16.10-rancher1-1, 1.17.6-rancher1-1
			">=1.15.12-rancher1-1 <1.16.1-rancher1-1": nginxIngressV11512,
			">=1.16.10-rancher1-1 <1.17.0-rancher1-1": nginxIngressV11512,
			">=1.17.6-rancher1-1":                     nginxIngressV11512,
		},
		kdm.Nodelocal: {
			">=1.15.11-rancher0 <1.16.0-alpha": nodelocalv115,
			">=1.16.8-rancher0 <1.17.0-alpha":  nodelocalv115,
			">=1.17.4-rancher0":                nodelocalv115,
		},
		kdm.TemplateKeys: getTemplates(),
	}
}

func getTemplates() map[string]string {
	return map[string]string{
		calicov113:           CalicoTemplateV113,
		calicov115:           CalicoTemplateV115,
		calicov11512:         CalicoTemplateV115,
		calicov115Privileged: CalicoTemplateV115Privileged,
		calicov116:           CalicoTemplateV116,
		calicov117:           CalicoTemplateV117,
		calicov117Privileged: CalicoTemplateV117Privileged,
		calicov18:            CalicoTemplateV112,
		calicov3160:          CalicoTemplateV3_16_0,
		calicov3165:          CalicoTemplateV3_16_5,
		calicov3171:          CalicoTemplateV3_17_1,

		flannelv115: FlannelTemplateV115,
		flannelv116: FlannelTemplateV116,
		flannelv18:  FlannelTemplate,

		canalv113:                     CanalTemplateV113,
		canalv18:                      CanalTemplateV112,
		canalv115:                     CanalTemplateV115,
		canalv11512:                   CanalTemplateV115,
		canalv115Privileged:           CanalTemplateV115Privileged,
		canalv115PrivilegedCalico3134: CanalTemplateV115PrivilegedCalico3134,
		canalv116:                     CanalTemplateV116,
		canalv117:                     CanalTemplateV117,
		canalv117Privileged:           CanalTemplateV117Privileged,
		canalv117PrivilegedCalico3134: CanalTemplateV117PrivilegedCalico3134,
		canalv3160:                    CanalTemplateV3_16_0,
		canalv3165:                    CanalTemplateV3_16_5,
		canalv3171:                    CanalTemplateV3_17_1,

		coreDnsv18:  CoreDNSTemplate,
		coreDnsv116: CoreDNSTemplateV116,
		coreDnsv117: CoreDNSTemplateV117,

		kubeDnsv18:  KubeDNSTemplate,
		kubeDnsv116: KubeDNSTemplateV116,

		metricsServerv18:  MetricsServerTemplate,
		metricsServerv120: MetricsServerTemplate_v0_4_1,

		weavev18:  WeaveTemplate,
		weavev116: WeaveTemplateV116,
		weavev120: WeaveTemplateV120,

		aciv500: AciTemplateV500,
		aciv513: AciTemplateV513,

		nginxIngressv18:    NginxIngressTemplate,
		nginxIngressV115:   NginxIngressTemplateV0251Rancher1,
		nginxIngressV11512: NginxIngressTemplateV0320Rancher1,

		nodelocalv115: NodelocalTemplateV115,
	}
}
