package pkg

import (
	solrkubernetesv1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/kubernetes/solrkubernetes/v1"
	"fmt"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/provider/kubernetes/kuberneteslabelkeys"
	"github.com/project-planton/solr-kubernetes-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	IngressCertClusterIssuerName string
	IngressCertSecretName        string
	IngressExternalHostname      string
	IngressHostnames             []string
	IngressInternalHostname      string
	KubePortForwardCommand       string
	KubeServiceFqdn              string
	KubeServiceName              string
	Namespace                    string
	SolrKubernetes               *solrkubernetesv1.SolrKubernetes
	Labels                       map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *solrkubernetesv1.SolrKubernetesStackInput) *Locals {
	locals := &Locals{}
	//assign value for the locals variable to make it available across the project
	locals.SolrKubernetes = stackInput.Target

	solrKubernetes := stackInput.Target

	locals.Labels = map[string]string{
		kuberneteslabelkeys.Environment:  stackInput.Target.Spec.EnvironmentInfo.EnvId,
		kuberneteslabelkeys.Organization: stackInput.Target.Spec.EnvironmentInfo.OrgId,
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceId:   stackInput.Target.Metadata.Id,
		kuberneteslabelkeys.ResourceKind: "solr_kubernetes",
	}

	//decide on the namespace
	locals.Namespace = solrKubernetes.Metadata.Id

	ctx.Export(outputs.Namespace, pulumi.String(locals.Namespace))

	locals.KubeServiceName = fmt.Sprintf("%s-solrcloud-common", solrKubernetes.Metadata.Name)

	//export kubernetes service name
	ctx.Export(outputs.Service, pulumi.String(locals.KubeServiceName))

	locals.KubeServiceFqdn = fmt.Sprintf(
		"%s.%s.svc.cluster.local", locals.KubeServiceName, locals.Namespace)

	//export kubernetes endpoint
	ctx.Export(outputs.KubeEndpoint, pulumi.String(locals.KubeServiceFqdn))

	locals.KubePortForwardCommand = fmt.Sprintf("kubectl port-forward -n %s service/%s 8080:8080",
		locals.Namespace, solrKubernetes.Metadata.Name)

	//export kube-port-forward command
	ctx.Export(outputs.KubePortForwardCommand, pulumi.String(locals.KubePortForwardCommand))

	if solrKubernetes.Spec.Ingress == nil ||
		!solrKubernetes.Spec.Ingress.IsEnabled ||
		solrKubernetes.Spec.Ingress.EndpointDomainName == "" {
		return locals
	}

	locals.IngressExternalHostname = fmt.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
		solrKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressInternalHostname = fmt.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
		solrKubernetes.Spec.Ingress.EndpointDomainName)

	locals.IngressHostnames = []string{
		locals.IngressExternalHostname,
		locals.IngressInternalHostname,
	}

	//export ingress hostnames
	ctx.Export(outputs.IngressExternalHostname, pulumi.String(locals.IngressExternalHostname))
	ctx.Export(outputs.IngressInternalHostname, pulumi.String(locals.IngressInternalHostname))

	//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
	//this is typically taken care of by the kubernetes cluster administrator.
	//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
	//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
	//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
	locals.IngressCertClusterIssuerName = solrKubernetes.Spec.Ingress.EndpointDomainName

	locals.IngressCertSecretName = fmt.Sprintf("cert-%s", solrKubernetes.Metadata.Id)

	return locals
}
