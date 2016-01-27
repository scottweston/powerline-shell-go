#!/bin/bash
go install -v -ldflags "-X main.build=`date -u +%Y%m%d.%H%M%S`"
