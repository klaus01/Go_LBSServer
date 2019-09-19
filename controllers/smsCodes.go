package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/klaus01/Go_LBSServer/database"
	"github.com/klaus01/Go_LBSServer/models"
	"github.com/klaus01/Go_LBSServer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

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
	sig := fmt.Sprintf("SMS%sCODE%dS", body.PhoneNumber, body.Time)
	sig = utils.Md5(sig)
	if sig != body.Sign {
		utils.ResultClientError(c, "缺少参数3")
		return
	}

	collection := database.GetDB().Collection(models.TableNameSmsCode)
	filter := bson.M{"phoneNumber": body.PhoneNumber}
	update := bson.M{"$set": bson.M{"createAt": time.Now()}}
	var smsCode models.SmsCode
	err := collection.FindOneAndUpdate(database.Context(), filter, update).Decode(&smsCode)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Println("[ERROR]", "查询", body.PhoneNumber, "短信失败", err)
			utils.ResultServerError(c, err.Error())
			return
		}

		smsCode.PhoneNumber = body.PhoneNumber
		smsCode.Code = fmt.Sprintf("%d", 1000+utils.RandomInt(8999))
		smsCode.CreateAt = time.Now()
		_, err := collection.InsertOne(database.Context(), smsCode)
		if err != nil {
			log.Println("[ERROR]", "插入短信失败", err)
			utils.ResultServerError(c, err.Error())
			return
		}
	}

	err = models.SendVerificationCode(smsCode.PhoneNumber, smsCode.Code)
	if err != nil {
		utils.ResultServerError(c, err.Error())
		return
	}

	utils.ResultSuccessMsg(c, "获取验证码成功")
}
