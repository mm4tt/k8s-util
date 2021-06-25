package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/klog/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	clientset kubernetes.Interface
)

func main() {
	updateConfigMap()
}

func updateConfigMap() {

	mustBuildClientSet()

	updater := func(instanceID int) {
		cfg, err := clientset.CoreV1().ConfigMaps("default").Get(context.Background(), "tstcfg", metav1.GetOptions{})
		if err != nil {
			klog.Fatalf("error getting: %v", err)
		}

		if cfg.Data == nil {
			cfg.Data = map[string]string{}
		}
		cfg.Data[fmt.Sprint(instanceID)] = cfg.ResourceVersion
		klog.Infof("%d see RV = %v", instanceID, cfg.ResourceVersion)

		if cfg2, err := clientset.CoreV1().ConfigMaps("default").Update(context.Background(), cfg, metav1.UpdateOptions{}); err != nil {
			klog.Errorf("%d: error updating: %v", instanceID, err)
		} else {
			klog.Infof("%d: successfuly updated, newRV = %v", instanceID, cfg2.ResourceVersion)
		}
	}

	requestsWindow := 100 * time.Millisecond
	for {
		go updater(1)
		time.Sleep(requestsWindow)
		go updater(2)
		time.Sleep(requestsWindow)
		go updater(3)
		time.Sleep(5 * time.Second)
	}
}

func mustBuildClientSet() {
	var err error
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home := os.Getenv("HOME")
		kubeconfig = home + "/.kube/config"
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}