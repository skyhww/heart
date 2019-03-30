package service

import (
	"time"
	"github.com/garyburd/redigo/redis"
	"heart/sms"
	"heart/helper"
	"heart/entity"
	"crypto/md5"
	"heart/service/common"
	"github.com/satori/go.uuid"
)

//全局唯一
type Token struct {
	Token  string
	UserId int64
}

type TokenService interface {
	CreateToken(user int64) (*Token, *base.Info)
	GetToken(token string) (*Token, *base.Info)
}

type SimpleTokenService struct {
	Pool   *redis.Pool
	expire time.Duration
}

func (simpleTokenService *SimpleTokenService) CreateToken(user int64) (*Token, *base.Info) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, base.ServerError
	}
	token := &Token{uid.String(), user}
	count, err := simpleTokenService.Pool.Get().Do("setnx", token.Token, user, "EX", simpleTokenService.expire.Seconds())
	if err != nil {
		return nil, base.ServerError
	}
	if count.(int) == 0 {
		return nil, base.ServerError
	}
	return token, base.Success

}
func (simpleTokenService *SimpleTokenService) GetToken(token string) (*Token, *base.Info) {
	userId, err := simpleTokenService.Pool.Get().Do("get", token)
	if err != nil {
		return nil, base.ServerError
	}
	return &Token{token, userId.(int64)}, base.Success
}

type Security interface {
	Login(mobile, password string) *base.Info
	SendSmsCode(mobile string) *base.Info
	Regist(mobile, password, smsCode string) *base.Info
}

type SimpleSecurity struct {
	Pool         *redis.Pool
	SmsClient    sms.Sms
	UserPersist  entity.UserPersist
	tokenService *SimpleTokenService
}

func (security *SimpleSecurity) Login(mobile, password string) *base.Info {
	user, err := security.UserPersist.Get(mobile)
	if err != nil {
		return base.GetUserInfoFailed
	}
	if user.Id == 0 {
		return base.UsernameOrPasswordError
	}
	b:=md5.Sum([]byte(password))
	psd := string(b[:])
	if psd != user.Password {
		return base.UsernameOrPasswordError
	}
	t, in := security.tokenService.CreateToken(user.Id)
	if !in.IsSuccess() {
		return in
	}
	return base.NewSuccess(&User{User: user, Token: t})
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
	//用户已经注册完成
	t, in := security.tokenService.CreateToken(user.Id)
	if !in.IsSuccess() {
		return base.ServerError
	}
	return base.NewSuccess(&User{User: user, Token: t})
}
