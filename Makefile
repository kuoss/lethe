VERSION := v0.2.0
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
