package controller

import (
	"github.com/astaxie/beego"
	"heart/controller/common"
	"heart/service/common"
	"heart/service"
)

type UserCollectorController struct {
	beego.Controller
	TokenHolder      *common.TokenHolder
	CollectorService service.CollectorService
}

func (userCollectorController *UserCollectorController) Put() {
	info := base.Success
	defer func() {
		userCollectorController.Data["json"] = info
		userCollectorController.ServeJSON()
	}()
	postsId, err := userCollectorController.GetInt64(":posts_id", -1)
	if err != nil || postsId == -1 {
		info = common.IllegalRequestDataFormat
		return
	}
	t, info := userCollectorController.TokenHolder.GetToken(&userCollectorController.Controller)
	if !info.IsSuccess() {
		return
	}
	info = userCollectorController.CollectorService.Add(t, postsId)
}
func (userCollectorController *UserCollectorController) Delete() {
	info := base.Success
	defer func() {
		userCollectorController.Data["json"] = info
		userCollectorController.ServeJSON()
	}()
	postsId, err := userCollectorController.GetInt64("posts_id", -1)
	if err != nil || postsId == -1 {
		info = common.IllegalRequestDataFormat
		return
	}
	t, info := userCollectorController.TokenHolder.GetToken(&userCollectorController.Controller)
	if !info.IsSuccess() {
		return
	}
	info = userCollectorController.CollectorService.Remove(t, postsId)
}
func (userCollectorController *UserCollectorController) Get() {
	info := base.Success
	defer func() {
		userCollectorController.Data["json"] = info
		userCollectorController.ServeJSON()
	}()
	t, info := userCollectorController.TokenHolder.GetToken(&userCollectorController.Controller)
	if !info.IsSuccess() {
		return
	}
	size, err := userCollectorController.GetInt("page_size", 5)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	no, err := userCollectorController.GetInt("page_no", 1)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	p := &base.Page{PageSize: size, PageNo: no}
	info = userCollectorController.CollectorService.Get(t, p)
}
