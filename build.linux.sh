#!/bin/bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o build/powerline.linux.amd64
GOOS=linux GOARCH=386   CGO_ENABLED=0 go build -v -o build/powerline.linux.386
