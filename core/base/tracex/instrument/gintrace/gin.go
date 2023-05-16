package gintrace

import (
	"github.com/bang-go/crab/core/base/tracex"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Middleware(service string) gin.HandlerFunc {
	return otelgin.Middleware(service, otelgin.WithPropagators(tracex.Propagator()), otelgin.WithTracerProvider(tracex.Provider()))
}
