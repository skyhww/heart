package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
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

var PassWordLengthNotEnough = &service.Info{"000001", "密码长度必须大于6"}
var RequestDataRequired = &service.Info{"000002", "上送的数据为空"}
var PasswordRequired = &service.Info{"000003", "密码不能为空"}

func (loginInput *LoginInput) validatePassword() (*service.Info) {
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
	info := input.validatePassword()
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
	info = token.service.Regist(&service.Token{Token: input.Token}, input.Mobile, input.Password, input.SmsCode)
	if !info.IsSuccess() {
		token.Data["json"] = info
		return
	}
}
