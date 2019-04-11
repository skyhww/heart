package entity

import (
	"database/sql"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type UserExtraInfo struct {
	UserId int `db:"user_id"`
	//关注用户数
	FollowUserCount int `db:"follow_user_count"`
	//收藏帖子数量
	Collections int `db:"collections"`
	//粉丝数
	FollowedUserCount int `db:"followed_user_count"`
	//未读消息数--普通消息
	UnreadMessageCount int `db:"unread_message_count"`
	//未读消息数--站点消息
	UnreadAdminMessageCount int `db:"unread_admin_message_count"`
}

type UserExtraInfoPersist interface {
	Save(user *UserExtraInfo) bool
	Update(user *UserExtraInfo) (*UserExtraInfo, error)
	GetById(id int64) (*UserExtraInfo, error)
}

type UserExtraInfoDao struct {
	db *sqlx.DB
}

func NewUserExtraInfoPersist(db *sqlx.DB) UserExtraInfoPersist {
	return &UserExtraInfoDao{db: db}
}

func (userDao *UserExtraInfoDao) Save(user *UserExtraInfo) bool {
	tx := userDao.db.MustBegin()
	//r, err := tx.NamedExec("INSERT INTO user_extra_info (user_id) VALUES (:user_id)", user)
	_, err := tx.NamedExec("INSERT INTO user_extra_info (user_id) VALUES (:user_id)", user)
	if err != nil {
		logs.Error(err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		logs.Error(err)
		return false
	}
	/*
		id, _ := r.LastInsertId()
		user.Id = id
	*/
	return true
}

func (userDao *UserExtraInfoDao) GetById(id int64) (*UserExtraInfo, error) {
	user := &UserExtraInfo{}

	err := userDao.db.Get(user, "select id,follow_user_count, collections, followed_user_count, unread_message_count, unread_admin_message_count from user_extra_info where id=? ", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (userDao *UserExtraInfoDao) Update(user *UserExtraInfo) (*UserExtraInfo, error) {
	tx := userDao.db.MustBegin()
	_, err := tx.NamedExec("update  user set follow_user_count=:follow_user_count, collections=:collections, followed_user_count=:followed_user_count, unread_message_count=:unread_message_count, unread_admin_message_count=:unread_admin_message_count where id=:id", user)
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
