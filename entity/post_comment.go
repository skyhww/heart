package entity

import "time"

type PostComment struct {
	Id         int        `db:"id"`
	UserId     int        `db:"user_id"`
	CreateTime *time.Time `db:"create_time"`
	Enable     int        `db:"enable"`
	Content    []byte     `db:"content"`
	PostId     int        `db:"post_id"`
	ReplyId   int        `db:"reply_id"`
}
