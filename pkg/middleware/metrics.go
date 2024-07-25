package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"mygo/template/pkg/infra/metric"
)

// Metrics ...
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())
		metric.RequestCount.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.FullPath(),
			"status": status,
		}).Inc()

		// request duration, in ms
		metric.RequestDuration.With(prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.FullPath(),
			"status": status,
		}).Observe(float64(duration) / float64(time.Microsecond))
	}
}
