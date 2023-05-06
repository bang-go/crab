package logx_test

import (
	"github.com/bang-go/crab/core/base/logx"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	logx.Debug("hello1")
	logx.Debug("hello2")
	logx.Info("hello3", logx.Int("num", 1))
	logger := logx.New(logx.WithDefaultConfig(logx.DefaultConfigKindProd))
	logger.Info("hello4")
	logx.Info("hello5")
	logx.Build(logx.WithDefaultConfig(logx.DefaultConfigKindProd))
	logx.Info("hello6")
	logger.Info("hello7")
	if err := logger.Sync(); err != nil {
		log.Fatal(err)
	}
}

func TestLogFile(t *testing.T) {
	logx.Build(logx.WithOutStdout(), logx.WithOutFile(&lumberjack.Logger{
		Filename: "./logx.log",
	}))
	logx.Info("log file 1")
	logx.Info("log file 2", logx.String("name", "bale"))
}
