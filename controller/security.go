package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/helper"
)

type Token struct {
	beego.Controller
	service service.Security
}

func (token *Token) Get() {
	token.ServeJSON()
}

type LoginInput struct {
	Token    string `form:"token" valid:"Required"`
	Mobile   string `form:"mobile" valid:"Required;Mobile"`
	Password string `form:"password" valid:"Required;MinSize(6);MaxSize(20)"`
	SmsCode  string `form:"smsCode" valid:"Required;Length(6)"`
}

var PassWordLengthNotEnough = &service.Info{Code: "000001", Message: "密码长度必须大于6"}
var RequestDataRequired = &service.Info{Code: "000002", Message: "上送的数据为空"}
var PasswordRequired = &service.Info{Code: "000003", Message: "密码不能为空"}
var IllegalMobileFormat = &service.Info{Code: "000004", Message: "手机号非法"}

func (loginInput *LoginInput) validateLoginPassword() (*service.Info) {
	if len(loginInput.Password) == 0 {
		return PasswordRequired
	}
	return service.Success
}
func (loginInput *LoginInput) validateRegistPassword() (*service.Info) {
	if len(loginInput.Password) == 0 {
		return PasswordRequired
	}
	if len(loginInput.Password) < 6 {
		return PassWordLengthNotEnough
	}
	return service.Success
}

//登录
func (token *Token) Post() {
	defer token.ServeJSON()
	input := &LoginInput{}
	token.ParseForm(input)
	info := input.validateLoginPassword()
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
	info = token.service.Login(&service.Token{Token: input.Token}, input.Mobile, input.Password)
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
}

type SmsController struct {
	beego.Controller
}

func (sms *SmsController) Get() {
	defer sms.ServeJSON()
	mobile := sms.Ctx.Input.Param(":mobile")
	if !helper.MobileRegexp.MatchString(mobile) {
		sms.Data["json"] = IllegalMobileFormat
	}

}
