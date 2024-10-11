FROM golang:1.21-alpine AS builder
ARG VERSION
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /build/bin/lethe ./cmd/lethe/

FROM alpine:3.20.3
COPY --from=builder /build/bin/lethe /app/bin/lethe
COPY --from=builder /build/etc /app/etc
RUN apk add --no-cache coreutils util-linux curl grep

WORKDIR /app
ENTRYPOINT ["/app/bin/lethe"]
