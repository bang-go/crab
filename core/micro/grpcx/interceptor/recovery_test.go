package interceptor_test

import (
	"context"
	"github.com/bang-go/crab/core/micro/grpcx/interceptor"
	"google.golang.org/grpc"
	"testing"
)

func TestRecovery(t *testing.T) {
	custom := func(ctx context.Context, p any) (err error) {
		return nil
	}
	grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryServerRecoveryInterceptor(custom)))
}
