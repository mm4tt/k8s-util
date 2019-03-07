// junit-parser - util for fetching and analyzing junit files from k8s ci tests.
// Example useage:
//   go run $GOPATH/src/github.com/mm4tt/k8s-util/test-summary-parser/test-summary-parser.go --test-id=ci-kubernetes-e2e-gce-scale-correctness/208
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

var (
	testId = flag.String("test-id", "ci-kubernetes-e2e-gce-scale-correctness/208", "Id of the test to analyze in form <name>/<run_number>")
)

func main() {
	flag.Parse()

	tcChan := make(chan Testcase, 1000)
	var wg sync.WaitGroup

	wg.Add(1)
	go process("/artifacts/junit_runner.xml", tcChan, &wg)

	// TODO(mm4tt): Figure out something more sophisticated
	for i := 1; i < 100; i++ {
		wg.Add(1)
		go process(fmt.Sprintf("/artifacts/junit_%02d.xml", i), tcChan, &wg)
	}

	wg.Wait()
	close(tcChan)

	testcases := make([]Testcase, 0, len(tcChan))
	for tc := range tcChan {
		testcases = append(testcases, tc)
	}

	sort.Slice(testcases, func(i, j int) bool {
		return testcases[i].Duration > testcases[j].Duration
	})

	fmt.Println("Top 50 longest testcases")
	for i, tc := range testcases {
		fmt.Println(tc)
		if i == 50 {
			break
		}
	}
}

func process(relativePath string, outputChan chan<- Testcase, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := read(relativePath)
	if err != nil {
		klog.Warning(err)
		return
	}
	var suite xmlTestsuite
	if err := xml.Unmarshal(data, &suite); err != nil {
		klog.Warning(err)
		return
	}
	for _, tc := range suite.Testcases {
		if tc.Skipped != (xmlSkipped{}) {
			continue
		}

		floatDuration, err := strconv.ParseFloat(tc.Duration, 64)
		if err != nil {
			klog.Fatal(err)
		}
		duration, err := time.ParseDuration(fmt.Sprintf("%fs", floatDuration))
		if err != nil {
			klog.Fatal(err)
			continue
		}

		outputChan <- Testcase{
			Classname: tc.Classname,
			Name:      tc.Name,
			Duration:  duration,
		}
	}
}

func read(relativePath string) ([]byte, error) {
	resp, err := http.Get("http://gcsweb.k8s.io/gcs/kubernetes-jenkins/logs/" + *testId + relativePath)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Testcase struct {
	Classname string
	Name      string
	Duration  time.Duration
}

func (t Testcase) String() string {
	return fmt.Sprintf("%v: [%s] %s", t.Duration, t.Classname, t.Name)
}

type xmlTestsuite struct {
	XMLName   xml.Name      `xml:"testsuite"`
	Testcases []xmlTestcase `xml:"testcase"`
}

type xmlTestcase struct {
	XMLName xml.Name `xml:"testcase"`

	Classname string `xml:"classname,attr"`
	Name      string `xml:"name,attr"`
	Duration  string `xml:"time,attr"`

	Skipped xmlSkipped `xml:"skipped"`
}

type xmlSkipped struct {
	XMLName xml.Name `xml:"skipped"`
}
