#!/bin/bash

set -e
set -u

export IMAGE_NAME=landzero/cutee
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

binfs views public locales > cmd/cutee/binfs.gen.go

go build -o cutee.out ./cmd/cutee
go build -o cutee-maint.out ./cmd/cutee-maint

docker build -t $IMAGE_NAME .
