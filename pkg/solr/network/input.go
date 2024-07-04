package network

import (
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type input struct {
	isIngressEnabled   bool
	endpointDomainName string
}

func extractInput(ctx *pulumi.Context) *input {
	var ctxConfig = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		isIngressEnabled:   ctxConfig.Spec.IsIngressEnabled,
		endpointDomainName: ctxConfig.Spec.EndpointDomainName,
	}
}
