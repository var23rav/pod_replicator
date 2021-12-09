package resource

import (
	"fmt"

	prclient "github.com/var23rav/pod_replicator/client"
)

type MockResource struct {
	*prclient.ReplicationClient
}

func (dummy MockResource) Name() string {
	return fmt.Sprintf("Example: %s/%s", dummy.GetResourceType(), dummy.GetResourceName())
}

func (dummy MockResource) Migrate(enableForceOverride bool) error {
	if enableForceOverride {
		fmt.Printf("This will overide the resource(%s) regardless one already exist or not \n", dummy.Name())
	} else {
		fmt.Printf("This will update an existing resource(%s), if targeted resource missing will throw error \n", dummy.Name())
	}

	fmt.Println("-----------------------")
	fmt.Println("Thank you ! This was a mock representation of k8s resource.")
	fmt.Println("You can include migration logic for other k8s resource, take this sample resource as base.")
	fmt.Println("PRs are more than welcome to enrich this cli tool.")
	return nil
}
