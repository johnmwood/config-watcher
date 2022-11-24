package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("Starting main")
	clientset := setKubeConfig()

	initWatcher("default", clientset)

	for {
		time.Sleep(1234 * time.Hour)
	}
}

func initWatcher(namespace string, clientset *kubernetes.Clientset) {
	for {
		watcher, err := clientset.CoreV1().ConfigMaps(namespace).Watch(context.TODO(), metav1.SingleObject(
			metav1.ObjectMeta{
				Name:      "app-config",
				Namespace: namespace,
			}),
		)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			panic("Unable to init watcher")
		}
		startWatcher(watcher.ResultChan())
	}
}

func startWatcher(eventChan <-chan watch.Event) {
	for {
		fmt.Println("Start watcher please?")
		event, open := <-eventChan
		if open {
			switch event.Type {
			case watch.Modified:
				if updatedMap, ok := event.Object.(*corev1.ConfigMap); ok {
					fmt.Printf("Updated map: %v\n", updatedMap)
				}
			case watch.Deleted:
				fmt.Printf("Deleted map: %v\n", event.Object.(*corev1.ConfigMap))
			default:
				// Do nothing
			}
		} else {
			// If eventChannel is closed, it means the server has closed the connection
			return
		}
	}
}

func setKubeConfig() *kubernetes.Clientset {
	clientCfg, err := rest.InClusterConfig()
	if err != nil {
		panic("Unable to get our client configuration")
	}
	clientset, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		panic("Unable to create our clientset")
	}
	return clientset
}
