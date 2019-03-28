package entity

import "time"

type UserCollectionInfo struct {
	Id         int        `db:"id"`
	PostId     int        `db:"post_id"`
	UserId     int        `db:"user_id"`
	CreateTime *time.Time `db:"create_time"`
	Enable     int        `db:"enable"`
}

type UserCollectionInfoPersist interface {
	Save(userCollectionInfo *UserCollectionInfo)
	Get(id int) *UserCollectionInfo
}
