package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/solrkubernetes/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	NamespaceOutputName               = "namespace"
	ServiceOutputName                 = "service"
	PortForwardCommandOutputName      = "port-forward-command"
	KubeEndpointOutputName            = "kube-endpoint"
	IngressExternalHostnameOutputName = "ingress-external-hostname"
	IngressInternalHostnameOutputName = "ingress-internal-hostname"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.SolrKubernetesStackInput) *model.SolrKubernetesStackOutputs {
	return &model.SolrKubernetesStackOutputs{
		Namespace:          autoapistackoutput.GetVal(pulumiOutputs, NamespaceOutputName),
		KubeEndpoint:       autoapistackoutput.GetVal(pulumiOutputs, KubeEndpointOutputName),
		Service:            autoapistackoutput.GetVal(pulumiOutputs, ServiceOutputName),
		PortForwardCommand: autoapistackoutput.GetVal(pulumiOutputs, PortForwardCommandOutputName),
		ExternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressExternalHostnameOutputName),
		InternalHostname:   autoapistackoutput.GetVal(pulumiOutputs, IngressInternalHostnameOutputName),
	}
}
