package v1

import (
	"net/http"
	"simple-rest/pkg/app"
	"simple-rest/pkg/message"
	"simple-rest/service/records_service"

	"github.com/gin-gonic/gin"
)

func GetRecords(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, ok := c.MustGet("UserID").(int)
	if !ok {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	rs := records_service.ActionRequest{
		UserID: userID,
	}
	records, err := rs.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, message.SUCCESS, records)
}

func DeleteRecordByID(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, ok := c.MustGet("UserID").(int)
	if !ok {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	recordID := c.Param("record")

	rs := records_service.ActionRequest{
		UserID:   userID,
		RecordID: recordID,
	}
	err := rs.DeleteByRecordID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, message.SUCCESS, nil)
}
