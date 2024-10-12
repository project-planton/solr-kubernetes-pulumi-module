package main

import (
	solrkubernetesv1 "buf.build/gen/go/project-planton/apis/protocolbuffers/go/project/planton/provider/kubernetes/solrkubernetes/v1"
	"github.com/pkg/errors"
	"github.com/project-planton/pulumi-module-golang-commons/pkg/stackinput"
	"github.com/project-planton/solr-kubernetes-pulumi-module/pkg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &solrkubernetesv1.SolrKubernetesStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return pkg.Resources(ctx, stackInput)
	})
}
