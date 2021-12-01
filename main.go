package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	prclient "github.com/var23rav/pod_replicator/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func panicOnErr(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("Err.. %s, %s", msg, err))
	}
}

func main() {
	var debug bool
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("----------------------")
			fmt.Println("Failed to complete!")
			if debug {
				fmt.Printf("\nError: %s \n", r)
			}
		}
	}()

	defaultConfigPath, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprint("Failed to detect current path", err))
	}
	defaultConfigPath = filepath.Join(defaultConfigPath, ".kubeconfig")
	flag.BoolVar(&debug, "debug", false, "Enable debug logs")
	srcConfigPath := flag.String("srcconfig", defaultConfigPath, "Enter the source config file path")
	destConfigPath := flag.String("destconfig", defaultConfigPath, "Enter the destination config file path")

	var srcnamespace string
	flag.StringVar(&srcnamespace, "srcnamespace", "default", "Enter the source namespace")
	var destnamespace string
	flag.StringVar(&destnamespace, "destnamespace", "default", "Enter the source namespace")

	var deployment string
	flag.StringVar(&deployment, "deployment", "vote", "Enter the deploymentset")
	if deployment == "" {
		panicOnErr(fmt.Errorf("%s, current %s", "set --deployment flag to duplicate", deployment), "No Deployment flag")
	}

	flag.Parse()

	replicator := prclient.New()
	replicator.SetSourceCluster(*srcConfigPath)
	replicator.SetTargetCluster(*destConfigPath)
	replicator.SetSourceNamespace(srcnamespace)
	replicator.SetTargetNamespace(destnamespace)
	replicator.SetDeploymentName(deployment)
	replicaClient, err := replicator.Build()
	panicOnErr(err, "Replication client creation failed")

	srcNSDeployment := replicaClient.SrcClientSet.AppsV1().Deployments(srcnamespace)
	destNSDeployment := replicaClient.DestClientSet.AppsV1().Deployments(destnamespace)

	// Initial try
	_, err = srcNSDeployment.Get(context.Background(), deployment, metav1.GetOptions{})
	panicOnErr(err, fmt.Sprintf("While reading deployment/%s DeploymentSet from source cluster", deployment))
	var newReplicaCount, oldReplicaCount int32
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		srcDeployment, err := srcNSDeployment.Get(context.Background(), deployment, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("%s While reading deployment/%s from source cluster", err, deployment)
		}

		destDeployment, err := destNSDeployment.Get(context.Background(), deployment, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("%s While reading deployment/%s from destination cluster", err, deployment)
		}

		newReplicaCount = *srcDeployment.Spec.Replicas
		oldReplicaCount = *destDeployment.Spec.Replicas
		destDeployment.Spec.Replicas = &newReplicaCount
		_, err = destNSDeployment.Update(context.Background(), destDeployment, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("%s While updating replica count update %d->%d %s Deploy", err, oldReplicaCount, newReplicaCount, deployment)
		}

		return nil
	})
	panicOnErr(retryErr, fmt.Sprintf("Replication failed for deployment/%s", deployment))
	fmt.Printf("Replication(%d -> %d) of deployment/%s successful\n", oldReplicaCount, newReplicaCount, deployment)
}
