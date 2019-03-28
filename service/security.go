package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"heart/sms"
	"heart/math"
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
	SendSmsCode(token *Token, mobile, smsCode, use string) *Info
	Regist(token *Token, mobile, password, smsCode string) *Info
}

var SmsSendFailure = &Info{"000100", "短信验证码发送失败"}
var SmsFindFailure = &Info{"000101", "短信验证异常"}
var SmsExpired = &Info{"000102", "短信验证码已经过期"}
var SmsNotMatched = &Info{"000103", "短信验证码匹配失败"}

type SimpleSecurity struct {
	pool      *redis.Pool
	smsClient sms.Sms
}

func (security *SimpleSecurity) Login(token *Token, mobile, password string) *Info {
	return Success
}
func (security *SimpleSecurity) SendSmsCode(token *Token, mobile, smsCode, use string) *Info {
	code := math_helper.CreateCaptcha()
	_, err := security.pool.Get().Do("setnx", use+"_"+token.Token, code, "EX", 60)
	if err != nil {
		return SmsSendFailure
	}
	b := security.smsClient.Send(mobile, smsCode)
	if !b {
		return SmsSendFailure
	}
	return Success

}
func (security *SimpleSecurity) Regist(token *Token, mobile, password, smsCode string) *Info {
	code, err := security.pool.Get().Do("get", "loginSms_"+token.Token)
	if err != nil {
		return SmsFindFailure
	}
	if code == nil {
		return SmsExpired
	}
	if code.(string) != smsCode {
		return SmsNotMatched
	}
	return Success
}
