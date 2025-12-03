package gintrace

import (
	"net/http"

	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/opt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Filter = otelgin.Filter

type Options struct {
	filters []Filter
}

func Middleware(service string, opts ...opt.Option[Options]) gin.HandlerFunc {
	o := &Options{
		filters: []Filter{DefaultHealthCheckFilter},
	}
	opt.Each(o, opts...)
	return otelgin.Middleware(service, otelgin.WithPropagators(tracex.Propagator()), otelgin.WithTracerProvider(tracex.Provider()), otelgin.WithFilter(o.filters...))
}

func WithFilter(f ...Filter) opt.Option[Options] {
	return opt.OptionFunc[Options](func(o *Options) {
		o.filters = append(o.filters, f...)
	})
}

func DefaultHealthCheckFilter(req *http.Request) bool {
	if req.URL.Path == "/probe/health" {
		return false
	}
	return true
}
