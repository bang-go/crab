package crab_test

import (
	"context"
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/micro/grpcx"
	"google.golang.org/grpc"
	"log"
	"testing"
)
import pb "github.com/bang-go/crab/examples/proto/echo"

type greeterWrapper struct {
	pb.UnimplementedEchoServer
}

func (g *greeterWrapper) SayHello(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: "hello " + req.Message}, nil
}
func TestGrpcServer(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.LogEncodeJson), crab.WithLogAllowLevel(logx.LevelInfo))

	cmder := cmd.New(&cmd.Config{CmdUse: "grpc-server", CmdShort: "serve grpc server"})
	cmder.SetRun(func(args []string) {
		server := grpcx.NewServer(&grpcx.ServerConfig{Addr: ":8081"})
		_ = server.Start(func(server *grpc.Server) {
			pb.RegisterEchoServer(server, &greeterWrapper{})
		})
	})
	crab.RegisterCmd(cmder)
	err = crab.Start()
	defer crab.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrpcClient(t *testing.T) {
	crab.Build(crab.WithLogEncoding(logx.LogEncodeJson), crab.WithLogAllowLevel(logx.LevelInfo))
	cmder := cmd.NewWithRunFunc(&cmd.Config{CmdUse: "grpc-client", CmdShort: "serve grpc client"}, func(args []string) {
		var err error
		client := grpcx.NewClient(&grpcx.ClientConfig{Addr: ":8081"})
		defer client.Close()
		reply, err := client.DialWithCall(func(conn *grpc.ClientConn) (any, error) {
			client := pb.NewEchoClient(conn)
			return client.UnaryEcho(context.Background(), &pb.EchoRequest{Message: "crab"})
		})
		if err != nil {
			return
		}
		log.Println(reply.(*pb.EchoResponse))
	})
	crab.RegisterCmd(cmder)
	err := crab.Start()
	defer crab.Close()
	if err != nil {
		log.Fatal(err)
	}
}
