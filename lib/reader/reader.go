// resource-usage-summary-parser - util for analyzing resource-usage-summary-parsers
// Example useage:
//   go run $GOPATH/src/github.com/mm4tt/k8s-util/resource-usage-summary-parser/resource-usage-summary-parser.go --resource-usage-summary-path /tmp/cl2/coredns_after/ResourceUsageSummary_load_2019-03-28T17\:13\:59+01\:00.txt --name-pattern=coredns --stderrthreshold=INFO
package reader

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Read reads data from the provided path, which can be a local filesystem path or url.
func Read(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return readHttp(path)
	}
	return ioutil.ReadFile(path)
}


func readHttp(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
