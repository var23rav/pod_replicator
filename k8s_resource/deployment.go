package resource

import (
	"context"
	"fmt"

	prclient "github.com/var23rav/pod_replicator/client"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
)

type Deployment struct {
	*prclient.ReplicationClient
}

func panicOnErr(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("Err.. %s, %s", msg, err))
	}
}

func (deploy Deployment) Name() string {
	return fmt.Sprintf("%s/%s", deploy.GetResourceType(), deploy.GetResourceName())
}

func (deploy Deployment) Migrate(enableForceOverride bool) error {
	replicaClient := deploy.ReplicationClient
	srcnamespace := replicaClient.GetSrcNamespace()
	destnamespace := replicaClient.GetTargetNamespace()
	deploymentName := replicaClient.GetResourceName()

	srcNSDeployment := replicaClient.GetSrcClientSet().AppsV1().Deployments(srcnamespace)
	destNSDeployment := replicaClient.GetTargetClientSet().AppsV1().Deployments(destnamespace)

	// Initial try
	_, err := srcNSDeployment.Get(context.Background(), deploymentName, metav1.GetOptions{})
	panicOnErr(err, fmt.Sprintf("While reading deployment/%s DeploymentSet from source cluster", deploymentName))
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		srcDeployment, err := srcNSDeployment.Get(context.Background(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("%s While reading deployment/%s from source cluster", err, deploymentName)
		}

		destDeploymentTemp := srcDeployment.DeepCopy()
		destDeploymentTemp.Namespace = destnamespace

		isResourceExist := true
		destDeployment, err := destNSDeployment.Get(context.Background(), deploymentName, metav1.GetOptions{})
		if err != nil {
			if !enableForceOverride {
				return fmt.Errorf("%s While reading deployment/%s from destination cluster", err, deploymentName)
			}
			isResourceExist = false
		}

		if enableForceOverride {
			return Replace(destNSDeployment, destDeploymentTemp, isResourceExist)
		} else {
			destDeploymentTemp.SetResourceVersion(destDeployment.GetResourceVersion())
			destDeploymentTemp.SetUID(destDeployment.GetUID())
			_, err = destNSDeployment.Update(context.Background(), destDeploymentTemp, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("%s While updating deployment/%s", err, deploymentName)
			}
		}

		return nil
	})

	return retryErr
}

func Replace(nsDeployment typedappsv1.DeploymentInterface, deployment *appsv1.Deployment, isResourceExist bool) error {
	if isResourceExist {
		_ = nsDeployment.Delete(context.Background(), deployment.Name, metav1.DeleteOptions{})
		// if err != nil {
		// fmt.Printf("\nDelete Failed\n%+v\n\n", err)
		// return fmt.Errorf("Replace Deploy/%s step 1 Delete existing failed, Err %s", deployment.Name, err)
		// }
	}

	err := Create(nsDeployment, deployment)
	if err != nil {
		// fmt.Printf("\nDelete Failed\n%+v\n\n", err)
		return fmt.Errorf("Replace Deploy/%s step 1 Delete existing failed, Err %s", deployment.Name, err)
	}

	fmt.Println("Overide is successful")
	return nil
}

func Create(nsDeployment typedappsv1.DeploymentInterface, deployment *appsv1.Deployment) error {
	deployment.CreationTimestamp.Reset()
	deployment.SetResourceVersion("")
	deployment.SetUID("")
	deployment.SetSelfLink("")
	_, err := nsDeployment.Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Replace Deploy/%s step 2 Create new failed, Err %s", deployment.Name, err)
	}

	return nil
}
