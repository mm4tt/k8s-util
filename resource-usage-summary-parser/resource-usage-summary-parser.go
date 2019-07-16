// resource-usage-summary-parser - util for analyzing resource-usage-summary-parsers
// Example useage:
//   go run $GOPATH/src/github.com/mm4tt/k8s-util/resource-usage-summary-parser/resource-usage-summary-parser.go --resource-usage-summary-path /tmp/cl2/coredns_after/ResourceUsageSummary_load_2019-03-28T17\:13\:59+01\:00.txt --name-pattern=coredns --stderrthreshold=INFO
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"k8s.io/klog"
	"sort"
	"strconv"
	"strings"

	"github.com/mm4tt/k8s-util/lib/reader"
)

var (
	resourceUsageSummaryPath = flag.String("resource-usage-summary-path", "", "Path to the resource usage summary path")
	namePattern              = flag.String("name-pattern", "", "Pattern (substring) of the component to analyze")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	verifyFlags()

	jsonMap, err := readJson(*resourceUsageSummaryPath)
	if err != nil {
		klog.Fatal(err)
	}

	err = computeSummary(jsonMap, *namePattern)
	if err != nil {
		klog.Fatal(err)
	}
}

func verifyFlags() {
	if *resourceUsageSummaryPath == "" {
		klog.Fatalf("--resource-usage-summary-path must be set")
	}
	if *namePattern == "" {
		klog.Fatalf("--name-pattern must be set")
	}
}

func readJson(path string) (map[string]interface{}, error) {
	data, err := reader.Read(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func computeSummary(jsonMap map[string]interface{}, namePattern string) error {
	for pctlString, v := range jsonMap {
		pctl, err := strconv.Atoi(pctlString)
		if err != nil {
			return err
		}
		var cpus, mems []float64

		measurements, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("unable to cast value to measurments")
		}
		for _, mRaw := range measurements {
			m, ok := mRaw.(map[string]interface{})
			if !ok {
				return fmt.Errorf("unable to parse measurement %v", mRaw)
			}
			if strings.Contains(m["Name"].(string), namePattern) {
				cpus = append(cpus, m["Cpu"].(float64))
				mems = append(mems, m["Mem"].(float64))
			}
		}

		sort.Float64s(cpus)
		sort.Float64s(mems)

		cpu := cpus[pctl*len(cpus)/100-1]
		mem := mems[pctl*len(mems)/100-1]

		klog.Infof("%d: Cpu: %f, Mem: %f\n", pctl, cpu, mem)
	}
	return nil
}
