package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/klaus01/Go_LBSServer/config"
	"github.com/klaus01/Go_LBSServer/utils"
)

// TableNameSmsCode SmsCode 表名
const TableNameSmsCode string = "smscodes"

// SmsCode 短信模型
type SmsCode struct {
	PhoneNumber string    `bson:"phoneNumber"`
	Code        string    `bson:"code"`
	CreateAt    time.Time `bson:"createAt"`
}

func postSendSMS(phoneNumber string, templateID string, parameters []string) error {
	c := config.GetConfig()

	timeNow := time.Now().Format("20060102030405")
	sig := fmt.Sprintf("%s%s%s", c.GetString("yuntongxun.accountSID"), c.GetString("yuntongxun.authToken"), timeNow)
	sig = strings.ToUpper(utils.Md5(sig))
	authorization := fmt.Sprintf("%s:%s", c.GetString("yuntongxun.accountSID"), timeNow)
	authorization = base64.URLEncoding.EncodeToString([]byte(authorization))
	postJSON := map[string]interface{}{"to": phoneNumber, "appId": c.GetString("yuntongxun.appID"), "templateId": templateID, "datas": parameters}
	postData, err := json.Marshal(postJSON)
	if err != nil {
		log.Println("[ERROR]", "JSON 序列化失败", err)
		return err
	}

	url := fmt.Sprintf("https://app.cloopen.com:8883/2013-12-26/Accounts/%s/SMS/TemplateSMS?sig=%s", c.GetString("yuntongxun.accountSID"), sig)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postData))
	if err != nil {
		log.Println("[ERROR]", "NewRequest 失败", err)
		return err
	}
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[ERROR]", "client.Do 失败", err)
		return err
	}
	defer resp.Body.Close()

	httpCode := strings.TrimSpace(resp.Status)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR]", "httpCode", httpCode, "read body error", err)
		return err
	}
	var bodyMap map[string]interface{}
	err = json.Unmarshal(b, &bodyMap)
	if err != nil {
		log.Println("[ERROR]", "json.Unmarshal", err, "json string:", string(b))
		return err
	}
	ytxCode := bodyMap["statusCode"].(string)
	if httpCode != "200" || ytxCode != "000000" {
		log.Println("[ERROR]", "send sms fail", string(b))
		message := "发送短信失败"
		if len(ytxCode) > 0 || bodyMap["statusMsg"] != nil {
			var ytxMsg string
			if bodyMap["statusMsg"] != nil {
				ytxMsg = bodyMap["statusMsg"].(string)
			}
			message = fmt.Sprintf("%s %s%s", message, ytxCode, ytxMsg)
		} else if len(httpCode) > 0 {
			message = fmt.Sprintf("%s code:%s", message, httpCode)
		}
		return errors.New(message)
	}
	log.Println("[INFO]", "send sms success", string(b))
	return nil
}

// SendVerificationCode 发送验证码短信
func SendVerificationCode(phoneNumber string, code string) error {
	c := config.GetConfig()
	expirationTime := c.GetInt("smsCodeExpireAfterSeconds") / 60
	return postSendSMS(phoneNumber, c.GetString("yuntongxun.templateIDs.verificationCode"), []string{code, strconv.Itoa(expirationTime)})
}

// SendOrderConfirmation 发送订单确认短信
func SendOrderConfirmation(phoneNumber string, name string, orderID string) error {
	c := config.GetConfig()
	return postSendSMS(phoneNumber, c.GetString("yuntongxun.templateIDs.orderConfirmation"), []string{name, orderID})
}

// SendShippingNotice 发送订单发货短信
func SendShippingNotice(phoneNumber string, courierCompany string, waybillNumber string) error {
	c := config.GetConfig()
	return postSendSMS(phoneNumber, c.GetString("yuntongxun.templateIDs.shippingNotice"), []string{courierCompany, waybillNumber})
}
