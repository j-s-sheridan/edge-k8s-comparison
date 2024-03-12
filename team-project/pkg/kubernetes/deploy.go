package kubernetes

import (
	"context"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const NAMESPACE = "test-ns"

func Start() *kubernetes.Clientset {
	args, err := GetCommandLineArgs()
	if err != nil {
		panic(err.Error())
	}
	kubeConnection, err := CreateKubeConnection(args.KubeConfigPath)
	if err != nil {
		panic(err.Error())
	}
	UninstallNginx(kubeConnection)
	_, err = deployNginx(kubeConnection)
	if err != nil {
		panic(err.Error())
	}
	return kubeConnection
}

func deployNginx(kubeConnection *kubernetes.Clientset) (*v1.Deployment, error) {
	deployment, err := kubeConnection.AppsV1().Deployments(NAMESPACE).Create(context.Background(), createNginxDeploymentSpec(1), metav1.CreateOptions{})
	fmt.Println(&deployment.Spec.Replicas)
	if err != nil {
		return &v1.Deployment{}, err
	}
	return PollNginxDeploymentStatus(kubeConnection, deployment, "nginx:1.14.2")
}

func UninstallNginx(kubeconnection *kubernetes.Clientset) error {
	if err := kubeconnection.AppsV1().Deployments(NAMESPACE).Delete(context.Background(), "nginx-deployment", metav1.DeleteOptions{}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to delete nginx deployment. Error was: %s", err.Error()))
	}
	fmt.Println("Successfully delete nginx deployment.")
	return nil
}

func createNginxContainerSpec() []v12.Container {
	var containerPorts []v12.ContainerPort
	port := v12.ContainerPort{
		Name:          "containerport",
		ContainerPort: 80,
	}
	containerPorts = append(containerPorts, port)

	var containers []v12.Container
	container := v12.Container{
		Name:  "nginx",
		Image: "nginx:1.14.2",
		Ports: containerPorts,
	}
	containers = append(containers, container)
	return containers
}

func createNginxDeploymentSpec(noReplicas int32) *v1.Deployment {
	selectorMap := make(map[string]string)
	selectorMap["app"] = "nginx"
	containers := createNginxContainerSpec()
	deployment := v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: v1.DeploymentSpec{
			Replicas: &noReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorMap,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selectorMap,
				},
				Spec: v12.PodSpec{
					Containers: containers,
				},
			},
			Strategy:                v1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
		Status: v1.DeploymentStatus{},
	}
	return &deployment
}
