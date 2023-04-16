VERSION := v0.2.0-beta.1
IMAGE := ghcr.io/kuoss/lethe:$(VERSION)

install-dev:
	go mod tidy
	which air || go install github.com/cosmtrek/air@latest

run-dev:
	air

test:
	hack/go-test-failfast.sh

test-win:
	.\hack\go-test-failfast.bat

checks:
	hack/checks.sh

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
