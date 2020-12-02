package sockets

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"rpiWater/config"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnS struct {
	connects *websocket.Conn
	sync.Mutex
}

var SocketConnect *ConnS // socket 连接对象

var addr = flag.String("addr", config.Cfg.SOCKETAddr, "http service address")
var upgrader = websocket.Upgrader{} // use default options

var clientF, _ = ioutil.ReadFile("client.html")
var homeTemplate = template.Must(template.New("").Parse(string(clientF)))

func socketRouter(w http.ResponseWriter, r *http.Request) {

	if SocketConnect != nil {
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func() {
		c.Close()
		SocketConnect = nil
	}()

	SocketConnect = &ConnS{connects: c}
	handler()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/socketRouter")
}

func Run() {

	go MsgChannelSender()

	flag.Parse()
	// log.SetFlags(0)
	http.HandleFunc("/socketRouter", socketRouter)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
