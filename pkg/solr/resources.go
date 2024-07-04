package solr

import (
	"github.com/pkg/errors"
	code2cloudv1deployslcstackk8smodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/solrkubernetes/stack/model"
	solrcontextstate "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/contextstate"
	solrnamespace "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/namespace"
	solrnetwork "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network"
	solroperator "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/operator"
	solroutputs "github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/outputs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	WorkspaceDir     string
	Input            *code2cloudv1deployslcstackk8smodel.SolrKubernetesStackInput
	KubernetesLabels map[string]string
}

func (resourceStack *ResourceStack) Resources(ctx *pulumi.Context) error {
	//load context config
	var ctxConfig, err = loadConfig(ctx, resourceStack)
	if err != nil {
		return errors.Wrap(err, "failed to initiate context config")
	}
	ctx = ctx.WithValue(solrcontextstate.Key, *ctxConfig)

	// Create the namespace resource
	ctx, err = solrnamespace.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace resource")
	}

	if err := solroperator.Resources(ctx); err != nil {
		return errors.Wrap(err, "failed to add solr-kubernetes resources")
	}

	ctx, err = solrnetwork.Resources(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to add solr network resources")
	}

	err = solroutputs.Export(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to export solr kubernetes outputs")
	}

	return nil
}
