package sockets

import (
	"fmt"
	"rpiWater/public"

	"github.com/gorilla/websocket"
)

func MsgChannelSender() {
	for {
		select {
		case msg := <-public.Chan_sendMsg:
			fmt.Println("websocket.TextMessage =>", websocket.TextMessage)
			if SocketConnect != nil {
				SocketConnect.Lock()
				err := SocketConnect.connects.WriteMessage(websocket.TextMessage, []byte(msg+"\r\n"))
				SocketConnect.Unlock()
				if err != nil {
					fmt.Println("send process msg err =>", err)
				}
			}
		}
	}
}
