package ingress

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/kubernetesworkloadingresstype"
	solristio "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/istio"
	solrloadbalancer "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/loadbalancer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context) (newCtx *pulumi.Context, err error) {
	i := extractInput(ctx)
	switch i.ingressType {
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_load_balancer:
		ctx, err = solrloadbalancer.Resources(ctx)
		if err != nil {
			return ctx, errors.Wrap(err, "failed to add load balancer resources")
		}
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_ingress_controller:
		if err = solristio.Resources(ctx); err != nil {
			return ctx, errors.Wrap(err, "failed to add istio resources")
		}
	}
	return ctx, nil
}
