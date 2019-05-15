package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/mm4tt/k8s-util/lib/logs"
	"k8s.io/klog"
)

var (
	logFile = flag.String("log-file", "", "Log file to parse")

	namespace = flag.String("namespace", "", "namespace to look for in the log")
	name      = flag.String("name", "", "name to look for in the log")

	smartCompact = flag.Bool("smart-compact", true, "Wheter to compact the same log lines")

	namespaceRegex, nameRegex, stepRegex *regexp.Regexp
)

func main() {
	klog.Info("I'm log-parser.")

	flag.Parse()
	validateFlags()

	namespaceRegex = regexp.MustCompile(*namespace + "[^0-9]")
	nameRegex = regexp.MustCompile(*name + "[^0-9]")
	stepRegex = regexp.MustCompile("Step \"")

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

	prettifier := logs.DefaultPrettifier()

	var hasSkippedLastLine bool
	var lastPrinted, prevLine string

	fmt.Println()
	fmt.Println()
	fmt.Println()
	for scanner.Scan() {
		nLines++
		if !matchesNamespaceAndName(scanner.Bytes()) {
			if stepRegex.Match(scanner.Bytes()) {
				fmt.Println(prettifier.Prettify(scanner.Text(), "Step"))
			}
			continue
		}
		nMatched++
		line := scanner.Text()

		if !*smartCompact {
			fmt.Println(prettifier.Prettify(line, *name))
			continue
		}

		withoutDate := line[len("I0410 02:34:51.241928"):]
		if lastPrinted == "" || withoutDate != lastPrinted {
			lastPrinted = withoutDate
			if hasSkippedLastLine {
				fmt.Println(prettifier.Prettify(prevLine, *name))
			}
			fmt.Println(prettifier.Prettify(line, *name))
		} else {
			hasSkippedLastLine = true
		}
		prevLine = line
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


