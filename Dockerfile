FROM golang:1.22-alpine AS builder
ARG VERSION
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /build/bin/lethe ./cmd/lethe/

FROM alpine:3.20.6
RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
    coreutils \
    curl \
    grep \
    util-linux \
  && rm -rf /var/cache/apk/*

COPY --from=builder /build/bin/lethe /app/bin/lethe
COPY --from=builder /build/etc       /app/etc

WORKDIR    /app
USER       nobody
EXPOSE     3030
ENTRYPOINT ["/app/bin/lethe"]
