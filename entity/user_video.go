package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
	"heart/service/common"
	"database/sql"
)

type UserVideo struct {
	Id         int64      `db:"id" json:"id"`
	UserId     int64      `db:"user_id"`
	Url        *string    `db:"url" json:"url"`
	StoreType  *string    `db:"store_type"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable"`
	Hash       *string    `db:"hash" json:"hash"`
	Content    *string    `db:"content" json:"content"`
}

type UserVideoPersist interface {
	Save(video *UserVideo) error
	Delete(video *UserVideo) error
	Get(video *UserVideo) error
	SelectByContent(userId int64,content string, page *base.Page) error
}

func NewUserVideoPersist(db *sqlx.DB) UserVideoPersist {
	return &UserVideoDao{db: db}
}

type UserVideoDao struct {
	db *sqlx.DB
}

func (userVideoDao *UserVideoDao) SelectByContent(userId int64,content string, page *base.Page) error {
	count := 0
	err := userVideoDao.db.Get(&count, "select count(1) from user_video where user_id=? and  content like '%"+content+"%' and enable=1 ",userId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	u := &[]UserVideo{}
	page.Count = count
	err = userVideoDao.db.Select(u, "select * from user_video where user_id=? and  content like '%"+content+"%' and enable=1 order by create_time desc limit ?,?",userId, (page.PageNo-1)*page.PageSize, page.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	page.Data = u
	return nil
}

func (userVideoDao *UserVideoDao) Save(video *UserVideo) error {
	tx := userVideoDao.db.MustBegin()
	r, err := tx.NamedExec("INSERT INTO user_video (user_id, url,hash,create_time,enable,content) VALUES (:user_id, :url, :hash,:create_time,1,:content)", video)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	id, _ := r.LastInsertId()
	video.Id = id
	return nil
}

func (userVideoDao *UserVideoDao) Delete(video *UserVideo) error {
	tx := userVideoDao.db.MustBegin()
	_, err := tx.NamedExec("update user_video set enbale=0 where id=:id and user_id=:user_id", video)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (userVideoDao *UserVideoDao) Get(video *UserVideo) error {
	tx := userVideoDao.db.MustBegin()
	_, err := tx.NamedExec("select id,url,hash,create_time from user_video where id=:id and user_id=:user_id", video)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
