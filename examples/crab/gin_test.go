package crab_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/micro/ginx"
	"github.com/gin-gonic/gin"
)

func TestGin(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.LogEncodeText), crab.WithLogAllowLevel(logx.LevelInfo))
	baseSetting()
	//do something
	log.Println("do something")
	//cmd为可选模式
	cmdGin := cmd.New(&cmd.Config{CmdUse: "gin"})
	cmdGin.SetRun(func(args []string) {
		//cmd gin logic
		server := ginx.New(&ginx.ServerConfig{Addr: ":8080", Mode: gin.ReleaseMode})
		setRoute(server)
		_ = server.Start()
	})
	crab.RegisterCmd(cmdGin)
	if err = crab.Start(); err != nil {
		log.Println(err)
	}
	defer crab.Close()
}

func setRoute(server ginx.Server) {
	gp := server.Group("/")
	gp.Handle(http.MethodGet, "/health", func(c *gin.Context) {
		c.JSON(200, "success")
	})
}
