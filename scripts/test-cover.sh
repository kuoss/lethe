#!/bin/bash
MIN_COVER=50.0

cd $(dirname $0)/..
go test ./... -coverprofile cover.out
COVER=$(go tool cover -func cover.out | tail -1 | grep -oP [0-9.]+)
if [[ $COVER < $MIN_COVER ]]; then
    echo
    echo "❌ FAIL ( test coverage: $COVER% < $MIN_COVER% )"
    echo
    exit 2
fi
echo
echo "✔️ PASS ( test coverage: $COVER% >= $MIN_COVER% )"
echo
