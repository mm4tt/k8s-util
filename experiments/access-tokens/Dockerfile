FROM golang:1.13.4 AS build-env

ARG gopkg=github.com/mm4tt/k8s-util/access-tokens

ADD ["cmd", "/go/src/$gopkg/cmd"]
#ADD ["pkg", "/go/src/$gopkg/pkg"]
ADD ["go.mod", "/go/src/$gopkg"]
ADD ["go.sum", "/go/src/$gopkg"]

WORKDIR /go/src/$gopkg
RUN CGO_ENABLED=0 go build -o /workspace/main.bin ./cmd

FROM golang:1.13.4-alpine
WORKDIR /workspace
COPY --from=build-env /workspace/main.bin .
ENTRYPOINT ["/workspace/main.bin"]

