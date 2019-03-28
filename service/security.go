package service

import "time"

//全局唯一
type Token struct {
	expire time.Duration
	token  string
	auth bool
}

func (token *Token) isExpired() bool {
	return token.expire <= 0
}

type TokenService interface {
	CreateToken() *Token
	GetToken(token string) *Token
}
type Security interface {
	Login(token *Token, userName, password string) (*User, string)
	SendSmsCode(token *Token, smsCode, use string) (bool, string)
	Regist(token *Token, userName, password, smsCode string) (bool, string)
}
