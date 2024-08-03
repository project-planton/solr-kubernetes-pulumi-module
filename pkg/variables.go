package pkg

var vars = struct {
	IstioIngressNamespace      string
	IstioIngressSelectorLabels map[string]string
	SolrCloudSolrModules       []string
}{
	IstioIngressNamespace: "istio-ingress",
	IstioIngressSelectorLabels: map[string]string{
		"app":   "gateway",
		"istio": "gateway",
	},
	SolrCloudSolrModules: []string{
		"jaegertracer-configurator",
		"ltr",
	},
}
