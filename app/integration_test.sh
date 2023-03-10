#!/bin/sh

go test $(go list ./integration_test/...) -cover