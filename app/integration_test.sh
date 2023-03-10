#!/bin/sh

go test $(go list ./... | grep /integration_test) -cover