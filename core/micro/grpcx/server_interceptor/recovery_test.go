package server_interceptor_test

import (
	"context"
	"github.com/bang-go/crab/core/micro/grpcx/server_interceptor"
	"google.golang.org/grpc"
	"testing"
)

func TestRecovery(t *testing.T) {
	custom := func(ctx context.Context, p any) {}
	grpc.NewServer(grpc.ChainUnaryInterceptor(server_interceptor.UnaryServerRecoveryInterceptor(custom)))
}
