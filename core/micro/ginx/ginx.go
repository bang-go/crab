package ginx

import (
	"context"
	"github.com/bang-go/crab/core/base/tracex/instrument/gintrace"
	"github.com/bang-go/crab/core/pub/graceful"
	"github.com/bang-go/crab/internal/vars"
	"github.com/bang-go/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server interface {
	Start() error
	Use(...gin.HandlerFunc)
	Engine() *http.Server
	GinEngine() *gin.Engine
	Group(relativePath string, handlers ...gin.HandlerFunc) RouterGroup
	Shutdown() error
}

type ServerConfig struct {
	Addr        string
	Mode        string
	Trace       bool
	TraceFilter gintrace.Filter
}
type ServerEntity struct {
	*ServerConfig
	ginEngine  *gin.Engine
	httpServer *http.Server
}

func New(conf *ServerConfig) Server {
	mode := util.If(conf.Mode != "", conf.Mode, gin.ReleaseMode)
	gin.SetMode(mode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())
	if conf.Trace {
		ginEngine.Use(gintrace.Middleware(vars.DefaultAppName.Load(), gintrace.WithFilter(conf.TraceFilter)))
	}
	return &ServerEntity{
		ServerConfig: conf,
		ginEngine:    ginEngine,
	}
}

func (s *ServerEntity) GinEngine() *gin.Engine {
	return s.ginEngine
}

func (s *ServerEntity) Engine() *http.Server {
	return s.httpServer
}

func (s *ServerEntity) Use(middlewares ...gin.HandlerFunc) {
	s.ginEngine.Use(middlewares...)
}

func (s *ServerEntity) Start() (err error) {
	s.httpServer = &http.Server{
		Addr:    s.Addr,
		Handler: s.ginEngine,
	}
	graceful.Register(s.Shutdown)
	err = s.httpServer.ListenAndServe()
	return
}

func (s *ServerEntity) Group(relativePath string, handlers ...gin.HandlerFunc) RouterGroup {
	return &routerGroupEntity{
		group: s.ginEngine.Group(relativePath, handlers...),
	}
}

func (s *ServerEntity) Shutdown() error {
	//cxt, cancel := context.WithTimeout(context.Background(), graceful.MaxWaitTime)
	//defer cancel()
	return s.httpServer.Shutdown(context.Background())
}
