#!/bin/bash

cd $(dirname $0)/../
for s in $(go list ./...); do
    echo =============== $s ===============
    go test -failfast -v $s || exit 1
done
