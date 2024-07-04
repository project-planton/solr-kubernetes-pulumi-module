package loadbalancer

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/enums/kubernetesprovider"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	kubeProvider kubernetesprovider.KubernetesProvider
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		kubeProvider: ctxConfig.Spec.EnvironmentInfo.KubernetesProvider,
	}
}
