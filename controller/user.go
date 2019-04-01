package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"encoding/json"
	"heart/controller/common"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
)

type User struct {
	beego.Controller
	Service      service.Security
	TokenHolder *common.TokenHolder
}

//注册
func (user *User) Put() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	passInput := &PassInput{}
	b,err:=ioutil.ReadAll(user.Ctx.Request.Body)
	if err!=nil{
		logs.Error(err)
		info=common.IllegalRequest
		return
	}
	if err:= json.Unmarshal(b, &passInput); err != nil {
		logs.Error(err)
		info=common.IllegalRequestDataFormat
		return
	}
	info = passInput.validateMobile()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateSmsCode()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateRegistPassword()
	if !info.IsSuccess() {
		return
	}
	info = user.Service.Regist(passInput.Mobile, passInput.Password, passInput.SmsCode)
}
