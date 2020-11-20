package sockets

import (
	"net"
	"rpiWater/config"
	"rpiWater/public"

	// "github.com/golang/protobuf/proto"
	"encoding/json"
	"fmt"
	"strings"

	// "github.com/gorilla/websocket"

	// "strconv"
	"bytes"
	"log"
)

var cfg = config.Cfg

var TCPServer *net.TCPAddr
var Listener *net.TCPListener

// data type
type RetData struct {
	Opa string `json:"opa"`
}

// ###############################################################################
// 初始化
func Run() {

	TCPServer, err := net.ResolveTCPAddr("tcp4", cfg.SOCKETAddr)
	if err != nil {
		fmt.Println("socket server err 1 =>", err.Error())
	}

	Listener, err = net.ListenTCP("tcp", TCPServer)
	fmt.Println("socket server running")

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println("socket server err 2 =>", err.Error())
			continue
		}
		go handler(conn)
	}
}

// ###############################################################################
// 处理器
func handler(conn net.Conn) {

	var chan_sendMsg_stop = make(chan int, 1) // 退出发送进程消息

	defer func() {
		fmt.Println("conn close =>")
		chan_sendMsg_stop <- 1 // 消息协程退出
		conn.Close()
	}()
	fmt.Println("handle ===>")

	// send process's msg
	go func(conn net.Conn) {
		for {
			select {
			case msg := <-public.Chan_sendMsg:
				conn.Write([]byte(msg + "\r\n"))
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

		buffer := make([]byte, 512)

		// 接收
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("socket read buffer err =>", err.Error())
			conn.Write([]byte("read buffer err\r\n"))
			goto end0
		}
		fmt.Println("n & buffer =>", n, string(buffer))
		// size
		var splitStr = strings.Split(strings.Trim(string(buffer), " "), ":")
		// is size data
		if splitStr[0] == "size" {
			// is data
		} else {
			// ----------------------------------------------------------------
			// handle ret
			var retData RetData
			err = json.Unmarshal(buffer[:bytes.IndexByte(buffer, 0)], &retData)
			if err != nil {
				conn.Write([]byte("json err " + err.Error() + "\r\n"))
				continue
			}

			/*
			   start || stop
			*/
			if retData.Opa == "start" {
				public.Chan_start <- 1
				conn.Write([]byte("received =>\r\n"))

			} else if retData.Opa == "stop" {
				public.Chan_stop <- 1 // 停止操作
				conn.Write([]byte("received =>\r\n"))
			}

			// ----------------------------------------------------------------
		}

	}

	// 出现连接断开时返回
end0:
	public.Chan_stop <- 1
	return
}
