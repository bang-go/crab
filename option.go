package crab

import (
	"github.com/bang-go/crab/core/base/logx"
	"github.com/bang-go/crab/internal/vars"
	"github.com/bang-go/opt"
)

type logOptions struct {
	allowLogLevel logx.Level //允许的log level -1:Debug info:0 1:warn 2:error 3:dpanic 4 panic 5 fatal
	logEncodeType uint       //日志编码方式
}
type options struct {
	logOptions
	appName string
}

func WithLogAllowLevel(logLevel logx.Level) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.allowLogLevel = logLevel
	})
}

func WithLogEncoding(logEncodeType uint) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.logEncodeType = logEncodeType
	})
}

func WithAppName(appName string) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.appName = appName
		vars.DefaultAppName.Store(appName)
	})
}
