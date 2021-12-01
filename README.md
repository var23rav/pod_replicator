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

`Replication within same cluster. Keep the cluster k8s config file in same folder`

> pod_replicator_executable  --srcnamespace vote  --destnamespace default  --deployment vote --debug

`Replication across cluster. Keep the cluster k8s config file in same folder`

> pod_replicator_executable  --srcconfig src_k8s_config --destconfig dest_k8s_config --srcnamespace source_namespace --destnamespace dest_namespace --deployment test_deployment --debug

## Docker Usage

```md
> Replication within same cluster.
docker run                                            \
      -v <host path to kube config>:/.kubeconfig      \
      --rm -it <image_name:tag>                       \
      --srcnamespace <source namespace>               \
      --destnamespace <destination namespace>         \
      --deployment <deployment name>

> Replication across cluster.
docker run                                                        \
      -v <host path to source kube config>:/src_k8s_config        \
      -v <host path to destination kube config>:/dest_k8s_config  \
      --rm -it <image_name:tag>                                   \
      --srcconfig /src_k8s_config                                 \
      --destconfig /dest_k8s_config                               \
      --srcnamespace <source namespace>                           \
      --destnamespace <destination namespace>                     \
      --deployment <deployment name>
```

