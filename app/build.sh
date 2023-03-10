#!/bin/sh

test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"
go test $(go list ./... | grep /integration_test) -cover
go build .