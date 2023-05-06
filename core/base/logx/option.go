package logx

import (
	"github.com/bang-go/opt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type options struct {
	callerSkip    int
	logOutType    uint
	logStdout     bool
	logFileConfig *lumberjack.Logger
	levelEnabler  zapcore.LevelEnabler
	zapOption     []zap.Option
	zapEncoder    zapcore.Encoder
}

func WithEncoder(encoder zapcore.Encoder) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.zapEncoder = encoder
	})
}

func WithCallerSkip(skip int) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.callerSkip = skip
	})
}

func WithDefaultConfig(kind int) opt.Option[options] {
	var encoder zapcore.Encoder
	switch kind {
	case DefaultConfigKindProd:
		encoder = NewDefaultProdEncoder()
	default:
		encoder = NewDefaultDevEncoder()
	}
	return WithEncoder(encoder)
}

func WithZapOption(opts ...zap.Option) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.zapOption = opts
	})
}

func WithOutStdout() opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.logStdout = true
		o.logOutType |= 1
	})
}

func WithOutFile(config *lumberjack.Logger) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.logFileConfig = config
		o.logOutType |= 2
	})
}
func WithLevelEnabler(levelEnabler zapcore.LevelEnabler) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.levelEnabler = levelEnabler
	})
}
