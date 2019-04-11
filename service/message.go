package service

import (
	"heart/entity"
	"heart/service/common"
	"time"
)

type Message struct {
	entity.Message
}
type MessageService interface {
	GetMessage(token *Token, id int64) *base.Info
	SendMessage(token *Token, message *Message) *base.Info
}

type SimpleMessageService struct {
	MessagePersist entity.MessagePersist
	UserPersist    entity.UserPersist
}

func (service *SimpleMessageService) GetMessage(token *Token, id int64) *base.Info {
	u, err := service.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	m, err := service.MessagePersist.GetUnreadMessage(token.UserId, id)
	if err != nil {
		return base.ServerError
	}
	if m == nil {
		return base.Success
	}
	affected, err := service.MessagePersist.SetRead(token.UserId, m.Id)
	if err != nil {
		return base.ServerError
	}
	//如果信息已经被读取，那么读取下一条
	if affected == 0 {
		return service.GetMessage(token, id)
	}
	return base.NewSuccess(m)
}

func (service *SimpleMessageService) SendMessage(token *Token, message *Message) *base.Info {
	u, err := service.UserPersist.GetById(token.UserId)
	if err != nil {
		return base.ServerError
	}
	if u == nil {
		return base.NoUserFound
	}
	u, err = service.UserPersist.GetById(message.ToUser)
	if err != nil {
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
		return base.ServerError
	}
	return base.Success
}
