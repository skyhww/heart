package entity

import (
	"database/sql"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	Id            int64      `db:"id" json:"id"`
	Name          *string    `db:"name" json:"name"`
	IconUrl       *string    `db:"icon_url" json:"icon_url"`
	Signature     *string    `db:"signature" json:"signature"`
	CreateTime    *time.Time `db:"create_time" json:"create_time"`
	Mobile        *string    `db:"mobile" json:"-"`
	Enable        int        `db:"enable" json:"-"`
	Password      *string    `db:"password" json:"-"`
	FollowCount   int        `json:"follow_count"`
	FollowedCount int        `json:"followed_count"`
}

type UserPersist interface {
	Save(user *User) bool
	Get(mobile string) (*User, error)
	GetByUserName(userName string) (*User, error)
	Update(user *User) (*User, error)
	GetById(id int64) (*User, error)
	UpdatePassword(id int64, password string) bool
}

type UserDao struct {
	db *sqlx.DB
}

func NewUserPersist(db *sqlx.DB) UserPersist {
	return &UserDao{db: db}
}
func (userDao *UserDao) UpdatePassword(id int64, password string) bool {
	tx := userDao.db.MustBegin()
	_, err := tx.Exec("update  user set password=? where id=?", password, id)
	if err != nil {
		tx.Rollback()
		return false
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err)
		return false
	}
	return true
}
func (userDao *UserDao) Save(user *User) bool {
	tx := userDao.db.MustBegin()
	r, err := tx.NamedExec("INSERT INTO user (name, mobile,password,create_time) VALUES (:name , :mobile, :password,:create_time)", user)
	if err != nil {
		logs.Error(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err)
		return false
	}
	id, _ := r.LastInsertId()
	user.Id = id
	return true
}
func (userDao *UserDao) Get(mobile string) (*User, error) {
	user := &User{}
	err := userDao.db.Get(user, "select id,name,icon_url,create_time,password,signature from user where mobile=? ", mobile)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
func (userDao *UserDao) GetById(id int64) (*User, error) {
	user := &User{}
	err := userDao.db.Get(user, "select id,name,icon_url,signature,create_time from user where id=? ", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (userDao *UserDao) GetByUserName(userName string) (*User, error) {
	user := &User{}
	err := userDao.db.Get(user, "select id,name,icon_url,signature,create_time from user where name=? ", userName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (userDao *UserDao) Update(user *User) (*User, error) {
	tx := userDao.db.MustBegin()
	_, err := tx.NamedExec("update  user set name=:name,icon_url=:icon_url,signature=:signature where id=:id", user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	return user, err
}
