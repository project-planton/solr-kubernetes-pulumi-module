package locals

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubernetes/solrkubernetes/model"
	"github.com/plantoncloud/solr-kubernetes-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	IngressCertClusterIssuerName string
	IngressCertSecretName        string
	IngressExternalHostname      string
	IngressHostnames             []string
	IngressInternalHostname      string
	KubePortForwardCommand       string
	KubeServiceFqdn              string
	KubeServiceName              string
	Namespace                    string
	SolrKubernetes               *model.SolrKubernetes
)

// Initializer will be invoked by the stack-job-runner sdk before the pulumi operations are executed.
func Initializer(ctx *pulumi.Context, stackInput *model.SolrKubernetesStackInput) {
	//assign value for the local variable to make it available across the project
	SolrKubernetes = stackInput.ApiResource

	solrKubernetes := stackInput.ApiResource

	//decide on the namespace
	Namespace = solrKubernetes.Metadata.Id

	KubeServiceName = solrKubernetes.Metadata.Name

	//export kubernetes service name
	ctx.Export(outputs.Service, pulumi.String(KubeServiceName))

	KubeServiceFqdn = fmt.Sprintf(
		"%s-solrcloud-common.%s.svc.cluster.local",
		solrKubernetes.Metadata.Name, Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpoint, pulumi.String(KubeServiceFqdn))

	KubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s 8080:8080",
		Namespace, solrKubernetes.Metadata.Name)

	//export kube-port-forward command
	ctx.Export(outputs.KubePortForwardCommand, pulumi.String(KubePortForwardCommand))

	if solrKubernetes.Spec.Ingress == nil ||
		!solrKubernetes.Spec.Ingress.IsEnabled ||
		solrKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return
	}

	IngressExternalHostname = fmt.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
		solrKubernetes.Spec.Ingress.EndpointDomainName)

	IngressInternalHostname = fmt.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
		solrKubernetes.Spec.Ingress.EndpointDomainName)

	IngressHostnames = []string{
		IngressExternalHostname,
		IngressInternalHostname,
	}

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostname, pulumi.String(IngressExternalHostname))
	ctx.Export(outputs.IngressInternalHostname, pulumi.String(IngressInternalHostname))

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	IngressCertClusterIssuerName = solrKubernetes.Spec.Ingress.EndpointDomainName

	IngressCertSecretName = solrKubernetes.Metadata.Id
}
