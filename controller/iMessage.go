package controller

import (
	"github.com/gorilla/websocket"
	"time"
	"github.com/astaxie/beego/logs"
	"fmt"
	"net/http"
)

type IMessageController struct {
}

func (iMessageController *IMessageController) ServeHTTP(w http.ResponseWriter, r *http.Request){
	conn, err :=websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		logs.Error(err)
		return
	}
	for {
		time.Sleep(time.Second * 1)
		fmt.Print(conn.WriteJSON("{\"a\":1}"))
	}
}