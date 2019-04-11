package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"k8s.io/klog"
	"github.com/mm4tt/k8s-util/lib/logs"
)

var (
	logFile = flag.String("log-file", "", "Log file to parse")

	namespace = flag.String("namespace", "", "namespace to look for in the log")
	name      = flag.String("name", "", "name to look for in the log")

	namespaceRegex, nameRegex *regexp.Regexp
)

func main() {
	klog.Info("I'm log-parser.")

	flag.Parse()
	validateFlags()

	namespaceRegex = regexp.MustCompile(*namespace + "[^0-9]")
	nameRegex = regexp.MustCompile(*name + "[^0-9]")

	klog.Infof("Processing file: %s", *logFile)
	if err := processLog(*logFile); err != nil {
		klog.Fatal(err)
	}
}

func validateFlags() {
	validateStringFlagNotEmpty(logFile, "log-file")
	validateStringFlagNotEmpty(name, "namespace")
	validateStringFlagNotEmpty(namespace, "name")
}

func validateStringFlagNotEmpty(flag *string, name string) {
	if *flag == "" {
		klog.Fatalf("--%s cannot be empty!", name)
	}
}

func processLog(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1000), 100*1024*1024)
	nLines, nMatched := 0, 0

	for scanner.Scan() {
		nLines++
		if matchesNamespaceAndName(scanner.Bytes()) {
			nMatched++
			fmt.Println(logs.Shorten(scanner.Text(), *name, 50))
		}
	}
	klog.Infof("Processed %d lines, matched %d lines", nLines, nMatched)
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func matchesNamespaceAndName(line []byte) bool {
	return namespaceRegex.Match(line) && nameRegex.Match(line)
}


