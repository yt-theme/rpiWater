package sockets

import (
	"net"
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

var TCPServer *net.TCPAddr
var Listener *net.TCPListener

// data type
type RetData struct {
	Opa string `json:"opa"`
}

// 处理器
func handler(conn *ConnS) {
    fmt.Println("handler =>")
	var chan_sendMsg_stop = make(chan int, 1) // 退出发送进程消息

	defer func() {
		fmt.Println("conn close =>")
		chan_sendMsg_stop <- 1 // 消息协程退出
        conn.Lock()
		conn.connects.Close()
        conn.Unlock()
	}()
    
	// send process's msg
	go func(conn *ConnS) {
		for {
			select {
			case msg := <-public.Chan_sendMsg:
                fmt.Println("websocket.TextMessage =>", websocket.TextMessage)
                conn.Lock()
				err := conn.connects.WriteMessage(websocket.TextMessage, []byte(msg + "\r\n"))
                conn.Unlock()
                if err != nil {
                    public.Chan_stop <- 1 // 停止操作
                    fmt.Println("send process msg err =>", err)
                    continue
                }
				break
			case <-chan_sendMsg_stop:
				goto end1
				break
			}
		}

    	end1:
    		return
	}(conn)

	// socket 通信
	for {
		fmt.Println("======================>")
		/*
		   TODO: socket通信规则
		*/
		// 接收
		_, message, err := conn.connects.ReadMessage()
		if err != nil {
            chan_sendMsg_stop <- 1
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
			conn.connects.WriteMessage(websocket.TextMessage, []byte("json err " + err.Error() + "\r\n"))
            conn.Unlock()
			continue
		}

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

	}

	// 出现连接断开时返回
end0:
	public.Chan_stop <- 1
	return
}
