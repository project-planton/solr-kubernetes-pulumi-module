package service

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
)

func GetKubeServiceNameFqdn(solrKubernetesName, namespace string) string {
	return fmt.Sprintf("%s.%s.%s", GetKubeServiceName(solrKubernetesName), namespace, dns.DefaultDomain)
}

func GetKubeServiceName(solrKubernetesName string) string {
	return fmt.Sprintf("%s-solrcloud-common", solrKubernetesName)
}
