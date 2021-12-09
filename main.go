package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	prclient "github.com/var23rav/pod_replicator/client"
	k8sresource "github.com/var23rav/pod_replicator/k8s_resource"
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

	var resourceType string
	flag.StringVar(&resourceType, "kind", "Deployment", "Enter the K8s ResourceType")
	if resourceType == "" {
		requiredErr := fmt.Errorf("%s, current %s", "set --kind flag to identify the targeted resource type", resourceType)
		fmt.Println(requiredErr)
		panicOnErr(requiredErr, "No ResourceType flag")
	}

	var resourceName string
	flag.StringVar(&resourceName, "name", "default", "Enter the Resource name")
	if resourceName == "" {
		requiredErr := fmt.Errorf("%s, current %s", "set --name flag to identify the targeted resource", resourceName)
		fmt.Println(requiredErr)
		panicOnErr(requiredErr, "No ResourceName flag")
	}

	var enableForceOverride bool
	flag.BoolVar(&enableForceOverride, "force", false, "Force override; Create in case missing replace when existing.")

	flag.Parse()

	replicator := prclient.New()
	replicator.SetSourceCluster(*srcConfigPath)
	replicator.SetTargetCluster(*destConfigPath)
	replicator.SetSourceNamespace(srcnamespace)
	replicator.SetTargetNamespace(destnamespace)
	replicator.SetResourceType(resourceType)
	replicator.SetResourceName(resourceName)
	replicaClient, err := replicator.Build()
	panicOnErr(err, "Replication client creation failed")

	resourceObj := k8sresource.ResourceType(resourceType)
	resourceMigrator, err := resourceObj.Build(replicaClient)
	panicOnErr(err, "Building Resource Object failed")

	err = resourceMigrator.Migrate(enableForceOverride)
	panicOnErr(err, fmt.Sprintf("Resource migration failed for %s", resourceMigrator.Name()))
	if err != nil {
		fmt.Printf("Resource migration failed for %s, %+v\n", resourceMigrator.Name(), err)
	}

	fmt.Println("Resource migration completed successfully")
}
