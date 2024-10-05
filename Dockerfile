FROM golang:1.20-alpine AS builder
ARG VERSION
WORKDIR /build
COPY . ./
RUN go mod download -x
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /build/bin/lethe ./cmd/lethe/

# 2023-11 latest
FROM alpine:3.18.4
COPY --from=builder /build/bin/lethe /app/bin/lethe
COPY --from=builder /build/etc /app/etc
RUN apk add --no-cache coreutils util-linux curl grep

WORKDIR /app
ENTRYPOINT ["/app/bin/lethe"]
