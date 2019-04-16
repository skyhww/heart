package common

import (
	"heart/service"
	"github.com/astaxie/beego"
	"heart/service/common"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
	"encoding/json"
)

type TokenHolder struct {
	Name     string
	Helper   *service.TokenHelper
	UserInfo *service.UserInfo
}

func (holder *TokenHolder) GetToken(controller *beego.Controller) (*service.Token, *base.Info) {
	token := controller.GetString(holder.Name)
	if token == "" {
		return nil, ReLogin
	}
	return holder.Helper.GetToken(token)
}

func (holder *TokenHolder) GetUser(controller *beego.Controller) (*service.User, *base.Info) {
	t, in := holder.GetToken(controller)
	if !in.IsSuccess() {
		return nil, in
	}
	u, err := holder.UserInfo.GetUser(t)
	if err != nil {
		return nil, base.GetUserInfoFailed
	}
	if u==nil{
		return nil,base.NoUserFound
	}
	return u, base.Success
}

func (holder *TokenHolder) ReadJsonBody(controller *beego.Controller, target interface{}) *base.Info {
	info := base.Success
	b, err := ioutil.ReadAll(controller.Ctx.Request.Body)
	if err != nil {
		logs.Error(err)
		info = IllegalRequest
		return info
	}
	if err := json.Unmarshal(b, target); err != nil {
		logs.Error(err)
		info = IllegalRequestDataFormat
		return info
	}
	return info
}

