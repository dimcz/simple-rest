package routers

import (
	"github.com/gin-gonic/gin"
	"simple-rest/middleware/jwt"
	"simple-rest/middleware/metric"
	"simple-rest/pkg/prometheus"
	"simple-rest/pkg/util"
	v1 "simple-rest/routers/v1"
)

func InitRouters(logger *util.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	prom, err := prometheus.NewPrometheus()
	if err != nil {
		logger.Fatalf("can't start prometheus: %v", err)
	}
	r.Use(metric.Prometheus(prom))
	r.GET("/metrics", metric.Metrics)

	r.POST("/auth", GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/records", v1.GetRecords)
		apiv1.DELETE("/record/:record", v1.DeleteRecordByID)
	}

	return r
}
