package solr

import (
	"github.com/pkg/errors"
	environmentblueprinthostnames "github.com/plantoncloud/environment-pulumi-blueprint/pkg/gcpgke/endpointdomains/hostnames"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/kubernetes/pulumikubernetesprovider"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	solrnetutilshostname "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/netutils/hostname"
	solrnetutilsservice "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/netutils/service"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func loadConfig(ctx *pulumi.Context, resourceStack *ResourceStack) (*solrcontextstate.ContextState, error) {

	kubernetesProvider, err := pulumikubernetesprovider.GetWithStackCredentials(ctx, resourceStack.Input.CredentialsInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup kubernetes provider")
	}

	var resourceId = resourceStack.Input.ResourceInput.Metadata.Id
	var resourceName = resourceStack.Input.ResourceInput.Metadata.Name
	var environmentInfo = resourceStack.Input.ResourceInput.Spec.EnvironmentInfo
	var isIngressEnabled = false

	if resourceStack.Input.ResourceInput.Spec.Ingress != nil {
		isIngressEnabled = resourceStack.Input.ResourceInput.Spec.Ingress.IsEnabled
	}

	var endpointDomainName = ""
	var envDomainName = ""
	var ingressType = kubernetesworkloadingresstype.KubernetesWorkloadIngressType_unspecified
	var internalHostname = ""
	var externalHostname = ""

	if isIngressEnabled {
		endpointDomainName = resourceStack.Input.ResourceInput.Spec.Ingress.EndpointDomainName
		envDomainName = environmentblueprinthostnames.GetExternalEnvHostname(environmentInfo.EnvironmentName, endpointDomainName)
		ingressType = resourceStack.Input.ResourceInput.Spec.Ingress.IngressType

		internalHostname = solrnetutilshostname.GetInternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
		externalHostname = solrnetutilshostname.GetExternalHostname(resourceId, environmentInfo.EnvironmentName, endpointDomainName)
	}

	return &solrcontextstate.ContextState{
		Spec: &solrcontextstate.Spec{
			KubeProvider:           kubernetesProvider,
			ResourceId:             resourceId,
			ResourceName:           resourceName,
			Labels:                 resourceStack.KubernetesLabels,
			WorkspaceDir:           resourceStack.WorkspaceDir,
			NamespaceName:          resourceId,
			EnvironmentInfo:        resourceStack.Input.ResourceInput.Spec.EnvironmentInfo,
			IsIngressEnabled:       isIngressEnabled,
			IngressType:            ingressType,
			EndpointDomainName:     endpointDomainName,
			EnvDomainName:          envDomainName,
			InternalHostname:       internalHostname,
			ExternalHostname:       externalHostname,
			KubeServiceName:        solrnetutilsservice.GetKubeServiceName(resourceName),
			KubeLocalEndpoint:      solrnetutilsservice.GetKubeServiceNameFqdn(resourceName, resourceId),
			SolrContainerSpec:      resourceStack.Input.ResourceInput.Spec.SolrContainer,
			ZookeeperContainerSpec: resourceStack.Input.ResourceInput.Spec.ZookeeperContainer,
		},
		Status: &solrcontextstate.Status{},
	}, nil
}
