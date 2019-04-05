package main

import (
	"flag"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

var (
	numWatches = flag.Int("num-watches", 5, "Number of watches to start")
)

func main() {
	flag.Parse()

	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal(err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err.Error())
	}

	for i := 0; i < *numWatches; i++ {
		go startWatch(i, client)
	}

	// Block main routine.
	<-make(chan bool)
}

func startWatch(id int, client kubernetes.Interface) {
	for {
		watch, err := client.CoreV1().Endpoints("").Watch(metav1.ListOptions{})
		if err != nil {
			klog.Warningf("Got error: %v", err)
			klog.Info("Sleeping 10s and retrying...")
			time.Sleep(10 * time.Second)
		}
		nReceived := 0
		for {
			<-watch.ResultChan()
			nReceived++
			if nReceived%10 == 0 {
				klog.Infof("Watch %d received %d events so far", id, nReceived)
			}
		}
	}
}
