package namespace

import (
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	namespaceName string
	labels        map[string]string
	kubeProvider  *pulumikubernetes.Provider
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		namespaceName: contextState.Spec.NamespaceName,
		labels:        contextState.Spec.Labels,
		kubeProvider:  contextState.Spec.KubeProvider,
	}
}
