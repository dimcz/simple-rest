package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type Service struct {
	httpRequestHistogram *prometheus.HistogramVec
}

func NewPrometheusService() (*Service, error) {
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "requests_duration_seconds",
		Help:      "The latency of HTTP requests",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &Service{
		httpRequestHistogram: http,
	}
	err := prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

func (s *Service) SaveHTTP(h *HTTP) {
	s.httpRequestHistogram.WithLabelValues(h.Handler, h.Method, strconv.Itoa(h.StatusCode)).Observe(h.Duration)
}
