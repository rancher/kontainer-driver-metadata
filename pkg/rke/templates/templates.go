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
	calicov319           = "calico-v3.19.0"
	calicov3211          = "calico-v3.21.1"
	calicov3225          = "calico-v3.22.5"
	calicov3225Rancher2  = "calico-v3.22.5-rancher2"
	calicov3241          = "calico-v3.24.1"
	calicov3250          = "calico-v3.25.0"
	calicov3261          = "calico-v3.26.1"

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
	canalv319                     = "canal-v3.19.0"
	canalv3211                    = "canal-v3.21.1"
	canalv3225                    = "canal-v3.22.5"
	canalv3225Rancher2            = "canal-v3.22.5-rancher2"
	canalv3241                    = "canal-v3.24.1"
	canalv3250                    = "canal-v3.25.0"
	canalv3261                    = "canal-v3.26.1"

	flannelv18           = "flannel-v1.8"
	flannelv115          = "flannel-v1.15"
	flannelv116          = "flannel-v1.16"
	flannelv0140         = "flannel-v0.14.0"
	flannelv0140Rancher2 = "flannel-v0.14.0-rancher2"
	flannelv0192         = "flannel-v0.19.2"
	flannelv0214         = "flannel-v0.21.4"

	coreDnsv18          = "coredns-v1.8"
	coreDnsv116         = "coredns-v1.16"
	coreDnsv117         = "coredns-v1.17"
	coreDnsv183         = "coredns-v1.8.3"
	coreDnsv183Rancher2 = "coredns-v1.8.3-rancher2"

	kubeDnsv18  = "kubedns-v1.8"
	kubeDnsv116 = "kubedns-v1.16"

	metricsServerv18  = "metricsserver-v1.8"
	metricsServerv120 = "metricsserver-v1.20"
	metricsServerv050 = "metricsserver-v0.5.0"
	metricsServerv052 = "metricsserver-v0.5.2"
	metricsServerv061 = "metricsserver-v0.6.1"

	weavev18  = "weave-v1.8"
	weavev116 = "weave-v1.16"
	weavev120 = "weave-v1.20"
	weavev122 = "weave-v1.22"
	/* aciv500 Supports k8s versions 1.17 and 1.18 */
	/* Versioning: va.b.c-<special-attr/base-if-none>-x.y.z where a.b.c is ACI version and x.y.z is the k8s version,
	if required
	*/
	aciv500  = "aci-v5.0.0"
	aciv523  = "aci-v5.2.3"
	aciv5234 = "aci-v5.2.3.4"
	aciv5235 = "aci-v5.2.3.5"
	aciv5236 = "aci-v5.2.3.6"
	aciv5271 = "aci-v5.2.7.1"

	nginxIngressv18          = "nginxingress-v1.8"
	nginxIngressV115         = "nginxingress-v1.15"
	nginxIngressV11512       = "nginxingress-v1.15.12"
	nginxIngressv0481        = "nginxingress-v0.48.1"
	nginxIngressv0493        = "nginxingress-v0.49.3"
	nginxIngressv110         = "nginxingress-v1.1.0"
	nginxIngressv110Rancher2 = "nginxingress-v1.1.0-rancher2"
	nginxIngressv110Rancher3 = "nginxingress-v1.1.0-rancher3"
	nginxIngressv110Rancher4 = "nginxingress-v1.1.0-rancher4"
	nginxIngressv120Rancher1 = "nginxingress-v1.2.0-rancher1"
	nginxIngressv120Rancher2 = "nginxingress-v1.2.0-rancher2"
	nginxIngressv151Rancher2 = "nginxingress-v1.5.1-rancher2"
	nginxIngressv170Rancher1 = "nginxingress-v1.7.0-rancher1"
	nginxIngressv181Rancher1 = "nginxingress-v1.8.1-rancher1"

	nodelocalv115 = "nodelocal-v1.15"
	nodelocalv121 = "nodelocal-v1.21"
)

var TemplateIntroducedRanges = map[string][]string{
	kdm.Nodelocal: {">=1.17.4-rancher1-1", ">=1.16.8-rancher1-1 <1.17.0-alpha", ">=1.15.11-rancher1-1 <1.16.0-alpha"},
}

