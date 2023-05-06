package env

import (
	"github.com/bang-go/opt"
)

type options struct {
	appKey string
	appEnv string
}

func WithAppKey(appKey string) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.appKey = appKey
	})
}
func WithAppEnv(appEnv string) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.appEnv = appEnv
	})
}
