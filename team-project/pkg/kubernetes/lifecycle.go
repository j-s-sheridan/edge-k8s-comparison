package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpgradeNginx(kubeConnection *kubernetes.Clientset) error {
	deployment, err := GetLatestNginxDeployment(kubeConnection)
	if err != nil {
		fmt.Println("Failed to get the latest nginx deployment")
		return err
	}
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:1.15.12"
	_, err = kubeConnection.AppsV1().Deployments(NAMESPACE).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Upgrade failed. Error was: %s", err.Error())
		return UpgradeNginx(kubeConnection)
	}
	if _, err := PollNginxDeploymentStatus(kubeConnection, deployment, "nginx:1.15.12"); err != nil {
		fmt.Println("Failed to poll nginx deployment status")
		return err
	}
	fmt.Println("Successfully upgraded nginx containers.")
	return nil
}

func DowngradeNginx(kubeConnection *kubernetes.Clientset) error {
	deployment, err := GetLatestNginxDeployment(kubeConnection)
	if err != nil {
		fmt.Println("Failed to get the latest nginx deployment")
		return err
	}
	deployment.Spec.Template.Spec.Containers[0].Image = "nginx:1.14.2"
	_, err = kubeConnection.AppsV1().Deployments(NAMESPACE).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Downgrade failed. Error was: %s", err.Error())
		return DowngradeNginx(kubeConnection)
	}
	if _, err := PollNginxDeploymentStatus(kubeConnection, deployment, "nginx:1.14.2"); err != nil {
		fmt.Println("Failed to poll nginx deployment status")
		return err
	}
	fmt.Println("Successfully downgraded nginx containers.")
	return err
}
