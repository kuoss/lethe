#!/bin/bash
MIN_COVER=50.0

cd $(dirname $0)/..

if [ ! -f cover.out ]; then
    echo
    echo "❌ cover.out file not found"
    echo "run 'make test-cover' first"
    echo
    exit 1
fi

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
