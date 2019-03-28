package service

import (
	"heart/entity"
)

type User struct {
	 entity.User
}

func (user *User) GetIcon() []byte {
	return nil
}


type UserService interface {
	GetExtraInfo() *entity.UserExtraInfo
	//获取
	GetFollowUsers(page *Page) []User
	//获取粉丝
	GetFollowedUser(page *Page) []User
	//获取发布的帖子
	GetPosts(page *Page) *Post
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


