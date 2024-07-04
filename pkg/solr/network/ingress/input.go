package ingress

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	ingressType kubernetesworkloadingresstype.KubernetesWorkloadIngressType
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		ingressType: ctxConfig.Spec.IngressType,
	}
}
