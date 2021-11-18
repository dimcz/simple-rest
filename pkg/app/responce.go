package app

import (
	"simple-rest/pkg/message"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Responce struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Responce{
		Code:    errCode,
		Message: message.GetMessage(errCode),
		Data:    data,
	})
}
