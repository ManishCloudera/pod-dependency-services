package main

import (
	"github.com/spf13/cast"
	"log"
	"os"
	"time"
)

var (
	noOfTrys     = 1
	namespace    = os.Getenv("NAMESPACE")
	podLabels    = os.Getenv("POD_LABELS")
	retryTimeOut = os.Getenv("RETRY_TIME_OUT_SECOND")
	maxRetry     = os.Getenv("MAX_RETRY")
)

func main() {
	checkIfPodsRunning()
}

func checkIfPodsRunning() {
	log.Println(`Checking for running pods try `, noOfTrys)
	podsRunning := getRunningPods(namespace, podLabels)
	if podsRunning == true {
		log.Println(`Found running pod with Name:`, podLabels)
		os.Exit(0)
	} else if noOfTrys <= cast.ToInt(maxRetry) {
		setTimeout()
	} else {
		log.Println(`Didn't find any running pod with label`, podLabels, ` after `, noOfTrys, ` try`)
		os.Exit(1)
	}
}

func setTimeout() {
	time.Sleep(time.Duration(cast.ToInt(retryTimeOut)) * time.Second)
	noOfTrys++
	checkIfPodsRunning()
}
