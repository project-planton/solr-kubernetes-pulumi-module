package outputs

import (
	pulumicommonsloadbalancerservice "github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/kubernetes/loadbalancer/service"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceId                    string
	resourceName                  string
	environmentName               string
	endpointDomainName            string
	namespaceName                 string
	externalLoadBalancerIpAddress string
	internalLoadBalancerIpAddress string
	internalHostname              string
	externalHostname              string
	kubeServiceName               string
	kubeLocalEndpoint             string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)
	var externalLoadBalancerIpAddress = ""
	var internalLoadBalancerIpAddress = ""

	if ctxConfig.Status.AddedResources.LoadBalancerExternalService != nil {
		externalLoadBalancerIpAddress = pulumicommonsloadbalancerservice.GetIpAddress(ctxConfig.Status.AddedResources.LoadBalancerExternalService)
	}

	if ctxConfig.Status.AddedResources.LoadBalancerInternalService != nil {
		internalLoadBalancerIpAddress = pulumicommonsloadbalancerservice.GetIpAddress(ctxConfig.Status.AddedResources.LoadBalancerExternalService)
	}

	return &input{
		resourceId:                    ctxConfig.Spec.ResourceId,
		resourceName:                  ctxConfig.Spec.ResourceName,
		environmentName:               ctxConfig.Spec.EnvironmentInfo.EnvironmentName,
		endpointDomainName:            ctxConfig.Spec.EndpointDomainName,
		namespaceName:                 ctxConfig.Spec.NamespaceName,
		externalLoadBalancerIpAddress: externalLoadBalancerIpAddress,
		internalLoadBalancerIpAddress: internalLoadBalancerIpAddress,
		internalHostname:              ctxConfig.Spec.InternalHostname,
		externalHostname:              ctxConfig.Spec.ExternalHostname,
		kubeServiceName:               ctxConfig.Spec.KubeServiceName,
		kubeLocalEndpoint:             ctxConfig.Spec.KubeLocalEndpoint,
	}
}
