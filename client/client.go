package client

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type replicatorClient struct {
	srcConfigPath    string
	targetConfigPath string
	srcNamespace     string
	targetNamespace  string
	deploymentName   string
}

type ReplicationClient struct {
	Client        *replicatorClient
	SrcClientSet  *kubernetes.Clientset
	DestClientSet *kubernetes.Clientset
}

func New() *replicatorClient {
	return &replicatorClient{}
}

func (replicator *replicatorClient) SetSourceCluster(configPath string) {
	replicator.srcConfigPath = configPath
}

func (replicator *replicatorClient) SetTargetCluster(configPath string) {
	replicator.targetConfigPath = configPath
}

func (replicator *replicatorClient) SetSourceNamespace(ns string) {
	replicator.srcNamespace = ns
}

func (replicator *replicatorClient) SetTargetNamespace(ns string) {
	replicator.targetNamespace = ns
}

func (replicator *replicatorClient) SetDeploymentName(deploymentName string) {
	replicator.deploymentName = deploymentName
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
		Client:        builder,
		SrcClientSet:  srcClientSet,
		DestClientSet: destClientSet,
	}, nil
}
