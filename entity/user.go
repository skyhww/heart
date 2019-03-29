package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int        `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	IconUrl    string     `db:"icon_url" json:"icon_url"`
	Signature  []byte     `db:"signature" json:"signature"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Mobile     string     `db:"mobile" json:"mobile"`
	Enable     int        `db:"enable" json:"enable"`
	Password   string     `db:"password" json:"password"`
}
type UserPersist interface {
	Save(user *User)
	Get(id int) *User
}

type UserDao struct {
	db *sqlx.DB
}

func NewUserPersist(db *sqlx.DB) UserPersist {
	return &UserDao{db: db}
}
func (userDao *UserDao) Save(user *User) {
	userDao.db.BeginTx()
}
func (userDao *UserDao) Get(id int) *User {

}
