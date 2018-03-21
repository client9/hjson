#!/bin/sh
set -ex

go build ./...
gometalinter \
    --vendor \
    --deadline=60s \
    --disable-all \
    --enable=vet \
    --enable=golint \
    --enable=gofmt \
    --enable=goimports \
    --enable=gosimple \
    --enable=staticcheck \
    --enable=ineffassign \
    ./...
go test ./...
