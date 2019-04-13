package service

import (
	"heart/entity"
	"heart/service/common"
	"time"
	"fmt"
	"github.com/astaxie/beego/logs"
)

type Message struct {
	entity.Message
}
type MessageService interface {
	GetMessage(token *Token, id int64) *base.Info
	SendTxtMessage(token *Token, message *Message) *base.Info
	SendBinaryMessage(token *Token, message *Message, content *[]byte, suffix string) *base.Info
	GetMessageAttach(token *Token, messageId int64) (*base.Info, *[]byte, string)
}

type SimpleMessageService struct {
	MessagePersist entity.MessagePersist
	UserPersist    entity.UserPersist
	StoreService   StoreService
}

func (service *SimpleMessageService) GetMessageAttach(token *Token, messageId int64) (*base.Info, *[]byte, string) {
	u, err := service.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if u == nil {
		logs.Error(err)
		return base.NoUserFound, nil, ""
	}
	m, err := service.MessagePersist.GetMessage(token.UserId, messageId)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	if m.Url == nil {
		logs.Error(err)
		return base.MessageAttachNotFound, nil, ""
	}
	b, name, err := service.StoreService.Get("message/send/"+fmt.Sprint(m.FromUser), *m.Url)
	if err != nil {
		logs.Error(err)
		return base.ServerError, nil, ""
	}
	return base.Success, &b, name
}
func (service *SimpleMessageService) GetMessage(token *Token, id int64) *base.Info {
	u, err := service.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	m, err := service.MessagePersist.GetUnreadMessage(token.UserId, id)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if m == nil {
		return base.Success
	}
	affected, err := service.MessagePersist.SetRead(token.UserId, m.Id)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	//如果信息已经被读取，那么读取下一条
	if affected == 0 {
		return service.GetMessage(token, id)
	}
	return base.NewSuccess(m)
}

func (service *SimpleMessageService) sendMessage(token *Token, message *Message) *base.Info {
	u, err := service.UserPersist.GetById(token.UserId)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	u, err = service.UserPersist.GetById(message.ToUser)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	if u == nil {
		return base.TargetUserNotFound
	}
	now := time.Now()
	message.CreateTime = &now
	message.FromUser = token.UserId
	err = service.MessagePersist.Save(&message.Message)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	return base.Success
}
func (service *SimpleMessageService) SendTxtMessage(token *Token, message *Message) *base.Info {
	return service.sendMessage(token, message)
}
func (service *SimpleMessageService) SendBinaryMessage(token *Token, message *Message, content *[]byte, suffix string) *base.Info {
	id, err := service.StoreService.Save("message/send/"+fmt.Sprint(token.UserId), content, suffix)
	if err != nil {
		logs.Error(err)
		return base.ServerError
	}
	message.Url = &id
	return service.sendMessage(token, message)
}
