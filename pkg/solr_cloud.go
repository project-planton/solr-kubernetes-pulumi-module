package pkg

import (
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) solrCloud(ctx *pulumi.Context, createdNamespace *kubernetescorev1.Namespace) error {
	return nil
}
