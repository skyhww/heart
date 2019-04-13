package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"heart/controller/common"
	"heart/entity"
	"io/ioutil"
	"path"
	"time"
)

//  user/posts
type UserPostsController struct {
	beego.Controller
	UserPostService service.UserPostService
	TokenHolder     *common.TokenHolder
	StoreService    service.StoreService
	Limit           int64
	MaxAttach       int
}

func (postsController *UserPostsController) Put() {
	info := base.Success
	defer func() {
		postsController.Data["json"] = info
		postsController.ServeJSON()
		postsController.Ctx.Request.MultipartForm.RemoveAll()
	}()
	t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}

	fs, err := postsController.GetFiles("attach")
	if err != nil {
		info = common.FileUploadFailed
		return
	}
	content := postsController.GetString("content")
	posts := &service.Post{}
	now := time.Now()
	posts.Content = &content
	len := len(fs)
	if len > postsController.MaxAttach {
		info = common.AttachTooMuch
		return
	}
	if len != 0 {
		attach := make([]entity.PostAttach, len)
		posts.PostAttach = &attach
		for i, v := range fs {
			if v.Size > (postsController.Limit << 20) {
				info = common.FileSizeUnbound
				return
			}
			f, err := v.Open()
			if err != nil {
				info = common.FileUploadFailed
				return
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				info = common.FileUploadFailed
				return
			}
			ext := path.Ext(v.Filename)
			id, err := postsController.StoreService.Save("post/attach", &b, ext)
			if err != nil {
				info = common.FileUploadFailed
				return
			}
			attach[i].Url = &id
			attach[i].No = i
			attach[i].CreateTime = &now
			attach[i].Enable = 1
		}
	} else if content == "" {
		info = common.RequestDataRequired
		return
	}
	info = postsController.UserPostService.PutPosts(t, posts)
}

func (postsController *UserPostsController) Delete() {
	info := base.Success
	defer func() {
		postsController.Data["json"] = info
		postsController.ServeJSON()
	}()
	id, err := postsController.GetInt64("id", -1)
	if err != nil || id == -1 {
		info = common.RequestDataRequired
		return
	}
	t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}
	info = postsController.UserPostService.DeletePosts(t, id)
}
func (postsController *UserPostsController) Get() {
	info := base.Success
	defer func() {
		postsController.Data["json"] = info
		postsController.ServeJSON()
	}()
	t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}
	pageSize, err := postsController.GetInt("page_size", 5)
	if err != nil {
		info = common.RequestDataRequired
		return
	}
	pageNo, err := postsController.GetInt("page_no", 1)
	if err != nil {
		info = common.RequestDataRequired
		return
	}
	page := &base.Page{PageSize: pageSize, PageNo: pageNo}
	info = postsController.UserPostService.GetPosts(t, page)
}

type CommentController struct {
	beego.Controller
	PostService service.UserPostService
}

func (commentController *CommentController) Put() {

}

func (commentController *CommentController) Delete() {

}
func (commentController *CommentController) Get() {

}

type PostAttachController struct {
	beego.Controller
	TokenHolder       *common.TokenHolder
	PostAttachService service.PostAttachService
}

func (postAttachController *PostAttachController) Get() {
	info := base.Success
	defer func() {
		postAttachController.Data["json"] = info
		postAttachController.ServeJSON()
	}()
	id, err := postAttachController.GetInt64("id", -1)
	if err != nil || id == -1 {
		info = common.RequestDataRequired
		return
	}
	t, info := postAttachController.TokenHolder.GetToken(&postAttachController.Controller)
	if !info.IsSuccess() {
		return
	}
	info, b, name := postAttachController.PostAttachService.GetAttach(t, id)
	if !info.IsSuccess() {
		return
	}
	output := postAttachController.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "must-revalidate")
	output.Header("Pragma", "public")
	postAttachController.Ctx.ResponseWriter.Write(*b)
}

type PostsController struct {
	beego.Controller
	PostService service.PostService
	TokenHolder *common.TokenHolder
}

func (postsController *PostsController) Get() {
	info := base.Success
	defer func() {
		postsController.Data["json"] = info
		postsController.ServeJSON()
	}()
	t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}
	pageSize, err := postsController.GetInt("page_size", 5)
	if err != nil {
		info = common.RequestDataRequired
		return
	}
	pageNo, err := postsController.GetInt("page_no", 1)
	if err != nil {
		info = common.RequestDataRequired
		return
	}
	page := &base.Page{PageSize: pageSize, PageNo: pageNo}
	info = postsController.PostService.GetPosts(t, page)
}
