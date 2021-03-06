package base

import (
	JSON "encoding/json"
)

type Json interface {
	toJsonString() string
}

type ByteChan struct {
	Content []byte
	Last    bool
}

type Page struct {
	PageSize int         `json:"page_size"`
	PageNo   int         `json:"page_no"`
	Count    int         `json:"count"`
	Data     interface{} `json:"data"`
}

type Info struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var Success = &Info{Code: "000000", Message: ""}

func NewSuccess(data interface{}) *Info {
	return &Info{Success.Code, Success.Message, data}
}

func (e *Info) toJsonString() string {
	b, _ := JSON.Marshal(e)
	return string(b)
}

func (e *Info) IsSuccess() bool {
	return e == Success || e.Code == Success.Code
}

var SmsSendFailure = &Info{Code: "000100", Message: "短信验证码发送失败"}
var SmsFindFailure = &Info{Code: "000101", Message: "短信验证异常"}
var SmsExpired = &Info{Code: "000102", Message: "短信验证码已经过期"}
var SmsNotMatched = &Info{Code: "000103", Message: "短信验证码匹配失败"}
var GetUserInfoFailed = &Info{Code: "000104", Message: "获取用户信息失败"}
var NonSignedUser = &Info{Code: "000105", Message: "用户未注册"}
var UsernameOrPasswordError = &Info{Code: "000106", Message: "用户名或密码错误"}
var SaveUserFailed = &Info{Code: "000107", Message: "保存用户失败"}
var ServerError = &Info{Code: "000108", Message: "服务器繁忙"}
var IllegalOperation = &Info{Code: "000109", Message: "非法操作"}
var NoUserFound = &Info{Code: "000120", Message: "用户不存在"}
var TokenExpired = &Info{Code: "000121", Message: "token已过期，请重新登录"}
var SignedUser = &Info{Code: "000122", Message: "用户已注册"}
var TargetUserNotFound = &Info{Code: "000123", Message: "接收用户不存在"}
var MessageAttachNotFound = &Info{Code: "000124", Message: "消息附件不存在"}
var CommentNotFound = &Info{Code: "000125", Message: "评论不存在"}
var CantFollowYourSelf = &Info{Code: "000126", Message: "不需要关注自己"}
var NoFollowUserFound = &Info{Code: "000127", Message: "关注的用户不存在"}
var IllegalWord = &Info{Code: "000128", Message: "敏感词汇"}