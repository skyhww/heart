package common

import (
	"heart/service/common"
	"github.com/astaxie/beego/logs"
)

type MessageRequest struct {
	Message *string `json:"message"`
	ToUser  int64   `json:"to_user"`
	Url     *string `db:"url" json:"url"`
}

func (messageRequest *MessageRequest) ValidateMessage() *base.Info {
	if messageRequest.ToUser == 0 {
		logs.Warn("接收人为空！")
		return IllegalRequestDataFormat
	}
	return base.Success
}
