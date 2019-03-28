package service

import (
	"heart/entity"
	"time"
)
type Post struct {

}

type PostService interface {
	GetComments() *Comment
	GetUser() User
	GetCreateTime() *time.Time
}

type Comment struct {
	comment *entity.PostComment
	reply   []Comment
}
