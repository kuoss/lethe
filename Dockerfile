FROM golang:1.20-alpine AS base
ARG VERSION
WORKDIR /temp/
COPY . ./
RUN go mod download -x
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /app/bin/lethe
RUN cp -a ./etc                                       /app/etc

FROM alpine:3.18
COPY --from=base /app /app
RUN set -x && apk add --no-cache coreutils util-linux curl grep

WORKDIR /app
ENTRYPOINT ["/app/bin/lethe"]
