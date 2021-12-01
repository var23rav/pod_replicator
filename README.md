# Pod Replicator

* Support across cluster
* Support across namespace

## Prerequisites

* Cluster is supposed to define with proper k8s resources like Namespaces and Deployment
* Keep source and destination config file in handy

## Usage

```md
Usage of pod_replicator:
  -debug
        Enable debug logs
  -deployment string
        Enter the deploymentset (default "vote")
  -destconfig string
        Enter the destination config file path (default "$pwd\.kubeconfig")
  -destnamespace string
        Enter the source namespace (default "default")
  -srcconfig string
        Enter the source config file path (default "$pwd\.kubeconfig")
  -srcnamespace string
        Enter the source namespace (default "default")
```
