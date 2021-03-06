package service

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"heart/entity"
	"heart/helper"
	"heart/service/common"
	"heart/sms"
	"time"
)

//全局唯一
type Token struct {
	Token  string `json:"token"`
	UserId int64  `json:"-"`
}

type TokenHelper struct {
	Rds *redis.Pool
}

func (h *TokenHelper) GetToken(token string) (*Token, *base.Info) {
	userId, err := h.Rds.Get().Do("get", token)
	if err != nil {
		return nil, base.ServerError
	}
	if userId == nil {
		return nil, base.TokenExpired
	}
	go h.Rds.Get().Do("expire", token, time.Second*60*60*24)
	return &Token{Token: token, UserId: helper.Int2int64(userId.([]uint8))}, base.Success
}

type TokenService interface {
	CreateToken(user int64) (*Token, *base.Info)
	GetToken(token string) (*Token, *base.Info)
	Expire(token *Token) *base.Info
}

type SimpleTokenService struct {
	Pool *redis.Pool
	Ex   time.Duration
}

func (simpleTokenService *SimpleTokenService) Expire(token *Token) *base.Info {
	if token == nil || token.UserId == 0 || token.Token == "" {
		return base.Success
	}
	userId, err := simpleTokenService.Pool.Get().Do("get", token.Token)
	if err != nil {
		return base.ServerError
	}
	if userId.(int64) != token.UserId {
		return base.IllegalOperation
	}
	return base.Success
}

func (simpleTokenService *SimpleTokenService) CreateToken(user int64) (*Token, *base.Info) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, base.ServerError
	}
	token := &Token{uid.String(), user}
	ok, err := simpleTokenService.Pool.Get().Do("set", token.Token, user, "EX", simpleTokenService.Ex.Seconds(), "NX")
	if err != nil {
		logs.Error(err)
		return nil, base.ServerError
	}
	if ok.(string) != "OK" {
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
	SendRestPasswordCode(mobile string) *base.Info
	SendRegistCode(mobile string) *base.Info
	Regist(mobile, password, smsCode string) *base.Info
	LogOut(token *Token) *base.Info
}

type SimpleSecurity struct {
	Pool                  *redis.Pool
	SmsClient             sms.Sms
	UserPersist           entity.UserPersist
	TokenService          *SimpleTokenService
	UserFollowInfoPersist entity.UserFollowInfoPersist
}

//用户登录时，用户信息验证成功后，失效这个用户对应的token
func (security *SimpleSecurity) Login(mobile, password string) *base.Info {
	user, err := security.UserPersist.Get(mobile)
	if err != nil {
		logs.Error(err)
		return base.GetUserInfoFailed
	}
	if user == nil {
		return base.NoUserFound
	}
	if user.Id == 0 {
		return base.UsernameOrPasswordError
	}
	psd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	if psd != *user.Password {
		return base.UsernameOrPasswordError
	}
	t, in := security.TokenService.CreateToken(user.Id)
	if !in.IsSuccess() {
		return in
	}
	user.FollowCount = security.UserFollowInfoPersist.GetFollowCount(user.Id)
	user.FollowedCount = security.UserFollowInfoPersist.GetFollowedCount(user.Id)
	return base.NewSuccess(&User{User: user, Token: t})
}

func (security *SimpleSecurity) SendRegistCode(mobile string) *base.Info {
	u, _ := security.UserPersist.Get(mobile)
	if u != nil {
		return base.SignedUser
	}
	code := helper.CreateCaptcha()
	_, err := security.Pool.Get().Do("set", "regist_"+mobile, code, "EX", 60000, "NX")
	if err != nil {
		logs.Error(err)
		return base.SmsSendFailure
	}
	b := security.SmsClient.SendRegistCode(mobile, code)
	if !b {
		return base.SmsSendFailure
	}
	return base.Success
}
func (security *SimpleSecurity) SendRestPasswordCode(mobile string) *base.Info {
	u, err := security.UserPersist.Get(mobile)
	if err != nil && u == nil {
		return base.NonSignedUser
	}
	code := helper.CreateCaptcha()
	_, err = security.Pool.Get().Do("set", "rest_"+mobile, code, "EX", 60000, "NX")
	if err != nil {
		return base.SmsSendFailure
	}
	b := security.SmsClient.SendRegistCode(mobile, code)
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
	strCode := helper.Int2str(code.([]uint8))
	if strCode != smsCode {
		return base.SmsNotMatched
	}
	now := time.Now()
	name := helper.Random(8)
	psd := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user := &entity.User{Name: &name, CreateTime: &now, Mobile: &mobile, Password: &psd}
	if !security.UserPersist.Save(user) {
		return base.SaveUserFailed
	}
	//用户已经注册完成
	t, in := security.TokenService.CreateToken(user.Id)
	if !in.IsSuccess() {
		return base.ServerError
	}
	return base.NewSuccess(&User{User: user, Token: t})
}

func (security *SimpleSecurity) LogOut(token *Token) *base.Info {
	_, err := security.Pool.Get().Do("del", token.Token)
	if err != nil {
		logs.Error(err)
		return base.SmsSendFailure
	}
	return base.Success
}
