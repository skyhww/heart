package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"heart/controller/common"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
	"path"
)

type VideoController struct {
	beego.Controller
	VideoService service.VideoService
	TokenHolder  *common.TokenHolder
	Limit        int64
}

func (videoController *VideoController) Put() {
	info := base.Success
	defer func() {
		videoController.Data["json"] = info
		videoController.ServeJSON()
		videoController.Ctx.Request.MultipartForm.RemoveAll()
	}()
	t, info := videoController.TokenHolder.GetToken(&videoController.Controller)
	if !info.IsSuccess() {
		return
	}
	content := videoController.GetString("content")
	f, h, err := videoController.GetFile("attach")
	if err != nil {
		logs.Error(err)
		info = common.FileUploadFailed
		return
	}
	defer f.Close()
	//字节
	if h.Size > (videoController.Limit << 20) {
		logs.Warn("size %d",h.Size)
		info = common.FileSizeUnbound
		return
	}
	ext := path.Ext(h.Filename)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logs.Error(err)
		info = common.FileUploadFailed
		return
	}
	if b == nil || len(b) == 0 {
		info = common.FileRequired
		return
	}
	info = videoController.VideoService.PushVideo(t, &b, ext, content)
}

//删除视频
func (videoController *VideoController) Delete() {
	info := base.Success
	defer func() {
		videoController.Data["json"] = info
		videoController.ServeJSON()

	}()
	t, info := videoController.TokenHolder.GetToken(&videoController.Controller)
	if !info.IsSuccess() {
		return
	}
	id, err := videoController.GetInt64("id", -1)
	if err != nil {
		info = common.IllegalRequestDataFormat
		return
	}
	info = videoController.VideoService.RemoveVideo(t, id)
}

func (videoController *VideoController) Get() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			videoController.Data["json"] = info
			videoController.ServeJSON()
		} else {
			videoController.Ctx.ResponseWriter.Flush()
		}
	}()
	id, err := videoController.GetInt64("id")
	if err != nil || id == 0 {
		info = common.IllegalRequest
		return
	}
	t, info := videoController.TokenHolder.GetToken(&videoController.Controller)
	if !info.IsSuccess() {
		return
	}
	info, b, name := videoController.VideoService.GetVideo(t, id)
	if !info.IsSuccess() {
		return
	}
	output := videoController.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "31536000")
	output.Header("Cache-Control", "public")
	output.Header("Pragma", "public")
	videoController.Ctx.ResponseWriter.Write(*b)
}

//视频检索
func (videoController *VideoController) Search() {
	info := base.Success
	defer func() {
		videoController.Data["json"] = info
		videoController.ServeJSON()
	}()
	t, info := videoController.TokenHolder.GetToken(&videoController.Controller)
	if !info.IsSuccess() {
		return
	}
	content := videoController.GetString("content")
	pageSize, err := videoController.GetInt("page_size", 5)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	pageNo, err := videoController.GetInt("page_no", 1)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	p := &base.Page{PageNo: pageNo, PageSize: pageSize}
	info = videoController.VideoService.SearchVideo(t, content, p)
}
