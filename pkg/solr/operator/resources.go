package operator

import (
	"fmt"
	"path/filepath"

	"github.com/apache/solr-operator/api/v1beta1"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	kubernetesv1model "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8sapimachineryv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func Resources(ctx *pulumi.Context) error {
	if err := addSolrKubernetes(ctx); err != nil {
		return errors.Wrap(err, "failed to add solr kubernetes")
	}
	return nil
}

func addSolrKubernetes(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	solrKubernetesObject, _ := buildSolrKubernetesObject(i)
	resourceName := fmt.Sprintf("solr-kubernetes-%s", solrKubernetesObject.Name)
	manifestPath := filepath.Join(i.workspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, solrKubernetesObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{
		File: manifestPath,
	}, pulumi.DependsOn([]pulumi.Resource{i.namespace}), pulumi.Parent(i.namespace))
	if err != nil {
		return errors.Wrap(err, "failed to add virtual-service manifest")
	}
	return nil
}

/*
apiVersion: solr.apache.org/v1beta1
kind: SolrKubernetes
metadata:

	name: solr-kubernetes-name
	namespace: solr-kubernetes-namespace

spec:

	dataStorage:
	  persistent:
	    reclaimPolicy: Delete
	    pvcTemplate:
	      spec:
	        resources:
	          requests:
	            storage: "1Gi"
	replicas: 1
	solrImage:
	  tag: 8.7.0
	solrJavaMem: "-Xms1g -Xmx3g"
	solrModules:
	  - jaegertracer-configurator
	  - ltr
	customSolrKubeOptions:
	  podOptions:
	    resources:
	      limits:
	        memory: "1G"
	      requests:
	        cpu: "65m"
	        memory: "156Mi"
	zookeeperRef:
	  provided:
	    chroot: "/this/will/be/auto/created"
	    persistence:
	      spec:
	        resources:
	          requests:
	            storage: "1Gi"
	    replicas: 1
	    zookeeperPodPolicy:
	      resources:
	        limits:
	          memory: "1G"
	        requests:
	          cpu: "65m"
	          memory: "156Mi"
	solrOpts: "-Dsolr.autoSoftCommit.maxTime=10000"
	solrGCTune: "-XX:SurvivorRatio=4 -XX:TargetSurvivorRatio=90 -XX:MaxTenuringThreshold=8"
*/
func buildSolrKubernetesObject(i *input) (*v1beta1.SolrCloud, error) {
	solrDiskSize, err := resource.ParseQuantity(i.solrContainerSpec.DiskSize)
	if err != nil {
		return nil, errors.Wrapf(err, "solr-disk-size value %s is invalid",
			i.solrContainerSpec.DiskSize)
	}
	zookeeperDiskSize, err := resource.ParseQuantity(i.zookeeperContainerSpec.DiskSize)
	if err != nil {
		return nil, errors.Wrapf(err, "solr-disk-size value %s is invalid",
			i.zookeeperContainerSpec.DiskSize)
	}
	solrContainerResources, err := getContainerResources(i.solrContainerSpec.Resources)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse solr container resources")
	}
	zookeeperContainerResources, err := getContainerResources(i.zookeeperContainerSpec.Resources)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse zookeeper container resources")
	}
	return &v1beta1.SolrCloud{
		TypeMeta: k8sapimachineryv1.TypeMeta{
			APIVersion: "solr.apache.org/v1beta1",
			Kind:       "SolrKubernetes",
		},
		ObjectMeta: k8sapimachineryv1.ObjectMeta{
			Name:      i.resourceName,
			Namespace: i.namespaceName,
			Labels:    i.labels,
		},
		Spec: v1beta1.SolrCloudSpec{
			Replicas: pointer.Int32(i.solrContainerSpec.Replicas),
			SolrImage: &v1beta1.ContainerImage{
				Repository: i.solrContainerSpec.Image.Repo,
				Tag:        i.solrContainerSpec.Image.Tag,
			},
			SolrJavaMem: i.solrContainerSpec.Config.JavaMem,
			SolrOpts:    i.solrContainerSpec.Config.Opts,
			SolrGCTune:  i.solrContainerSpec.Config.GarbageCollectionTuning,
			SolrModules: solrModules,
			CustomSolrKubeOptions: v1beta1.CustomSolrKubeOptions{
				PodOptions: &v1beta1.PodOptions{
					Resources: *solrContainerResources,
				},
			},
			StorageOptions: v1beta1.SolrDataStorageOptions{
				PersistentStorage: &v1beta1.SolrPersistentDataStorageOptions{
					VolumeReclaimPolicy: "Delete",
					PersistentVolumeClaimTemplate: v1beta1.PersistentVolumeClaimTemplate{
						Spec: k8scorev1.PersistentVolumeClaimSpec{
							Resources: k8scorev1.ResourceRequirements{
								Requests: k8scorev1.ResourceList{
									"storage": solrDiskSize,
								},
							},
						},
					},
				},
			},
			ZookeeperRef: &v1beta1.ZookeeperRef{
				ProvidedZookeeper: &v1beta1.ZookeeperSpec{
					Replicas: pointer.Int32(i.zookeeperContainerSpec.Replicas),
					Persistence: &v1beta1.ZKPersistence{
						PersistentVolumeClaimSpec: k8scorev1.PersistentVolumeClaimSpec{
							Resources: k8scorev1.ResourceRequirements{
								Requests: k8scorev1.ResourceList{
									"storage": zookeeperDiskSize,
								},
							},
						},
					},
					ZookeeperPod: v1beta1.ZookeeperPodPolicy{
						Resources: *zookeeperContainerResources,
					},
				},
			},
		},
	}, nil
}

func getContainerResources(solrKubernetesSpecInput *kubernetesv1model.ContainerResources) (*k8scorev1.ResourceRequirements, error) {
	cpuLimits, err := resource.ParseQuantity(solrKubernetesSpecInput.Limits.Cpu)
	if err != nil {
		return nil, errors.Wrapf(err, "cpu limits value %s is invalid",
			solrKubernetesSpecInput.Limits.Cpu)
	}
	memoryLimits, err := resource.ParseQuantity(solrKubernetesSpecInput.Limits.Memory)
	if err != nil {
		return nil, errors.Wrapf(err, "memory limits value %s is invalid",
			solrKubernetesSpecInput.Limits.Memory)
	}
	cpuRequests, err := resource.ParseQuantity(solrKubernetesSpecInput.Requests.Cpu)
	if err != nil {
		return nil, errors.Wrapf(err, "cpu requests value %s is invalid",
			solrKubernetesSpecInput.Requests.Cpu)
	}
	memoryRequests, err := resource.ParseQuantity(solrKubernetesSpecInput.Requests.Memory)
	if err != nil {
		return nil, errors.Wrapf(err, "memory requests value %s is invalid",
			solrKubernetesSpecInput.Requests.Memory)
	}
	return &k8scorev1.ResourceRequirements{
		Limits: k8scorev1.ResourceList{
			k8scorev1.ResourceCPU:    cpuLimits,
			k8scorev1.ResourceMemory: memoryLimits,
		},
		Requests: k8scorev1.ResourceList{
			k8scorev1.ResourceCPU:    cpuRequests,
			k8scorev1.ResourceMemory: memoryRequests,
		},
	}, nil
}
