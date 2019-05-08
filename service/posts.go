package service

import (
	"heart/entity"
	"heart/service/common"
	"time"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type Post struct {
	entity.UserPost `json:"posts"`
}

type PostService interface {
	GetPosts(keyword string, token *Token, page *base.Page) *base.Info
}
type SimplePostService struct {
	PostsPersist         entity.PostsPersist
	PostAttachPersist    entity.PostAttachPersist
	ElasticSearchService ElasticSearchService
}

func (simplePostService *SimplePostService) GetPosts(keyword string, token *Token, page *base.Page) *base.Info {
	key := make(map[string]interface{})
	key["content"] = keyword
	key["enable"] = true
	bs, err := simplePostService.ElasticSearchService.Query(key, "3dheart", "posts", page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if len(bs) == 0 {
		return base.NewSuccess(page)
	}
	data := make([]entity.UserPost,0)
	for _, b := range bs {
		tmp := &entity.UserPost{}
		err = json.Unmarshal(b, tmp)
		if err != nil {
			logs.Error(err)
			return base.ServerError
		}
		data=append(data, *tmp)
	}
	page.Data = &data

	if page.Data != nil {
		data, _ := page.Data.(*[]entity.UserPost)
		attachPage := &base.Page{PageNo: 1, PageSize: 5}
		for index := range *data {
			err = simplePostService.PostAttachPersist.Get((*data)[index].UserId, (*data)[index].Id, attachPage)
			if err != nil {
				logs.Error(err)
				return base.ServerError
			}
			(*data)[index].Attach = attachPage
		}
	}
	return base.NewSuccess(page)
}

type UserPostService interface {
	GetComments(postId int64) *base.Info
	//添加评论
	AddComment(token *Token, comment *Comment) *base.Info
	//添加回复
	AddReplay(token *Token, comment *Comment) *base.Info
	//获取回复
	GetReplay(comment *Comment, page *base.Page) *base.Info
	//发贴
	PutPosts(token *Token, post *Post) *base.Info
	//删除帖子,如果帖子被删除，则帖子对应的所有的评论也被删除
	DeletePosts(token *Token, id int64) *base.Info
	//
	GetPosts(token *Token, page *base.Page) *base.Info
	//删除评论
	DeleteComment(token *Token, id int64) *base.Info
}

type Comment struct {
	entity.PostComment
}

type SimpleUserPostService struct {
	PostCommentPersist entity.PostCommentPersist
	PostAttachPersist  entity.PostAttachPersist
	UserPersist        entity.UserPersist
	UserPostPersist    entity.UserPostPersist
}

func (simplePostService *SimpleUserPostService) GetPosts(token *Token, page *base.Page) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simplePostService.UserPostPersist.Get(token.UserId, page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if page.Data != nil {
		data, _ := page.Data.(*[]entity.UserPost)
		attachPage := &base.Page{PageNo: 1, PageSize: 5}
		for index := range *data {
			err = simplePostService.PostAttachPersist.Get(token.UserId, (*data)[index].Id, attachPage)
			if err != nil {
				logs.Error(err)
				return base.ServerError
			}
			(*data)[index].Attach = attachPage
		}
	}
	return base.NewSuccess(page)
}

func (simplePostService *SimpleUserPostService) GetComments(postId int64) *base.Info {
	c, err := simplePostService.PostCommentPersist.GetComments(&entity.UserPost{Id: postId})
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(c)

}

func (simplePostService *SimpleUserPostService) DeleteComment(token *Token, id int64) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simplePostService.PostCommentPersist.Delete(&entity.PostComment{Id: id})
	if err != nil {
		return base.ServerError
	}
	return base.Success

}
func (simplePostService *SimpleUserPostService) AddComment(token *Token, comment *Comment) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	comment.PostComment.UserId = u.Id
	now := time.Now()
	comment.PostComment.CreateTime = &now
	err = simplePostService.PostCommentPersist.Save(&comment.PostComment)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(comment)

}
func (simplePostService *SimpleUserPostService) GetReplay(comment *Comment, page *base.Page) *base.Info {
	now := time.Now()
	comment.PostComment.CreateTime = &now
	err := simplePostService.PostCommentPersist.GetReply(&comment.PostComment, page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(comment)
}
func (simplePostService *SimpleUserPostService) AddReplay(token *Token, comment *Comment) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	now := time.Now()
	comment.CreateTime = &now
	c, err := simplePostService.PostCommentPersist.Get(comment.ReplyId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if c == nil {
		return base.CommentNotFound
	}
	comment.UserId = token.UserId
	comment.PostId = c.PostId
	err = simplePostService.PostCommentPersist.Save(&comment.PostComment)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(comment)
}

func (simplePostService *SimpleUserPostService) PutPosts(token *Token, post *Post) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	now := time.Now()
	post.UserId = token.UserId
	post.CreateTime = &now
	if post.PostAttach != nil && len(*post.PostAttach) > 0 {
		for _, v := range *post.PostAttach {
			v.CreateTime = &now
		}
	}
	err = simplePostService.UserPostPersist.Save(&post.UserPost)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(post)
}
func (simplePostService *SimpleUserPostService) DeletePosts(token *Token, id int64) *base.Info {
	u, err := simplePostService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err=simplePostService.UserPostPersist.Delete(token.UserId,id)
	if err!=nil{
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}

type CommentService interface {
	//如果评论被删除，则这个评论下的所有回复都会被删除
	DeleteComment(token *Token, id int64) *base.Info
}

type SimpleCommentService struct {
	PostCommentPersist entity.PostCommentPersist
	UserPersist        entity.UserPersist
}

func (simpleCommentService *SimpleCommentService) DeleteComment(token *Token, id int64) *base.Info {
	u, err := simpleCommentService.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = simpleCommentService.PostCommentPersist.Delete(&entity.PostComment{UserId: token.UserId, Id: id})
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}

type PostAttachService interface {
	GetAttach(token *Token, attachId int64) (*base.Info, *[]byte, string)
}

type SimplePostAttachService struct {
	PostAttachPersist entity.PostAttachPersist
	StoreService      StoreService
}

func (simplePostAttachService *SimplePostAttachService) GetAttach(token *Token, attachId int64) (*base.Info, *[]byte, string) {

	a, err := simplePostAttachService.PostAttachPersist.GetAttach(attachId)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	b, name, err := simplePostAttachService.StoreService.Get("post/attach", *a.Url)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	return base.Success, &b, name
}
