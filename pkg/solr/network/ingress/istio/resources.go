package istio

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/istio/virtualservice"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) error {
	if err := virtualservice.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add virtual resources")
	}
	return nil
}
