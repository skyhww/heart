package service

import (
	"heart/entity"
	"heart/service/common"
	"time"
	"github.com/astaxie/beego/logs"
)

type Post struct {
	entity.UserPost
}

type PostService interface {
	GetComments(postId int64,page *base.Page) *base.Info
	//添加评论
	AddComment(token *Token, comment *Comment) *base.Info
	//添加回复
	AddReplay(token *Token, comment *Comment) *base.Info
	//获取回复
	GetReplay(comment *Comment, page base.Page) *base.Info
}

type Comment struct {
	entity.PostComment
}
type SimplePostService struct {
	PostCommentPersist entity.PostCommentPersist
	PostAttachPersist  entity.PostAttachPersist
	UserPersist      entity.UserPersist
}

func (simplePostService *SimplePostService)  GetComments(postId int64,page *base.Page) *base.Info{
	err:=simplePostService.PostCommentPersist.GetComments(&entity.UserPost{Id:postId},page)
	if err!=nil{
		return base.ServerError
	}
	return base.NewSuccess(page)

}
func (simplePostService *SimplePostService)  AddComment(token *Token, comment *Comment) *base.Info{
	u,err:=simplePostService.UserPersist.GetById(token.UserId)
	if err!=nil{
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	comment.PostComment.UserId=u.Id
	now:=time.Now()
	comment.PostComment.CreateTime=&now
	err=simplePostService.PostCommentPersist.Save(&comment.PostComment)
	if err!=nil{
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(comment)

}
func (simplePostService *SimplePostService)  GetReplay(comment *Comment, page *base.Page) *base.Info{
	now:=time.Now()
	comment.PostComment.CreateTime=&now
	err:=simplePostService.PostCommentPersist.GetReply(&comment.PostComment, page)
	if err!=nil{
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(comment)
}