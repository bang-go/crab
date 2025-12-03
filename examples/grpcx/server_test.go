package grpcx_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/bang-go/micro/throttle"
	"github.com/bang-go/network/grpcx"
	"github.com/bang-go/network/grpcx/server_interceptor"
	"github.com/bang-go/util"
	"google.golang.org/grpc"
)

import pb "github.com/bang-go/crab/examples/proto/echo"

type greeterWrapper struct {
	pb.UnimplementedEchoServer
}

func (g *greeterWrapper) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: "hello " + req.Message}, nil
}

func (g *greeterWrapper) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	var num = 100
	for i := 0; i < num; i++ {
		_ = stream.Send(&pb.EchoResponse{Message: util.IntToString(i)})
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestServer(t *testing.T) {
	greeter := &greeterWrapper{}
	rpcServer := grpcx.NewServer(&grpcx.ServerConfig{Addr: "127.0.0.1:8081", Trace: true})
	err := rpcServer.Start(func(server *grpc.Server) {
		pb.RegisterEchoServer(server, greeter)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestServerWithThrottleControl(t *testing.T) {
	var err error
	greeter := &greeterWrapper{}
	rpcServer := grpcx.NewServer(&grpcx.ServerConfig{Addr: "127.0.0.1:8081"})

	limiter := throttle.Limiter()
	_ = limiter.Build(throttle.WithLogLevel(logging.WarnLevel))
	err = limiter.Rule([]*flow.Rule{{
		Resource:               "some-test",
		TokenCalculateStrategy: flow.Direct,
		ControlBehavior:        flow.Reject,
		StatIntervalInMs:       1000,
		Threshold:              10,
	}})
	if err != nil {
		log.Fatal(err)
	}
	rpcServer.AddUnaryInterceptor(server_interceptor.UnaryServerThrottleInterceptor(func() bool {
		return limiter.Guard("some-test", func() error {
			log.Println("pass")
			return nil
		}, func() {
			log.Println("reject")
		})
	}))
	err = rpcServer.Start(func(server *grpc.Server) {
		pb.RegisterEchoServer(server, greeter)
	})
	if err != nil {
		log.Fatal(err)
	}
}
