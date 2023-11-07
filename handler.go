package crab

import (
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/tracex/aliyun_trace"
	"github.com/bang-go/crab/core/base/viperx"
	"github.com/bang-go/opt"
	"log"
)

func UseAppEnv(opts ...opt.Option[env.Options]) Handler {
	return Handler{
		Pre: func() error {
			return env.Build(opts...)
		},
	}
}

func UseViper(conf *viperx.Config) Handler {
	return Handler{
		Pre: func() error {
			return viperx.Build(conf)
		},
	}
}

func UseTraceByAliSLS(conf *aliyun_trace.Config) Handler {
	config, err := aliyun_trace.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	return Handler{
		Init: func() error {
			return config.Start()
		},
		Close: func() error {
			config.Stop()
			return nil
		},
	}
}

func UseAppLog(opts ...opt.Option[logx.Options]) Handler {
	return Handler{
		Init: func() error {
			logx.Build(opts...)
			return nil
		}, Close: func() error {
			return nil
		},
	}
}
