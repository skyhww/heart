package controller

import (
	"github.com/astaxie/beego"
	"heart/controller/common"
	"heart/service"
	"heart/service/common"
	"github.com/astaxie/beego/logs"
	"path"
	"io/ioutil"
)

type MessageRequest struct {
	Message *string `json:"message"`
	ToUser  int64   `json:"to_user"`
	Url     *string `db:"url" json:"url"`
}

func (messageRequest *MessageRequest) ValidateMessage() *base.Info {
	if messageRequest.ToUser == 0 {
		logs.Warn("接收人为空！")
		return common.IllegalRequestDataFormat
	}
	return base.Success
}

type IMessageController struct {
	beego.Controller
	TokenHolder    *common.TokenHolder
	MessageService service.MessageService
	//MB
	Limit int64
}

func (iMessageController *IMessageController) Get() {
	info := base.Success
	defer func() {
		iMessageController.Data["json"] = info
		iMessageController.ServeJSON()
	}()
	t, info := iMessageController.TokenHolder.GetToken(&iMessageController.Controller)
	if !info.IsSuccess() {
		return
	}
	id, err := iMessageController.GetInt64("id", -1)
	if err != nil {
		info = common.IllegalRequestDataFormat
		return
	}
	info = iMessageController.MessageService.GetMessage(t, id)
}

//文本消息
func (iMessageController *IMessageController) putText(t *service.Token, toUser int64) *base.Info {
	info := base.Success
	m := &MessageRequest{}
	info = iMessageController.TokenHolder.ReadJsonBody(&iMessageController.Controller, m)
	if !info.IsSuccess() {
		return info
	}
	m.ToUser = toUser
	info = m.ValidateMessage()
	if m.Message == nil {
		return common.MessageRequired
	}
	sm := &service.Message{}
	sm.Message.Message = m.Message
	sm.Message.ToUser = m.ToUser
	return iMessageController.MessageService.SendTxtMessage(t, sm)
}

//流消息
func (iMessageController *IMessageController) putBinary(t *service.Token, toUser int64) *base.Info {
	info := base.Success
	f, h, err := iMessageController.GetFile("attach")
	if err != nil {
		logs.Error(err)
		return common.FileUploadFailed
	}
	defer f.Close()
	//字节
	if h.Size > (iMessageController.Limit << 20) {
		return common.FileSizeUnbound
	}
	ext := path.Ext(h.Filename)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logs.Error(err)
		return common.FileUploadFailed
	}
	if b == nil || len(b) == 0 {
		return common.FileRequired
	}
	m := &service.Message{}
	m.ToUser = toUser
	iMessageController.MessageService.SendBinaryMessage(t, m, &b, ext)
	return info
}

func (iMessageController *IMessageController) Put() {
	info := base.Success
	defer func() {
		iMessageController.Data["json"] = info
		iMessageController.ServeJSON()
	}()
	t, info := iMessageController.TokenHolder.GetToken(&iMessageController.Controller)
	if !info.IsSuccess() {
		return
	}
	ty := iMessageController.GetString("type")
	toUser, err := iMessageController.GetInt64("to_user", -1)
	if err != nil || toUser == -1 {
		info = common.IllegalRequest
		return
	}
	if ty == "text" {
		info = iMessageController.putText(t, toUser)
	} else if ty == "binary" {
		info = iMessageController.putBinary(t, toUser)
	} else {
		info = common.IllegalRequestDataFormat
	}

}

//附件消息
type IMessageAttachController struct {
	beego.Controller
	TokenHolder    *common.TokenHolder
	MessageService service.MessageService
	//MB
	Limit int64
}

func (iMessageAttachController *IMessageAttachController) Get() {
	info := base.Success
	defer func() {
		if !info.IsSuccess() {
			iMessageAttachController.Data["json"] = info
			iMessageAttachController.ServeJSON()
		} else {
			iMessageAttachController.Ctx.ResponseWriter.Flush()
		}
	}()
	id, err := iMessageAttachController.GetInt64("id")
	if err != nil || id == 0 {
		info = common.IllegalRequest
		return
	}
	t, info := iMessageAttachController.TokenHolder.GetToken(&iMessageAttachController.Controller)
	if !info.IsSuccess() {
		return
	}
	info, b, name := iMessageAttachController.MessageService.GetMessageAttach(t, id)
	if !info.IsSuccess() {
		return
	}
	output := iMessageAttachController.Ctx.Output
	output.Header("Content-Disposition", "attachment; filename="+name)
	output.Header("Content-Description", "File Transfer")
	output.Header("Content-Type", "application/octet-stream")
	output.Header("Content-Transfer-Encoding", "binary")
	output.Header("Expires", "0")
	output.Header("Cache-Control", "must-revalidate")
	output.Header("Pragma", "public")
	iMessageAttachController.Ctx.ResponseWriter.Write(*b)
}

/*func NewIMessageController() *IMessageController {
	return &IMessageController{close: make(chan struct{}, 1)}
}
*/
/*type IMessageRequest struct {
	Id string `json:"id"`
	//确认报文，确认收到了哪个消息（确认以后才能收到下一个消息），
	// -1表示由服务端决定去获取哪一个消息，一般从第一个sending处的数据重新发送
	//服务端只有收到了ack，才会把该条消息置为已读
	//消息状态： unread-->sending-->read
	Ack int64 `json:"ack"`
	//token
	Token string `json:"token"`
	// 1-发送消息 2-Ack
	Type int64 `json:"type"`
	//type=1是生效
	Message service.Message `json:"message"`
}
type IMessageResponse struct {

	RequestId string `json:"request_id"`
	//确认报文，确认收到了哪个消息（确认以后才能收到下一个消息），

	Ack int64 `json:"ack"`
	//token
	Token string `json:"token"`
	// 1-发送消息 2-Ack
	Type int64 `json:"type"`
	//type=1是生效
	Message service.Message `json:"message"`
}
type Data struct {
	ToUser  int64   `json:"to_user"`
	Message *string `json:"message"`
	Url     *string `json:"url"`
}*/

/*func (iMessageController *IMessageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logs.Error(err)
		return
	}
	defer conn.Close()
	for {
		err = nil
		m, reader, err := conn.NextReader()
		if err != nil {
			break
		}
		if m != websocket.TextMessage {
			logs.Warn("非文本消息")
			continue
		}
		if reader == nil {
			continue
		}
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			continue
		}
		data := &IMessageRequest{}
		err = json.Unmarshal(b, data)
		if err != nil {

		}
	}
}

func (iMessageController *IMessageController) read() {

}
func (iMessageController *IMessageController) write() {

}
*/
