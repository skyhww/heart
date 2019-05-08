package entity

import (
	"github.com/jmoiron/sqlx"
	"time"
	"heart/service/common"
	"database/sql"
)

type UserPost struct {
	Id         int64         `db:"id" json:"id"`
	UserId     int64         `db:"user_id" json:"user_id"`
	Content    *string       `db:"content" json:"content"`
	CreateTime *time.Time    `db:"create_time" json:"create_time"`
	PostAttach *[]PostAttach `json:"post_attach"`
	Attach     interface{}   `json:"attach"`
	Enable     int           `db:"enable" json:"-"`
}

type PostAttach struct {
	Id         int64      `db:"id" json:"id"`
	PostId     int64      `db:"post_id" json:"post_id"`
	Url        *string    `db:"url" json:"-"`
	No         int        `db:"no" json:"no"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable" json:"-"`
}

//文本评论
type PostComment struct {
	Id         int64      `db:"id" json:"id"`
	UserId     int64      `db:"user_id" json:"user_id"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Enable     int        `db:"enable" json:"-"`
	Content    string     `db:"content" json:"content"`
	PostId     int64      `db:"post_id" json:"post_id"`
	ReplyId    int64      `db:"reply_id" json:"reply_id"`
	Attach     *string    `db:"attach" json:"-"`
}

type PostCommentPersist interface {
	Save(postComment *PostComment) error

	GetComments(post *UserPost) (*[]PostComment, error)
	//评论回复
	GetReply(comments *PostComment, page *base.Page) error
	//删除评论
	Delete(post *PostComment) error

	Get(id int64) (*PostComment, error)
}

type UserPostPersist interface {
	Save(post *UserPost) error
	Get(userId int64, page *base.Page) error
	Delete(id int64) error
}

type PostAttachPersist interface {
	Get(userId int64, postId int64, page *base.Page) error
	GetAttach(attach int64) (*PostAttach, error)
}

type PostsPersist interface {
	Get(keyword string, page *base.Page) error
}
type PostsDao struct {
	DB *sqlx.DB
}

func NewPostsPersist(db *sqlx.DB) PostsPersist {
	return &PostsDao{DB: db}
}

func (postsDao *PostsDao) Get(keyword string, page *base.Page) error {
	post := &[]UserPost{}
	count := 0
	page.Data = post
	err := postsDao.DB.Get(&count, "select count(id) from user_post where   enable=1")
	if err != nil {
		return err
	}
	page.Count = count
	if count != 0 {
		err = postsDao.DB.Select(post, "select id,user_id,content,create_time from user_post where  enable=1   order by create_time desc limit ?,?", (page.PageNo-1)*page.PageSize, page.PageSize)
		if err != nil {
			return err
		}
	}
	return nil
}

type PostCommentDao struct {
	DB *sqlx.DB
}

func NewPostCommentPersist(db *sqlx.DB) PostCommentPersist {
	return &PostCommentDao{DB: db}
}
func (postCommentDao *PostCommentDao) Get(id int64) (*PostComment, error) {
	p := &PostComment{}
	err := postCommentDao.DB.Get(p, "select id,user_id,create_time,content,post_id,reply_id from post_comment where id=? ", id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if p.Id==0{
		return nil,nil
	}
	return p, nil
}

func (postCommentDao *PostCommentDao) Save(postComment *PostComment) error {
	tx := postCommentDao.DB.MustBegin()
	r, err := tx.NamedExec("insert into post_comment(user_id,create_time,enable,content,post_id,reply_id) values(:user_id,:create_time,1,:content,:post_id,:reply_id)", postComment)
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
func (postCommentDao *PostCommentDao) GetComments(post *UserPost) (*[]PostComment, error) {
	p := &[]PostComment{}
	err := postCommentDao.DB.Select(p, "select * from post_comment where post_id=?   order by create_time desc", post.Id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return p, nil
}
func (postCommentDao *PostCommentDao) Delete(comment *PostComment) error {
	tx := postCommentDao.DB.MustBegin()
	_, err := tx.Exec("update post_comment  set enable=0 where id=? and user_id=?", comment.Id, comment.UserId)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (postCommentDao *PostCommentDao) GetReply(comments *PostComment, page *base.Page) error {
	count := 0
	err := postCommentDao.DB.Select(count, "select * from post_comment where reply_id=?   order by create_time desc limit  ?,?", comments.Id, page.PageSize*page.PageNo, page.PageSize)

	p := &[]PostComment{}
	err = postCommentDao.DB.Select(p, "select * from post_comment where reply_id=?   order by create_time desc limit  ?,?", comments.Id, page.PageSize*page.PageNo, page.PageSize)
	if err != nil {
		return err
	}
	page.Data = p
	return nil
}

type UserPostDao struct {
	DB *sqlx.DB
}

func NewUserPostPersist(db *sqlx.DB) UserPostPersist {
	return &UserPostDao{DB: db}
}

type PostAttachDao struct {
	DB *sqlx.DB
}

func NewPostAttachPersist(db *sqlx.DB) PostAttachPersist {
	return &PostAttachDao{DB: db}
}
func (postAttachDao *PostAttachDao) GetAttach(attach int64) (*PostAttach, error) {
	att := &PostAttach{}
	err := postAttachDao.DB.Get(att, "select * from user_post_attach where id=? and enable=1", attach)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return att, nil
}
func (userPostDao *UserPostDao) Save(post *UserPost) error {
	tx := userPostDao.DB.MustBegin()
	r, err := tx.NamedExec("insert into user_post(user_id,content,create_time,enable) values(:user_id,:content,:create_time,1)", post)
	if err != nil {
		tx.Rollback()
		return err
	}
	id, err := r.LastInsertId()
	post.Id = id
	if post.PostAttach != nil && len(*post.PostAttach) > 0 {
		for index := range *post.PostAttach {
			(*post.PostAttach)[index].PostId = id
			r2, err := tx.NamedExec("insert into user_post_attach(post_id,url,no,create_time,enable) values(:post_id,:url,:no,:create_time,1) ", (*post.PostAttach)[index])
			if err != nil {
				tx.Rollback()
				return err
			}
			(*post.PostAttach)[index].Id, err = r2.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}
func (userPostDao *UserPostDao) Get(userId int64, page *base.Page) error {
	post := &[]UserPost{}
	count := 0
	page.Data = post
	err := userPostDao.DB.Get(&count, "select count(id) from user_post where   user_id=? and enable=1", userId)
	if err != nil {
		return err
	}
	page.Count = count
	if count != 0 {
		err = userPostDao.DB.Select(post, "select id,user_id,content,create_time from user_post where user_id=? and enable=1 order by create_time desc limit ?,?", userId, (page.PageNo-1)*page.PageSize, page.PageSize)
		if err != nil {
			return err
		}
	}
	return nil
}
func (userPostDao *UserPostDao) Delete(id int64) error {
	tx := userPostDao.DB.MustBegin()
	_, err := tx.Exec("update user_post set enable=0 where id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (postAttachDao *PostAttachDao) Get(userId int64, postId int64, page *base.Page) error {
	count := 0
	postAttach := &[]PostAttach{}
	page.Data = postAttach
	err := postAttachDao.DB.Get(&count, "select count(user_post_attach.id) from user_post_attach join user_post on user_post.id=user_post_attach.post_id and  user_post.user_id=?  and  user_post_attach.post_id=? and user_post.enable=1 and user_post_attach.enable=1", userId, postId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if count == 0 {
		return nil
	}
	page.Count = count
	err = postAttachDao.DB.Select(postAttach, "select user_post_attach.id,user_post_attach.post_id,user_post_attach.url,user_post_attach.no,user_post_attach.create_time from user_post_attach join user_post on user_post.id=user_post_attach.post_id and  user_post.user_id=?  and  user_post_attach.post_id=? and user_post.enable=1 and user_post_attach.enable=1 limit ?,?", userId, postId, page.PageSize*(page.PageNo-1), page.PageSize)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
