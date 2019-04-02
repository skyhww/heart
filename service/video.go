package service

import (
	"heart/entity"
	"heart/service/common"
	"hash/crc32"
	"time"
)

type Video struct {
	entity.UserVideo
	StoreService     StoreService
	UserVideoPersist entity.UserVideoPersist
}

type VideoService interface {
	PushVideo(video *Video, content []byte) *base.Info
	RemoveVideo(video *Video) *base.Info
	GetVideos(page base.Page) []Video
	//数据迁移
	Move(video *Video, storeService StoreService) *Video
}

func (video *Video) PushVideo(v *Video, content []byte) *base.Info {
	//优先保存文件
	id, err := video.StoreService.Save("video/", &content)
	if err != nil {
		return base.ServerError
	}
	v.Url = id
	v.StoreType = video.StoreService.GetType()
	c := crc32.NewIEEE()
	hash := string(c.Sum(content))
	v.Hash = hash
	now := time.Now()
	err = video.UserVideoPersist.Save(&entity.UserVideo{UserId: video.UserId, Url: video.Url, StoreType: video.StoreType, Hash: v.Hash, CreateTime: &now})
	if err != nil {
		return base.ServerError
	}
	return base.Success
}
func (video *Video) RemoveVideo(v *Video) *base.Info {
	err := video.UserVideoPersist.Delete(&entity.UserVideo{Id: v.Id})
	if err != nil {
		return base.ServerError
	}
	return base.Success
}
func (video *Video) GetVideos(page base.Page) []Video {
	return nil
}

func (video *Video) Move(v *Video, storeService StoreService) *Video {
	return nil
}
