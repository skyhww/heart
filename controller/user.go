package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
)

type User struct {
	beego.Controller
	Service service.Security
}

//注册
func (user *User) Put() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	passInput := &PassInput{Mobile: user.GetString("mobile"), Password: user.GetString("password"), SmsCode: user.GetString("smsCode"), ConfirmPassword: user.GetString("confirm_password")}
	info = passInput.validateRegistPassword()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateSmsCode()
	if !info.IsSuccess() {
		return
	}
	info = user.Service.Regist(passInput.Mobile, passInput.Password, passInput.SmsCode)
}
