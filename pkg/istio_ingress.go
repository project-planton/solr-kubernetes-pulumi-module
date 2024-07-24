package pkg

import (
	"github.com/pkg/errors"
	certmanagerv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/certmanager/certmanager/v1"
	istiov1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/istio/networking/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	v1 "istio.io/api/networking/v1"
)

const (
	IstioIngressNamespace = "istio-ingress"
)

func (s *ResourceStack) istioIngress(ctx *pulumi.Context, createdNamespace *kubernetescorev1.Namespace) error {
	solrKubernetes := s.Input.ApiResource
	addedCertificate, err := certmanagerv1.NewCertificate(ctx,
		"ingress-certificate",
		&certmanagerv1.CertificateArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(solrKubernetes.Metadata.Id),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(s.Labels),
			},
			Spec: certmanagerv1.CertificateSpecArgs{
				DnsNames: pulumi.StringArray{
					pulumi.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
						solrKubernetes.Spec.Ingress.EndpointDomainName),
					pulumi.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
						solrKubernetes.Spec.Ingress.EndpointDomainName),
				},
				SecretName: nil,
				IssuerRef: certmanagerv1.CertificateSpecIssuerRefArgs{
					Kind: pulumi.String("ClusterIssuer"),
					//note: a ClusterIssuer resource should have already exist on the kubernetes-cluster.
					//this is typically taken care of by the kubernetes cluster administrator.
					//if the kubernetes-cluster is created using Planton Cloud, then the cluster-issuer name will be
					//same as the ingress-domain-name as long as the same ingress-domain-name is added to the list of
					//ingress-domain-names for the GkeCluster/EksCluster/AksCluster spec.
					Name: pulumi.String(solrKubernetes.Spec.Ingress.EndpointDomainName),
				},
			},
		})
	if err != nil {
		return errors.Wrap(err, "error creating certificate")
	}

	_, err = istiov1.NewGateway(ctx,
		solrKubernetes.Metadata.Id,
		&istiov1.GatewayArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(solrKubernetes.Metadata.Id),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(s.Labels),
			},
			Spec: istiov1.GatewaySpecArgs{
				//the selector labels map should match the desired istio-ingress deployment.
				Selector: pulumi.StringMap{
					"app":   pulumi.String("istio-ingress"),
					"istio": pulumi.String("ingress"),
				},
				Servers: istiov1.GatewaySpecServersArray{
					&istiov1.GatewaySpecServersArgs{
						Name: pulumi.String("solr-https"),
						Port: &istiov1.GatewaySpecServersPortArgs{
							Number:   pulumi.Int(443),
							Name:     pulumi.String("solr-https"),
							Protocol: pulumi.String("HTTPS"),
						},
						Hosts: pulumi.StringArray{
							pulumi.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
								solrKubernetes.Spec.Ingress.EndpointDomainName),
							pulumi.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
								solrKubernetes.Spec.Ingress.EndpointDomainName),
						},
						Tls: &istiov1.GatewaySpecServersTlsArgs{
							CredentialName: addedCertificate.Spec.SecretName(),
							Mode:           pulumi.String(v1.ServerTLSSettings_SIMPLE.String()),
						},
					},
					&istiov1.GatewaySpecServersArgs{
						Name: pulumi.String("solr-http"),
						Port: &istiov1.GatewaySpecServersPortArgs{
							Number:   pulumi.Int(80),
							Name:     pulumi.String("solr-http"),
							Protocol: pulumi.String("HTTP"),
						},
						Hosts: pulumi.StringArray{
							pulumi.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
								solrKubernetes.Spec.Ingress.EndpointDomainName),
							pulumi.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
								solrKubernetes.Spec.Ingress.EndpointDomainName),
						},
						Tls: &istiov1.GatewaySpecServersTlsArgs{
							HttpsRedirect: pulumi.Bool(true),
						},
					},
				},
			},
		})
	if err != nil {
		return errors.Wrap(err, "error creating gateway")
	}

	_, err = istiov1.NewVirtualService(ctx,
		solrKubernetes.Metadata.Id,
		&istiov1.VirtualServiceArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:      pulumi.String(solrKubernetes.Metadata.Id),
				Namespace: createdNamespace.Metadata.Name(),
				Labels:    pulumi.ToStringMap(s.Labels),
			},
			Spec: istiov1.VirtualServiceSpecArgs{
				Gateways: pulumi.StringArray{
					pulumi.Sprintf("%s/%s", IstioIngressNamespace,
						solrKubernetes.Metadata.Id),
				},
				Hosts: pulumi.StringArray{
					pulumi.Sprintf("%s.%s", solrKubernetes.Metadata.Id,
						solrKubernetes.Spec.Ingress.EndpointDomainName),
					pulumi.Sprintf("%s-internal.%s", solrKubernetes.Metadata.Id,
						solrKubernetes.Spec.Ingress.EndpointDomainName),
				},
				Http: istiov1.VirtualServiceSpecHttpArray{
					&istiov1.VirtualServiceSpecHttpArgs{
						Name: pulumi.String(solrKubernetes.Metadata.Id),
						Route: istiov1.VirtualServiceSpecHttpRouteArray{
							&istiov1.VirtualServiceSpecHttpRouteArgs{
								Destination: istiov1.VirtualServiceSpecHttpRouteDestinationArgs{
									Host: pulumi.String(""),
									Port: istiov1.VirtualServiceSpecHttpRouteDestinationPortArgs{
										Number: pulumi.Int(8080),
									},
								},
							},
						},
					},
				},
			},
			Status: nil,
		})
	if err != nil {
		return errors.Wrap(err, "error creating virtual-service")
	}
	return nil
}
