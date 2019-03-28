package service

import (
	"heart/entity"
	"time"
)

type Post struct {
	entity.UserPost
}

type PostService interface {
	GetComments() *Comment
	GetUser() *User
	GetCreateTime() *time.Time
	//添加评论
	AddComment(user *User, comment Comment)
}

type Comment struct {
	entity.PostComment
	reply   []Comment
}
