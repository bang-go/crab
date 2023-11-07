package tracex_test

import (
	"flag"
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/tracex/aliyun_trace"
	"log"
)

var (
	endpoint        = flag.String("endpoint", "", "")
	project         = flag.String("project", "", "")
	instanceID      = flag.String("instance_id", "", "")
	accessKeyID     = flag.String("access_key_id", "", "")
	accessKeySecret = flag.String("access_key_secret", "", "")
)

func InitFrame() {
	flag.Parse()
	crab.Build(crab.WithLogAllowLevel(logx2.WarnLevel))
	err := crab.Use(crab.UseAppEnv())
	if err != nil {
		log.Fatal(err)
	}
	crab.Use(crab.UseAppLog(),
		crab.UseTraceByAliSLS(&aliyun_trace.Config{
			ServiceName:           "crab-test-service",
			ServiceNamespace:      "ns",
			ServiceVersion:        "v1.0",
			TraceExporterEndpoint: *endpoint,
			SlsConfig: aliyun_trace.SlsConfig{
				Project:         *project,
				InstanceID:      *instanceID,
				AccessKeyID:     *accessKeyID,
				AccessKeySecret: *accessKeySecret,
			},
		}),
	)

	if err := crab.Start(); err != nil {
		log.Fatal(err)
	}

}
