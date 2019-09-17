package controllers

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化
func Init(e *gin.Engine) {
	var r *gin.RouterGroup

	smsCodes := new(smsCodesController)
	r = e.Group("/smscodes")
	r.POST("/", smsCodes.post)

	orders := new(ordersController)
	r = e.Group("/orders")
	r.POST("/", orders.post)
}
