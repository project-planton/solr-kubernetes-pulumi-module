package virtualservice

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/network/dns/zone"
	"github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/netutils/port"
	"github.com/plantoncloud/solr-kubernetes-pulumi-blueprint/pkg/solr/network/ingress/netutils/service"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	ingressnamespace "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(ctx *pulumi.Context) error {
	if err := addVirtualService(ctx); err != nil {
		return errors.Wrap(err, "failed to add external virtual service")
	}
	return nil
}

func addVirtualService(ctx *pulumi.Context) error {
	i := extractInput(ctx)
	var virtualServiceObject = buildVirtualServiceObject(i)

	resourceName := fmt.Sprintf("virtual-service-%s", virtualServiceObject.Name)
	manifestPath := filepath.Join(i.workspaceDir, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, virtualServiceObject); err != nil {
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

func buildVirtualServiceObject(i *input) *v1beta1.VirtualService {
	gatewayName := getGatewayName(i.environmentName, i.endpointDomainName)
	return &v1beta1.VirtualService{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "VirtualService",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      i.resourceName,
			Namespace: i.namespaceName,
		},
		Spec: networkingv1beta1.VirtualService{
			Gateways: []string{fmt.Sprintf("%s/%s", ingressnamespace.Name, gatewayName)},
			Hosts:    []string{i.externalHostname},
			Http: []*networkingv1beta1.HTTPRoute{{
				Name: i.resourceId,
				Route: []*networkingv1beta1.HTTPRouteDestination{
					{
						Destination: &networkingv1beta1.Destination{
							Host: service.GetKubeServiceNameFqdn(i.resourceName, i.namespaceName),
							Port: &networkingv1beta1.PortSelector{
								Number: port.SolrCloudCommonServicePort,
							},
						},
					},
				},
			}},
		},
	}
}

func getGatewayName(environmentName, endpointDomainName string) string {
	envDomain := fmt.Sprintf("%s.%s", environmentName, endpointDomainName)
	return fmt.Sprintf(zone.GetZoneName(envDomain))
}
