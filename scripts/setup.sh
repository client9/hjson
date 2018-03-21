#!/bin/sh
set -ex
./scripts/godownloader-goreleaser.sh
go get -u github.com/alecthomas/gometalinter && gometalinter --install
