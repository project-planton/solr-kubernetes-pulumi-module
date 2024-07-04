package virtualservice

import (
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	resourceId         string
	resourceName       string
	namespaceName      string
	workspaceDir       string
	namespace          *kubernetescorev1.Namespace
	externalHostname   string
	internalHostname   string
	kubeEndpoint       string
	environmentName    string
	endpointDomainName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxState = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		resourceId:         ctxState.Spec.ResourceId,
		resourceName:       ctxState.Spec.ResourceName,
		workspaceDir:       ctxState.Spec.WorkspaceDir,
		namespaceName:      ctxState.Spec.NamespaceName,
		namespace:          ctxState.Status.AddedResources.Namespace,
		externalHostname:   ctxState.Spec.ExternalHostname,
		internalHostname:   ctxState.Spec.InternalHostname,
		kubeEndpoint:       ctxState.Spec.KubeLocalEndpoint,
		environmentName:    ctxState.Spec.EnvironmentInfo.EnvironmentName,
		endpointDomainName: ctxState.Spec.EndpointDomainName,
	}
}
