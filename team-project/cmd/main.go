package main

import (
	"fmt"
	"team-project/pkg/kubernetes"
)

func main() {
	fmt.Println("Starting...")
	kubernetes.Start()
	kubeConnection := kubernetes.Start()
	fmt.Println("nginx deployed.")
	fmt.Println("Beginning upgrade / Downgrade endurance test")
	enduranceTestResults := kubernetes.RunUpgradeDowngradeEnduranceTest(kubeConnection, 2)
	fmt.Println(fmt.Sprintf("Performed a total of %v upgrades", len(enduranceTestResults.UpgradeTimes)))
	fmt.Println(fmt.Sprintf("Performed a total of %v downgrades", len(enduranceTestResults.DowngradeTimes)))
}
