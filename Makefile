COVERAGE_THRESHOLD = 65
PROMETHEUS_VERSION := v2.42.0

install-dev:
	go mod tidy
	which air || go install github.com/cosmtrek/air@latest

run-dev:
	air

.PHONY: test
test:
	go test -race --failfast ./...


## checks
.PHONY: checks
checks: cover build lint licenses

.PHONY: cover
cover: test
	@echo "Running tests and checking coverage..."
	@go test -coverprofile=coverage.out ./...
	@cat coverage.out | grep -v yaccpar > coverage2.out
	@go tool cover -func=coverage2.out | grep total | awk '{print $$3}' | sed 's/%//' | \
	awk -v threshold=$(COVERAGE_THRESHOLD) '{ if ($$1 < threshold) { print "Coverage is below threshold: " $$1 "% < " threshold "%"; exit 1 } else { print "Coverage is sufficient: " $$1 "%" } }'

.PHONY: build
build:
	go build -o bin/lethe

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

.PHONY: licenses
licenses: go-licenses
	$(GO_LICENSES) check .


## docker
docker: 
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) .


#### parser
.PHONY: parser-clone
parser-clone:
	git clone -b $(PROMETHEUS_VERSION) --depth=1 https://github.com/prometheus/prometheus.git
	rm -rf letheql/parser
	mv prometheus/promql/parser letheql/
	rm -rf prometheus

.PHONY: parser-generate
parser-generate: goyacc
	$(GOYACC) -o letheql/parser/generated_parser.y.go letheql/parser/generated_parser.y
	rm -f letheql/parser/y.output

.PHONY: parser-test
parser-test:
	go test -failfast github.com/kuoss/lethe/letheql/parser
	go test -failfast github.com/kuoss/lethe/letheql/parser_test


##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOYACC ?= $(LOCALBIN)/goyacc
GO_LICENSES ?= $(LOCALBIN)/go-licenses

## Tool Versions
GOLANGCI_LINT_VERSION ?= v1.60.2
GOYACC_VERSION ?= v0.3.0
GO_LICENSES_VERSION ?= v1.6.0

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT)
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: goyacc
goyacc: $(GOYACC)
$(GOYACC): $(LOCALBIN)
	$(call go-install-tool,$(GOYACC),golang.org/x/tools/cmd/goyacc,$(GOYACC_VERSION))

.PHONY: go-licenses
go-licenses: $(GO_LICENSES)
$(GO_LICENSES): $(LOCALBIN)
	$(call go-install-tool,$(GO_LICENSES),github.com/google/go-licenses,$(GO_LICENSES_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef
