LETHE_VERSION=v0.1.4

run-dev:
	air

git-push:
	git add -A; git commit -am ${LETHE_VERSION}; git push

docker-build-dev:
	docker build -t ghcr.io/kuoss/lethe:dev -f Dockerfile.dev . && docker push ghcr.io/kuoss/lethe:dev

docker-build:
	docker build -t ghcr.io/kuoss/lethe:${LETHE_VERSION} --build-arg LETHE_VERSION=${LETHE_VERSION} . && docker push ghcr.io/kuoss/lethe:${LETHE_VERSION}
