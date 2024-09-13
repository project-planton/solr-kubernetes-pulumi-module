
# Apache Solr on Kubernetes Pulumi Module

This Pulumi module automates the deployment of **Apache Solr** on **Kubernetes** using the [Solr Operator](https://github.com/apache/solr-operator). It provides an easy-to-use, customizable approach to manage Solr cloud instances with built-in support for Kubernetes Ingress, TLS management, and Zookeeper.

## Features

- Deploys a SolrCloud resource using the Solr Operator.
- Manages Solr pods, Zookeeper pods, and Persistent Volumes for both.
- Supports custom configuration for JVM tuning and garbage collection.
- Automatic creation of Kubernetes namespaces, ingress routes, and certificates using cert-manager.
- Compatible with Kubernetes Gateway API for managing ingress with TLS termination.

## Prerequisites

To use this module, ensure that the following are installed and properly configured:

1. **Pulumi CLI**: Follow the installation guide [here](https://www.pulumi.com/docs/get-started/install/).
2. **Kubernetes cluster**: The module requires access to a Kubernetes cluster with sufficient resources.
3. **Solr Operator**: Install the [Solr Operator](https://github.com/apache/solr-operator) in your Kubernetes cluster.
4. **Cert-Manager**: Used for creating TLS certificates.

## Module Architecture

This module manages the following resources:

- **SolrCloud Custom Resource**: A Solr cloud cluster managed by the Solr Operator.
- **Zookeeper Pods**: Required for the SolrCloud, managed as a Kubernetes resource.
- **Persistent Volume Claims (PVCs)**: For storing Solr and Zookeeper data.
- **Kubernetes Namespace**: A namespace to isolate the resources in the cluster.
- **Ingress with TLS**: Secure ingress to the Solr service with HTTPS termination via cert-manager.

## Usage

### Basic Usage

#### Setup Input

```yaml
apiVersion: code2cloud.planton.cloud/v1
kind: SolrKubernetes
metadata:
  name: main
  id: solk8s-my-org-prod-main
spec:
  solrContainer:
    config:
      garbageCollectionTuning: -XX:SurvivorRatio=4 -XX:TargetSurvivorRatio=90 -XX:MaxTenuringThreshold=8
      javaMem: -Xms1g -Xmx3g
      opts: -Dsolr.autoSoftCommit.maxTime=10000
    diskSize: 1Gi
    image:
      repo: solr
      tag: 8.7.0
    replicas: 1
    resources:
      limits:
        cpu: 2000m
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 250Mi
  zookeeperContainer:
    diskSize: 1Gi
    replicas: 1
    resources:
      limits:
        cpu: 2000m
        memory: 2Gi
      requests:
        cpu: 50m
        memory: 250Mi
  ingress:
    isEnabled: true
    endpointDomainName: my-company.com
```

#### Deploy

```bash
planton pulumi up --input <path-to-stack-input-yaml>
```

#### Destroy

```bash
planton pulumi destroy --input <path-to-stack-input-yaml>
```

This will tear down the SolrCloud and all related resources from your Kubernetes cluster.

### Inputs

#### Solr Kubernetes Spec

- `kubernetes_cluster_credential_id`: Kubernetes cluster credential ID to set up the Kubernetes provider in the stack job.
- `solr_container`: Configuration for the Solr container, including replica count, image, resource requests/limits, disk size, and Solr JVM tuning options.
- `zookeeper_container`: Configuration for the Zookeeper container, including replica count, resource requests/limits, and persistent disk size.

https://github.com/plantoncloud/planton-cloud-apis/blob/84d2058812f72c939b0341c2ddaefdd85ea67a21/cloud/planton/apis/code2cloud/v1/kubernetes/solrkubernetes/spec.proto#L1-L109

### Outputs

This module will create the following outputs:

- **SolrCloud resource**: A running SolrCloud instance in Kubernetes.
- **Zookeeper cluster**: Zookeeper pods supporting the SolrCloud instance.
- **Namespace**: A dedicated Kubernetes namespace for the SolrCloud deployment.
- **Ingress and Certificate**: A secure HTTPS ingress with a certificate issued via cert-manager (if enabled).

## Configuration Options

- **Replica Count**: You can adjust the number of Solr and Zookeeper pods via the `replicas` parameter in the spec.
- **Resource Limits**: Adjust CPU and memory limits for both Solr and Zookeeper containers.
- **Persistent Volumes**: Customize the size of persistent volumes for Solr and Zookeeper using the `disk_size` parameter.

## License

This project is licensed under the Apache 2.0 License

## Contributing

We welcome contributions!
