package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"strings"
)

func connectToK8sInClusterConfig() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Failed to create K8s client")
		return nil, err
	}

	return clientSet, nil
}

func getRunningPods(namespace, podLabels string) bool {
	client, err := connectToK8sInClusterConfig()
	if err != nil {
		return false
	}
	ns := client.CoreV1().Pods(namespace)
	pods, errPods := ns.List(context.TODO(), metav1.ListOptions{})
	if errPods != nil {
		log.Println("Error Get list pods:", errPods)
		return false
	}
	podLabelArray := strings.Split(podLabels, ",")
	podPhaseMapping := map[string]v1.PodPhase{}
	for _, pod := range pods.Items {
		podPhaseMapping[pod.Labels["app"]] = pod.Status.Phase
	}

	for _, podLabel := range podLabelArray {
		if podPhaseMapping[podLabel] == "Running" {
			log.Println(podLabel+" pod phase: ", podPhaseMapping[podLabel])
		} else {
			return false
		}
	}

	return true
}
