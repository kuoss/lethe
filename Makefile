LETHE_VERSION=v0.1.5

run-dev:
	air

docker-build-dev:
	docker build -t ghcr.io/kuoss/lethe:dev -f Dockerfile.dev . && docker push ghcr.io/kuoss/lethe:dev

docker-build:
	docker build -t ghcr.io/kuoss/lethe:${LETHE_VERSION} --build-arg LETHE_VERSION=${LETHE_VERSION} . && docker push ghcr.io/kuoss/lethe:${LETHE_VERSION}

check: fmt vet staticcheck test

fmt:
	go fmt ./...

vet:
	go vet ./...

staticcheck:
	# go install honnef.co/go/tools/cmd/staticcheck@latest
	/root/go/bin/staticcheck ./...

test:
	go test ./... -failfast -cover

test-all:
	scripts/go_test_all_packages_failfast.sh

test-win:
	.\scripts\go_test_all_packages_failfast.bat


go-licenses: go-licenses-report go-licenses-check

go-licenses-report:
	# go install github.com/google/go-licenses@latest
	go-licenses report github.com/kuoss/lethe | tee docs/go-licenses.csv

go-licenses-check:
	go-licenses check github.com/kuoss/lethe && echo OK

mock:
	./scripts/mock/restart.sh

mock-status:
	./scripts/mock/status.sh

mock-logs:
	./scripts/mock/logs.sh

mock-delete:
	./scripts/mock/delete.sh

