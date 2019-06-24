package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

var commit = flag.String("commit", "master", "Commit to use in the fileUrl address")
var fileUrl = flag.String("file-url", "https://raw.githubusercontent.com/kubernetes/test-infra/$commit/config/jobs/kubernetes/sig-scalability/sig-scalability-presets.yaml", "URL of the yaml file with presets to read")
var presetName = flag.String("preset-name", "preset-e2e-scalability-common", "Name of the preset to load")

func main() {
	flag.Parse()

	url := strings.ReplaceAll(*fileUrl, "$commit", *commit)
	resp, err := http.Get(url)
	if err != nil {	panic(err) }
	yamlFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {	panic(err) }

	config := new(conf)
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}

	for _, preset := range config.Presets {
		name := preset.getName()

		if name == *presetName {
			for _, env := range preset.Env {
				fmt.Printf("export %s=\"%s\"\n", env.Name, env.Value)
			}
		}
	}
}

type conf struct {
	Presets []preset `yaml:presets`
}

type preset struct {
	Labels map[string]string `yaml:labels`
	Env []env `yaml:env`
}

func (p preset) getName() string {
	for name, _ := range p.Labels {
		return name
	}
	panic("Empty Labels!")
}

type env struct {
	Name string `yaml:name`
	Value string `yaml:value`
}