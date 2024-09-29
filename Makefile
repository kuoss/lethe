PROMETHEUS_VER := v2.42.0
VERSION := v0.2.1
IMAGE := ghcr.io/kuoss/lethe:$(VERSION)

install-dev:
	go mod tidy
	which air || go install github.com/cosmtrek/air@latest

run-dev:
	air

test:
	hack/test-failfast.sh

test-win:
	.\hack\test-failfast.bat

cover:
	hack/test-cover.sh

misspell:
	hack/misspell.sh

gocyclo:
	hack/gocyclo.sh

checks:
	hack/checks.sh

# go build
build:
	go build -ldflags="-X 'main.Version=$(VERSION)'" -o bin/lethe

# docker build & push
docker: 
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) . && docker push $(IMAGE)


mock:
	hack/mock/restart.sh

mock-status:
	hack/mock/status.sh

mock-logs:
	hack/mock/logs.sh

mock-delete:
	hack/mock/delete.sh

## letheql
.PHONY: parser
parser:
	git clone -b $(PROMETHEUS_VER) --depth=1 https://github.com/prometheus/prometheus.git
	rm -rf letheql/parser
	mv prometheus/promql/parser letheql/
	rm -rf prometheus

.PHONY: goyacc
goyacc:
	which goyacc || go install golang.org/x/tools/cmd/goyacc@v0.6.0
	goyacc -o letheql/parser/generated_parser.y.go letheql/parser/generated_parser.y
	rm -f letheql/parser/y.output
	go test -failfast letheql/parser/
