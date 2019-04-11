package controller

import (
	"github.com/astaxie/beego"
	"heart/controller/common"
	"heart/service"
	"heart/service/common"
	"github.com/astaxie/beego/logs"
)

type MessageRequest struct {
	Message *string `json:"message"`
	ToUser  int64   `json:"to_user"`
	Url     *string `db:"url" json:"url"`
}

func (messageRequest *MessageRequest) ValidateMessage() *base.Info {
	if messageRequest.Message == nil && messageRequest.Url == nil {
		return common.MessageRequired
	}
	if messageRequest.Message != nil && messageRequest.Url != nil {
		return common.IllegalRequestDataFormat
	}
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
	m := &MessageRequest{}
	info=iMessageController.TokenHolder.ReadJsonBody(&iMessageController.Controller, m)
	if !info.IsSuccess(){
		return
	}
	info=m.ValidateMessage()
	if !info.IsSuccess(){
		return
	}
	sm:=&service.Message{}
	sm.Message.Message=m.Message
	sm.Message.Url=m.Url
	sm.Message.ToUser=m.ToUser
	info=iMessageController.MessageService.SendMessage(t,sm)
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
