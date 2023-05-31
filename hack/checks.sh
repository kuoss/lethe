#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go mod tidy
go fmt ./...
./hack/misspell.sh
go vet ./...
which goimports || go install golang.org/x/tools/cmd/goimports@latest
goimports -local -v -w .
which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --timeout 5m
./hack/test-cover.sh
./hack/go-licenses.sh
