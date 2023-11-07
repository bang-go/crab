package logx2_test

import (
	"github.com/bang-go/crab/core/base/logx2"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	logx2.Debug("hello1")
	logx2.Debug("hello2")
	logx2.Info("hello3", logx2.Int("num", 1))
	logger := logx2.New(logx2.WithDefaultConfig(logx2.DefaultConfigKindProd))
	logger.Info("hello4")
	logx2.Info("hello5")
	logx2.Build(logx2.WithDefaultConfig(logx2.DefaultConfigKindProd))
	logx2.Info("hello6")
	logger.Info("hello7")
	if err := logger.Sync(); err != nil {
		log.Fatal(err)
	}
}

func TestLogFile(t *testing.T) {
	logx2.Build(logx2.WithOutStdout(), logx2.WithOutFile(&lumberjack.Logger{
		Filename: "./logx.log",
	}))
	logx2.Info("log file 1")
	logx2.Info("log file 2", logx2.String("name", "bale"))
}
