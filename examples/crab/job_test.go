package crab_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/viperx"
	"log"
	"testing"
)

func TestJob(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.EncodeConsole), crab.WithLogAllowLevel(logx.InfoLevel))
	baseSetting()
	crab.Use([]crab.Handler{
		{
			Name: "log",
			Init: func() error {
				logx.Build(logx.WithDefaultConfig(logx.DefaultConfigKindDev))
				return nil
			}, Close: func() error {
				return logx.Sync()
			}}})
	//do something
	log.Println("do something")
	//cmd为可选模式
	cmdJob := cmd.New(&cmd.Config{CmdUse: "job"})
	cmdJob.SetRun(func(args []string) error {
		//cmd job logic
		return nil
	})
	if err = crab.Start(); err != nil {
		log.Fatal(err)
	}
	defer crab.Exit()
}

func baseSetting() {
	if err := crab.Exec([]crab.Handler{
		{
			Name: "env",
			Init: func() error {
				return env.Build()
			},
		}, {
			Name: "viper",
			Init: func() error {
				return viperx.Build(&viperx.Config{ConfigFormat: viperx.FileFormatYaml, ConfigPaths: []string{"./"}, ConfigNames: []string{"app"}})
			},
		},
	}); err != nil {
		log.Fatal(err)
	}
}
