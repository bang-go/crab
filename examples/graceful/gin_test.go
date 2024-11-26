package graceful_test

import (
	"github.com/bang-go/crab"
	"github.com/bang-go/crab/cmd"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/network/ginx"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
)

func TestGin(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.LogEncodeJson), crab.WithLogAllowLevel(logx.LevelInfo))
	//do something
	log.Println("do something")
	if err = crab.Start(); err != nil {
		log.Fatal(err)
	}
	server := ginx.New(&ginx.ServerConfig{Addr: ":8080", Mode: gin.ReleaseMode})
	setRoute(server)
	if err = server.Start(); err != nil {
		log.Println(err)
	}
	select {}
}

func TestGinWithCmd(t *testing.T) {
	var err error
	crab.Build(crab.WithLogEncoding(logx.LogEncodeJson), crab.WithLogAllowLevel(logx.LevelInfo))
	//do something
	log.Println("do something")
	cmdGin := cmd.New(&cmd.Config{CmdUse: "gin"})
	cmdGin.SetRun(func(args []string) {
		server := ginx.New(&ginx.ServerConfig{Addr: ":8080", Mode: gin.ReleaseMode})
		setRoute(server)
		_ = server.Start()
		//cmd gin logic
	})
	crab.RegisterCmd(cmdGin)
	if err = crab.Start(); err != nil {
		log.Println(err)
	}
}
func setRoute(server ginx.Server) {
	gp := server.Group("/")
	gp.RouterHandle(http.MethodGet, "/health", func(c *gin.Context) {
		crab.Done()
		c.JSON(200, "success")
	})
}
