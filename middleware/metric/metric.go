package metric

import (
	"simple-rest/pkg/prometheus"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Prometheus(s *prometheus.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		http := prometheus.NewHTTP(c.Request.URL.Path, c.Request.Method)
		http.Started()
		c.Next()
		http.Finished()
		http.StatusCode = strconv.Itoa(c.Writer.Status())
		s.SaveHTTP(http)
	}
}

func Metrics(c *gin.Context) {
	h := promhttp.Handler()
	h.ServeHTTP(c.Writer, c.Request)
}
