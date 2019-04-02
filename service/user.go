package service

import (
	"heart/entity"
	"heart/service/common"
	"github.com/astaxie/beego/logs"
)

type User struct {
	*entity.User                   `json:"user"`
	Token       *Token             `json:"token"`
	UserPersist entity.UserPersist `json:"-"`
}

func (user *User) GetIcon() []byte {
	return nil
}

type UserInfoService interface {
	UpdateName(token *Token, name *string) (*base.Info)
	UpdateSignature(token *Token, signature *string) (*base.Info)
	UpdateIcon(token *Token, icon *[]byte) (*base.Info)
}

type UserInfo struct {
	UserPersist  entity.UserPersist
	StoreService StoreService
}

func (user *UserInfo) GetUserByName(name *string) (*entity.User, *base.Info) {
	u, err := user.UserPersist.GetByUserName(*name)
	if err != nil {
		return nil, base.ServerError
	}
	return u, base.Success
}

func (user *UserInfo) UpdateName(token *Token, name *string) (*base.Info) {
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
		return base.GetUserInfoFailed
	}
	if u == nil {
		return base.NoUserFound
	}
	u.Name = name
	u, err = user.UserPersist.Update(u)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	return base.NewSuccess(u)
}

func (user *UserInfo) UpdateSignature(token *Token, signature *string) (*base.Info) {
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.GetUserInfoFailed
	}
	if u == nil {
		return base.NoUserFound
	}
	u.Signature = signature
	u, err = user.UserPersist.Update(&entity.User{Id: token.UserId, Signature: signature})
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	return base.NewSuccess(u)
}
func (user *UserInfo) UpdateIcon(token *Token, icon *[]byte) (*base.Info) {
	u, err := user.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	url, err := user.StoreService.Save("icon", icon)
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

type UserService interface {
	GetExtraInfo() *entity.UserExtraInfo
	//获取
	GetFollowUsers(page *base.Page) []User
	//获取粉丝
	GetFollowedUser(page *base.Page) []User
	//获取发布的帖子
	GetPosts(page *base.Page) *Post
	//获取未读的消息
	GetUnreadMessages() []Message
	//发布视频
	CommitVideo(video *Video) bool
	//发贴
	CommitPosts(post *Post) bool
	//关注
	Follow(user *User) bool
	//收藏帖子
	Collect(post *Post) bool
}
