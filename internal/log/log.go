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
		frameLogger = logx.New(logx.WithLevel(logx.LevelWarn), logx.WithOutStdout())
		m.Unlock()
	}
	return frameLogger
}

func SetFrameLogger(logLevel logx.Level, logEncode uint) {
	opts := []opt.Option[logx.Options]{logx.WithLevel(logLevel), logx.WithSource(true)}
	switch logEncode {
	case logx.LogEncodeJson:
		opts = append(opts, logx.WithEncodeJson())
		break
	default:
		opts = append(opts, logx.WithEncodeText())
	}
	frameLogger = logx.New(opts...)
}

//func InitLog(AllowLogLevel logx2.Level, LogEncoding string) {
//	encodeConfig := zapcore.EncoderConfig{
//		TimeKey:        "time",
//		LevelKey:       "level",
//		NameKey:        "logger",
//		CallerKey:      "caller",
//		FunctionKey:    zapcore.OmitKey,
//		MessageKey:     "msg",
//		StacktraceKey:  "stacktrace",
//		LineEnding:     zapcore.DefaultLineEnding,
//		EncodeLevel:    zapcore.LowercaseLevelEncoder,
//		EncodeTime:     zapcore.ISO8601TimeEncoder,
//		EncodeDuration: zapcore.SecondsDurationEncoder,
//		EncodeCaller:   zapcore.ShortCallerEncoder,
//	}
//	var encode zapcore.Encoder
//	switch LogEncoding {
//	case logx2.EncodeJson:
//		encode = zapcore.NewJSONEncoder(encodeConfig)
//	default:
//		encode = zapcore.NewConsoleEncoder(encodeConfig)
//	}
//	frameLogger = logx2.New(logx2.WithLevelEnabler(AllowLogLevel), logx2.WithEncoder(encode), logx2.WithCallerSkip(0))
//}
