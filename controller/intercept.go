package controller

import (
	"github.com/astaxie/beego"
	"heart/controller/common"
)

type NN struct {
	beego.Controller
}

func (test *NN) Post() {
	test.ServeJSON()
	common.GetUser(&test.Controller)
}
func (test *NN) Get() {
	test.ServeJSON()
	common.GetUser(&test.Controller)
}
func (test *NN) Put() {
	test.ServeJSON()
	common.GetUser(&test.Controller)
}
