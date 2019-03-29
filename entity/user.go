package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int64      `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	IconUrl    string     `db:"icon_url" json:"icon_url"`
	Signature  []byte     `db:"signature" json:"signature"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Mobile     string     `db:"mobile" json:"mobile"`
	Enable     int        `db:"enable" json:"enable"`
	Password   string     `db:"password" json:"password"`
}
type UserPersist interface {
	Save(user *User) bool
	Get(id int64) *User
}

type UserDao struct {
	db *sqlx.DB
}

func NewUserPersist(db *sqlx.DB) UserPersist {
	return &UserDao{db: db}
}
func (userDao *UserDao) Save(user *User) bool {
	tx := userDao.db.MustBegin()
	r, err := tx.NamedExec("INSERT INTO user (name, mobile,password,create_time) VALUES (:name, , :mobile, :password,create_time)", user)
	if err != nil {
		return false
	}
	err = tx.Commit()
	if err != nil {
		return false
	}
	id, _ := r.LastInsertId()
	user.Id = id
	return true
}
func (userDao *UserDao) Get(id int64) *User {
	return nil
}
