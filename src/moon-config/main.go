package main

import (
    "flag"
    "fmt"

    "moon/cfg"
)

func main() {
    key := flag.String("key", "", "agent authentication key")
    path := flag.String("conf", "/etc/moon/moon.cfg", "config file")
    flag.Parse()

    if len(*key) == 0 {
        fmt.Printf("Error: invalid key %s\n", *key)
        return
    }

    c, err := cfg.Parse(*path)
    if err != nil {
        fmt.Printf("Error: unable to parse %s, %v\n", path, err)
        return
    }

    if err := cfg.ResetKey(c.KeyFile, *key); err != nil {
        fmt.Printf("Error: unable to reset key, %v\n", err)
        return
    }
}
