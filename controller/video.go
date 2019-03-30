package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"heart/controller/common"
	"io/ioutil"
)

type VideoController struct {
	beego.Controller
	VideoService service.VideoService
}

func (videoController *VideoController) Put() {
	info := base.Success
	defer func() {
		videoController.Data["json"] = info
		videoController.ServeJSON()
	}()
	f, _, err := videoController.GetFile("video")
	if err != nil {
		info = common.UploadFailed
		return
	}
	defer f.Close()
	_,err= ioutil.ReadAll(f)
	if err != nil {
		info = common.UploadFailed
		return
	}
	//video:= &service.Video{}

	videoController.GetString("token")
}
