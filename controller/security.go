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
	Service service.Security
}

type PassInput struct {
	Token           string `form:"token" valid:"Required"`
	Mobile          string `form:"mobile" valid:"Required;Mobile"`
	Password        string `form:"password" valid:"Required;MinSize(6);MaxSize(20)"`
	ConfirmPassword string `form:"confirm_password" valid:"Required;MinSize(6);MaxSize(20)"`
	SmsCode         string `form:"smsCode" valid:"Required;Length(6)"`
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

//登录
func (token *Token) Get() {
	info := base.Success
	defer func() {
		token.Data["json"] = info
		token.ServeJSON()
	}()
	passInput := &PassInput{Mobile: token.GetString("mobile"), Password: token.GetString("password")}
	info = passInput.validateLoginPassword()
	if !info.IsSuccess() {
		return
	}
	info = token.Service.Login(passInput.Mobile, passInput.Password)
	if !info.IsSuccess() {
		return
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
	info = sms.Security.SendSmsCode(in.Mobile)
}


