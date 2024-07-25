package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mygo_api_requests_total",
			Help: "How many HTTP requests processed.",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mygo_api_request_duration_milliseconds",
			Help:    "How long it took to process the request.",
			Buckets: []float64{100, 300, 1000, 5000},
		},
		[]string{"method", "path", "status"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(
		RequestCount,
		RequestDuration,
	)
}
