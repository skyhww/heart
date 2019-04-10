package sms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"fmt"
	"heart/cfg"
)

type Sms interface {
	SendRegistCode(phoneNo string, content string) bool
	SendResetPasswordCode(phoneNo string, content string) bool
}

type AliYun struct {
	AliYunConfig *cfg.AliYunConfig
	Client       *sdk.Client
}

func (aliYun *AliYun) getRequest(phoneNo string, code string,c *cfg.SendSmsConfig) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	request.Method = c.Method
	request.Scheme = c.Scheme
	request.Domain = c.Domain
	request.Version = c.Version
	request.ApiName = c.ApiName
	request.QueryParams["PhoneNumbers"] = phoneNo
	request.QueryParams["SignName"] = c.SignName
	request.QueryParams["TemplateCode"] = c.Code
	request.QueryParams["TemplateParam"] = "{\"code\":\"" + code + "\"}"
	return request
}

func (aliYun *AliYun) SendRegistCode(phoneNo string, content string) bool {
	response, err := aliYun.Client.ProcessCommonRequest(aliYun.getRequest(phoneNo, content,aliYun.AliYunConfig.SmsConfig.Regist))
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
	return true
}

func (aliYun *AliYun) SendResetPasswordCode(phoneNo string, content string) bool {
	response, err := aliYun.Client.ProcessCommonRequest(aliYun.getRequest(phoneNo, content,aliYun.AliYunConfig.SmsConfig.ResetPassword))
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
	return true
}
