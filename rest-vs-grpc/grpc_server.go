// gRPCのサーバー
package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "path/to/your/protofile" // Protocol Buffersの定義
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Message: "Hello, " + in.Name}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
