package entity

import "time"

type UserFollowInfo struct {
	Id           int        `db:"id"`
	UserId       int        `db:"user_id"`
	FollowUserId int        `db:"follow_user_id"`
	CreateTime   *time.Time `db:"create_time"`
	Enable       int        `db:"enable"`
}

