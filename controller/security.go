package controller

import "github.com/astaxie/beego"

type Token struct {
	beego.Controller
}

func (token *Token) Get() {
	token.ServeJSON()
}

func (token *Token) Post() {

}
