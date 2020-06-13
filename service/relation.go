package service

import (
	"github.com/astaxie/beego/logs"
	"heart/entity"
	"heart/service/common"
	"time"
)

type UserFollowService interface {
	//关注
	Follow(token *Token, user int64) *base.Info
	//取消关注
	UnFollow(token *Token, user int64) *base.Info

	GetFollowUser(token *Token, page *base.Page) *base.Info
	GetFollowedUser(token *Token, page *base.Page) *base.Info
}

type SimpleUserFollowService struct {
	UserPersist           entity.UserPersist
	UserFollowInfoPersist entity.UserFollowInfoPersist
}

func (simpleUserFollowService *SimpleUserFollowService) Follow(token *Token, user int64) *base.Info {
	if token.UserId == user {
		return base.CantFollowYourSelf
	}
	u, err := simpleUserFollowService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	n := time.Now()
	f := &entity.UserFollowInfo{UserId: token.UserId, FollowUserId: user, CreateTime: &n}
	followUser, err := simpleUserFollowService.UserPersist.GetById(user)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if followUser == nil || followUser.Id == 0 {
		return base.NoFollowUserFound
	}
	err = simpleUserFollowService.UserFollowInfoPersist.Save(f)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}
func (simpleUserFollowService *SimpleUserFollowService) UnFollow(token *Token, user int64) *base.Info {
	u, err := simpleUserFollowService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simpleUserFollowService.UserFollowInfoPersist.DeleteUserFollowInfo(token.UserId, user)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}
func (simpleUserFollowService *SimpleUserFollowService) GetFollowUser(token *Token, page *base.Page) *base.Info {
	u, err := simpleUserFollowService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simpleUserFollowService.UserFollowInfoPersist.GetFollowUsers(token.UserId, page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(page)
}
func (simpleUserFollowService *SimpleUserFollowService) GetFollowedUser(token *Token, page *base.Page) *base.Info {
	u, err := simpleUserFollowService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simpleUserFollowService.UserFollowInfoPersist.GetFollowed(token.UserId, page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(page)
}
