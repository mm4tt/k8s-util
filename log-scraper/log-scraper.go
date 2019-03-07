// log_scraper - util for fetching and greping scalability test logs.
// Example usage:
// go run log-scraper/log-scraper.go --log-url=http://gcsweb.k8s.io/gcs/kubernetes-jenkins/logs/ci-kubernetes-e2e-gce-scale-performance/294/ --pod-name=density-latency-pod-3787-5rt98
//
// Requirements:
//  * go get github.com/nhooyr/color
package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nhooyr/color"
)

var (
	logURL  = flag.String("log-url", "http://gcsweb.k8s.io/gcs/kubernetes-jenkins/logs/ci-kubernetes-e2e-gce-scale-performance/294/", "URL of the root directory of the test run results.")
	podName = flag.String("pod-name", "density-latency-pod-3779-qvspw", "Name of the pod to look for.")

	charsAround = flag.Int("chars-around", 50, "How many chars should be printed around pattern in long lines")
	linesAround = flag.Int("lines-around", 0, "How many lines should be printed around (i.e. before and after) pattern")

	cacheRootDir = flag.String("cache-root-dir", "/tmp/log_cache/", "Root dir of the log cache dir.")
)

func main() {
	flag.Parse()

	fmt.Println("Running pod log tracker...")
	fmt.Printf("\tLog url is: %s\n", *logURL)
	fmt.Printf("\tLooking for pod: %s\n", *podName)
	fmt.Println()

	store := newLogStore(*cacheRootDir, *logURL, *linesAround)

	buildLog := store.Grep("build-log.txt", *podName)

	masterDir := store.List("artifacts", "master")[0]

	schedulerLog := store.Grep(masterDir+"kube-scheduler.log", *podName)
	controllerManagerLog := store.Grep(masterDir+"kube-controller-manager.log", *podName)

	apiserverLogs := []string{}
	for _, apiserverLog := range store.List(masterDir, "kube-apiserver.log") {
		log := store.Grep(apiserverLog, *podName)
		apiserverLogs = append(apiserverLogs, log...)
	}

	node := getNode(schedulerLog)
	var kubeletLog []string
	if node != "" {
		fmt.Printf("Pod's node is: %s\n\n", node)
		kubeletLog = store.Grep("artifacts/"+node+"/kubelet.log", *podName)
	} else {
		fmt.Printf("Node not found!!!")
		kubeletLog = []string{}
	}

	printAll(buildLog, apiserverLogs, schedulerLog, controllerManagerLog, kubeletLog)

	fmt.Println()
}

func getNode(schedulerLog []string) string {
	schedulerLog = grep(schedulerLog, *podName, 0 /* linesAround */)
	re := regexp.MustCompile("bound successfully on node ([^,]+),")
	for _, line := range schedulerLog {
		groups := re.FindStringSubmatch(line)
		if len(groups) > 1 {
			return groups[1]
		}
	}
	return ""
}

func shorten(line string, pattern string) string {
	i := strings.Index(line, pattern)

	type rangeType struct{ a, b int }

	ranges := []rangeType{
		{0, 100},
		{i - *charsAround, i + len(pattern) + *charsAround},
		{len(line) - 50, len(line)},
	}

	// Compact ranges.
	curr := ranges[0]
	ranges2 := []rangeType{}
	for i := 1; i < len(ranges); i++ {
		this := ranges[i]
		if curr.b < this.a {
			ranges2 = append(ranges2, curr)
			curr = this
		} else {
			curr = rangeType{curr.a, this.b}
		}
	}
	ranges2 = append(ranges2, curr)

	ret := []string{}
	for _, r := range ranges2 {
		ret = append(ret, line[r.a:r.b])
	}

	return strings.Join(ret, " ... ")
}

func printAll(logs ...[]string) {
	allLogs := []string{}
	for _, log := range logs {
		allLogs = append(allLogs, log...)
	}
	sort.Strings(allLogs)
	print(allLogs, *podName)
}

func print(log []string, pattern string) {
	coloredPattern := "%h[fgGreen]" + pattern + "%r"
	for _, line := range log {
		if len(line) > 1000 {
			line = shorten(line, pattern)
		}
		line = strings.Replace(line, "%", "%%", -1)
		line = strings.Replace(line, pattern, coloredPattern, -1)
		color.Printf(line + "\n")
	}
}

// ------ Log Store ------- //

// LogStore is an util for reading log files from GCS. It provides a cache mechanism, backed by
// filesystem, so log files don't have to be fetched more than once, even across different runs.
// In addition it handles gzip compression in a transparent way.
type LogStore interface {
	Read(relativePath string) []string
	Grep(relativeDir, pattern string) []string
	List(relativeDir, pattern string) []string
}

