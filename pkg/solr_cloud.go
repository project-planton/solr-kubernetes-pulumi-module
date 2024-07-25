package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/solroperator/solr/v1beta1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) solrCloud(ctx *pulumi.Context,
	createdNamespace *kubernetescorev1.Namespace) error {
	//create a variable with descriptive name for the api-resource from input
	solrKubernetes := s.Input.ApiResource

	_, err := v1beta1.NewSolrCloud(ctx, "solr-cloud", &v1beta1.SolrCloudArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String(solrKubernetes.Metadata.Id),
			Namespace: createdNamespace.Metadata.Name(),
			Labels:    pulumi.ToStringMap(s.Labels),
		},
		Spec: v1beta1.SolrCloudSpecArgs{
			Replicas: pulumi.Int(solrKubernetes.Spec.SolrContainer.Replicas),
			SolrImage: v1beta1.SolrCloudSpecSolrImageArgs{
				Repository: pulumi.String(solrKubernetes.Spec.SolrContainer.Image.Repo),
				Tag:        pulumi.String(solrKubernetes.Spec.SolrContainer.Image.Tag),
			},
			SolrJavaMem: pulumi.String(solrKubernetes.Spec.SolrContainer.Config.JavaMem),
			SolrOpts:    pulumi.String(solrKubernetes.Spec.SolrContainer.Config.Opts),
			SolrGCTune:  pulumi.String(solrKubernetes.Spec.SolrContainer.Config.GarbageCollectionTuning),
			SolrModules: pulumi.ToStringArray([]string{
				"jaegertracer-configurator",
				"ltr",
			}),
			CustomSolrKubeOptions: v1beta1.SolrCloudSpecCustomSolrKubeOptionsArgs{
				PodOptions: v1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsArgs{
					Resources: v1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsResourcesArgs{
						Limits: pulumi.ToMap(map[string]interface{}{
							"cpu":    solrKubernetes.Spec.SolrContainer.Resources.Limits.Cpu,
							"memory": solrKubernetes.Spec.SolrContainer.Resources.Limits.Memory,
						}),
						Requests: pulumi.ToMap(map[string]interface{}{
							"cpu":    solrKubernetes.Spec.SolrContainer.Resources.Requests.Cpu,
							"memory": solrKubernetes.Spec.SolrContainer.Resources.Requests.Memory,
						}),
					},
				},
			},
			DataStorage: v1beta1.SolrCloudSpecDataStorageArgs{
				Ephemeral: nil,
				Persistent: v1beta1.SolrCloudSpecDataStoragePersistentArgs{
					ReclaimPolicy: pulumi.String("Delete"),
					PvcTemplate: v1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateArgs{
						Spec: v1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateSpecArgs{
							Resources: v1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateSpecResourcesArgs{
								Requests: pulumi.ToMap(map[string]interface{}{
									"storage": solrKubernetes.Spec.SolrContainer.DiskSize,
								}),
							},
						},
					},
				},
			},
			ZookeeperRef: v1beta1.SolrCloudSpecZookeeperRefArgs{
				Provided: v1beta1.SolrCloudSpecZookeeperRefProvidedArgs{
					Replicas: pulumi.Int(solrKubernetes.Spec.ZookeeperContainer.Replicas),
					Persistence: v1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceArgs{
						Spec: v1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceSpecArgs{
							Resources: v1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceSpecResourcesArgs{
								Requests: pulumi.ToMap(map[string]interface{}{
									"storage": solrKubernetes.Spec.ZookeeperContainer.DiskSize,
								}),
							},
						},
					},
					ZookeeperPodPolicy: v1beta1.SolrCloudSpecZookeeperRefProvidedZookeeperPodPolicyArgs{
						Resources: v1beta1.SolrCloudSpecZookeeperRefProvidedZookeeperPodPolicyResourcesArgs{
							Limits: pulumi.ToMap(map[string]interface{}{
								"cpu":    solrKubernetes.Spec.ZookeeperContainer.Resources.Limits.Cpu,
								"memory": solrKubernetes.Spec.ZookeeperContainer.Resources.Limits.Memory,
							}),
							Requests: pulumi.ToMap(map[string]interface{}{
								"cpu":    solrKubernetes.Spec.ZookeeperContainer.Resources.Requests.Cpu,
								"memory": solrKubernetes.Spec.ZookeeperContainer.Resources.Requests.Memory,
							}),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to create solr-cloud resource")
	}
	return nil
}
