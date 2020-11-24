package sockets

import (
    "rpiWater/config"
    "flag"
    "html/template"
    "log"
    "net/http"
    "io/ioutil"
    "sync"
    "github.com/gorilla/websocket"
)

type ConnS struct {
    connects *websocket.Conn
    sync.Mutex
}

var addr = flag.String("addr", config.Cfg.SOCKETAddr, "http service address")
var upgrader = websocket.Upgrader{} // use default options

var clientF, _ = ioutil.ReadFile("client.html")
var homeTemplate = template.Must(template.New("").Parse(string(clientF)))

func socketRouter(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()
        
    connObj := &ConnS{connects: c}
    handler(connObj)
}

func home(w http.ResponseWriter, r *http.Request) {
    homeTemplate.Execute(w, "ws://"+r.Host+"/socketRouter")
}

func Run() {
    flag.Parse()
    // log.SetFlags(0)
    http.HandleFunc("/socketRouter", socketRouter)
    http.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(*addr, nil))
}

