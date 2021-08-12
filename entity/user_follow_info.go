package entity

import (
	"database/sql"
	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
	"heart/service/common"
	"time"
)

type UserFollowInfo struct {
	Id           int64      `db:"id"`
	UserId       int64      `db:"user_id"`
	FollowUserId int64      `db:"follow_user_id"`
	CreateTime   *time.Time `db:"create_time"`
	Enable       int        `db:"enable"`
	//互相关注，逻辑字段
	FollowEachOther bool `db:"follow_each_other"`
}

type UserFollowInfoPersist interface {
	Save(userFollowInfo *UserFollowInfo) error
	//获取关注的人
	GetFollowUsers(userId int64, page *base.Page) error
	//取消关注
	DeleteUserFollowInfo(from, to int64) error
	//获取粉丝
	GetFollowed(userId int64, page *base.Page) error

	GetFollowedCount(userId int64) int
	GetFollowCount(userId int64) int
}

type UserFollowInfoDao struct {
	DB *sqlx.DB
}

func NewUserFollowInfoPersist(db *sqlx.DB) UserFollowInfoPersist {
	return &UserFollowInfoDao{DB: db}
}

func (userFollowInfoDao *UserFollowInfoDao) Save(userFollowInfo *UserFollowInfo) error {
	tx := userFollowInfoDao.DB.MustBegin()
	r, err := tx.NamedExec("update user_follow_info set enable=1,create_time=:create_time  where user_id=:user_id and follow_user_id=:follow_user_id", userFollowInfo)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 1 {
		return tx.Commit()
	} else {
		_, err := tx.NamedExec("insert into user_follow_info(user_id,follow_user_id,create_time,enable) values(:user_id,:follow_user_id,:create_time,1)", userFollowInfo)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}
}
func (userFollowInfoDao *UserFollowInfoDao) GetFollowUsers(userId int64, page *base.Page) error {
	info := &[]UserFollowInfo{}
	count := 0
	err := userFollowInfoDao.DB.Get(&count, "select count(id) from user_follow_info where user_id=? and enable=1 ", userId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count == 0 {
		return nil
	}
	err = userFollowInfoDao.DB.Select(info, "select user_id,follow_user_id,create_time from user_follow_info where user_id=? and enable=1 order by create_time desc limit ?,?", userId, page.PageSize*(page.PageNo-1), page.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	page.Count = count / page.PageSize
	page.Data = info
	return nil
}
func (userFollowInfoDao *UserFollowInfoDao) DeleteUserFollowInfo(from, to int64) error {
	tx := userFollowInfoDao.DB.MustBegin()
	_, err := tx.Exec("update user_follow_info set enable=0 where user_id=? and follow_user_id=?", from, to)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}
func (userFollowInfoDao *UserFollowInfoDao) GetFollowed(userId int64, page *base.Page) error {
	info := &[]UserFollowInfo{}
	count := 0
	err := userFollowInfoDao.DB.Get(&count, "select count(id) from user_follow_info where follow_user_id=? and enable=1 ", userId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count == 0 {
		return nil
	}
	err = userFollowInfoDao.DB.Select(info, "select user_id,follow_user_id,create_time from user_follow_info where follow_user_id=? and enable=1 order by create_time desc limit ?,?", userId, page.PageSize*(page.PageNo-1), page.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	page.Count = count
	page.Data = info
	return nil
}
func (userFollowInfoDao *UserFollowInfoDao) GetFollowedCount(userId int64) int {
	count := 0
	err := userFollowInfoDao.DB.Get(&count, "select count(id) from user_follow_info where follow_user_id=? and enable=1 ", userId)
	if err != nil && err != sql.ErrNoRows {
		logs.Error(err)
		return 0
	}

	return count
}
func (userFollowInfoDao *UserFollowInfoDao) GetFollowCount(userId int64) int {
	count := 0
	err := userFollowInfoDao.DB.Get(&count, "select count(id) from user_follow_info where user_id=? and enable=1 ", userId)
	if err != nil && err != sql.ErrNoRows {
		logs.Error(err)

		return 0
	}
	return count

}
