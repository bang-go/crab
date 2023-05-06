package grpcx_test

import (
	"context"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/bang-go/crab/core/micro/grpcx"
	"github.com/bang-go/crab/core/micro/grpcx/interceptor"
	"github.com/bang-go/crab/core/micro/throttle_control"
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
func TestUnaryServer(t *testing.T) {
	greeter := &greeterWrapper{}
	rpcServer := grpcx.NewServer(&grpcx.ServerConfig{Addr: "127.0.0.1:8081"})
	err := rpcServer.Start(func(server *grpc.Server) {
		pb.RegisterGreeterServer(server, greeter)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestServerWithThrottleControl(t *testing.T) {
	var err error
	greeter := &greeterWrapper{}
	rpcServer := grpcx.NewServer(&grpcx.ServerConfig{Addr: "127.0.0.1:8081"})

	limiter := throttle_control.Limiter()
	_ = limiter.Build(throttle_control.WithLogLevel(logging.WarnLevel))
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
	rpcServer.AddUnaryInterceptor(interceptor.UnaryServerThrottleInterceptor(func() bool {
		return limiter.Guard("some-test", func() error {
			log.Println("pass")
			return nil
		}, func() {
			log.Println("reject")
		})
	}))
	err = rpcServer.Start(func(server *grpc.Server) {
		pb.RegisterGreeterServer(server, greeter)
	})
	if err != nil {
		log.Fatal(err)
	}
}
