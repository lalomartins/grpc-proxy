FROM grpc/go:1.0

# forked from https://github.com/devsu/docker-grpc-proxy to build a more
# recent version which proxies the metadata (otherwise we can't authenticate)

WORKDIR "/go/src/github.com"
RUN mkdir mwitkow \
 && git clone https://github.com/mwitkow/grpc-proxy.git mwitkow/grpc-proxy \
 && cd mwitkow/grpc-proxy \
 && git checkout 97396d94749c00db659393ba5123f707062f829f

WORKDIR "/go/src/google.golang.org/grpc/"
RUN git fetch && git checkout v1.3.x

WORKDIR "/go/src/github.com/devsu/grpc-proxy"
COPY . .
RUN go get -d -v ./... \
 && go install -v ./...
 # && cd \
 # && rm -rf /go/src

RUN mkdir -p "/opt/grpc-proxy"
VOLUME ["/opt/grpc-proxy"]
WORKDIR "/opt/grpc-proxy"

ENTRYPOINT ["grpc-proxy"]
