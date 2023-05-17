package graceful_test

import (
	"context"
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/micro/grpcx"
	pb "github.com/bang-go/crab/examples/proto/echo"
	"github.com/bang-go/util"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

type greeterWrapper struct {
	pb.UnimplementedEchoServer
}

func (g *greeterWrapper) SayHello(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Println(req.Message)
	return &pb.EchoResponse{Message: "hello " + req.Message}, nil
}

func (g *greeterWrapper) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	for {
		_ = stream.Send(&pb.EchoResponse{Message: util.IntToString(util.IntRandRange(1, 100))})
		time.Sleep(1 * time.Second)
	}
	return nil
}
func TestGrpcServer(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.EncodeConsole), crab.WithLogAllowLevel(logx.InfoLevel))
	cmder := cmd.New(&cmd.Config{CmdUse: "grpc-server", CmdShort: "serve grpc server"})
	cmder.SetRun(func(args []string) {
		server := grpcx.NewServer(&grpcx.ServerConfig{Addr: ":8081"})
		_ = server.Start(func(server *grpc.Server) {
			pb.RegisterEchoServer(server, &greeterWrapper{})
		})
	})
	crab.RegisterCmd(cmder)
	err = crab.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrpcClient(t *testing.T) {
	crab.Build(crab.WithLogEncoding(logx.EncodeConsole), crab.WithLogAllowLevel(logx.InfoLevel))
	err := crab.Start()
	if err != nil {
		log.Fatal(err)
	}
	client := grpcx.NewClient(&grpcx.ClientConfig{Addr: ":8081"})
	defer client.Close()
	streamReply, err := client.DialWithCall(func(conn *grpc.ClientConn) (any, error) {
		client := pb.NewEchoClient(conn)
		return client.ServerStreamingEcho(context.Background(), &pb.EchoRequest{Message: "crab"})
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	stream := streamReply.(pb.Echo_ServerStreamingEchoClient)
	for {
		res, err := stream.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(res)
	}
}