func LoadK8sVersionedTemplates() map[string]map[string]string {
	return map[string]map[string]string{
		kdm.Calico: {
			">=1.27.0-rancher1-1":                      calicov3261,
			">=1.26.0-rancher1-1 <1.27.0-rancher1":     calicov3250,
			">=1.25.0-rancher1-1 <1.26.0-rancher1-1":   calicov3241,
			">=1.24.10-rancher1-1 <1.25.0-rancher1-1":  calicov3225Rancher2,
			">=1.24.9-rancher1-1 <1.24.10-rancher1-1":  calicov3225,
			">=1.24.0-rancher1-1 <1.24.9-rancher1-1":   calicov3211,
			">=1.23.16-rancher1-1 <1.24.0-rancher1-1":  calicov3225Rancher2,
			">=1.23.15-rancher1-1 <1.23.16-rancher1-1": calicov3225,
			">=1.23.0-rancher1-1 <1.23.15-rancher1-1":  calicov3211,
			">=1.22.17-rancher1-2 <1.23.0-rancher1-1":  calicov3225Rancher2,
			">=1.22.17-rancher1-1 <1.22.17-rancher1-2": calicov3225,
			">=1.22.0-rancher1-1 <1.22.17-rancher1-1":  calicov3211,
			">=1.21.0-rancher1-1 <1.22.0-rancher1-1":   calicov319,
			">=1.20.4-rancher1-1 <1.21.0-rancher1-1":   calicov3171,
			">=1.19.4-rancher1-2 <1.20.4-rancher1-1":   calicov3165,
			">=1.19.0-rancher0 <1.19.4-rancher1-2":     calicov3160,
			">=1.17.4-rancher0 <1.19.0-rancher0":       calicov117Privileged,
			">=1.17.0-rancher0 <1.17.4-rancher0":       calicov117,
			">=1.16.8-rancher0 <1.17.0-rancher0":       calicov117Privileged,
			">=1.16.4-rancher1 <1.16.8-rancher0":       calicov117,
			">=1.16.0-alpha <1.16.4-rancher1":          calicov116,

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
			">=1.27.0-rancher1-1":                      canalv3261,
			">=1.26.0-rancher1-1 <1.27.0-rancher1":     canalv3250,
			">=1.25.0-rancher1-1 <1.26.0-rancher1-1":   canalv3241,
			">=1.24.10-rancher1-1 <1.25.0-rancher1-1":  canalv3225Rancher2,
			">=1.24.9-rancher1-1 <1.24.10-rancher1-1":  canalv3225,
			">=1.24.0-rancher1-1 <1.24.9-rancher1-1":   canalv3211,
			">=1.23.16-rancher1-1 <1.24.0-rancher1-1":  canalv3225Rancher2,
			">=1.23.15-rancher1-1 <1.23.16-rancher1-1": canalv3225,
			">=1.23.0-rancher1-1 <1.23.15-rancher1-1":  canalv3211,
			">=1.22.17-rancher1-2 <1.23.0-rancher1-1":  canalv3225Rancher2,
			">=1.22.17-rancher1-1 <1.22.17-rancher1-2": canalv3225,
			">=1.22.0-rancher1-1 <1.22.17-rancher1-1":  canalv3211,
			">=1.21.0-rancher1-1 <1.22.0-rancher1-1":   canalv319,
			">=1.20.4-rancher1-1 <1.21.0-rancher1-1":   canalv3171,
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
			">=1.26.0-rancher1-1":                     flannelv0214,
			">=1.25.0-rancher1-1 <1.26.0-rancher1-1":  flannelv0192,
			">=1.24.10-rancher1-1 <1.25.0-rancher1-1": flannelv0140Rancher2,
			">=1.24.0-rancher1-1 <1.24.10-rancher1-1": flannelv0140,
			">=1.23.16-rancher1-1 <1.24.0-rancher1-1": flannelv0140Rancher2,
			">=1.21.0-alpha <1.23.16-rancher1-1":      flannelv0140,
			">=1.16.0-alpha <1.21.0-alpha":            flannelv116,
			">=1.15.0-rancher0 <1.16.0-alpha":         flannelv115,
			">=1.8.0-rancher0 <1.15.0-rancher0":       flannelv18,
		},
		kdm.CoreDNS: {
			">=1.21.9-rancher1-2":                     coreDnsv183Rancher2,
			">=1.21.0-rancher1-1 <1.21.9-rancher1-2":  coreDnsv183,
			">=1.20.15-rancher1-2 <1.21.0-rancher1-1": coreDnsv183Rancher2,
			">=1.17.0-alpha <1.20.15-rancher1-2":      coreDnsv117,
			">=1.16.0-alpha <1.17.0-alpha":            coreDnsv116,
			">=1.8.0-rancher0 <1.16.0-alpha":          coreDnsv18,
		},
		kdm.KubeDNS: {
			">=1.16.0-alpha":                 kubeDnsv116,
			">=1.8.0-rancher0 <1.16.0-alpha": kubeDnsv18,
		},
		kdm.MetricsServer: {
			">=1.23.3-rancher1-1":                      metricsServerv061,
			">=1.22.15-rancher1-1 <1.23.3-rancher1-1":  metricsServerv052,
			">=1.20.14-rancher2-1 <1.22.15-rancher1-1": metricsServerv050,
			">=1.20.4-rancher1-1 <1.20.14-rancher2-1":  metricsServerv120,
			">=1.8.0-rancher0 <1.20.4-rancher1-1":      metricsServerv18,
		},
		kdm.Weave: {
			">=1.22.0-rancher1-1":                    weavev122,
			">=1.20.4-rancher1-1 <1.22.0-rancher1-1": weavev120,
			">=1.16.0-alpha <1.20.4-rancher1-1":      weavev116,
			">=1.8.0-rancher0 <1.16.0-alpha":         weavev18,
		},
		kdm.Aci: {
			">=1.17.0-alpha <1.20.15-rancher2-2":       aciv500,
			">=1.20.15-rancher2-2 <1.21.0-rancher0":    aciv523,
			">=1.21.0-rancher0 <1.21.14-rancher1-1":    aciv500,
			">=1.21.14-rancher1-1 <1.22.0-rancher0":    aciv523,
			">=1.22.0-rancher0 <1.22.11-rancher1-1":    aciv500,
			">=1.22.11-rancher1-1 <1.22.16-rancher1-1": aciv523,
			">=1.22.16-rancher1-1 <1.22.17-rancher1-1": aciv5234,
			">=1.22.17-rancher1-1 <1.23.0-rancher0":    aciv5235,
			">=1.23.0-rancher0 <1.23.8-rancher1-1":     aciv500,
			">=1.23.8-rancher1-1 <1.23.14-rancher1-1":  aciv523,
			">=1.23.14-rancher1-1 <1.23.15-rancher1-1": aciv5234,
			">=1.23.15-rancher1-1 <1.23.16-rancher2-2": aciv5235,
			">=1.23.16-rancher2-2 <1.23.16-rancher2-3": aciv5236,
			">=1.23.16-rancher2-3 <1.24.0-rancher0":    aciv5271,
			">=1.24.0-rancher0 <1.24.8-rancher1-1":     aciv523,
			">=1.24.8-rancher1-1 <1.24.9-rancher1-1":   aciv5234,
			">=1.24.9-rancher1-1 <1.24.13-rancher2-1":  aciv5235,
			">=1.24.13-rancher2-1 <1.24.13-rancher2-2": aciv5236,
			">=1.24.13-rancher2-2 <1.25.0-rancher0":    aciv5271,
			">=1.25.0-rancher0 <1.25.9-rancher2-1":     aciv5235,
			">=1.25.9-rancher2-1 <1.25.9-rancher2-2":   aciv5236,
			">=1.25.9-rancher2-2":                      aciv5271,
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
			">=1.15.12-rancher1-1 <1.16.1-rancher1-1":  nginxIngressV11512,
			">=1.16.10-rancher1-1 <1.17.0-rancher1-1":  nginxIngressV11512,
			">=1.17.6-rancher1-1 <1.19.16-rancher1-1":  nginxIngressV11512,
			">=1.19.16-rancher1-1 <1.19.16-rancher1-3": nginxIngressv0493,
			">=1.19.16-rancher1-3 <1.19.16-rancher1-4": nginxIngressv110Rancher2,
			">=1.19.16-rancher1-4 <1.19.16-rancher1-5": nginxIngressv110Rancher4,
			">=1.19.16-rancher1-5 <1.20.0-rancher0":    nginxIngressv120Rancher1,
			">=1.20.0-rancher0 <1.20.12-rancher1-1":    nginxIngressV11512,
			">=1.20.12-rancher1-1 <1.20.14-rancher2-2": nginxIngressv0493,
			">=1.20.14-rancher2-2 <1.20.15-rancher1-2": nginxIngressv110Rancher2,
			">=1.20.15-rancher1-2 <1.20.15-rancher1-3": nginxIngressv110Rancher4,
			">=1.20.15-rancher1-3 <1.21.0-rancher0":    nginxIngressv120Rancher1,
			">=1.21.0-rancher0 <1.21.6-rancher1-1":     nginxIngressv0481,
			">=1.21.6-rancher1-1 <1.21.8-rancher2-2":   nginxIngressv0493,
			">=1.21.8-rancher2-2 <1.21.9-rancher1-2":   nginxIngressv110Rancher2,
			">=1.21.9-rancher1-2 <1.21.12-rancher1-1":  nginxIngressv110Rancher4,
			">=1.21.12-rancher1-1 <1.22.0-rancher1-1":  nginxIngressv120Rancher1,
			">=1.22.0-rancher1-1 <1.22.5-rancher2-2":   nginxIngressv110,
			">=1.22.5-rancher2-2 <1.22.9-rancher1-1":   nginxIngressv110Rancher3,
			">=1.22.9-rancher1-1 <1.23.0-rancher1-1":   nginxIngressv120Rancher2,
			">=1.23.0-rancher1-1 <1.23.6-rancher1-1":   nginxIngressv110Rancher3,
			">=1.23.6-rancher1-1 <1.23.16-rancher1-1":  nginxIngressv120Rancher2,
			">=1.23.16-rancher1-1 <1.24.0-rancher1-1":  nginxIngressv151Rancher2,
			">=1.24.2-rancher1-1 <1.24.10-rancher1-1":  nginxIngressv120Rancher2,
			">=1.24.10-rancher1-1 <1.26.0-rancher1-1":  nginxIngressv151Rancher2,
			">=1.26.0-rancher1-1 <1.27.0-rancher1-1":   nginxIngressv170Rancher1,
			">=1.27.0-rancher1-1":                      nginxIngressv181Rancher1,
		},
		kdm.Nodelocal: {
			">=1.15.11-rancher0 <1.16.0-alpha":     nodelocalv115,
			">=1.16.8-rancher0 <1.17.0-alpha":      nodelocalv115,
			">=1.17.4-rancher0 <1.21.0-rancher1-1": nodelocalv115,
			">=1.21.0-rancher1-1":                  nodelocalv121,
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
		calicov319:           CalicoTemplateV3_19_0,
		calicov3211:          CalicoTemplateV3_21_1,
		calicov3225:          CalicoTemplateV3_22_5,
		calicov3225Rancher2:  CalicoTemplateV3_22_5Rancher2,
		calicov3241:          CalicoTemplateV3_24_1,
		calicov3250:          CalicoTemplateV3_25_0,
		calicov3261:          CalicoTemplateV3_26_1,

		flannelv115:          FlannelTemplateV115,
		flannelv116:          FlannelTemplateV116,
		flannelv0140:         FlannelTemplateV0_14_0,
		flannelv0140Rancher2: FlannelTemplateV0_14_0Rancher2,
		flannelv0192:         FlannelTemplateV0_19_2,
		flannelv0214:         FlannelTemplateV0_21_4,
		flannelv18:           FlannelTemplate,

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
		canalv319:                     CanalTemplateV3_19_0,
		canalv3211:                    CanalTemplateV3_21_1,
		canalv3225:                    CanalTemplateV3_22_5,
		canalv3225Rancher2:            CanalTemplateV3_22_5Rancher2,
		canalv3241:                    CanalTemplateV3_24_1,
		canalv3250:                    CanalTemplateV3_25_0,
		canalv3261:                    CanalTemplateV3_26_1,

		coreDnsv18:          CoreDNSTemplate,
		coreDnsv116:         CoreDNSTemplateV116,
		coreDnsv117:         CoreDNSTemplateV117,
		coreDnsv183:         CoreDNSTemplateV183,
		coreDnsv183Rancher2: CoreDNSTemplateV183Rancher2,

		kubeDnsv18:  KubeDNSTemplate,
		kubeDnsv116: KubeDNSTemplateV116,

		metricsServerv18:  MetricsServerTemplate,
		metricsServerv120: MetricsServerTemplateV0_4_1,
		metricsServerv050: MetricsServerTemplateV0_5_0,
		metricsServerv052: MetricsServerTemplateV0_5_2,
		metricsServerv061: MetricsServerTemplateV0_6_1,

		weavev18:  WeaveTemplate,
		weavev116: WeaveTemplateV116,
		weavev120: WeaveTemplateV120,
		weavev122: WeaveTemplateV122,

		aciv500:  AciTemplateV500,
		aciv523:  AciTemplateV523,
		aciv5234: AciTemplateV5234,
		aciv5235: AciTemplateV5235,
		aciv5236: AciTemplateV5236,
		aciv5271: AciTemplateV5271,

		nginxIngressv18:          NginxIngressTemplate,
		nginxIngressV115:         NginxIngressTemplateV0251Rancher1,
		nginxIngressV11512:       NginxIngressTemplateV0320Rancher1,
		nginxIngressv0481:        NginxIngressTemplateV0481Rancher1,
		nginxIngressv0493:        NginxIngressTemplateV0493Rancher1,
		nginxIngressv110:         NginxIngressTemplateV110Rancher1,
		nginxIngressv110Rancher2: NginxIngressTemplateV110Rancher2,
		nginxIngressv110Rancher3: NginxIngressTemplateV110Rancher3,
		nginxIngressv110Rancher4: NginxIngressTemplateV110Rancher4,
		nginxIngressv120Rancher1: NginxIngressTemplateV120Rancher1,
		nginxIngressv120Rancher2: NginxIngressTemplateV120Rancher2,
		nginxIngressv151Rancher2: NginxIngressTemplateV151Rancher2,
		nginxIngressv170Rancher1: NginxIngressTemplateV170Rancher1,
		nginxIngressv181Rancher1: NginxIngressTemplateV181Rancher1,

		nodelocalv115: NodelocalTemplateV115,
		nodelocalv121: NodelocalTemplateV121,
	}
}