func newLogStore(cacheRootDir, logURL string, linesAround int) LogStore {
	if !strings.HasSuffix(logURL, "/") {
		logURL += "/"
	}

	gcsDir := (logURL)[len("http://gcsweb.k8s.io/gcs"):]

	urlParts := strings.Split(logURL[:len(logURL)-1], "/")
	cacheDir := cacheRootDir + strings.Join(urlParts[len(urlParts)-2:], "/") + "/"

	return &logStore{
		logURL:      logURL,
		gcsDir:      gcsDir,
		cacheDir:    cacheDir,
		linesAround: linesAround,
	}
}

const (
	// Separator used in cached files, used to store info about pattern and number of lines around.
	Separator = "_-_"
	// DirFile is an artificial file, representing fetched directory listing, for caching purposes.
	DirFile = "/__dir"
)

type logStore struct {
	logURL      string
	gcsDir      string
	cacheDir    string
	linesAround int
}

func (l *logStore) Read(relativePath string) []string {
	fmt.Printf("Reading %s...\n", relativePath)

	if d := l.readFromCache(relativePath); d != nil {
		return d
	}

	var reader io.Reader
	if isGzip(relativePath) {
		reader = l.fetchGzip(relativePath)
	} else {
		if strings.HasSuffix(relativePath, DirFile) {
			reader = l.fetch(relativePath[:len(relativePath)-6])
		} else {
			reader = l.fetch(relativePath)
		}
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	l.writeToCache(relativePath, data)

	return strings.Split(string(data), "\n")
}

func (l *logStore) Grep(relativePath, pattern string) []string {
	pathWithPattern := relativePath + Separator + pattern
	if l.linesAround > 0 {
		pathWithPattern = fmt.Sprintf("%s%s%d", pathWithPattern, Separator, l.linesAround)
	}

	if d := l.readFromCache(pathWithPattern); d != nil {
		return d
	}

	lines := l.checkIfMoreLinesAroundInCache(relativePath, pattern)
	// Unfortunately no, read whole file.
	if lines == nil {
		lines = l.Read(relativePath)
	}

	lines = grep(lines, pattern, l.linesAround)

	l.writeToCache(pathWithPattern, []byte(strings.Join(lines, "\n")))

	return lines
}

func (l *logStore) List(relativeDir, pattern string) []string {
	lines := l.Read(relativeDir + DirFile)
	lines = grep(lines, l.gcsDir, 0 /* linesAround */)
	lines = grep(lines, pattern, 0 /* linesAround */)
	ret := []string{}
	for _, line := range lines {
		i := strings.Index(line, l.gcsDir)
		line = line[i+len(l.gcsDir):]
		i = strings.Index(line, "\"")
		if i == -1 {
			continue
		}
		line = line[:i]
		ret = append(ret, line)
	}
	return ret
}

func (l *logStore) checkIfMoreLinesAroundInCache(relativePath, pattern string) []string {
	path := l.cacheDir + relativePath + Separator + pattern + Separator
	files, err := filepath.Glob(path + "*")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		n, _ := strconv.Atoi(file[len(path):])
		if n > l.linesAround {
			return l.readFromCache(file[len(l.cacheDir):])
		}
	}

	return nil
}

func (l *logStore) writeToCache(relativePath string, data []byte) {
	path := l.getCachePath(relativePath)
	dir := filepath.Dir(path)

	os.MkdirAll(dir, 0755)

	f, err := os.Create(l.getCachePath(relativePath))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		panic(err)
	}
}

func (l *logStore) fetch(relativePath string) io.Reader {
	resp, err := http.Get(l.logURL + relativePath)
	if err != nil {
		panic(err)
	}
	return resp.Body
}

func (l *logStore) fetchGzip(relativePath string) io.Reader {
	body := l.fetch(relativePath)
	gr, err := gzip.NewReader(body)
	if err != nil {
		panic(err)
	}
	return gr
}

func (l *logStore) readFromCache(relativePath string) []string {
	path := l.getCachePath(relativePath)

	if _, err := os.Stat(path); err == nil {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		return strings.Split(string(data), "\n")
	}
	return nil
}

func (l *logStore) getCachePath(relativePath string) string {
	return l.cacheDir + relativePath
}

// ------ Utils -------- //

func isGzip(path string) bool {
	return strings.HasSuffix(path, ".gz")
}

func grep(lines []string, pattern string, linesAround int) []string {
	lineNumbers := map[int]bool{}

	for i, line := range lines {
		if strings.Contains(line, pattern) {
			start := i - linesAround
			if start < 0 {
				start = 0
			}
			for j := start; j <= i+linesAround && j < len(lines); j++ {
				lineNumbers[j] = true
			}
		}
	}

	i := 0
	ret := make([]string, len(lineNumbers), len(lineNumbers))
	for l := range lineNumbers {
		ret[i] = lines[l]
		i++
	}
	return ret
}