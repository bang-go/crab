package httptrace

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var Transport = otelhttp.NewTransport(http.DefaultTransport)
