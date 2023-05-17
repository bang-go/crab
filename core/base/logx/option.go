package logx

import (
	"github.com/bang-go/opt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Options struct {
	callerSkip    int
	logOutType    uint
	logStdout     bool
	logFileConfig *lumberjack.Logger
	levelEnabler  zapcore.LevelEnabler
	zapOption     []zap.Option
	zapEncoder    zapcore.Encoder
}

func WithEncoder(encoder zapcore.Encoder) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.zapEncoder = encoder
	})
}

func WithCallerSkip(skip int) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.callerSkip = skip
	})
}

func WithDefaultConfig(kind int) opt.Option[Options] {
	var encoder zapcore.Encoder
	switch kind {
	case DefaultConfigKindProd:
		encoder = NewDefaultProdEncoder()
	default:
		encoder = NewDefaultDevEncoder()
	}
	return WithEncoder(encoder)
}

func WithZapOption(opts ...zap.Option) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.zapOption = opts
	})
}

func WithOutStdout() opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logStdout = true
		o.logOutType |= 1
	})
}

func WithOutFile(config *lumberjack.Logger) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.logFileConfig = config
		o.logOutType |= 2
	})
}
func WithLevelEnabler(levelEnabler zapcore.LevelEnabler) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.levelEnabler = levelEnabler
	})
}
