#!/bin/bash
cd $(dirname $0)/..
set -x

PROMPARSER=letheql/promql/promparser

set -x
rm -rf        $PROMPARSER
cp -a $PARSER $PROMPARSER
sed 's|package parser|package promparser|g' -i $PROMPARSER/*.*

which goyacc || go install golang.org/x/tools/cmd/goyacc@v0.6.0
goyacc -o $PROMPARSER/generated_parser.y.go $PROMPARSER/generated_parser.y
rm -f y.output

cd $PROMPARSER
go test -failfast .

# grep PIPE_ generated_parser.y.go
# go mod tidy
# go test -failfast .
# cd ../
# go test -failfast .
