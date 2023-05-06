package ginx

import (
	"github.com/bang-go/util"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Start() error
	Use(...gin.HandlerFunc)
	Engine() *gin.Engine
	Group(relativePath string, handlers ...gin.HandlerFunc) RouterGroup
}

type ServerConfig struct {
	Addr string
	Mode string
}
type ServerWrapper struct {
	*ServerConfig
	engine *gin.Engine
}

func New(conf *ServerConfig) Server {
	mode := util.If(conf.Mode != "", conf.Mode, gin.ReleaseMode)
	gin.SetMode(mode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &ServerWrapper{
		ServerConfig: conf,
		engine:       engine,
	}
}

func (s *ServerWrapper) Engine() *gin.Engine {
	return s.engine
}
func (s *ServerWrapper) Use(middlewares ...gin.HandlerFunc) {
	s.engine.Use(middlewares...)
}

func (s *ServerWrapper) Start() error {
	return s.engine.Run(s.Addr)
}

func (s *ServerWrapper) Group(relativePath string, handlers ...gin.HandlerFunc) RouterGroup {
	return &RouterGroupWrapper{
		Group: s.engine.Group(relativePath, handlers...),
	}
}
