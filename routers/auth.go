package routers

import (
	"net/http"
	"simple-rest/pkg/app"
	"simple-rest/pkg/message"
	"simple-rest/pkg/util"
	"simple-rest/service/auth_service"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	username := c.PostForm("username")
	password := c.PostForm("password")

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		appG.Response(http.StatusBadRequest, message.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	userID, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	if userID == 0 {
		appG.Response(http.StatusUnauthorized, message.ERROR, nil)
		return
	}

	token, err := util.GenerateToken(userID)
	if err != nil {
		appG.Response(http.StatusInternalServerError, message.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, message.SUCCESS, map[string]string{
		"token": token,
	})
}
