package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-rest/pkg/message"
	"simple-rest/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = message.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = message.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			c.Set("UserID", claims.ID)
			if err != nil {
				switch jwt.ValidationErrorExpired {
				default:
					code = message.ERROR
				}
			}
		}

		if code != message.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message.GetMessage(code),
				"data":    data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
