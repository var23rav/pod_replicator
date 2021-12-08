package client

import (
	"fmt"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type replicatorClient struct {
	srcConfigPath    string
	targetConfigPath string
	srcNamespace     string
	targetNamespace  string
	resourceType     string
	resourceName     string
}

type ReplicationClient struct {
	*replicatorClient
	srcClientSet    *kubernetes.Clientset
	targetClientSet *kubernetes.Clientset
}

func New() *replicatorClient {
	return &replicatorClient{}
}

func (replicator *replicatorClient) SetSourceCluster(configPath string) {
	replicator.srcConfigPath = strings.Trim(configPath, " ")
}

func (replicator *replicatorClient) SetTargetCluster(configPath string) {
	replicator.targetConfigPath = strings.Trim(configPath, " ")
}

func (replicator *replicatorClient) SetSourceNamespace(ns string) {
	replicator.srcNamespace = strings.Trim(ns, " ")
}

func (replicator *replicatorClient) SetTargetNamespace(ns string) {
	replicator.targetNamespace = strings.Trim(ns, " ")
}

func (replicator *replicatorClient) SetResourceType(resourceType string) {
	replicator.resourceType = strings.Trim(resourceType, " ")
}

func (replicator *replicatorClient) SetResourceName(resourceName string) {
	replicator.resourceName = strings.Trim(resourceName, " ")
}

func (builder *replicatorClient) Build() (*ReplicationClient, error) {
	if builder.srcConfigPath == builder.targetConfigPath && builder.srcNamespace == builder.targetNamespace {
		return nil, fmt.Errorf("Source and Destination config point to same cluster, if you intended replicate across namespace specify different src, dest namespace")
	}

	srcK8sConfig, err := clientcmd.BuildConfigFromFlags("", builder.srcConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get k8sConfig for source cluster(%s), Error: %s", builder.srcConfigPath, err)
	}

	destK8sConfig, err := clientcmd.BuildConfigFromFlags("", builder.targetConfigPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get k8sConfig for source cluster(%s), Error: %s", builder.srcConfigPath, err)
	}

	srcClientSet, err := kubernetes.NewForConfig(srcK8sConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create clientset for source cluster, Error: %s", err)
	}

	destClientSet, err := kubernetes.NewForConfig(destK8sConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create clientset for destination cluster, Error: %s", err)
	}

	return &ReplicationClient{
		replicatorClient: builder,
		srcClientSet:     srcClientSet,
		targetClientSet:  destClientSet,
	}, nil
}

func (replicationClient ReplicationClient) GetSrcClientSet() *kubernetes.Clientset {
	return replicationClient.srcClientSet
}

func (replicationClient ReplicationClient) GetTargetClientSet() *kubernetes.Clientset {
	return replicationClient.targetClientSet
}

func (replicationClient ReplicationClient) GetSrcNamespace() string {
	if replicationClient.srcNamespace == "" {
		return replicationClient.targetNamespace
	}

	return replicationClient.srcNamespace
}

func (replicationClient ReplicationClient) GetTargetNamespace() string {
	if replicationClient.targetNamespace == "" {
		return replicationClient.srcNamespace
	}

	return replicationClient.targetNamespace
}

func (replicationClient ReplicationClient) GetResourceType() string {
	return replicationClient.resourceType
}

func (replicationClient ReplicationClient) GetResourceName() string {
	return replicationClient.resourceName
}
