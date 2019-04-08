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
//文本评论
type PostComment struct {
	Id         int64      `db:"id"`
	UserId     int64      `db:"user_id"`
	CreateTime *time.Time `db:"create_time"`
	Enable     int        `db:"enable"`
	Content    string     `db:"content"`
	PostId     int64      `db:"post_id"`
	ReplyId    int64      `db:"reply_id"`
}

type PostCommentPersist interface {
	Save(postComment *PostComment) error
	//一级评论
	GetComments(post *UserPost, page *base.Page) error
	//评论回复
	GetReply(comments *PostComment, page *base.Page) error
	//删除评论
	Delete(post *PostComment)error
}

type UserPostPersist interface {
	Save(post *UserPost) error
	Get(userId int64, page base.Page) error
	Delete(id int64) error
}

type PostAttachPersist interface {
	Get(userId int64, postId int64, page base.Page) error
}
type PostCommentDao struct {
	DB *sqlx.DB
}

func (postCommentDao *PostCommentDao) Save(postComment *PostComment) error {
	tx := postCommentDao.DB.MustBegin()
	r, err := tx.Exec("insert into post_comment(user_id,create_time,enable,content,post_id,reply_id) values(:user_id,:create_time,:enable,:content,:post_id,:reply_id)", postComment)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	postComment.Id, _ = r.LastInsertId()
	return nil
}
func (postCommentDao *PostCommentDao) GetComments(post *UserPost, page *base.Page) error {
	p:=&[]PostComment{}
	err:=postCommentDao.DB.Select(p,"select * from post_comment where post_id=? and reply_id is null order by create_time desc limit  ?,?",post.Id,page.PageSize*page.PageNo,page.PageSize)
	if err!=nil{
		return err
	}
	return nil
}
func (postCommentDao *PostCommentDao) Delete(comment *PostComment)error {
	tx:=postCommentDao.DB.MustBegin()
	_,err:=tx.Exec("update post_comment  set enable=0 where id=? and user_id=?",comment.Id,comment.UserId)
	if err!=nil{
		tx.Rollback()
		return err
	}
	err=tx.Commit()
	if err!=nil{
		return err
	}
	return nil
}
func (postCommentDao *PostCommentDao) GetReply(comments *PostComment, page *base.Page) error {
	count:=0
	err:=postCommentDao.DB.Select(count,"select * from post_comment where reply_id=?   order by create_time desc limit  ?,?",comments.Id,page.PageSize*page.PageNo,page.PageSize)

	p:=&[]PostComment{}
	err=postCommentDao.DB.Select(p,"select * from post_comment where reply_id=?   order by create_time desc limit  ?,?",comments.Id,page.PageSize*page.PageNo,page.PageSize)
	if err!=nil{
		return err
	}
	page.Data=p
	return nil
}
type UserPostDao struct {
	DB *sqlx.DB
}
type PostAttachDao struct {
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
	count := 0
	err := userPostDao.DB.Select(&count, "select count(id) from user_post where   user_id=? and enable=1", userId)
	if err != nil {
		return err
	}
	err = userPostDao.DB.Select(post, "select id,user_id,content,create_time from user_post where user_id=? and enable=1 order by create_time desc limit ?,?", userId, page.PageSize*page.PageNo, page.PageSize)
	if err != nil {
		return err
	}
	page.Count = count
	page.Data = post
	return nil
}
func (userPostDao *UserPostDao) Delete(id int64) error {
	_, err := userPostDao.DB.Exec("update user_post set enable=0 where id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func (postAttachDao *PostAttachDao) Get(userId int64, postId int64, page base.Page) error {
	count := 0
	err := postAttachDao.DB.Select(&count, "select count(user_post_attach.id) from user_post_attach join user_post on user_post.id=user_post_attach.post_id and  user_post.user_id=?  and  user_post_attach.post_id=? and user_post.enable=1 and user_post_attach.enable=1 limit ?,?", userId, postId, page.PageSize*page.PageNo, page.PageSize)
	if err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	postAttach := &[]PostAttach{}
	err = postAttachDao.DB.Select(&postAttach, "select user_post_attach.id,user_post_attach.post_id,user_post_attach.url,user_post_attach.no,user_post_attach.create_time from user_post_attach join user_post on user_post.id=user_post_attach.post_id and  user_post.user_id=?  and  user_post_attach.post_id=? and user_post.enable=1 and user_post_attach.enable=1 limit ?,?", userId, postId, page.PageSize*page.PageNo, page.PageSize)
	if err != nil {
		return err
	}
	return nil
}
