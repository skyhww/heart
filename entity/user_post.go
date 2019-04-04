package entity

import (
	"github.com/jmoiron/sqlx"
	"time"
	"heart/service/common"
)

type UserPost struct {
	Id         int64      `db:"id"`
	UserId     int64      `db:"user_id"`
	Content    *string    `db:"content"`
	CreateTime *time.Time `db:"create_time"`
	PostAttach *[]PostAttach
	Enable     int        `db:"enable"`
}

type PostAttach struct {
	Id         int64      `db:"id"`
	PostId     int64      `db:"post_id"`
	Url        *string    `db:"url"`
	No         int        `db:"no"`
	CreateTime *time.Time `db:"create_time"`
	Enable     int        `db:"enable"`
}

type UserPostPersist interface {
	Save(post *UserPost) error
	Get(userId int64, page base.Page) error
	Delete(id int64) error
}

type UserPostDao struct {
	DB *sqlx.DB
}

func (userPostDao *UserPostDao) Save(post *UserPost) error {
	tx := userPostDao.DB.MustBegin()
	r, err := tx.Exec("insert into user_post(user_id,content,create_time,enable) values(:user_id,:content,:create_time,1)", post)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := r.LastInsertId()
	post.Id = id
	if post.PostAttach != nil && len(*post.PostAttach) > 0 {
		for _, v := range *post.PostAttach {
			r, err = tx.Exec("insert into user_post_attach(post_id,url,no,create_time,enable) values(post_id,url,no,create_time,1) ", v)
			if err != nil {
				tx.Rollback()
				return err
			}
			v.Id, _ = r.LastInsertId()
		}
	}
	return tx.Commit()
}
func (userPostDao *UserPostDao) Get(userId int64, page base.Page) error {
	post := &[]UserPost{}
	userPostDao.DB.Select(post, "select id,user_id,content,create_time from user_post where user_id=? and enable=1 order by create_time desc limit ?,?", userId, page.PageSize*page.PageNo, page.PageSize)
	return nil
}
func (userPostDao *UserPostDao) Delete(id int64) error {
	return nil
}
