package tracex_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/crab/core/micro/httpx"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"net/http"
	"testing"
	"time"
)

// TestMicroPropagate task=>gin server=>grpc server1 =>grpc server2
func TestMicroPropagate(t *testing.T) {
	InitFrame()
	defer crab.Close()
	//task
	tracer := tracex.Tracer("crab-tracex-propagate")
	ctx, span := tracer.NewSpan("step1-1")
	defer span.End()
	span.SetAttributes(attribute.String("num", "1"))
	_, span2 := tracer.ChildSpan(ctx, "step1-2")
	span2.SetAttributes(attribute.String("num", "2"))
	span2.End()
	httpClient := httpx.New(httpx.Config{
		Trace:   true,
		Timeout: 20 * time.Second,
	})
	res, err := httpClient.Send(ctx, &httpx.Request{Url: "http://127.0.0.1:8080/health", Method: http.MethodGet})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(res.Content))
}
