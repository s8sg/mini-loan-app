#!/bin/sh

go test $(go list ./... | grep -v /vendor/ | grep -v /integration_test) -cover