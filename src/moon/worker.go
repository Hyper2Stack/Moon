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
)

type Message struct {
    Type    string      `json:"type"`
    Content interface{} `json:"content"`
}

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
    u := url.URL{Scheme: "ws", Host: "localhost", Path: "/api/v1/agent"}
    header := make(http.Header)
    header.Set("Moon-Authentication", auth)

    for {
        log.Printf("Trying to connect %s\n", u.String())

        c, res, err := websocket.DefaultDialer.Dial(u.String(), header)
        if err != nil {
            log.Printf("Error: websocket dail, %v\n", err)
            errBody, _ := ioutil.ReadAll(res.Body)
            log.Printf("Error: websocket body, %s\n", string(errBody))
            time.Sleep(3 * time.Second)
            continue
        }

        log.Printf("Successfully to create websocket connectiton to %s\n", u.String())
        conn = c
        connected <- struct{}{}
        <- reconnect
    }
}

func process() {
    <- connected
    for {
        message := new(Message)
        if err := conn.ReadJSON(message); err != nil {
            log.Printf("Error: read message, %v\n", err)
            conn.Close()
            reconnect <- struct{}{}
            <- connected
            continue
        }

        log.Printf("Recv message of type %s\n", message.Type)
        job, err := createJob(message)
        if err != nil {
            log.Printf("Error: create job, %v\n", err)
            continue
        }

        result := job.Run()
        if err := conn.WriteJSON(result); err != nil {
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
