package service

import (
	"heart/entity"
	"heart/service/common"
	"hash/crc32"
	"time"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type Video struct {
	entity.UserVideo
}
type VideoService interface {
	PushVideo(token *Token, content *[]byte, suffix, title string) *base.Info
	RemoveVideo(token *Token, id int64) *base.Info
	GetVideo(token *Token, id int64) (*base.Info, *[]byte, string)
	SearchVideo(token *Token, content string, page *base.Page) *base.Info
}
type SimpleVideoService struct {
	StoreService     StoreService
	UserVideoPersist entity.UserVideoPersist
	UserPersist      entity.UserPersist
	ElasticSearchService ElasticSearchService
}

func (video *SimpleVideoService) PushVideo(token *Token, content *[]byte, suffix, title string) *base.Info {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	id, err := video.StoreService.Save("video", content, suffix)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	//优先保存文件
	now := time.Now()
	c := crc32.NewIEEE()
	hash := string(c.Sum(*content))
	userVideo := &entity.UserVideo{UserId: u.Id, CreateTime: &now, Url: &id, Hash: &hash, Content: &title}
	err = video.UserVideoPersist.Save(userVideo)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.NewSuccess(userVideo)
}
func (video *SimpleVideoService) RemoveVideo(token *Token, id int64) *base.Info {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	err = video.UserVideoPersist.Delete(&entity.UserVideo{Id: id, UserId: token.UserId})
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}

func (video *SimpleVideoService) SearchVideo(token *Token, content string, page *base.Page) *base.Info {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound
	}
	key := make(map[string]interface{})
	key["content"] = content
	key["enable"] = true
	key["user_id"]=u.Id
	bs, err := video.ElasticSearchService.Query(key, "3dheart", "video", page)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if len(bs) == 0 {
		return base.NewSuccess(page)
	}
	data := make([]entity.UserVideo,0)
	for _, b := range bs {
		tmp := &entity.UserVideo{}
		err = json.Unmarshal(b, tmp)
		if err != nil {
			logs.Error(err)
			return base.ServerError
		}
		data=append(data, *tmp)
	}
	page.Data=data
	return base.NewSuccess(page)
}

func (video *SimpleVideoService) GetVideo(token *Token, id int64) (*base.Info, *[]byte, string) {
	u, err := video.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound, nil, ""
	}
	v := &entity.UserVideo{Id: id, UserId: u.Id}
	err = video.UserVideoPersist.Get(v)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	b, n, err := video.StoreService.Get("video", *v.Url)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if b == nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	return base.Success, &b, n
}
