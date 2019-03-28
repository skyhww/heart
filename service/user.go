package service

import (
	"heart/entity"
	"time"
)

type User struct {
	entity *entity.User
}

func (user *User) GetId() int {
	return user.entity.Id
}
func (user *User) GetName() string {
	return user.entity.Name
}
func (user *User) GetIcon() []byte {
	return nil
}
func (user *User) GetSignature() []byte {
	return nil
}
func (user *User) GetCreateTime() *time.Time {
	return user.entity.CreateTime
}
func (user *User) GetMobile() string {
	return user.entity.Mobile
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
