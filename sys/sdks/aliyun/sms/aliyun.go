package sms

import (
	"fmt"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func newSys(options Options) (sys *SMS, err error) {
	sys = &SMS{options.key, options.KeySecret, options.SignName}
	return
}

type SMS struct {
	key       string
	keysecret string
	signName  string
}

//发送短信验证吗
func (this *SMS) SendCaptcha(mobile string, captcha string) error {
	client, err := sdk.NewClientWithAccessKey("default", this.key, this.keysecret)
	if err != nil {
		return err
	}
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = this.signName
	request.QueryParams["TemplateCode"] = "SMS_137655450"
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + captcha + "\"}"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	var xx struct {
		Message   string `json:"Message"`
		RequestID string `json:"RequestId"`
		BizID     string `json:"BizId"`
		Code      string `json:"Code"`
	}

	rspStr := response.GetHttpContentString()
	err = json.Unmarshal([]byte(rspStr), &xx)
	if err != nil {
		return err
	}
	if xx.Message == "OK" && xx.Code == "OK" {
		return nil
	}
	return fmt.Errorf("aliyun sdk err %v", xx)
}
