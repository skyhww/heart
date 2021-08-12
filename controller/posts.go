package controller

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"heart/controller/common"
	"heart/entity"
	"heart/service"
	"heart/service/common"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strconv"
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
		if postsController.Ctx.Request.MultipartForm != nil {
			postsController.Ctx.Request.MultipartForm.RemoveAll()
		}
	}()
	t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}
	var fs []*multipart.FileHeader
	var err error
	if postsController.Ctx.Request.MultipartForm != nil {
		fs, err = postsController.GetFiles("attach")
		if err != nil {
			info = common.FileUploadFailed
			return
		}
	}

	content := postsController.GetString("content", "")
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
	id, err := postsController.GetInt64(":id", -1)
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
	TokenHolder *common.TokenHolder
}

//发评论
func (commentController *CommentController) Put() {
	info := base.Success
	defer func() {
		commentController.Data["json"] = info
		commentController.ServeJSON()
	}()
	t, info := commentController.TokenHolder.GetToken(&commentController.Controller)
	if !info.IsSuccess() {
		return
	}
	id := commentController.Ctx.Input.Param(":post_id")
	if id == "" {
		info = common.IllegalRequest
		return
	}
	replayId := commentController.Ctx.Input.Param("id")
	msg := common.MessageRequest{}
	info = commentController.TokenHolder.ReadJsonBody(&commentController.Controller, &msg)
	if !info.IsSuccess() {
		return
	}
	if msg.Message == nil {
		info = common.MessageRequired
		return
	}
	c := service.Comment{}
	c.Content = *msg.Message

	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		info = common.IllegalRequestDataFormat
		return
	}
	c.PostId = postId
	if replayId != "" {
		c.ReplyId, err = strconv.ParseInt(replayId, 10, 64)
	}
	info = commentController.PostService.AddComment(t, &c)
}

func (commentController *CommentController) Replay() {
	info := base.Success
	defer func() {
		commentController.Data["json"] = info
		commentController.ServeJSON()
	}()
	t, info := commentController.TokenHolder.GetToken(&commentController.Controller)
	if !info.IsSuccess() {
		return
	}
	id, err := commentController.GetInt64(":id", -1)
	if err != nil || id == -1 {
		info = common.IllegalRequest
		return
	}
	m := &common.MessageRequest{}
	info = commentController.TokenHolder.ReadJsonBody(&commentController.Controller, m)
	if !info.IsSuccess() {
		return
	}
	if m.Message == nil || (*m.Message) == "" {
		info = common.MessageRequired
		return
	}
	c := service.Comment{}
	c.ReplyId = id
	c.Content = *m.Message
	info = commentController.PostService.AddReplay(t, &c)
}

func (commentController *CommentController) Delete() {
	info := base.Success
	defer func() {
		commentController.Data["json"] = info
		commentController.ServeJSON()
	}()
	t, info := commentController.TokenHolder.GetToken(&commentController.Controller)
	if !info.IsSuccess() {
		return
	}
	id, err := commentController.GetInt64("id", -1)
	if err != nil || id == -1 {
		info = common.IllegalRequest
		return
	}
	info = commentController.PostService.DeleteComment(t, id)
}
func (commentController *CommentController) Get() {
	info := base.Success
	defer func() {
		commentController.Data["json"] = info
		commentController.ServeJSON()
	}()
	id, err := commentController.GetInt64("post_id", -1)
	if err != nil || id == -1 {
		info = common.IllegalRequest
		return
	}
	info = commentController.PostService.GetComments(id)
}

type PostAttachController struct {
	beego.Controller
	TokenHolder       *common.TokenHolder
	PostAttachService service.PostAttachService
}

func (postAttachController *PostAttachController) Get() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			postAttachController.Data["json"] = info
			postAttachController.ServeJSON()
		} else {
			postAttachController.Ctx.ResponseWriter.Flush()
		}
	}()
	id, err := postAttachController.GetInt64("id", -1)
	if err != nil || id == -1 {
		info = common.RequestDataRequired
		return
	}
	/*t, info := postAttachController.TokenHolder.GetToken(&postAttachController.Controller)
	if !info.IsSuccess() {
		return
	}*/
	info, b, name := postAttachController.PostAttachService.GetAttach(nil, id)
	if !info.IsSuccess() {
		return
	}
	output := postAttachController.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	//帖子附件一般不会变更，可缓存
	output.Header("Expires", "31536000")
	output.Header("Cache-Control", "public")
	output.Header("Pragma", "public")
	postAttachController.Ctx.ResponseWriter.Write(*b)
}

func (postAttachController *PostAttachController) GetPage() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			postAttachController.Data["json"] = info
			postAttachController.ServeJSON()
		} else {
			postAttachController.Ctx.ResponseWriter.Flush()
		}
	}()
	id, err := postAttachController.GetInt64(":posts_id", -1)
	if err != nil || id == -1 {
		info = common.RequestDataRequired
		return
	}

	info, b, name := postAttachController.PostAttachService.GetAttach(nil, id)
	if !info.IsSuccess() {
		return
	}
	output := postAttachController.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	//帖子附件一般不会变更，可缓存
	output.Header("Expires", "31536000")
	output.Header("Cache-Control", "public")
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
	/*t, info := postsController.TokenHolder.GetToken(&postsController.Controller)
	if !info.IsSuccess() {
		return
	}*/
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
	keyword := postsController.GetString("keyword", "")
	page := &base.Page{PageSize: pageSize, PageNo: pageNo}
	a := postsController.PostService.GetPosts(keyword, nil, page)
	logs.Info(a)
	info = a
}
