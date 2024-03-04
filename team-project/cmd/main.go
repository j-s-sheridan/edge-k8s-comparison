package main

import (
	"fmt"
	"team-project/pkg/kubernetes"
)

func main() {
	fmt.Println("Starting...")
	kubeConnection, deployment := kubernetes.Start()
	fmt.Println("nginx deployed.")
	fmt.Println("Upgrading nginx.")
	if err := kubernetes.UpgradeNginx(kubeConnection, deployment); err != nil {
		panic(fmt.Sprintf("Upgrade failed. Error was: %s", err.Error()))
	}
	if err := kubernetes.DowngradeNginx(kubeConnection, deployment); err != nil {
		panic(fmt.Sprintf("Downgrade failed. Error was: %s", err.Error()))
	}
}
