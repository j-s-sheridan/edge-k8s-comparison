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
	"time"
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

func PollReplicasReady(clientset *kubernetes.Clientset, deployment *v1.Deployment, timeoutInSeconds int) error {
	startTime := time.Now()
	timeoutInMillis := int64(timeoutInSeconds * 1000)
	for {
		deployment, _ = clientset.AppsV1().Deployments(NAMESPACE).Update(context.Background(), deployment, metav1.UpdateOptions{})
		if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
			fmt.Println(fmt.Sprintf("Deployment: %s successfully deployed. Wanted replicas: %v. Have replicas: %v", deployment.Name, deployment.Status.Replicas, deployment.Status.ReadyReplicas))
			return nil
		}
		if time.Since(startTime).Milliseconds() > timeoutInMillis {
			return errors.New("Deployment did not deploy ready in time.")
		}
	}
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
