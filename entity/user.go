package entity

import "time"

type User struct {
	Id         int        `db:"id"`
	Name       string     `db:"name"`
	IconUrl    string     `db:"icon_url"`
	Signature  []byte     `db:"signature"`
	CreateTime *time.Time `db:"create_time"`
	Mobile     string     `db:"mobile"`
	Enable     int        `db:"enable"`
	Password   string     `db:"password"`
}
type UserPersist interface {
	Save(user *User)
	Get(id int) *User
}
