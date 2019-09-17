package controllers

import (
	"github.com/klaus01/Go_LBSServer/utils"

	"github.com/gin-gonic/gin"
)

type smsCodesController struct{}

func (x smsCodesController) post(c *gin.Context) {
	type RequestBody struct {
		PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
		Time        int    `json:"time" form:"time"`
		Sign        string `json:"sign" form:"sign"`
	}
	var body RequestBody
	c.Bind(&body)

	if len(body.PhoneNumber) <= 0 {
		utils.ResultClientError(c, "缺少手机号")
		return
	}
	if body.Time <= 0 {
		utils.ResultClientError(c, "缺少参数1")
		return
	}
	if len(body.Sign) <= 0 {
		utils.ResultClientError(c, "缺少参数2")
		return
	}

	utils.ResultSuccessData(c, body)
}
