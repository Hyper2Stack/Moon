package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"

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

var (
    ActionNodeInfo      = "get-node-info"
    ActionAgentInfo     = "get-agent-info"
    ActionExecShell     = "exec-shell-script"
    StatusOK            = "ok"
    StatusBadRequest    = "bad-request"
    StatusError         = "error"
    StatusInternalError = "internal-error"
)

var conn *websocket.Conn = nil

var (
    reconnect = make(chan struct{})
    connected = make(chan struct{})
)

func getAuth() (string, bool) {
    path := "/etc/moon/key.cfg"
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        log.Printf("Error: read file %s, %v\n", path, err)
        return "", false
    }

    conf := struct {
        Key  string `json:"key"`
        Uuid string `json:"uuid"`
    } {}

    err = json.Unmarshal(bytes, &conf)
    if err != nil {
        log.Printf("Error: parse key, %v\n", err)
        return "", false
    }

    return fmt.Sprintf("%s,%s", conf.Key, conf.Uuid), true
}

func dail(auth string) {
    u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/api/v1/agent"}
    header := make(http.Header)
    header.Set("Moon-Authentication", auth)

    for {
        log.Printf("Trying to connect %s\n", u.String())

        c, res, err := websocket.DefaultDialer.Dial(u.String(), header)
        if err != nil {
            log.Printf("Error: websocket dail, %v\n", err)
            if res != nil {
                errBody, _ := ioutil.ReadAll(res.Body)
                log.Printf("Error: websocket body, %s\n", string(errBody))
            }
            time.Sleep(3 * time.Second)
            continue
        }

        log.Printf("Successfully to create websocket connectiton to %s\n", u.String())
        conn = c
        connected <- struct{}{}
        <- reconnect
    }
}

func decodeRequest(msg []byte) *Request {
    ss := strings.SplitN(string(msg), "\r\n", 2)
    return &Request{Action: ss[0], Content: []byte(ss[1])}
}

func encodeResponse(res *Response) []byte {
    return []byte(fmt.Sprintf("%s\r\n%s", res.Status, res.Content))
}

func process() {
    <- connected
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error: read message, %v\n", err)
            conn.Close()
            reconnect <- struct{}{}
            <- connected
            continue
        }

        req := decodeRequest(msg)
        log.Printf("Recv message of action %s\n", req.Action)

        res := new(Response)
        job, err := createJob(req)
        if err != nil {
            res.Status = StatusBadRequest
            res.Content, _ = json.Marshal(map[string]string{"message": err.Error()})
            log.Printf("Error: create job, %v\n", err)
        } else {
            res = job.Run()
        }

        if err := conn.WriteMessage(websocket.TextMessage, encodeResponse(res)); err != nil {
            log.Printf("Error: write message, %v\n", err)
        }
    }
}

func worker() {
    auth, ok := getAuth()
    if !ok {
        log.Fatalf("Error: missing authentication key\n")
    }

    go dail(auth)
    go process()

    <- interrupt
    // TBD
    done <- struct{}{}
}
