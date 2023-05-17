package aliyun_trace_test

import (
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/crab/core/base/tracex/aliyun_trace"
	"go.opentelemetry.io/otel/attribute"
	"log"
	"testing"
)

func TestTrace(t *testing.T) {
	var err error
	provider, err := aliyun_trace.New(&aliyun_trace.Config{
		ServiceName:           "testCrab",
		ServiceNamespace:      "ns",
		ServiceVersion:        "v1.0",
		TraceExporterEndpoint: aliyun_trace.TraceExporterEndpointStdout,
		//TraceExporterEndpoint: "",
		SlsConfig: aliyun_trace.SlsConfig{
			Project:         "",
			InstanceID:      "",
			AccessKeyID:     "",
			AccessKeySecret: "",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = provider.Start()
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Stop()
	tracer := tracex.Tracer("crab.com/basic")
	ctx, span := tracer.NewSpan("span11")
	span.SetStatus(404, "undefined")
	defer span.End()
	ctx2, span2 := tracer.ChildSpan(ctx, "span22")
	span2.SetAttributes(attribute.String("time", "111"))
	defer span2.End()
	_, span3 := tracer.ChildSpan(ctx2, "span33")
	defer span3.End()
}
