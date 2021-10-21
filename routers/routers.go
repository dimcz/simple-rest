package routers

import (
	"github.com/gin-gonic/gin"
	"simple-rest/model"
)

const RecordsUrl = "/records"
const RecordUrl = "/records/:record"

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET(RecordsUrl, func(ctx *gin.Context) {
		model.SelectAll(ctx)
	})

	r.GET(RecordUrl, func(ctx *gin.Context) {
		model.SelectRecordByID(ctx)
	})

	r.POST(RecordsUrl, func(ctx *gin.Context) {
		model.CreateUser(ctx)
	})

	r.DELETE(RecordUrl, func(ctx *gin.Context) {
		model.DeleteUser(ctx)
	})

	r.PUT(RecordUrl, func(ctx *gin.Context) {
		model.UpdateUser(ctx)
	})

	return r
}
