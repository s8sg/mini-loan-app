#!/bin/sh

test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"
go build .