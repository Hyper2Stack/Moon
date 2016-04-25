package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"

    "github.com/satori/go.uuid"
)

var keyPath string = "/etc/moon/key.cfg"

type Auth struct {
    Key  string `json:"key"`
    Uuid string `json:"uuid"`
}

func getAuth() *Auth {
    bytes, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return nil
    }

    auth := new(Auth)
    err = json.Unmarshal(bytes, auth)
    if err != nil {
        return nil
    }

    return auth
}

func generateUuid() string {
    return uuid.NewV4().String()
}

func main() {
    key := flag.String("key", "", "agent authentication key")
    flag.Parse()

    if len(*key) == 0 {
        fmt.Printf("Error: invalid key %s\n", *key)
        return
    }

    auth := getAuth()
    if auth != nil && auth.Key == *key {
        return
    }

    auth = new(Auth)
    auth.Key = *key
    auth.Uuid = generateUuid()

    content, _ := json.Marshal(auth)
    if err := ioutil.WriteFile(keyPath, content, 0644); err != nil {
        fmt.Printf("Error: write %s, %v\n", keyPath, err)
        return
    }
}
