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
		//TraceExporterEndpoint: "sf-test-ops.cn-hangzhou.log.aliyuncs.com:10010",
		SlsConfig: aliyun_trace.SlsConfig{
			Project:         "sf-test-ops",
			InstanceID:      "sf-test-ops",
			AccessKeyID:     "LTAI5tJCK1D4L5XuzVTNHTfw",
			AccessKeySecret: "SLWAmseH0bLnm0R8TzgFZEPGVyckYA",
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
	tracer := tracex.Tracer("ex.com/basic")
	ctx, span := tracer.NewSpan("span11")
	span.SetStatus(404, "undefined")
	defer span.End()
	ctx2, span2 := tracer.ChildSpan(ctx, "span22")
	span2.SetAttributes(attribute.String("time", "111"))
	defer span2.End()
	_, span3 := tracer.ChildSpan(ctx2, "span33")
	defer span3.End()
}