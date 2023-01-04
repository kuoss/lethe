FROM golang:1.18-alpine AS go
WORKDIR /temp/
COPY . ./
RUN go mod download -x
RUN go build -o           /app/bin/lethe
RUN cp -a ./etc           /app/etc
RUN cd cli && go build -o /app/bin/lethetool

FROM alpine:3.15
COPY --from=go /app /app

ARG LETHE_VERSION
ENV LETHE_VERSION=$LETHE_VERSION

RUN echo "http://nexus.sdsdev.co.kr:8081/repository/apk-alpine/v3.11/main" | tee /etc/apk/repositories \
    echo "http://nexus.sdsdev.co.kr:8081/repository/apk-alpine/v3.11/community" | tee -a /etc/apk/repositories

RUN set -x \
&& apk add --no-cache coreutils util-linux

WORKDIR  /app
ENTRYPOINT ["/app/bin/lethe"]
