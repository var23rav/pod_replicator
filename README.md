# Pod Replicator

* Support across cluster
* Support across namespace
* Support migration of Resources(currently deployment supported)

## Prerequisites

* Cluster is supposed to define with proper k8s resources like Namespaces
* Keep source and destination config file in handy

## Usage

```md
Usage of ./pod_replicator:
  -debug
        Enable debug logs
  -destconfig string
        Enter the destination config file path (default "$pwd/.kubeconfig")
  -destnamespace string
        Enter the source namespace (default "default")
  -force
        Force override; Create in case missing replace when existing
  -kind string
        Enter the K8s ResourceType (default "Deployment")
  -name string
        Enter the Resource name (default "default")
  -srcconfig string
        Enter the source config file path (default "$pwd/.kubeconfig")
  -srcnamespace string
        Enter the source namespace (default "default")
```

`Replication within same cluster for existing resource. Keep the cluster k8s config file in same folder`

> pod_replicator_executable  --srcnamespace vote  --destnamespace default  --name vote --debug

`Replication within same cluster for overriding existing or create new resources. Keep the cluster k8s config file in same folder`

> pod_replicator_executable  --srcnamespace vote  --destnamespace default  --name vote --debug --force

`Replication across cluster. Keep the cluster k8s config file in same folder`

> pod_replicator_executable  --srcconfig src_k8s_config --destconfig dest_k8s_config --srcnamespace source_namespace --destnamespace dest_namespace --kind deployment --name test_deployment --debug

## Docker Usage

```md
> Replication within same cluster.
docker run                                            \
      -v <host path to kube config>:/.kubeconfig      \
      --rm -it <image_name:tag>                       \
      --srcnamespace <source namespace>               \
      --destnamespace <destination namespace>         \
      --kind <resource type>                          \
      --name <resource name>                          \
      -- force                                        \
      -- debug

> Replication across cluster.
docker run                                                        \
      -v <host path to source kube config>:/src_k8s_config        \
      -v <host path to destination kube config>:/dest_k8s_config  \
      --rm -it <image_name:tag>                                   \
      --srcconfig /src_k8s_config                                 \
      --destconfig /dest_k8s_config                               \
      --srcnamespace <source namespace>                           \
      --destnamespace <destination namespace>                     \
      --name <resource name>
```

## Attention Developers (PRs expected)

`Sample resource file created to help the onboarding. Try mock resource migration`

* Update the supported resources under [ResourceType.Build()](./k8s_resource/k8s_resource.go)
* Write the migration logic for the same [k8s_resource/mock_resource.go](./k8s_resource/mock_resource.go)

> pod_replicator_executable --srcnamespace source --destnamespace  target --kind mock_resource --name test --debug
