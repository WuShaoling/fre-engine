#!/bin/bash

# build free
docker run -it --rm -v "$PWD":/go/src golang:1.14 \
  bash -c "cd /go/src && go build -o free main.go"
scp free root@server:/free/
rm -f free

# build pause
docker run -it --rm -v "$PWD":/go/src golang:1.14 \
  bash -c "cd /go/src && go build -o pause pause.go"
scp pause root@server:/free/workspace/runtime/python3.7/bin/
rm -f pause
