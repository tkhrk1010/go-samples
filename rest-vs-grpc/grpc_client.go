package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    pb "path/to/your/protofile" // Protocol Buffersの定義
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
		// gRPCのクライアント
		// gRPCの場合は、どういうinterfaceを持つかは、Protocol Buffersの定義による
		// .protoファイルに定義されており、これをコンパイルして生成されたコードを利用する
    c := pb.NewGreeterClient(conn)

    // gRPCのリクエスト
		// SayHelloはserverで定義されている。これをclinetで呼び出すのがRESTとの違いとして大きい
    r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
    if err != nil {
        log.Fatalf("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.Message)
}
