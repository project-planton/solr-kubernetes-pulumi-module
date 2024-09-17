# Apache Solr on Kubernetes Pulumi Module

This Pulumi module automates the deployment of **Apache Solr** on **Kubernetes** using the [Solr Operator](https://github.com/apache/solr-operator). It provides an easy-to-use, customizable approach to manage Solr cloud instances with built-in support for Kubernetes Ingress, TLS management, and Zookeeper.

## Features

- Deploys a SolrCloud resource using the Solr Operator.
- Manages Solr pods, Zookeeper pods, and Persistent Volumes for both.
- Supports custom configuration for JVM tuning and garbage collection.
- Automatic creation of Kubernetes namespaces, ingress routes, and certificates using cert-manager.
- Compatible with Kubernetes Gateway API for managing ingress with TLS termination.
