#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
go run $DIR/resource-usage-summary-parser.go --stderrthreshold=INFO $@
