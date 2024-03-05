package kubernetes

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"time"
)

type EnduranceTestResults struct {
	UpgradeTimes   []time.Duration
	DowngradeTimes []time.Duration
}

func RunUpgradeDowngradeEnduranceTest(kubeConnection *kubernetes.Clientset, runtimeInMinutes int64) EnduranceTestResults {
	var upgradeTimes []time.Duration
	var downgradeTimes []time.Duration
	start := time.Now()
	endTime := time.Minute * time.Duration(runtimeInMinutes)
	for {
		if time.Since(start) > endTime {
			break
		}
		upgradeStartTime := time.Now()
		if err := UpgradeNginx(kubeConnection); err != nil {
			panic(err)
		}
		upgradeTimes = append(upgradeTimes, time.Since(upgradeStartTime))
		downgradeStartTime := time.Now()
		if err := DowngradeNginx(kubeConnection); err != nil {
			panic(err)
		}
		downgradeTimes = append(downgradeTimes, time.Since(downgradeStartTime))
	}
	fmt.Println("Upgrade / downgrade test run completed.")
	return EnduranceTestResults{
		UpgradeTimes:   upgradeTimes,
		DowngradeTimes: downgradeTimes,
	}
}
