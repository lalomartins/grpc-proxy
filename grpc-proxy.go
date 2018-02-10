package main

import (
  "github.com/mwitkow/grpc-proxy/proxy"
  "google.golang.org/grpc"
  "log"
  "net"
  "net/http"
  "fmt"
  "os"
  "github.com/devsu/grpc-proxy/extras"
  "google.golang.org/grpc/credentials"
  "google.golang.org/grpc/grpclog"
  "golang.org/x/net/trace"
)

func main() {
  configurationFile := "./config.json"

  args := os.Args[1:]
  if len(args) > 0 {
    configurationFile = args[0]
  }

  config := extras.GetConfiguration(configurationFile)

  if (config.Trace) {
    grpc.EnableTracing = true
    // Start the default http server and explicitly bind it to listen on localhost for security purposes
    // Accessing http://localhost:6060/debug/events or http://localhost:6060/debug/requests will show the
    // currently gathered traces.
    trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
      return true, true
    }

    go http.ListenAndServe("0.0.0.0:6060", nil)
  }

  listen := ":50051"
  if config.Listen != "" {
    listen = config.Listen
  }

  lis, err := net.Listen("tcp", listen)

  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  fmt.Printf("Proxy running at %q\n", listen)

  server := GetServer(config)

  if err := server.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}

func GetServer (config extras.Config) *grpc.Server {
  var opts []grpc.ServerOption

  opts = append(opts, grpc.CustomCodec(proxy.Codec()),
    grpc.UnknownServiceHandler(proxy.TransparentHandler(extras.GetDirector(config))))

  if config.CertFile != "" && config.KeyFile != "" {
    creds, err := credentials.NewServerTLSFromFile(config.CertFile, config.KeyFile)
    if err != nil {
      grpclog.Fatalf("Failed to generate credentials %v", err)
    }
    opts = append(opts, grpc.Creds(creds))
  }

  return grpc.NewServer(opts...)
}
