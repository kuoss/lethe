#!/bin/bash
cd $(dirname $0)/..

[ -f ./bin/misspell ] || curl -L https://git.io/misspell | bash

find . -type f -name '*.*' | xargs ./bin/misspell -error
if [[ $? != 0 ]]; then
    echo "❌ FAIL - misspell found"
    exit 1
fi
echo "✔️ OK - misspell not found"
