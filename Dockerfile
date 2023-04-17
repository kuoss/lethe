FROM golang:1.19-alpine AS base
ARG VERSION
WORKDIR /temp/
COPY . ./
RUN go mod download -x
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /app/bin/lethe
RUN cp -a ./etc                                       /app/etc
RUN cd cli && go build -o                             /app/bin/lethetool

FROM alpine:3.15
COPY --from=base /app /app
RUN set -x \
&& apk add --no-cache coreutils util-linux

WORKDIR  /app
ENTRYPOINT ["/app/bin/lethe"]
