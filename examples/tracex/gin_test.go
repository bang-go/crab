package tracex_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/tracex"
	pb "github.com/bang-go/crab/examples/proto/echo"
	"github.com/bang-go/network/ginx"
	"github.com/bang-go/network/grpcx"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func TestGinServer(t *testing.T) {
	defer crab.Close()
	var err error
	server := ginx.New(&ginx.ServerConfig{Addr: ":8080", Trace: true})
	route(server)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func route(server ginx.Server) {
}

func healthHandle(c *gin.Context) {
	log.Println("receive :", c.Request)
	c.String(200, "success")
	tracer := tracex.Tracer("crab-gin-server")
	ctx, span := tracer.NewSpanWithContext(c.Request.Context(), "step-gin-1")
	span.SetAttributes(attribute.String("health handler", "1"))
	defer span.End()
	time.Sleep(10 * time.Second)
	log.Println(c.Request.Header)
	// call grpc
	callGrpcLogic(ctx)
}

func callGrpcLogic(ctx context.Context) {
	rpcClient := grpcx.NewClient(&grpcx.ClientConfig{Addr: "127.0.0.1:8081", Trace: true})
	defer rpcClient.Close()
	conn, err := rpcClient.Dial()
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewEchoClient(conn)
	unaryReply, err := client.UnaryEcho(ctx, &pb.EchoRequest{Message: "jk"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(unaryReply.Message)
	tracer := tracex.Tracer("crab-grpc-client")
	ctx, span := tracer.NewSpanWithContext(ctx, "step-grpc-1")
	span.SetAttributes(attribute.String("k", "v"))
	defer span.End()
	time.Sleep(2 * time.Second)

}
