package sms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"fmt"
	"heart/cfg"
)

type Sms interface {
	SendSmsCode(phoneNo string, content string) bool
}

type AliYun struct {
	AliYunConfig *cfg.AliYunConfig
	Client       *sdk.Client
}

func (aliYun *AliYun) getRequest(phoneNo string, code string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Method = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.Method
	request.Scheme = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.Scheme
	request.Domain = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.Domain
	request.Version = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.Version
	request.ApiName = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.ApiName
	request.QueryParams["PhoneNumbers"] = phoneNo
	request.QueryParams["SignName"] = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.SignName
	request.QueryParams["TemplateCode"] = aliYun.AliYunConfig.SmsConfig.SendSmsConfig.TemplateCode
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + code + "\"}"
	return request
}

func (aliYun *AliYun) SendSmsCode(phoneNo string, content string) bool {
	response, err := aliYun.Client.ProcessCommonRequest(aliYun.getRequest(phoneNo, content))
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
	return true
}
