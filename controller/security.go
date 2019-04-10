package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/helper"
	"heart/controller/common"
	"heart/service/common"
)

type Token struct {
	beego.Controller
	Service     service.Security
	TokenHolder *common.TokenHolder
}

type PassInput struct {
	Token           string `json:"token"`
	Mobile          string `json:"mobile"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	SmsCode         string `json:"sms_code"`
}

func (loginInput *PassInput) validateLoginPassword() (*base.Info) {
	if len(loginInput.Password) == 0 {
		return common.PasswordRequired
	}
	return base.Success
}
func (loginInput *PassInput) validateMobile() (*base.Info) {
	if len(loginInput.Mobile) == 0 {
		return common.MobileRequired
	}
	if !helper.MobileRegexp.MatchString(loginInput.Mobile) {
		return common.IllegalMobileFormat
	}
	return base.Success
}

func (loginInput *PassInput) validateSmsCode() (*base.Info) {
	if len(loginInput.SmsCode) == 0 {
		return common.SmsCodeRequired
	}
	if !helper.SmsCodeRegexp.MatchString(loginInput.SmsCode) {
		return common.IllegalSmsCodeFormat
	}
	return base.Success
}

func (loginInput *PassInput) validateRegistPassword() (*base.Info) {
	if len(loginInput.Password) == 0 {
		return common.PasswordRequired
	}
	if len(loginInput.Password) < 6 {
		return common.PassWordLengthNotEnough
	}
	if loginInput.Password != loginInput.ConfirmPassword {
		return common.ConfirmPasswordNotMatched
	}
	return base.Success
}

func (token *Token) Get() {
	info := base.Success
	defer func() {
		token.Data["json"] = info
		token.ServeJSON()
	}()
	passInput := &PassInput{Mobile: token.GetString("mobile"), Password: token.GetString("password")}
	info = passInput.validateMobile()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateLoginPassword()
	if !info.IsSuccess() {
		return
	}
	info = token.Service.Login(passInput.Mobile, passInput.Password)
	return
}

//登出
func (token *Token) Delete() {
	info := base.Success
	defer func() {
		token.Data["json"] = info
		token.ServeJSON()
	}()
	t, info := token.TokenHolder.GetToken(&token.Controller)
	if !info.IsSuccess() {
		return
	}
	if t != nil {
		info = token.Service.LogOut(t)
	}
}

type SmsController struct {
	beego.Controller
	Security service.Security
}

func (sms *SmsController) Get() {
	info := base.Success
	defer func() {
		sms.Data["json"] = info
		sms.ServeJSON()
	}()
	in := &PassInput{Mobile: sms.Ctx.Input.Param(":mobile")}
	info = in.validateMobile()
	if !info.IsSuccess() {
		return
	}

	str := sms.GetString("use")
	if str == "regist" {
		info = sms.Security.SendRegistCode(in.Mobile)
	} else if str == "reset_password" {
		info = sms.Security.SendRestPasswordCode(in.Mobile)
	}
}
