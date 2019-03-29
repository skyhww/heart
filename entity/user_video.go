package entity

import "time"

type UserVideo struct {
	Id         int        `db:"id" json:"id"`
	UserId     int        `db:"user_id"`
	Url        string     `db:"url" json:"url"`
	StoreType  string     `db:"store_type"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable"`
	Hash       string     `db:"hash" json:"hash"`
}
