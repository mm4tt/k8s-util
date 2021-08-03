package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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
	mustBuildClientSet()

	if err := compareCacheAndEtcdRv("TODO-node-name"); err != nil {
		log.Fatalf("got error: %v", err)
	}
	log.Println("Finished successfully!")
}

func watchNodesFromRV(RV string) error {
	ctx := context.Background()
	timeoutSec := int64(60)

	watch, err := clientset.CoreV1().Nodes().Watch(ctx, metav1.ListOptions{
		ResourceVersion: RV,
		TimeoutSeconds:  &timeoutSec,
	})
	if err != nil {
		return err
	}

	log.Printf("Opened watch from RV=%v", RV)

	for r := range watch.ResultChan() {
		log.Printf("Got watch event: %v", r)
	}

	return nil
}

func compareCacheAndEtcdRv(nodeName string) error {
	ctx := context.Background()

	for {
		cacheNode, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{
			ResourceVersion: "0",
		})
		if err != nil {
			return err
		}

		dbNode, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{
			ResourceVersion: "",
		})
		if err != nil {
			return err
		}
		cacheRV, _ := strconv.ParseInt(cacheNode.ResourceVersion, 10, 64)
		dbRV, _ := strconv.ParseInt(dbNode.ResourceVersion, 10, 64)
		diff := time.Unix(0, int64(time.Microsecond)*dbRV).Sub(time.Unix(0, int64(time.Microsecond)*cacheRV))
		log.Printf("Node %v db RV: %v, cache RV: %v, time diff: %v", cacheNode.Name, dbRV, cacheRV, diff)
		time.Sleep(time.Second)
	}
}

func updateNodeStatus(nodeName string) error {
	ctx := context.Background()
	node, err := clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{
		ResourceVersion: "0",
	})
	if err != nil {
		return err
	}
	node.Status.NodeInfo.OSImage += " TEST"
	n2, err := clientset.CoreV1().Nodes().UpdateStatus(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error while updating node status: %w", err)
	}
	log.Printf("UpdateStatus - success. Node %v RV is now: %v", n2.Name, n2.ResourceVersion)
	return nil
}

func updateConfigMap() {
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
