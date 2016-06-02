package cfg

import (
    "encoding/json"
    "io/ioutil"

    "github.com/satori/go.uuid"
    "github.com/shizeeg/gcfg"
)

type Cfg struct {
    PidFile string `gcfg:"pid_file"`
    LogFile string `gcfg:"log_file"`
    KeyFile string `gcfg:"key_file"`
}

type Auth struct {
    Key  string `json:"key"`
    Uuid string `json:"uuid"`
}

func Parse(cfgPath string) (*Cfg, error) {
    cfg := struct {
        Moon Cfg
    }{}

    if err := gcfg.ReadFileInto(&cfg, cfgPath); err != nil {
        return nil, err
    }

    if cfg.Moon.PidFile == "" {
        cfg.Moon.PidFile = "/var/run/moon.pid"
    }

    if cfg.Moon.LogFile == "" {
        cfg.Moon.LogFile = "/var/log/moon/moon.log"
    }

    if cfg.Moon.KeyFile == "" {
        cfg.Moon.KeyFile = "/etc/moon/key.json"
    }

    return &cfg.Moon, nil
}

func ParseKey(keyPath string) (*Auth, error) {
    bytes, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return nil, err
    }

    auth := new(Auth)
    if err := json.Unmarshal(bytes, auth); err != nil {
        return nil, err
    }

    return auth, nil
}

func ResetKey(keyPath string, key string) error {
    auth, err := ParseKey(keyPath)
    if err == nil && auth.Key == key {
        return nil
    }

    auth = new(Auth)
    auth.Key = key
    auth.Uuid = GenerateUuid()

    content, _ := json.Marshal(auth)
    if err := ioutil.WriteFile(keyPath, content, 0600); err != nil {
        return err
    }

    return nil
}

func GenerateUuid() string {
    return uuid.NewV4().String()
}
