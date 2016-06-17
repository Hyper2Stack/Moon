package protocol

import (
    "errors"
    "fmt"
    "regexp"
    "strings"
)

const (
    MsgHeader = "MOON/"
)

func Version() string {
    return "0.1"
}

var ErrInvalidMsg = errors.New("protocol: invalid message format")
var ErrInvalidVersion = errors.New("protocol: version mismatch")

// Message format:
// <version> <type> <uuid> <method>\r\n<payload>
// version: MOON/0.1
// type: [0-9a-zA-Z_.-]+
// uuid: [0-9a-zA-Z_.-]+
// method: [0-9a-zA-Z_.-]+
// payload: can be omitted

type Msg struct {
    Type    string
    Uuid    string
    Method  string
    Payload []byte
}

func (msg *Msg) Valid() bool {
    re := regexp.MustCompile("^[a-zA-Z0-9_.-]+$")
    return re.MatchString(msg.Type) && re.MatchString(msg.Uuid) && re.MatchString(msg.Method)
}

func Decode(body []byte) (*Msg, error) {
    ss1 := strings.SplitN(string(body), "\r\n", 2)
    ss2 := strings.Split(ss1[0], " ")
    if len(ss2) != 4 {
        return nil, ErrInvalidMsg
    }

    prefix := ss2[0]
    if len(prefix) < 5 || prefix[0:5] != MsgHeader {
        return nil, ErrInvalidMsg
    }

    if prefix[5:] != Version() {
        return nil, ErrInvalidVersion
    }

    msg := new(Msg)
    msg.Type = ss2[1]
    msg.Uuid = ss2[2]
    msg.Method = ss2[3]
    msg.Payload = []byte(ss1[1])

    return msg, nil
}

func Encode(msg *Msg) ([]byte, error) {
    if !msg.Valid() {
        return nil, ErrInvalidMsg
    }

    return []byte(fmt.Sprintf(
        "%s%s %s %s %s\r\n%s",
        MsgHeader, Version(), msg.Type, msg.Uuid, msg.Method, msg.Payload,
    )), nil
}
