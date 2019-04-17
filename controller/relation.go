package controller

import (
	"github.com/astaxie/beego"
	"heart/controller/common"
	"heart/service"
	"heart/service/common"
)

type RelationController struct {
	beego.Controller
	TokenHolder       *common.TokenHolder
	UserFollowService service.UserFollowService
}

func (relationController *RelationController) Put() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			relationController.Data["json"] = info
			relationController.ServeJSON()
		}
	}()
	userId, err := relationController.GetInt64("user_id", -1)
	if err != nil || userId == -1 {
		info = common.IllegalRequestDataFormat
		return
	}
	t, info := relationController.TokenHolder.GetToken(&relationController.Controller)
	if !info.IsSuccess() {
		return
	}
	info = relationController.UserFollowService.Follow(t, userId)

}
func (relationController *RelationController) Delete() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			relationController.Data["json"] = info
			relationController.ServeJSON()
		}
	}()
	userId, err := relationController.GetInt64("user_id", -1)
	if err != nil || userId == -1 {
		info = common.IllegalRequestDataFormat
		return
	}
	t, info := relationController.TokenHolder.GetToken(&relationController.Controller)
	if !info.IsSuccess() {
		return
	}
	info = relationController.UserFollowService.UnFollow(t, userId)
}
func (relationController *RelationController) Get() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			relationController.Data["json"] = info
			relationController.ServeJSON()
		}
	}()
	t, info := relationController.TokenHolder.GetToken(&relationController.Controller)
	if !info.IsSuccess() {
		return
	}
	use := relationController.GetString("use", "follow")
	size, err := relationController.GetInt("page_size", 5)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	no, err := relationController.GetInt("page_no", 1)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	p := &base.Page{PageSize: size, PageNo: no}
	if use == "follow" {
		info = relationController.UserFollowService.GetFollowUser(t, p)
	} else if use == "followed" {
		info = relationController.UserFollowService.GetFollowedUser(t, p)
	} else {
		info = common.IllegalRequest
	}
}
