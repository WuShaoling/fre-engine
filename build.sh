#!/bin/bash

# build free (fre engine)
docker run -it --rm -v "$PWD":/go/src golang:1.14 \
  bash -c "cd /go/src && go build -o free main.go"

scp free root@server:/usr/local/bin

rm -f free
