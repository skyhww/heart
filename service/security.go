package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"heart/sms"
	"heart/helper"
	"heart/entity"
	"crypto/md5"
	"heart/service/common"
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
	Login(mobile, password string) *base.Info
	SendSmsCode(mobile string) *base.Info
	Regist(mobile, password, smsCode string) *base.Info
}

type SimpleSecurity struct {
	Pool        *redis.Pool
	SmsClient   sms.Sms
	UserPersist entity.UserPersist
}

func (security *SimpleSecurity) Login(mobile, password string) *base.Info {
	user, err := security.UserPersist.Get(mobile)
	if err != nil {
		return base.GetUserInfoFailed
	}
	if user.Id == 0 {
		return base.UsernameOrPasswordError
	}
	psd := string(md5.Sum([]byte(password))[:])
	if psd != user.Password {
		return base.UsernameOrPasswordError
	}
	return base.NewSuccess(user)
}
func (security *SimpleSecurity) SendSmsCode(mobile string) *base.Info {
	code := helper.CreateCaptcha()
	_, err := security.Pool.Get().Do("setnx", "regist_"+mobile, code, "EX", 60)
	if err != nil {
		return base.SmsSendFailure
	}
	b := security.SmsClient.SendSmsCode(mobile, code)
	if !b {
		return base.SmsSendFailure
	}
	return base.Success

}
func (security *SimpleSecurity) Regist(mobile, password, smsCode string) *base.Info {
	code, err := security.Pool.Get().Do("get", "regist_"+mobile)
	if err != nil {
		return base.SmsFindFailure
	}
	if code == nil {
		return base.SmsExpired
	}
	if code.(string) != smsCode {
		return base.SmsNotMatched
	}
	now := time.Now()
	user := &entity.User{Name: helper.Random(8), CreateTime: &now, Mobile: mobile, Password: password}
	if !security.UserPersist.Save(user) {
		return base.SaveUserFailed
	}
	return base.NewSuccess(user)
}
