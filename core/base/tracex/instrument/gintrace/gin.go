package gintrace

import (
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/bang-go/opt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
)

type Filter = otelgin.Filter

type options struct {
	filters []Filter
}

func Middleware(service string, opts ...opt.Option[options]) gin.HandlerFunc {
	o := &options{
		filters: []Filter{DefaultHealthCheckFilter},
	}
	opt.Each(o, opts...)
	return otelgin.Middleware(service, otelgin.WithPropagators(tracex.Propagator()), otelgin.WithTracerProvider(tracex.Provider()), otelgin.WithFilter(o.filters...))
}

func WithFilter(f ...Filter) opt.Option[options] {
	return opt.OptionFunc[options](func(o *options) {
		o.filters = append(o.filters, f...)
	})
}

func DefaultHealthCheckFilter(req *http.Request) bool {
	if req.URL.Path == "/probe/health" {
		return false
	}
	return true
}
