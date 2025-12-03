package logz_test

import (
	"log"
	"testing"

	"github.com/bang-go/crab/core/base/logz"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLogger(t *testing.T) {
	logz.Debug("hello1")
	logz.Debug("hello2")
	logz.Info("hello3", logz.Int("num", 1))
	logger := logz.New(logz.WithDefaultConfig(logz.DefaultConfigKindProd))
	logger.Info("hello4")
	logz.Info("hello5")
	logz.Build(logz.WithDefaultConfig(logz.DefaultConfigKindProd))
	logz.Info("hello6")
	logger.Info("hello7")
	if err := logger.Sync(); err != nil {
		log.Fatal(err)
	}
}

func TestLogFile(t *testing.T) {
	logz.Build(logz.WithOutStdout(), logz.WithOutFile(&lumberjack.Logger{
		Filename: "./logx.log",
	}))
	logz.Info("log file 1")
	logz.Info("log file 2", logz.String("name", "bale"))
}
