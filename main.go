package main

import (
	"github.com/astaxie/beego"
	"heart/controller"
)

func main() {
	beego.Router("/token", &controller.Token{})
	beego.Run()
}

