package main

import (
    "fmt"
    "io/ioutil"
    "flag"
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/websocket"
)

type Request struct {
    Action  string
    Content []byte
}

type Response struct {
    Status  string
    Content []byte
}

func decodeResponse(msg []byte) *Response {
    ss := strings.SplitN(string(msg), "\r\n", 2)
    return &Response{Status: ss[0], Content: []byte(ss[1])}
}

func encodeRequest(req *Request) []byte {
    return []byte(fmt.Sprintf("%s\r\n%s", req.Action, req.Content))
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}

var conn *websocket.Conn
var key string

var (
    disconnect = make(chan struct{})
)

func ws(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error: websocket upgrade, %v\n", err)
        return
    }
    defer c.Close()

    key = r.Header.Get("Moon-Authentication")
    conn = c
    <- disconnect
}

func agentInfo(w http.ResponseWriter, r *http.Request) {
    req := new(Request)
    req.Action = "get-agent-info"
    if err := conn.WriteMessage(websocket.TextMessage, encodeRequest(req)); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    _, result, err := conn.ReadMessage()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(result)
}

func nodeInfo(w http.ResponseWriter, r *http.Request) {
    req := new(Request)
    req.Action = "get-node-info"
    if err := conn.WriteMessage(websocket.TextMessage, encodeRequest(req)); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    _, result, err := conn.ReadMessage()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(result)
}

func shell(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    req := new(Request)
    req.Action = "exec-shell-script"
    req.Content, _ = ioutil.ReadAll(r.Body)

    if err := conn.WriteMessage(websocket.TextMessage, encodeRequest(req)); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    
    _, result, err := conn.ReadMessage()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Write(result)
}

func getKey(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(key))
}

func done(w http.ResponseWriter, r *http.Request) {
    disconnect <- struct {}{}
}

func main() {
    flag.Parse()
    http.HandleFunc("/api/v1/agent", ws)
    http.HandleFunc("/test/agent-info", agentInfo)
    http.HandleFunc("/test/node-info", nodeInfo)
    http.HandleFunc("/test/shell", shell)
    http.HandleFunc("/key", getKey)
    http.HandleFunc("/done", done)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
