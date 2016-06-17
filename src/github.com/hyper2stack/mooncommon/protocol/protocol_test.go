package protocol

import (
    "fmt"
    "testing"
)

func Test_Decode_NormalWithPayload(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s REQ 000 GET\r\nxyz", MsgHeader, Version()))
    msg, err := Decode(body)
    if err != nil {
        t.Error(err)
    }

    if msg.Type != "REQ" {
        t.Errorf("msg.Type %s != REQ\n", msg.Type)
    }

    if msg.Uuid != "000" {
        t.Errorf("msg.Uuid %s != 000\n", msg.Uuid)
    }

    if msg.Method != "GET" {
        t.Errorf("msg.Method %s != GET\n", msg.Method)
    }

    if string(msg.Payload) != "xyz" {
        t.Errorf("msg.Payload %s != xyz\n", msg.Payload)
    }
}

func Test_Decode_NormalWithoutPayload(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s RES 100 post\r\n", MsgHeader, Version()))
    msg, err := Decode(body)
    if err != nil {
        t.Error(err)
    }

    if msg.Type != "RES" {
        t.Errorf("msg.Type %s != RES\n", msg.Type)
    }

    if msg.Uuid != "100" {
        t.Errorf("msg.Uuid %s != 100\n", msg.Uuid)
    }

    if msg.Method != "post" {
        t.Errorf("msg.Method %s != post\n", msg.Method)
    }

    if string(msg.Payload) != "" {
        t.Errorf("msg.Payload '%s' should be empty\n", msg.Payload)
    }
}

func Test_Decode_InvalidPrefix(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s REQ 000 GET\r\nxyz", "XXX", Version()))
    if _, err := Decode(body); err == nil {
        t.Errorf("should be error\n")
    }
}

func Test_Decode_LessFields(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s REQ 000\r\n", MsgHeader, Version()))
    if _, err := Decode(body); err == nil {
        t.Errorf("should be error\n")
    }
}

func Test_Decode_MoreFields(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s REQ 000 GET f1\r\nxyz", MsgHeader, Version()))
    if _, err := Decode(body); err == nil {
        t.Errorf("should be error\n")
    }
}

func Test_Decode_InvalidVersion(t *testing.T) {
    body := []byte(fmt.Sprintf("%s%s REQ 000 GET\r\n", MsgHeader, "XXX"))
    if _, err := Decode(body); err == nil {
        t.Errorf("should be error\n")
    }
}

func Test_Encode_NormalWithPayload(t *testing.T) {
    msg := new(Msg)
    msg.Type = "RES"
    msg.Uuid = "011"
    msg.Method = "200"
    msg.Payload = []byte("xyz")

    body, err := Encode(msg)
    if err != nil {
        t.Error(err)
    }

    target := fmt.Sprintf("%s%s RES 011 200\r\nxyz", MsgHeader, Version())
    if string(body) != target {
        t.Errorf("msg body %s != %s\n", body, target)
    }
}

func Test_Encode_NormalWithoutPayload(t *testing.T) {
    msg := new(Msg)
    msg.Type = "RES"
    msg.Uuid = "011"
    msg.Method = "200"

    body, err := Encode(msg)
    if err != nil {
        t.Error(err)
    }

    target := fmt.Sprintf("%s%s RES 011 200\r\n", MsgHeader, Version())
    if string(body) != target {
        t.Errorf("msg body %s != %s\n", body, target)
    }
}

func Test_Encode_InvalidField(t *testing.T) {
    msg := new(Msg)
    msg.Type = "REQ"
    msg.Uuid = "111"
    msg.Method = "GET /abc"

    if _, err := Encode(msg); err == nil {
        t.Errorf("should be error\n")
    }
}
