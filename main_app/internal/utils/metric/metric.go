package metric

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func WriteRequestDuration(start time.Time, met *prometheus.HistogramVec, method string) {
	elapsed := time.Since(start).Seconds()
	met.WithLabelValues(method).Observe(elapsed)
}
