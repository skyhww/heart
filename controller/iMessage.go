package controller

import (
	"github.com/astaxie/beego"
)

type IMessageController struct {
	beego.Controller
} 

func (iMessageController *IMessageController) Get(){
	/*conn, err := websocket.Upgrader{}.Upgrade(iMessageController.Ctx.ResponseWriter, iMessageController.Ctx.Request, nil)
	if err != nil {
		iMessageController.Ctx.ResponseWriter.Status=406
		iMessageController.Ctx.ResponseWriter.Flush()
		return
	}
	fmt.Print(conn)*/
}