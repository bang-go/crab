package env

import (
	"github.com/bang-go/opt"
)

type Options struct {
	appKey string
	appEnv string
}

func WithAppKey(appKey string) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.appKey = appKey
	})
}
func WithAppEnv(appEnv string) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.appEnv = appEnv
	})
}
