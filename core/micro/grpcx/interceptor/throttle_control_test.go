package interceptor_test

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/bang-go/crab/core/micro/grpcx/interceptor"
	"github.com/bang-go/crab/core/micro/throttle_control"
	"google.golang.org/grpc"
	"testing"
)

func TestThrottleBreaker(t *testing.T) {
	var err error
	breaker := throttle_control.Breaker()
	_ = breaker.Build(throttle_control.WithLogLevel(logging.InfoLevel), throttle_control.WithBreakerListener(throttle_control.DefaultStateChangeListener()))
	err = breaker.Rule([]*circuitbreaker.Rule{{
		Resource:                     "some-test",
		Strategy:                     circuitbreaker.ErrorCount,
		RetryTimeoutMs:               3000,
		MinRequestAmount:             10,
		StatIntervalMs:               5000,
		StatSlidingWindowBucketCount: 10,
		Threshold:                    50,
	}})
	checkError(err)
	handler := func() bool {
		return breaker.Guard("some-test", func() error {
			return nil
		}, func() {

		})
	}
	grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryServerThrottleInterceptor(handler)))
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
