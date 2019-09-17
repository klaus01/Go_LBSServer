package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ordersController struct{}

func (x ordersController) post(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "朝秦暮楚",
	})
}
