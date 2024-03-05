package kubernetes

import (
	"context"
	"errors"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
)

func CreateKubeConnection(filepath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath)
	if err != nil {
		fmt.Println("Error getting connection.")
		return nil, errors.New(fmt.Sprintf("Failed to get connection: %s", err))
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating clientset")
		return nil, errors.New(fmt.Sprintf("Failed to create the connection: %s", err))
	}
	return clientset, nil
}

func PollNginxDeploymentStatus(kubeConnection *kubernetes.Clientset, deployment *v1.Deployment) (*v1.Deployment, error) {
	for {
		fmt.Println("Checking deployment status")
		deployed, err := CheckNginxDeploymentIsReady(kubeConnection)
		if err != nil {
			fmt.Println("Successfully deployed nginx named: %s", deployment.Name)
			return deployment, err
		}
		if deployed {
			fmt.Println("Nginx deployed successfully")
			return deployment, nil
		}
	}
}

func CheckNginxDeploymentIsReady(kubeconnection *kubernetes.Clientset) (bool, error) {
	nginxDeployment, err := kubeconnection.AppsV1().Deployments(NAMESPACE).Get(context.Background(), "nginx-deployment", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if nginxDeployment != nil &&
		nginxDeployment.Spec.Replicas != nil &&
		*nginxDeployment.Spec.Replicas == nginxDeployment.Status.ReadyReplicas {
		return true, nil
	}
	return false, nil
}

func GetLatestNginxDeployment(kubeConnection *kubernetes.Clientset) (*v1.Deployment, error) {
	return kubeConnection.AppsV1().Deployments(NAMESPACE).Get(context.Background(), "nginx-deployment", metav1.GetOptions{})
}

type Args struct {
	KubeConfigPath string
}

func GetCommandLineArgs() (Args, error) {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 1 {
		defaultPath := path.Join(".", "config")
		return Args{KubeConfigPath: defaultPath}, nil
	}
	kubePath := argsWithoutProg[0]
	if _, err := os.Stat(kubePath); err != nil {
		return Args{}, err
	}
	return Args{KubeConfigPath: kubePath}, nil
}
