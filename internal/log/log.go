package log

import (
	"sync"

	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/opt"
)

var (
	frameLogger *logx.Logger
	m           sync.Mutex
)

func DefaultFrameLogger() *logx.Logger {
	if frameLogger == nil {
		m.Lock()
		if frameLogger == nil {
			frameLogger = logx.New(logx.WithLevel(logx.LevelWarn), logx.WithOutStdout())
		}
		m.Unlock()
	}
	return frameLogger
}

func SetFrameLogger(logLevel logx.Level, logEncode uint) {
	opts := []opt.Option[logx.Options]{logx.WithLevel(logLevel), logx.WithSource(true)}
	if logEncode == logx.LogEncodeJson {
		opts = append(opts, logx.WithEncodeJson())
	} else {
		opts = append(opts, logx.WithEncodeText())
	}
	frameLogger = logx.New(opts...)
}
