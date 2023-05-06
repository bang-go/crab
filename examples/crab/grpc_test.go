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
import pb "github.com/bang-go/crab/examples/grpcx/helloworld"

type greeterWrapper struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterWrapper) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + req.Name}, nil
}
func TestGrpcServer(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.EncodeConsole), crab.WithLogAllowLevel(logx.InfoLevel))

	cmder := cmd.New(&cmd.Config{CmdUse: "grpc-server", CmdShort: "serve grpc server"})
	cmder.SetRun(func(args []string) error {
		server := grpcx.NewServer(&grpcx.ServerConfig{Addr: ":8081"})
		return server.Start(func(server *grpc.Server) {
			pb.RegisterGreeterServer(server, &greeterWrapper{})
		})
	})
	crab.AddCmd(cmder)
	err = crab.Start()
	defer crab.Exit()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrpcClient(t *testing.T) {
	crab.Build(crab.WithLogEncoding(logx.EncodeConsole), crab.WithLogAllowLevel(logx.InfoLevel))
	cmder := cmd.NewWithRunFunc(&cmd.Config{CmdUse: "grpc-client", CmdShort: "serve grpc client"}, func(args []string) error {
		var err error
		client := grpcx.NewClient(&grpcx.ClientConfig{Addr: ":8081"})
		defer client.Close()
		reply, err := client.DialWithCall(func(conn *grpc.ClientConn) (any, error) {
			client := pb.NewGreeterClient(conn)
			return client.SayHello(context.Background(), &pb.HelloRequest{Name: "crab"})
		})
		if err != nil {
			return err
		}
		log.Println(reply.(*pb.HelloReply))
		return nil
	})
	crab.AddCmd(cmder)
	err := crab.Start()
	defer crab.Exit()
	if err != nil {
		log.Fatal(err)
	}
}
