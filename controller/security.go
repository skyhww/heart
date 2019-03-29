package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/helper"
)

type Token struct {
	beego.Controller
	Service service.Security
}

func (token *Token) Get() {
	token.ServeJSON()
}

type PassInput struct {
	Token           string `form:"token" valid:"Required"`
	Mobile          string `form:"mobile" valid:"Required;Mobile"`
	Password        string `form:"password" valid:"Required;MinSize(6);MaxSize(20)"`
	ConfirmPassword string `form:"confirm_password" valid:"Required;MinSize(6);MaxSize(20)"`
	SmsCode         string `form:"smsCode" valid:"Required;Length(6)"`
}

var PassWordLengthNotEnough = &service.Info{Code: "000001", Message: "密码长度必须大于6"}
var RequestDataRequired = &service.Info{Code: "000002", Message: "上送的数据为空"}
var PasswordRequired = &service.Info{Code: "000003", Message: "密码不能为空"}
var IllegalMobileFormat = &service.Info{Code: "000004", Message: "手机号非法"}
var ConfirmPasswordNotMatched = &service.Info{Code: "000005", Message: "确认密码不匹配"}
var SmsCodeRequired = &service.Info{Code: "000006", Message: "短信验证码不能为空"}
var IllegalSmsCodeFormat = &service.Info{Code: "000008", Message: "短信验证码格式不正确"}
var MobileRequired = &service.Info{Code: "000009", Message: "短信验证码格式不正确"}

func (loginInput *PassInput) validateLoginPassword() (*service.Info) {
	if len(loginInput.Password) == 0 {
		return PasswordRequired
	}
	return service.Success
}
func (loginInput *PassInput) validateMobile() (*service.Info) {
	if len(loginInput.Mobile) == 0 {
		return MobileRequired
	}
	if !helper.MobileRegexp.MatchString(loginInput.Mobile) {
		return IllegalMobileFormat
	}
	return service.Success
}

func (loginInput *PassInput) validateSmsCode() (*service.Info) {
	if len(loginInput.SmsCode) == 0 {
		return SmsCodeRequired
	}
	if !helper.SmsCodeRegexp.MatchString(loginInput.SmsCode) {
		return IllegalSmsCodeFormat
	}
	return service.Success
}

func (loginInput *PassInput) validateRegistPassword() (*service.Info) {
	if len(loginInput.Password) == 0 {
		return PasswordRequired
	}
	if len(loginInput.Password) < 6 {
		return PassWordLengthNotEnough
	}
	if loginInput.Password != loginInput.ConfirmPassword {
		return ConfirmPasswordNotMatched
	}
	return service.Success
}

//登录
func (token *Token) Post() {
	defer token.ServeJSON()
	input := &PassInput{}
	token.ParseForm(input)
	info := input.validateLoginPassword()
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
	info = token.Service.Login(&service.Token{Token: input.Token}, input.Mobile, input.Password)
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
}
func (token *Token) Put() {
	info := service.Success
	defer func() {
		token.Data["json"] = info
		token.ServeJSON()
	}()
	passInput := &PassInput{Mobile: token.GetString("mobile"), Password: token.GetString("password"), SmsCode: token.GetString("smsCode"), ConfirmPassword: token.GetString("confirm_password")}
	info = passInput.validateRegistPassword()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateSmsCode()
	if !info.IsSuccess() {
		return
	}
	info = token.Service.Regist(nil, passInput.Mobile, passInput.Password, passInput.SmsCode)
}

type SmsController struct {
	beego.Controller
	Security service.Security
}

func (sms *SmsController) Get() {
	info := service.Success
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
