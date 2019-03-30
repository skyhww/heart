package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
)

type UserVideo struct {
	Id         int64      `db:"id" json:"id"`
	UserId     int64      `db:"user_id"`
	Url        string     `db:"url" json:"url"`
	StoreType  string     `db:"store_type"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable"`
	Hash       string     `db:"hash" json:"hash"`
}

type UserVideoPersist interface {
	Save(video *UserVideo) error
	Delete(video *UserVideo) error
}

type UserVideoDao struct {
	db *sqlx.DB
}

func (userVideoDao *UserVideoDao) Save(video *UserVideo) error {
	tx := userVideoDao.db.MustBegin()
	r, err := tx.NamedExec("INSERT INTO user_video (user_id, url,hash,CreateTime) VALUES (:user_id, :url, :hash,:create_time)", video)
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
	_, err := tx.NamedExec("update user_video set enbale=0 where id=:id", video)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
