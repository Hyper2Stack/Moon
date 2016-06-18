package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "time"

    "github.com/gorilla/websocket"
    "github.com/hyper2stack/mooncommon/protocol"
    "moon/cfg"
)

var conn *websocket.Conn = nil

var (
    reconnect = make(chan struct{})
    connected = make(chan struct{})
)

func getAuth() (string, bool) {
    conf, err := cfg.ParseKey(config.KeyFile)
    if err != nil {
        log.Printf("Error: parse key %s, %v\n", config.KeyFile, err)
        return "", false
    }

    return fmt.Sprintf("%s,%s", conf.Key, conf.Uuid), true
}

func dail(auth string) {
    scheme := "wss"
    if config.Ssl != "on" {
        scheme = "ws"
    }

    u := url.URL{Scheme: scheme, Host: config.Server, Path: "/api/v1/agent"}
    header := make(http.Header)
    header.Set("Moon-Authentication", auth)
    header.Set("Moon-Version", version)

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

func readLoop() {
    <- connected
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error: read message, %v\n", err)
            conn.Close()
            conn = nil
            reconnect <- struct{}{}
            <- connected
            continue
        }

        go process(msg)
    }
}

func process(body []byte) {
    msg, err := protocol.Decode(body)
    if err != nil {
        // ignore those messgaes can not be decoded
        log.Printf("Error: decode messge, %v\n", err)
        return
    }

    if msg.Type != protocol.Req {
        // ignore thos messages are not REQ type
        log.Printf("Error: type not supported, %s\n", msg.Type)
        return
    }

    log.Printf("Recv request of method %s\n", msg.Method)

    payload, err := handle(msg.Method, msg.Payload)
    if err != nil {
        log.Printf("Error: exec request, %v\n", err)
        responseErr(msg.Uuid, protocol.StatusError, err)
        return
    }

    response(msg.Uuid, protocol.StatusOK, payload)
}

func responseErr(uuid, status string, err error) {
    payload, _ := json.Marshal(map[string]string{"error": err.Error()})
    response(uuid, status, payload)
}

func response(uuid, status string, payload []byte) {
    msg := new(protocol.Msg)
    msg.Type = protocol.Res
    msg.Uuid = uuid
    msg.Method = status
    msg.Payload = payload

    body, err := protocol.Encode(msg)
    if err != nil {
        log.Printf("Error: encode messge, %v\n", err)
        return
    }

    for {
        if conn == nil {
            time.Sleep(3 * time.Second)
            continue
        }
        if err := conn.WriteMessage(websocket.TextMessage, body); err != nil {
            log.Printf("Error: write message, %v\n", err)
        }
        break
    }
}

func worker() {
    auth, ok := getAuth()
    if !ok {
        log.Fatalf("Error: missing authentication key\n")
    }

    go dail(auth)
    go readLoop()

    <- interrupt
    // TBD
    done <- struct{}{}
}
