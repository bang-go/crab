package httptrace

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

var Transport = otelhttp.NewTransport(http.DefaultTransport)
