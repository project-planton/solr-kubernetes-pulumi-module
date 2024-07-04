package solr

import (
	"context"
	"github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/outputs"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	solrkubernetesstatemodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/solrkubernetes/model"
	solrkubernetesstackmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/solrkubernetes/stack/model"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *solrkubernetesstackmodel.SolrKubernetesStackInput) (*solrkubernetesstatemodel.SolrKubernetesStatusStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}

	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *solrkubernetesstackmodel.SolrKubernetesStackInput) *solrkubernetesstatemodel.SolrKubernetesStatusStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &solrkubernetesstatemodel.SolrKubernetesStatusStackOutputs{}
	}
	return &solrkubernetesstatemodel.SolrKubernetesStatusStackOutputs{
		Namespace:          backend.GetVal(stackOutput, outputs.GetNamespaceNameOutputName()),
		KubeEndpoint:       backend.GetVal(stackOutput, outputs.GetKubeEndpointOutputName()),
		Service:            backend.GetVal(stackOutput, outputs.GetKubeServiceNameOutputName()),
		PortForwardCommand: backend.GetVal(stackOutput, outputs.GetKubePortForwardCommandOutputName()),
		ExternalHostname:   backend.GetVal(stackOutput, outputs.GetExternalClusterHostnameOutputName()),
		InternalHostname:   backend.GetVal(stackOutput, outputs.GetInternalClusterHostnameOutputName()),
	}
}
