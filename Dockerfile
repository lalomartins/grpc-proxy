FROM golang:1.8-alpine

# forked from https://github.com/devsu/docker-grpc-proxy to build a more
# recent version which proxies the metadata (otherwise we can't authenticate)

RUN apk update && apk upgrade && \
    apk add --no-cache git && \
    rm -rf /var/cache/apk/*

WORKDIR "/go/src/github.com/devsu/grpc-proxy"
COPY . .
RUN go get -d -v ./... \
 && go install -v ./... \
 && cd \
 && rm -rf /go/src

RUN mkdir -p "/opt/grpc-proxy"
VOLUME ["/opt/grpc-proxy"]
WORKDIR "/opt/grpc-proxy"

ENTRYPOINT ["grpc-proxy"]
