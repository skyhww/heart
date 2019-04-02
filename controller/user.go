package controller

import (
	"github.com/astaxie/beego"
	"heart/service"
	"heart/service/common"
	"encoding/json"
	"heart/controller/common"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"fmt"
)

type User struct {
	beego.Controller
	Service     service.Security
	TokenHolder *common.TokenHolder
}

type UserInfoInput struct {
	//用户名
	Name string
	//头像
	Icon *[]byte
	//签名
	Signature *string
}

func (userInfoInput *UserInfoInput) ValidateName() *base.Info {
	if len(userInfoInput.Name) == 0 {
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

type Icon struct {
	beego.Controller
	TokenHolder *common.TokenHolder
	UserInfo    *service.UserInfo
	Limit  int64
}


//上传头像
func (user *Icon) Post() {
	info := base.Success
	defer func() {
		user.Data["json"] = info
		user.ServeJSON()
	}()
	defer user.Ctx.Request.MultipartForm.RemoveAll()
	/*err:=user.Ctx.Request.ParseMultipartForm(user.Limit)
	if err!=nil{
		logs.Error(err)
		info =common.IllegalRequestDataFormat
		return
	}

	form,err:=r.ReadForm(user.Limit)
	if err!=nil{
		logs.Error(err)
		info =common.IllegalRequestDataFormat
		return
	}
	head:=form.File["icon"]
	if  len(head)==0{
		info =common.IconRequired
		return
	}
	if len(head)>1{
		info =common.MultiIcon
		return
	}
	if head[0].Size>user.Limit{
		info =common.FileSizeUnbound
		return
	}
	f,err:=head[0].Open()
	defer f.Close()
	if err!=nil{
		logs.Error(err)
		info =common.FileUploadFailed
		return
	}*/
	name:=user.GetString("fileName")
	for e := range user.Ctx.Request.Form {
		fmt.Println(e)
	}
	f,h,err:=user.GetFile(name)
	if err!=nil{
		logs.Error(err)
		info =common.FileUploadFailed
		return
	}
	fmt.Print(h.Filename)
	b,err:= ioutil.ReadAll(f)
	if err!=nil{
		logs.Error(err)
		info =common.FileUploadFailed
		return
	}
	in := &UserInfoInput{Icon:&b}
	t, info := user.TokenHolder.GetToken(&user.Controller)
	if !info.IsSuccess() {
		return
	}
	info=user.UserInfo.UpdateSignature(t,in.Signature)
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
	info=user.UserInfo.UpdateSignature(t,in.Signature)
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
	info=user.UserInfo.UpdateName(t, &in.Name)
}
