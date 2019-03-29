package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"heart/sms"
	"heart/helper"
	"heart/entity"
)

//全局唯一
type Token struct {
	Expire time.Duration
	Token  string
	Auth   bool
}

func (token *Token) isExpired() bool {
	return token.Expire <= 0
}

type TokenService interface {
	CreateToken() *Token
	GetToken(token string) *Token
}
type Security interface {
	Login(token *Token, mobile, password string) *Info
	SendSmsCode(mobile string) *Info
	Regist(token *Token, mobile, password, smsCode string) *Info
}

var SmsSendFailure = &Info{Code: "000100", Message: "短信验证码发送失败"}
var SmsFindFailure = &Info{Code: "000101", Message: "短信验证异常"}
var SmsExpired = &Info{Code: "000102", Message: "短信验证码已经过期"}
var SmsNotMatched = &Info{Code: "000103", Message: "短信验证码匹配失败"}

type SimpleSecurity struct {
	Pool        *redis.Pool
	SmsClient   sms.Sms
	UserPersist entity.UserPersist
}

func (security *SimpleSecurity) Login(token *Token, mobile, password string) *Info {
	return Success
}
func (security *SimpleSecurity) SendSmsCode( mobile string) *Info {
	code := helper.CreateCaptcha()
	_, err := security.Pool.Get().Do("setnx", "regist_"+mobile, code, "EX", 60)
	if err != nil {
		return SmsSendFailure
	}
	b := security.SmsClient.SendSmsCode(mobile, code)
	if !b {
		return SmsSendFailure
	}
	return Success

}
func (security *SimpleSecurity) Regist(token *Token, mobile, password, smsCode string) *Info {
	code, err := security.Pool.Get().Do("get", "regist_"+mobile)
	if err != nil {
		return SmsFindFailure
	}
	if code == nil {
		return SmsExpired
	}
	if code.(string) != smsCode {
		return SmsNotMatched
	}
	now := time.Now()
	user := &entity.User{Name: helper.Random(8), CreateTime: &now, Mobile: mobile, Password: password}
	if security.UserPersist.Save(user) {
		return NewSuccess(user)
	}
	return Success
}
