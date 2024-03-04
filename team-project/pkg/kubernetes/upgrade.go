package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func UpgradeNginx(kubeConnection *kubernetes.Clientset, deployment *v1.Deployment) error {
	containers := deployment.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == "nginx" {
			container.Image = "nginx:1.15.12"
		}
	}
	_, err := kubeConnection.AppsV1().Deployments(NAMESPACE).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Upgrade failed. Error was: %s", err.Error())
	}
	err = PollReplicasReady(kubeConnection, deployment, 10)
	if err != nil {
		fmt.Println("Upgrade failed to be ready in time.")
	}
	return err
}

func DowngradeNginx(kubeConnection *kubernetes.Clientset, deployment *v1.Deployment) error {
	containers := deployment.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == "nginx" {
			container.Image = "nginx:1.14.2"
		}
	}
	_, err := kubeConnection.AppsV1().Deployments(NAMESPACE).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		fmt.Println("Downgrade failed. Error was: %s", err.Error())
	}
	err = PollReplicasReady(kubeConnection, deployment, 10)
	if err != nil {
		fmt.Println("Upgrade failed to be ready in time.")
	}
	return err
}
