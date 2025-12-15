package logx

import (
	"log/slog"

	"github.com/bang-go/opt"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	source        bool
	logLevel      slog.Leveler
	logEncodeType uint               // log encode type: text ,json
	logOutType    uint               // log out type: stdout ,file
	logFileConfig *lumberjack.Logger //default stdout
	callerSkip    int                // caller skip 层数，用于正确显示调用位置
}

func WithOutStdout() opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logOutType = LogOutByStdout
	})
}

func WithOutFile(config *lumberjack.Logger) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logFileConfig = config
		o.logOutType = LogOutByFile
	})
}
func WithEncodeText() opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logEncodeType = LogEncodeText
	})
}
func WithEncodeJson() opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logEncodeType = LogEncodeJson
	})
}

func WithLevel(level Level) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logLevel = level
	})
}

func WithSource(source bool) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.source = source
	})
}

// WithCallerSkip 设置 caller skip 层数，用于调整日志显示的调用位置
// 默认值为 3，适用于直接使用 logx.Info() 等函数
// 如果你封装了自己的日志函数，需要增加这个值
func WithCallerSkip(skip int) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.callerSkip = skip
	})
}
