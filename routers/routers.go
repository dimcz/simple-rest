package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"simple-rest/logging"
	"simple-rest/metric"
	"simple-rest/middleware"
	"simple-rest/model"
)

const RecordsUrl = "/records"
const RecordUrl = "/records/:record"

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func Setup(logger *logging.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	metricServer, err := metric.NewPrometheusService()
	if err != nil {
		logger.Fatalf("routers.Setup error: %v", err)
	}
	r.Use(middleware.Metric(metricServer))

	r.GET("/metrics", prometheusHandler())

	r.GET(RecordsUrl, func(ctx *gin.Context) {
		records, err := model.SelectAll()
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, records)
	})

	r.GET(RecordUrl, func(ctx *gin.Context) {
		id := ctx.Param("record")
		record, err := model.SelectRecordByID(id)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, record)
	})

	r.POST(RecordsUrl, func(ctx *gin.Context) {
		var record model.Records
		if err := ctx.BindJSON(&record); err != nil {
			logger.Error(err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := model.CreateUser(record)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		} else {
			result := struct {
				ID int `json:"id"`
			}{
				ID: result,
			}
			ctx.JSON(http.StatusOK, &result)
		}
	})

	r.DELETE(RecordUrl, func(ctx *gin.Context) {
		id := ctx.Param("record")
		err := model.DeleteUser(id)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		} else {
			ctx.Writer.WriteHeader(http.StatusNoContent)
		}
	})

	r.PUT(RecordUrl, func(ctx *gin.Context) {
		id := ctx.Param("record")
		var record model.Records
		if err := ctx.BindJSON(&record); err != nil {
			logger.Error(err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := model.UpdateUser(id, record)
		if err != nil {
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
		} else {
			ctx.Writer.WriteHeader(http.StatusNoContent)
		}
	})

	return r
}
