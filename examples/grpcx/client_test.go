package grpcx_test

import (
	"context"
	"fmt"
	pb "github.com/bang-go/crab/examples/proto/echo"
	"github.com/bang-go/network/grpcx"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

func TestUnaryClient(t *testing.T) {
	rpcClient := grpcx.NewClient(&grpcx.ClientConfig{Addr: "127.0.0.1:8081"})
	defer rpcClient.Close()
	var concurrency = 100
	for i := 0; i < concurrency; i++ {
		reply, err := rpcClient.DialWithCall(func(conn *grpc.ClientConn) (any, error) {
			client := pb.NewEchoClient(conn)
			return client.UnaryEcho(context.Background(), &pb.EchoRequest{Message: "jk"})
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply.(*pb.EchoResponse))
		time.Sleep(1 * time.Second)
	}
}

func TestClientKeepalive(t *testing.T) {
	rpcClient := grpcx.NewClient(&grpcx.ClientConfig{Addr: "127.0.0.1:8081"})
	defer rpcClient.Close()
	conn, err := rpcClient.Dial()
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewEchoClient(conn)
	streamReply, err := client.ServerStreamingEcho(context.Background(), &pb.EchoRequest{Message: "jk"})
	if err != nil {
		log.Fatal(err)
	}
	for {
		reply, err := streamReply.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reply)
		time.Sleep(time.Minute)
	}
}
