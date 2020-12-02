package sockets

import (
	// "net"
	"rpiWater/config"
	"rpiWater/public"

	// "github.com/golang/protobuf/proto"
	"encoding/json"
	"fmt"

	// "strings"

	"github.com/gorilla/websocket"
	// "strconv"
	// "bytes"
	// "log"
)

var cfg = config.Cfg

// data type
type RetData struct {
	Opa string `json:"opa"`
	Tok string `json:"tok"`
}

// 处理器
func handler() {
	fmt.Println("handler =>")

	// conn
	conn := SocketConnect

	defer func() {
		fmt.Println("conn close =>")
		conn.Lock()
		conn.connects.Close()
		conn.Unlock()
	}()

	// socket 通信
	for {
		fmt.Println("socket ReadMsg ======================>")
		/*
		   TODO: socket通信规则
		*/
		// 接收
		_, message, err := conn.connects.ReadMessage()
		if err != nil {
			conn.Lock()
			conn.connects.WriteMessage(websocket.TextMessage, []byte("read buffer err\r\n"))
			conn.Unlock()
			goto end0
		}

		// ----------------------------------------------------------------
		// handle ret
		var retData RetData
		err = json.Unmarshal(message, &retData)
		if err != nil {
			conn.Lock()
			conn.connects.WriteMessage(websocket.TextMessage, []byte("json err "+err.Error()+"\r\n"))
			conn.Unlock()
			continue
		}

		// 判断token 才能操作
		if retData.Tok != "" && retData.Tok == config.Cfg.PublicConnectToken {
			/*
			   start || stop
			*/
			if retData.Opa == "start" {
				public.Chan_start <- 1
				conn.Lock()
				conn.connects.WriteMessage(websocket.TextMessage, []byte("received =>\r\n"))
				conn.Unlock()

			} else if retData.Opa == "stop" {
				public.Chan_stop <- 1 // 停止操作
				conn.Lock()
				conn.connects.WriteMessage(websocket.TextMessage, []byte("received =>\r\n"))
				conn.Unlock()
			}
			// ----------------------------------------------------------------
		} else {
			conn.Lock()
			conn.connects.WriteMessage(websocket.TextMessage, []byte("token err =>\r\n"))
			conn.Unlock()

		}

	}

	// 出现连接断开时返回
end0:
	public.Chan_stop <- 1
	SocketConnect = nil
	return
}
