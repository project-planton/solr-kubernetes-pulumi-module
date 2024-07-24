package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/solrkubernetes/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/kubernetes/enums/kubernetesworkloadingresstype"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/pulumikubernetesprovider"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmetav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *model.SolrKubernetesStackInput
	Labels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	//create kubernetes-provider from the credential in the stack-input
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesClusterCredential(ctx,
		s.Input.KubernetesClusterCredential)
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	//create a new descriptive variable for the api-resource in the input.
	solrKubernetes := s.Input.ApiResource

	//decide on the name of the namespace
	namespaceName := solrKubernetes.Metadata.Id

	//create namespace resource
	createdNamespace, err := kubernetescorev1.NewNamespace(ctx,
		namespaceName,
		&kubernetescorev1.NamespaceArgs{
			Metadata: kubernetesmetav1.ObjectMetaPtrInput(
				&kubernetesmetav1.ObjectMetaArgs{
					Name:   pulumi.String(namespaceName),
					Labels: pulumi.ToStringMap(s.Labels),
				}),
		},
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s namespace", namespaceName)
	}

	//export name of the namespace
	ctx.Export(NamespaceOutputName, createdNamespace.Metadata.Name())

	//install the solr helm-chart
	if err := s.solrCloud(ctx, createdNamespace); err != nil {
		return errors.Wrap(err, "failed to create helm-chart resources")
	}

	//export kube-port-forward command
	ctx.Export(PortForwardCommandOutputName, pulumi.Sprintf(
		"kubectl port-forward -n %s service/%s 8080:8080",
		namespaceName, solrKubernetes.Metadata.Name))

	//no ingress resources required when ingress is not enabled
	if !solrKubernetes.Spec.Ingress.IsEnabled || solrKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return nil
	}

	//depending on the ingress-type in the input, create either istio-ingress resources or
	//create load-balancer resources
	switch solrKubernetes.Spec.Ingress.IngressType {
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_load_balancer:
		if err := s.loadBalancerIngress(ctx, createdNamespace); err != nil {
			return errors.Wrap(err, "failed to create load-balancer ingress resources")
		}
	case kubernetesworkloadingresstype.KubernetesWorkloadIngressType_ingress_controller:
		if err := s.istioIngress(ctx, createdNamespace); err != nil {
			return errors.Wrap(err, "failed to create istio ingress resources")
		}
	}

	//export ingress hostnames
	ctx.Export(IngressExternalHostnameOutputName, pulumi.Sprintf("%s.%s",
		solrKubernetes.Metadata.Id, solrKubernetes.Spec.Ingress.EndpointDomainName))
	ctx.Export(IngressInternalHostnameOutputName, pulumi.Sprintf("%s-internal.%s",
		solrKubernetes.Metadata.Id, solrKubernetes.Spec.Ingress.EndpointDomainName))

	return nil
}
