#!/bin/bash
cd $(dirname $0)/../

which gocyclo || go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gocyclo -over 15 -ignore letheql/parser .
if [[ $? != 0 ]]; then
    echo "❌ FAIL"
    exit 1
fi
echo "✔️ OK"
