package common

import (
	"heart/service"
	"github.com/astaxie/beego"
	"fmt"
)

func  GetUser(controller *beego.Controller) *service.User{
	token:=controller.GetString("token")
	fmt.Print(controller.Ctx.Request.Body)
	fmt.Print("1111")
	fmt.Print(token)
	return nil
}
