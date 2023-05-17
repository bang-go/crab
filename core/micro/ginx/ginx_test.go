package ginx_test

import (
	"github.com/bang-go/crab/core/micro/ginx"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"testing"
)

func TestGin(t *testing.T) {
	ginServer := ginx.New(&ginx.ServerConfig{Addr: ":8080", Mode: gin.ReleaseMode})
	gp := ginServer.Group("/bot")
	gp.RouterHandle(http.MethodGet, "/health", func(c *gin.Context) {
		c.JSON(200, "success")
	})
	err := ginServer.Start()
	if err != nil {
		log.Fatal(err)
	}
}
