#!/bin/bash
cd $(dirname $0)/..

which go-licenses || go install github.com/google/go-licenses@v1.6.0

echo go-licenses report
go-licenses report . | tee docs/go-licenses.csv

echo go-licenses check
go-licenses check .
if [[ $? != 0 ]]; then
    echo "❌ FAIL"
    exit 1
fi
echo "✔️ OK"
