package service

import (
	"heart/entity"
	"heart/service/common"
	"hash/crc32"
	"time"
)

type Video struct {
	entity.UserVideo
}
type VideoService interface {
	PushVideo(token *Token, content *[]byte, suffix string) *base.Info
	RemoveVideo(token *Token, id int64) *base.Info
	GetVideo(token *Token, id int64) (*base.Info, *[]byte, string)
}
type SimpleVideoService struct {
	StoreService     StoreService
	UserVideoPersist entity.UserVideoPersist
	UserPersist      entity.UserPersist
}

func (video *SimpleVideoService) PushVideo(token *Token, content *[]byte, suffix string) *base.Info {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	id, err := video.StoreService.Save("video", content, suffix)
	if err != nil {
		return base.ServerError
	}
	//优先保存文件
	now := time.Now()
	c := crc32.NewIEEE()
	hash := string(c.Sum(*content))
	userVideo := &entity.UserVideo{UserId: u.Id, StoreType: video.StoreService.GetType(), CreateTime: &now, Url: id, Hash: hash}
	err = video.UserVideoPersist.Save(userVideo)
	if err != nil {
		return base.ServerError
	}
	return base.NewSuccess(userVideo)
}
func (video *SimpleVideoService) RemoveVideo(token *Token, id int64) *base.Info {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	err = video.UserVideoPersist.Delete(&entity.UserVideo{Id: id})
	if err != nil {
		return base.ServerError
	}
	return base.Success
}
func (video *SimpleVideoService) GetVideo(token *Token, id int64) (*base.Info, *[]byte, string) {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError, nil, ""
	}
	if u == nil {
		return base.NoUserFound, nil, ""
	}
	v := &entity.UserVideo{Id: id, UserId: u.Id}
	err = video.UserVideoPersist.Get(v)
	if err != nil {
		return base.ServerError, nil, ""
	}
	b, n, err := video.StoreService.Get("video", v.Url)
	if err != nil {
		return base.ServerError, nil, ""
	}
	if b == nil {
		return base.ServerError, nil, ""
	}
	return base.Success, &b, n
}
