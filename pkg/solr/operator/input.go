package operator

import (
	code2cloudv1deployslcmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/solrkubernetes/model"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	pulk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var solrModules = []string{
	"jaegertracer-configurator",
	"ltr",
}

type input struct {
	workspaceDir           string
	namespaceName          string
	namespace              *pulk8scv1.Namespace
	solrContainerSpec      *code2cloudv1deployslcmodel.SolrKubernetesSpecSolrContainerSpec
	zookeeperContainerSpec *code2cloudv1deployslcmodel.SolrKubernetesSpecZookeeperContainerSpec
	labels                 map[string]string
	resourceName           string
}

func extractInput(ctx *pulumi.Context) *input {
	var contextState = ctx.Value(solrcontextstate.Key).(solrcontextstate.ContextState)

	return &input{
		namespaceName:          contextState.Spec.NamespaceName,
		labels:                 contextState.Spec.Labels,
		workspaceDir:           contextState.Spec.WorkspaceDir,
		namespace:              contextState.Status.AddedResources.Namespace,
		solrContainerSpec:      contextState.Spec.SolrContainerSpec,
		zookeeperContainerSpec: contextState.Spec.ZookeeperContainerSpec,
		resourceName:           contextState.Spec.ResourceName,
	}
}
