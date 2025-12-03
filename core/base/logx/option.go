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
