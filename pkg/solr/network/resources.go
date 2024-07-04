package network

import (
	"github.com/pkg/errors"
	solringress "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (*pulumi.Context, error) {
	i := extractInput(ctx)
	if !i.isIngressEnabled || i.endpointDomainName == "" {
		return ctx, nil
	}
	if ctx, err := solringress.Resources(ctx); err != nil {
		return ctx, errors.Wrap(err, "failed to add network resources")
	}
	return ctx, nil
}
