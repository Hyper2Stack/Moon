package main

import (
    "encoding/json"
    "flag"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/hyper2stack/mooncommon/protocol"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}

var conn *websocket.Conn

var clientKey string
var clientVersion string

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

    clientKey = r.Header.Get("Moon-Authentication")
    clientVersion = r.Header.Get("Moon-Version")

    log.Println("Connected.")
    conn = c
    <- disconnect
}

func test(w http.ResponseWriter, r *http.Request) {
    msgs := make([][]byte, 0)
    msgs = append(msgs, genReqNodeInfo("1001"))
    msgs = append(msgs, genReqAgentInfo("1002"))
    msgs = append(msgs, genReqExecScriptGood("1003"))
    msgs = append(msgs, genReqExecScriptBad("1004"))
    msgs = append(msgs, genReqCreateFile("1005"))

    done := make(chan struct{})

    result := make(map[string]bool)
    go func() {
        for i := 0; i < len(msgs); i++ {
            _, body, err := conn.ReadMessage()
            if err != nil {
                w.WriteHeader(http.StatusInternalServerError)
                return
            }
    
            msg, _ := protocol.Decode(body)
            switch msg.Uuid {
            case "1001":
                result["node-info"] = checkNodeInfo(msg.Method, msg.Payload)
            case "1002":
                result["agent-info"] = checkAgentInfo(msg.Method, msg.Payload)
            case "1003":
                result["exec-script-good"] = checkExecScriptGood(msg.Method, msg.Payload)
            case "1004":
                result["exec-script-bad"] = checkExecScriptBad(msg.Method, msg.Payload)
            case "1005":
                result["create-file"] = checkCreateFile(msg.Method, msg.Payload)
            }
        }
        done <- struct{}{}
    }()

    for _, msg := range msgs {
        if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
    }

    <- done

    payload, _ := json.Marshal(result)
    w.Write(payload)
}

func genReqNodeInfo(uuid string) []byte {
    msg, _ := protocol.Encode(&protocol.Msg{
        Type: protocol.Req,
        Uuid: uuid,
        Method: protocol.MethodNodeInfo,
    })

    return msg
}

func checkNodeInfo(status string, payload []byte) bool {
    log.Printf("node info: %s", payload)
    return status == protocol.StatusOK
}

func genReqAgentInfo(uuid string) []byte {
    msg, _ := protocol.Encode(&protocol.Msg{
        Type: protocol.Req,
        Uuid: uuid,
        Method: protocol.MethodAgentInfo,
    })

    return msg
}

func checkAgentInfo(status string, payload []byte) bool {
    log.Printf("agent info: %s", payload)
    return status == protocol.StatusOK 
}

func genReqExecScriptGood(uuid string) []byte {
    script := &protocol.Script{
        Commands: []*protocol.Command{
            &protocol.Command{
                Command: "ls",
                Args: []string{"/tmp/1/not-exist"},
                Restrict: false,
            },
            &protocol.Command{
                Command: "ls",
                Args: []string{"/usr"},
                Restrict: true,
            },
        },
    }

    payload, _ := json.Marshal(script)
    msg, _ := protocol.Encode(&protocol.Msg{
        Type: protocol.Req,
        Uuid: uuid,
        Method: protocol.MethodExecScript,
        Payload: payload,
    })

    return msg
}

func checkExecScriptGood(status string, payload []byte) bool {
    log.Printf("exec script good: %s", payload)

    if status != protocol.StatusOK {
        return false
    }

    result := new(protocol.ScriptResult)
    if err := json.Unmarshal(payload, result); err != nil {
        return false
    }

    return result.Ok
}

func genReqExecScriptBad(uuid string) []byte {
    script := &protocol.Script{
        Commands: []*protocol.Command{
            &protocol.Command{
                Command: "ls",
                Args: []string{"/tmp/1/not-exist"},
                Restrict: true,
            },
            &protocol.Command{
                Command: "ls",
                Args: []string{"/usr"},
                Restrict: true,
            },
        },
    }

    payload, _ := json.Marshal(script)
    msg, _ := protocol.Encode(&protocol.Msg{
        Type: protocol.Req,
        Uuid: uuid,
        Method: protocol.MethodExecScript,
        Payload: payload,
    })

    return msg
}

func checkExecScriptBad(status string, payload []byte) bool {
    log.Printf("exec script bad: %s", payload)

    if status != protocol.StatusOK {
        return false
    }

    result := new(protocol.ScriptResult)
    if err := json.Unmarshal(payload, result); err != nil {
        return false
    }

    return !result.Ok 
} 

func genReqCreateFile(uuid string) []byte {
    file := &protocol.File{
        Path: "/tmp/abcdefg",
        Mode: "755",
        Content: "abcdefg",
    }

    payload, _ := json.Marshal(file)
    msg, _ := protocol.Encode(&protocol.Msg{
        Type: protocol.Req,
        Uuid: uuid,
        Method: protocol.MethodCreateFile,
        Payload: payload,
    })

    return msg
}

func checkCreateFile(status string, payload []byte) bool {
    log.Printf("create file: %s", payload)
    return status == protocol.StatusOK
}

func clientInfo(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(clientKey+"\n"+clientVersion+"\n"))
}

func done(w http.ResponseWriter, r *http.Request) {
    disconnect <- struct {}{}
}

func main() {
    flag.Parse()
    http.HandleFunc("/api/v1/agent", ws)
    http.HandleFunc("/test", test)
    http.HandleFunc("/client", clientInfo)
    http.HandleFunc("/done", done)
    log.Fatal(http.ListenAndServe(*addr, nil))
}
