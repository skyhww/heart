package service

import (
	"heart/service/common"
	"heart/entity"
	"time"
	"github.com/astaxie/beego/logs"
)

type CollectorService interface {
	Add(token *Token, postId int64) *base.Info
	Remove(token *Token, postId int64) *base.Info
	Get(token *Token, page *base.Page) *base.Info
	GetCount(token *Token) *base.Info
}

type SimpleCollectorService struct {
	UserCollectionInfoPersist entity.UserCollectionInfoPersist
	UserPersist               entity.UserPersist
}

func (collectorService *SimpleCollectorService) Add(token *Token, postId int64) *base.Info {
	u, err := collectorService.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	n := time.Now()
	c := entity.UserCollectionInfo{PostId: postId, UserId: token.UserId, CreateTime: &n}
	err = collectorService.UserCollectionInfoPersist.Save(&c)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}
func (collectorService *SimpleCollectorService) Remove(token *Token, postId int64) *base.Info {
	u, err := collectorService.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = collectorService.UserCollectionInfoPersist.Delete(postId, token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}
func (collectorService *SimpleCollectorService) Get(token *Token, page *base.Page) *base.Info {
	u, err := collectorService.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = collectorService.UserCollectionInfoPersist.Get(token.UserId, page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(page)
}
func (collectorService *SimpleCollectorService) GetCount(token *Token) *base.Info {
	u, err := collectorService.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	c, err := collectorService.UserCollectionInfoPersist.GetCount(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(c)
}
