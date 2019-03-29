package service

import (
	"heart/entity"
	"heart/service/common"
)

type Video struct {
	entity.UserVideo
	StoreService StoreService
}

type VideoService interface {
	PushVideo(video *Video, content []byte) *base.Info
	RemoveVideo(video *Video) *base.Info
	GetVideos(page base.Page) []Video
}

func (video *Video) PushVideo(v *Video) *base.Info {
	//优先保存文件

}
func (video *Video) RemoveVideo(v *Video) *base.Info {

}
func (video *Video) GetVideos(page base.Page) []Video {

}
