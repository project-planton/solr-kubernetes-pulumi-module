package hostname

import (
	"fmt"
)

func GetInternalHostname(locustKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s-internal.%s", locustKubernetesId, environmentName, endpointDomainName)
}

func GetExternalHostname(locustKubernetesId, environmentName, endpointDomainName string) string {
	return fmt.Sprintf("%s.%s.%s", locustKubernetesId, environmentName, endpointDomainName)
}
