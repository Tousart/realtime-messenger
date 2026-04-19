package httpmw

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tousart/messenger/internal/metrics"
)

func InstrumentHandler(next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(metrics.HTTPRequestsTotal, next)
}
