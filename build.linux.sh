#!/bin/bash
go build -v -o build/powerline-shell-go -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`"
