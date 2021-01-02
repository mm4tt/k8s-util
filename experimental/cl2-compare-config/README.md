# CL2 Compare Config

Tool for comparing CL2 load test config changes introduced in a given PR against
the master branch. The tool purpose is to ease the migration of CL2 configs to 
modules. The tool expands CL2 config by flatenning the modules and them compares
list of executable steps.

## How to install

```
go get -d github.com/mm4tt/k8s-util/experimental/cl2-compare-config
```

## How to run

```
$GOPATH/src/github.com/mm4tt/k8s-util/experimental/cl2-compare-config/run.sh <pr-num, e.g. 1640>
```