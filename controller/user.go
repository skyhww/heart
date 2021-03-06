package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"encoding/json"
	"heart/controller/common"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"path"
)

type User struct {
	beego.Controller
	Service         service.Security
	TokenHolder     *common.TokenHolder
	UserInfoService service.UserInfoService
}

type UserInfoInput struct {
	//用户名
	Name *string `json:"name"`
	//头像
	Icon *string `json:"icon"`
	//签名
	Signature *string `json:"signature"`
}

func (userInfoInput *UserInfoInput) ValidateName() *base.Info {
	if userInfoInput.Name == nil {
		return common.UserNameRequired
	}
	return base.Success
}

//验证签名长度
func (userInfoInput *UserInfoInput) ValidateSignature() *base.Info {
	return base.Success
}

//注册
func (user *User) Put() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	passInput := &PassInput{}
	b, err := ioutil.ReadAll(user.Ctx.Request.Body)
	if err != nil {
		logs.Error(err)
		info = common.IllegalRequest
		return
	}
	if err := json.Unmarshal(b, &passInput); err != nil {
		logs.Error(err)
		info = common.IllegalRequestDataFormat
		return
	}
	info = passInput.validateMobile()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateSmsCode()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateRegistPassword()
	if !info.IsSuccess() {
		return
	}
	info = user.Service.Regist(passInput.Mobile, passInput.Password, passInput.SmsCode)
}

//修改密码
func (user *User) Post() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	passInput := &PassInput{}
	b, err := ioutil.ReadAll(user.Ctx.Request.Body)
	if err != nil {
		logs.Error(err)
		info = common.IllegalRequest
		return
	}
	if err := json.Unmarshal(b, &passInput); err != nil {
		logs.Error(err)
		info = common.IllegalRequestDataFormat
		return
	}
	info = passInput.validateMobile()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateSmsCode()
	if !info.IsSuccess() {
		return
	}
	info = passInput.validateRegistPassword()
	if !info.IsSuccess() {
		return
	}
	//info = user.Service(passInput.Mobile, passInput.Password, passInput.SmsCode)
}
//获取用户信息
func (user *User) Get() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	id, err := user.GetInt64("id")
	if err != nil {
		info = common.IllegalRequest
		return
	}
	info = user.UserInfoService.GetUserInfo(id)
}

type Icon struct {
	beego.Controller
	TokenHolder *common.TokenHolder
	UserInfo    *service.UserInfo
	Limit       int64
}

func (user *Icon) Get() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			user.Data["json"] = info
			user.ServeJSON()
		}else{
			user.Ctx.ResponseWriter.Flush()
		}
	}()

	id, err := user.GetInt64(":id", -1)
	if err != nil {
		info = common.IllegalRequest
		return
	}
	if id == -1 {
		t, _ := user.TokenHolder.GetToken(&user.Controller)
		if t == nil {
			info = common.IllegalRequest
			return
		}
		id = t.UserId
	}
	info, b, name := user.UserInfo.ReadIcon(id)
	if !info.IsSuccess() {
		return
	}

	output := user.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "2300000")
	output.Header("Cache-Control", "public")
	output.Header("Pragma", "public")
	user.Ctx.ResponseWriter.Write(*b)
}

//上传头像
func (user *Icon) Post() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
		user.Ctx.Request.MultipartForm.RemoveAll()
	}()
	t, info := user.TokenHolder.GetToken(&user.Controller)
	if !info.IsSuccess() {
		return
	}
	f, h, err := user.GetFile("icon")
	if err != nil {
		logs.Error(err)
		info = common.FileUploadFailed
		return
	}
	defer f.Close()
	//字节
	if h.Size > (user.Limit << 20) {
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

	info = user.UserInfo.UpdateIcon(t, &b, ext)
}

type Signature struct {
	beego.Controller
	TokenHolder *common.TokenHolder
	UserInfo    *service.UserInfo
}

//修改签名
func (user *Signature) Post() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	in := &UserInfoInput{}
	info = user.TokenHolder.ReadJsonBody(&user.Controller, in)
	if !info.IsSuccess() {
		return
	}
	t, info := user.TokenHolder.GetToken(&user.Controller)
	if !info.IsSuccess() {
		return
	}
	info = user.UserInfo.UpdateSignature(t, in.Signature)
}

type Name struct {
	beego.Controller
	TokenHolder *common.TokenHolder
	UserInfo    *service.UserInfo
}

//修改用户名
func (user *Name) Post() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	in := &UserInfoInput{}
	info = user.TokenHolder.ReadJsonBody(&user.Controller, in)
	if !info.IsSuccess() {
		return
	}
	info = in.ValidateName()
	if !info.IsSuccess() {
		return
	}
	t, info := user.TokenHolder.GetToken(&user.Controller)
	if !info.IsSuccess() {
		return
	}
	info = user.UserInfo.UpdateName(t, in.Name)
}
