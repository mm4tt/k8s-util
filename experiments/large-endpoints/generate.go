package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	fmt.Println(`
apiVersion: v1
kind: Endpoints
metadata:
  namespace: test
  name: large-endpoint
  labels:
    my-label: my-value
subsets:
  - ports:
      - name: https
        port: 443
      - name: http
        port: 80
      - name: http-proxy
        port: 8080
      - name: pulp
        port: 9090
    addresses:`)

	// https://groups.google.com/d/msg/golang-nuts/zlcYA4qk-94/TWRFHeXJCcYJ
	ip, ipnet, err := net.ParseCIDR("10.40.10.0/19")
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		fmt.Println("      - ip: ", ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}