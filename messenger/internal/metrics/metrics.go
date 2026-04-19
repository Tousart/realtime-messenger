package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests partitioned by method and status code.",
	},
	[]string{"method", "code"},
)

func init() {
	prometheus.MustRegister(HTTPRequestsTotal)
}

func InstrumentHandler(h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(HTTPRequestsTotal, h)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
