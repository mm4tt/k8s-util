package pingserver

import (
	"flag"
	"net/http"

	"k8s.io/klog"
)

var (
	pingServerBindAddress = flag.String("ping-server-bind-address", "", "The address to bind for ping server")
)

// Run runs the ping server.
func Run() {
	if *pingServerBindAddress == "" {
		klog.Fatal("--ping-server-bind-address not set!")
	}

	klog.Infof("Listening on %s \n", *pingServerBindAddress)
	http.HandleFunc("/", pong)
	klog.Fatal(http.ListenAndServe(*pingServerBindAddress, nil))
}

func pong(w http.ResponseWriter, r *http.Request) {
	klog.Infof("pong -> %s\n", r.RemoteAddr)
	w.Write([]byte("pong"))
}
