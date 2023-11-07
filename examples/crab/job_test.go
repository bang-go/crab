package crab_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"log"
	"testing"
)

func TestJob(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.LogEncodeJson), crab.WithLogAllowLevel(logx.LevelInfo))
	baseSetting()
	crab.Use(crab.UseAppLog(logx.WithEncodeJson()))
	//do something
	log.Println("do something")
	//cmd为可选模式
	cmdJob := cmd.New(&cmd.Config{CmdUse: "job"})
	cmdJob.SetRun(func(args []string) {
		//cmd job logic
	})
	if err = crab.Start(); err != nil {
		log.Fatal(err)
	}
	defer crab.Close()
}
