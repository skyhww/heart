package service

import (
	"heart/entity"
	"time"
)

type Post struct {
}

type PostService interface {
	GetComments() *Comment
	GetUser() *User
	GetCreateTime() *time.Time
	//添加评论
	AddComment(user *User, comment Comment)
}

type Comment struct {
	comment *entity.PostComment
	reply   []Comment
}
