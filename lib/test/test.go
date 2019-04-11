package main

import (
	"fmt"
	"github.com/mm4tt/k8s-util/lib/logs"
)

func main() {
	prettifier := logs.DefaultPrettifier()

	fmt.Println(prettifier.Prettify("ala ma kota", "ma"))
}
