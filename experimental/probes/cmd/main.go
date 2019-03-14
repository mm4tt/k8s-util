package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/mm4tt/k8s-util/experimental/probes/pkg/ping/client"
	"github.com/mm4tt/k8s-util/experimental/probes/pkg/ping/server"
	"k8s.io/klog"
)

var (
	metricAddress         = flag.String("metric-bind-address", "0.0.0.0:8080", "The address to serve the Prometheus metrics on.")
	mode                  = flag.String("mode", "", "Mode that should be run. Supported values: ping-server, ping-client")
)

func main() {
	flag.Parse()
	verifyFlags()

	klog.Infof("I'm probes.")
	klog.Infof("Mode is: %s\n", *mode)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		klog.Infof("Serving metrics on %s\n", *metricAddress)
		klog.Fatal(http.ListenAndServe(*metricAddress, nil))
	}()

	switch *mode {
	case "ping-client":
		pingclient.Run()
	case "ping-server":
		pingserver.Run()
	default:
		klog.Fatalf("Unrecognized mode: '%s'", *mode)
	}
}

func verifyFlags() {
	if *metricAddress == "" {
		klog.Fatal("--metric-bind-address not set!")
	}
	if *mode == "" {
		klog.Fatal("--mode not set!")
	}
}