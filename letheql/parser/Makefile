goyacc:
	set -x;\
	which goyacc || go install golang.org/x/tools/cmd/goyacc@v0.6.0;\
	goyacc -o generated_parser.y.go generated_parser.y;\
	rm -f y.output;\
	go test -failfast .
