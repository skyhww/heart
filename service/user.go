package service

import (
	"github.com/astaxie/beego/logs"
	"heart/entity"
	"heart/service/common"
)

type User struct {
	*entity.User `json:"user"`
	Token        *Token             `json:"token"`
	UserPersist  entity.UserPersist `json:"-"`
}

type UserInfoService interface {
	UpdateName(token *Token, name *string) *base.Info
	UpdateSignature(token *Token, signature *string) *base.Info
	UpdateIcon(token *Token, icon *[]byte, name string) *base.Info
	ReadIcon(userId int64) (*base.Info, *[]byte, string)
	GetUserInfo(userId int64) *base.Info
}

type UserInfo struct {
	UserPersist          entity.UserPersist
	StoreService         StoreService
	UserExtraInfoPersist entity.UserExtraInfoPersist
}
func (user *UserInfo) GetUserInfo(userId int64) *base.Info {
	u, err := user.UserPersist.GetById(userId)
	if err != nil {
		logs.Error(err)
		return  base.ServerError
	}
	return base.NewSuccess(u)
}

func (user *UserInfo) ReadIcon(userId int64) (*base.Info, *[]byte, string) {
	u, err := user.UserPersist.GetById(userId)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound, nil, ""
	}
	b, name, err := user.StoreService.Get("icon", *u.IconUrl)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if b == nil || len(b) == 0 {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	return base.Success, &b, name
}

func (user *UserInfo) GetUser(token *Token) (*User, error) {
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return &User{User: u, Token: token, UserPersist: user.UserPersist}, nil
}

func (user *UserInfo) GetUserByName(name *string) (*entity.User, *base.Info) {
	u, err := user.UserPersist.GetByUserName(*name)
	if err != nil {
		logs.Error(err)
		return nil, base.ServerError
	}
	return u, base.Success
}

func (user *UserInfo) UpdateName(token *Token, name *string) *base.Info {
	u, f := user.GetUserByName(name)
	if !f.IsSuccess() {
		return f
	}
	//未更改
	if u != nil && u.Id == token.UserId {
		return base.NewSuccess(u)
	}
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.GetUserInfoFailed
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	u.Name = name
	u, err = user.UserPersist.Update(u)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	return base.NewSuccess(u)
}

func (user *UserInfo) UpdateSignature(token *Token, signature *string) *base.Info {
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.GetUserInfoFailed
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	u.Signature = signature
	u, err = user.UserPersist.Update(u)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	return base.NewSuccess(u)
}
func (user *UserInfo) UpdateIcon(token *Token, icon *[]byte, suffix string) *base.Info {
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	url, err := user.StoreService.Save("icon", icon, suffix)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	u.IconUrl = &url
	u, err = user.UserPersist.Update(u)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(u)
}
