package service

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/kubernetes/network/dns"
)

func GetKubeServiceNameFqdn(locustKubernetesName, namespace string) string {
	return fmt.Sprintf("%s.%s.%s", GetKubeServiceName(locustKubernetesName), namespace, dns.DefaultDomain)
}

func GetKubeServiceName(locustKubernetesName string) string {
	return fmt.Sprintf(locustKubernetesName)
}
