package service

import (
	"heart/entity"
	"heart/service/common"
)

type User struct {
	*entity.User `json:"user"`
	Token *Token `json:"token"`
}

func (user *User) GetIcon() []byte {
	return nil
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
