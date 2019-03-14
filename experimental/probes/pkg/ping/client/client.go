package pingclient

import (
	"flag"
	"net/http"
	"time"

	"github.com/mm4tt/k8s-util/experimental/probes/pkg/metrics"
	"k8s.io/klog"
)

var (
	pingServerAddress = flag.String("ping-server-address", "", "The address of the ping server")
	pingSleepDuration = flag.Duration("ping-sleep-duration", 1*time.Second, "Duration of the sleep between pings")
)

// Run runs the ping client probe that periodically pings the ping server and exports latency metric.
func Run() {
	if *pingServerAddress == "" {
		klog.Fatal("--ping-server-address not set!")
	}

	for {
		time.Sleep(*pingSleepDuration)
		klog.Infof("ping -> %s...\n", *pingServerAddress)
		start := time.Now()
		if err := ping(*pingServerAddress); err != nil {
			klog.Infof("Got error: %v", err)
			// TODO(mm4tt): Increment server not available gauge metric.
			continue
		}
		end := time.Now()
		latency := end.Sub(start)
		klog.Infof("Request took: %v\n", latency)
		metrics.InClusterNetworkLatency.Observe(latency.Seconds())
	}
}

func ping(serverAddress string) error {
	_, err := http.Get("http://" + serverAddress)
	return err
}
