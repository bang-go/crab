package crab

import (
	"context"
	"net/http"
	"time"

	"github.com/bang-go/crab/core/base/env"
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/core/base/viperx"
	"github.com/bang-go/crab/core/db/mysqlx"
	"github.com/bang-go/crab/core/db/redisx"
	"github.com/bang-go/opt"
	"github.com/gin-gonic/gin"
)

// UseLogx 辅助函数：初始化日志
func UseLogx(opts ...opt.Option[logx.Options]) error {
	return Use(Handler{
		Pre: func() error {
			logx.Build(opts...)
			return nil
		},
	})
}

// UseEnv 辅助函数：初始化环境变量
func UseEnv(opts ...opt.Option[env.Options]) error {
	return Use(Handler{
		Pre: func() error {
			return env.Build(opts...)
		},
	})
}

// UseViper 辅助函数：初始化配置
func UseViper(conf *viperx.Config) error {
	return Use(Handler{
		Pre: func() error {
			return viperx.Build(conf)
		},
	})
}

// UseRedis 辅助函数：初始化 Redis 客户端
func UseRedis(opts *redisx.Options, clients ...*redisx.Client) error {
	if len(clients) == 0 {
		return nil
	}

	return Use(Handler{
		Init: func() error {
			client, err := redisx.New(opts)
			if err != nil {
				return err
			}
			*clients[0] = *client
			return nil
		},
		Close: func() error {
			if len(clients) > 0 && clients[0] != nil {
				return clients[0].Close()
			}
			return nil
		},
	})
}

// UseMySQL 辅助函数：初始化 MySQL 客户端
func UseMySQL(conf *mysqlx.ClientConfig, clients ...*mysqlx.Client) error {
	if len(clients) == 0 {
		return nil
	}

	return Use(Handler{
		Init: func() error {
			client, err := mysqlx.New(conf)
			if err != nil {
				return err
			}
			*clients[0] = *client
			return nil
		},
		Close: func() error {
			if len(clients) > 0 && clients[0] != nil {
				return clients[0].Close()
			}
			return nil
		},
	})
}

// UseGin 辅助函数：初始化 Gin HTTP 服务器
func UseGin(addr string, setupRouter func(*gin.Engine)) error {
	if setupRouter == nil {
		return nil
	}

	var server *http.Server

	return Use(Handler{
		Init: func() error {
			engine := gin.Default()
			setupRouter(engine)

			server = &http.Server{
				Addr:    addr,
				Handler: engine,
			}

			return server.ListenAndServe()
		},
		Close: func() error {
			if server != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return server.Shutdown(ctx)
			}
			return nil
		},
	})
}

// UseHTTP 辅助函数：初始化 HTTP 服务器
func UseHTTP(server *http.Server) error {
	if server == nil {
		return nil
	}

	return Use(Handler{
		Init: func() error {
			return server.ListenAndServe()
		},
		Close: func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		},
	})
}
