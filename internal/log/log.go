package log

import (
	"github.com/bang-go/crab/core/base/logx"
)

var FrameLogger *logx.Logger

// 默认初始化
func init() {
	FrameLogger = logx.New(logx.WithLevel(logx.LevelWarn), logx.WithOutStdout())
}

//func InitLog(logLevel logx.Level, logEncode uint){
//	FrameLogger=logx.New()
//}

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
//	FrameLogger = logx2.New(logx2.WithLevelEnabler(AllowLogLevel), logx2.WithEncoder(encode), logx2.WithCallerSkip(0))
//}
