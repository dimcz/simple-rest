package middleware

import (
	"github.com/gin-gonic/gin"
	"simple-rest/metric"
)

func Metric(mService metric.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		m := metric.NewHTTP(c.Request.URL.Path, c.Request.Method)
		m.Started()
		c.Next()
		m.Finished()
		c.Next()
		m.StatusCode = c.Writer.Status()
		mService.SaveHTTP(m)
	}
}
