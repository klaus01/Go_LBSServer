package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResultSuccess(c *gin.Context, msg string, data interface{}) {
	var message *string
	if len(msg) > 0 {
		message = &msg
	}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": message, "data": data})
}

func ResultSuccessMsg(c *gin.Context, msg string) {
	ResultSuccess(c, msg, nil)
}

func ResultSuccessData(c *gin.Context, data interface{}) {
	ResultSuccess(c, "", data)
}

func ResultClientError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": msg})
	c.Abort()
}

func ResultServerError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": msg})
	c.Abort()
}
