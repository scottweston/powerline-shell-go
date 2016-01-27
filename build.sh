#!/bin/bash
cwd=`pwd`
app=${cwd#*${GOPATH}/}
docker run --rm -it -v "$(pwd)":/go/${app} -w /go/${app} golang:1.4.0-cross ./build.linux.sh
CC=gcc GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -v -o build/powerline.darwin.amd64 -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`"
CC=gcc GOOS=darwin GOARCH=386 CGO_ENABLED=1 go build -v -o build/powerline.darwin.386 -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`"
