#!/bin/bash

cd $(dirname $0)/../

for s in $(go list ./...); do
    if ! go test -failfast -v $s; then
        break
    fi
done
