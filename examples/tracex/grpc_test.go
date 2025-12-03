package tracex_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/tracex"
	pb "github.com/bang-go/crab/examples/proto/echo"
	"github.com/bang-go/micro/grpcx"
	"github.com/bang-go/micro/grpcx/metadatax"
	"github.com/bang-go/util"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

type greeterEntity struct {
	pb.UnimplementedEchoServer
}

func (g *greeterEntity) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	tracer := tracex.Tracer("grpc-unary-echo")
	ctx, span := tracer.NewSpanWithContext(ctx, "server-echo-step1")
	span.SetAttributes(attribute.String("k", "v"))
	defer span.End()
	time.Sleep(2 * time.Second)
	//todo: 打印grpc metadata 获取到传播key-value
	md := metadatax.ExtractIncoming(ctx)
	log.Println(md)
	return &pb.EchoResponse{Message: "hello " + req.Message}, nil
}

func (g *greeterEntity) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	var num = 100
	for i := 0; i < num; i++ {
		_ = stream.Send(&pb.EchoResponse{Message: util.IntToString(i)})
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestGrpcServer(t *testing.T) {
	defer crab.Close()
	greeter := &greeterEntity{}
	rpcServer := grpcx.NewServer(&grpcx.ServerConfig{Addr: ":8081", Trace: true})
	err := rpcServer.Start(func(server *grpc.Server) {
		pb.RegisterEchoServer(server, greeter)
	})
	if err != nil {
		log.Fatal(err)
	}
}
