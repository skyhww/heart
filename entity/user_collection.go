package entity

import (
	"time"
	"heart/service/common"
	"github.com/jmoiron/sqlx"
	"database/sql"
)

type UserCollectionInfo struct {
	Id         int64      `db:"id" json:"-"`
	PostId     int64      `db:"post_id" json:"post_id"`
	UserId     int64      `db:"user_id" json:"user_id"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable"`
}

type UserCollectionInfoPersist interface {
	Save(userCollectionInfo *UserCollectionInfo) error
	Get(userId int64, page *base.Page) error
	Delete(postId int64, userId int64) error
	GetCount(userId int64) (int, error)
}

type UserCollectionInfoDao struct {
	DB *sqlx.DB
}

func NewUserCollectionInfoPersist(db *sqlx.DB)UserCollectionInfoPersist  {
	return &UserCollectionInfoDao{DB:db}
}

func (userCollectionInfoDao *UserCollectionInfoDao) Save(userCollectionInfo *UserCollectionInfo) error {
	tx := userCollectionInfoDao.DB.MustBegin()
	r, err := tx.NamedExec("update user_collection_info set enable=1 ,create_time=:create_time where user_id=:user_id and post_id=:post_id", userCollectionInfo)
	if err != nil {
		tx.Rollback()
		return err
	}
	c, err := r.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if c > 1 {
		return tx.Commit()
	}
	_, err = tx.NamedExec("insert into user_collection_info(post_id,user_id,create_time,enable) values(:post_id,:user_id,:create_time,1)", userCollectionInfo)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
func (userCollectionInfoDao *UserCollectionInfoDao) Get(userId int64, page *base.Page) error {
	r := &[]UserCollectionInfo{}
	err := userCollectionInfoDao.DB.Select(r, "select * from user_collection_info where user_id=? and enable=1 order by create_time desc limit ?,?", userId, page.PageSize*(page.PageNo-1), page.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	page.Data = r
	return nil
}
func (userCollectionInfoDao *UserCollectionInfoDao) GetCount(userId int64) (int, error) {
	count := 0
	err := userCollectionInfoDao.DB.Get(&count, "select count(1) from user_collection_info where user_id=? and enable=1 ", userId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return count, nil
}
func (userCollectionInfoDao *UserCollectionInfoDao) Delete(postId int64, userId int64) error {
	tx := userCollectionInfoDao.DB.MustBegin()
	_, err := tx.Exec("update user_collection_info set enable=0  where post_id=? and user_id=? ", postId, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
