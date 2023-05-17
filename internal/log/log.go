package log

import (
	"github.com/bang-go/crab/core/base/logx"
	"go.uber.org/zap/zapcore"
)

var FrameLogger *logx.Logger

// 默认初始化
func init() {
	InitLog(logx.WarnLevel, logx.EncodeConsole)
}

func InitLog(AllowLogLevel logx.Level, LogEncoding string) {
	encodeConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var encode zapcore.Encoder
	switch LogEncoding {
	case logx.EncodeJson:
		encode = zapcore.NewJSONEncoder(encodeConfig)
	default:
		encode = zapcore.NewConsoleEncoder(encodeConfig)
	}
	FrameLogger = logx.New(logx.WithLevelEnabler(AllowLogLevel), logx.WithEncoder(encode), logx.WithCallerSkip(0))
}
