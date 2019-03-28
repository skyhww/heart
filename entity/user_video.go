package entity

import "time"

type UserVideo struct {
	Id         int        `db:"id"`
	UserId     int        `db:"user_id"`
	Url        string     `db:"url"`
	StoreType  string     `db:"store_type"`
	CreateTime *time.Time `db:"create_time"`
	Enable     int        `db:"enable"`
	Hash       string     `db:"hash"`
}
