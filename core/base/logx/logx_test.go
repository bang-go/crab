package logx_test

import (
	"testing"

	"github.com/bang-go/crab/core/base/logx"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger(t *testing.T) {
	logx.Debug("hello1")
	logx.Info("hello2", "ddd", "ddd")
	logx.SetLoggerLevel(logx.LevelDebug)
	logx.Debug("hello3")
	logx.Build(logx.WithEncodeJson())
	logx.Debug("hello4")
	logx.Build(logx.WithEncodeText())
	logx.Debug("hello5")

}
func TestLogFile(t *testing.T) {
	logx.Build(logx.WithOutStdout(), logx.WithOutFile(&lumberjack.Logger{
		Filename: "./logx.log",
	}))
	logx.Info("log file 1")
	logx.Info("log file 2", "name", "bale")
}
