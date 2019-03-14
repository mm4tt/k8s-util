package metrics

import (
	"k8s.io/contrib/diurnal/Godeps/_workspace/src/github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "probes"
)

var (
	InClusterNetworkLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "in_cluster_network_latency",
		Buckets:   prometheus.ExponentialBuckets(0.000001, 2, 26), // from 1us up to ~1min
		Help:      "Histogram of the time (in seconds) it took to ping a ping-server instance.",
	})
)

func init() {
	prometheus.MustRegister(InClusterNetworkLatency)
}
