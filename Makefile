LETHE_VERSION=v0.2.0-dev.1

install-dev:
	go mod tidy

run-dev:
	air

docker-build:
	docker build -t ghcr.io/kuoss/lethe:${LETHE_VERSION} --build-arg LETHE_VERSION=${LETHE_VERSION} . && docker push ghcr.io/kuoss/lethe:${LETHE_VERSION}



test:
	go test ./... -failfast

test-all:
	scripts/go_test_all_packages_failfast.sh

test-win:
	.\scripts\go_test_all_packages_failfast.bat

test-cover:
	@./scripts/test-cover.sh



pre-checks:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/google/go-licenses@latest

checks: fmt vet staticcheck golangci-lint go-licenses-check test-cover

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

golangci-lint:
	golangci-lint run

go-licenses-check:
	go-licenses check  github.com/kuoss/lethe 2> /dev/null
	go-licenses report github.com/kuoss/lethe 2> /dev/null | tee docs/go-licenses.csv

mock:
	./scripts/mock/restart.sh

mock-status:
	./scripts/mock/status.sh

mock-logs:
	./scripts/mock/logs.sh

mock-delete:
	./scripts/mock/delete.sh

