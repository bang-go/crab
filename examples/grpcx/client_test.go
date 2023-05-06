package grpcx_test

import (
	"context"
	"fmt"
	"github.com/bang-go/crab/core/micro/grpcx"
	pb "github.com/bang-go/crab/examples/grpcx/helloworld"
	"google.golang.org/grpc"
	"testing"
)

func TestUnaryClient(t *testing.T) {
	rpcClient := grpcx.NewClient(&grpcx.ClientConfig{Addr: "127.0.0.1:8081"})
	defer rpcClient.Close()
	var concurrency = 100
	for i := 0; i < concurrency; i++ {
		reply, err := rpcClient.DialWithCall(func(conn *grpc.ClientConn) (any, error) {
			client := pb.NewGreeterClient(conn)
			return client.SayHello(context.Background(), &pb.HelloRequest{Name: "jk"})
		})
		if err != nil {
			//log.Fatal(err)
		}
		fmt.Println(reply.(*pb.HelloReply))
	}
}
