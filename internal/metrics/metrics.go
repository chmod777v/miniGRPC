package metrics

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(addres string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("Metrics listening", "Host", addres+"/metrics")
	return http.ListenAndServe(addres, mux)
}

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "server",
	Subsystem:  "http",
	Name:       "request",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"method", "status"})

func ObserveRequest(method string, d time.Duration, status int) {
	requestMetrics.WithLabelValues(method, strconv.Itoa(status)).Observe(d.Seconds())
}
