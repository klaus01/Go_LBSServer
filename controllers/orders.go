package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/klaus01/Go_LBSServer/utils"
)

type ordersController struct{}

func (x ordersController) getList(c *gin.Context) {
	utils.ResultSuccessData(c, nil)
}

func (x ordersController) createOrder(c *gin.Context) {
	utils.ResultSuccessData(c, nil)
}

func (x ordersController) postDeliveryInfo(c *gin.Context) {
	utils.ResultSuccessData(c, nil)
}
