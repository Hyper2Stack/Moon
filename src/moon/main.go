package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "syscall"

    "github.com/sevlyar/go-daemon"
    "github.com/shizeeg/gcfg"
)

var (
    interrupt = make(chan struct{})
    done = make(chan struct{})
)

var config Config

type Config struct {
    PidFile string `gcfg:"pid_file"`
    LogFile string `gcfg:"log_file"`
}

func parseConfig() {
    cfg := struct {
        Moon Config
    }{}

    err := gcfg.ReadFileInto(&cfg, "/etc/moon/moon.cfg")
    if err != nil {
        log.Fatalf("Failed to parse moon.cfg: %s", err)
    }

    config = cfg.Moon
    if config.PidFile == "" {
        config.PidFile = "/var/run/moon.pid"
    }

    if config.LogFile == "" {
        config.LogFile = "/var/log/moon/moon.log"
    }
}

func termHandler(sig os.Signal) error {
    log.Println("Terminating...")
    interrupt <- struct{}{}
    <-done
    return daemon.ErrStop
}

func printVersion() {
    fmt.Println(version)
}

func main() {
    version := flag.Bool("v", false, "print version")
    signal := flag.String("s", "", `send signal to the daemon
        quit — graceful shutdown`)
    flag.Parse()
    parseConfig()

    if *version {
        printVersion()
        return
    }

    daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)

    cntxt := &daemon.Context{
        PidFileName: config.PidFile,
        PidFilePerm: 0644,
        LogFileName: config.LogFile,
        LogFilePerm: 0644,
        WorkDir:     "/",
        Umask:       027,
        Args:        []string{"[/usr/sbin/moon]"},
    }

    if len(daemon.ActiveFlags()) > 0 {
        d, err := cntxt.Search()
        if err != nil {
            log.Fatalln("Error: unable to send signal to the daemon,", err)
        }
        if d != nil {
            daemon.SendCommands(d)
        }
        return
    }

    d, err := cntxt.Reborn()
    if err != nil {
        log.Fatalln("Error: unable to reborn daemon process,", err)
    }
    if d != nil {
        return
    }
    defer cntxt.Release()

    log.Println("Daemon started")

    go worker()

    err = daemon.ServeSignals()
    if err != nil {
        log.Println("Error: unable to serve signals,", err)
    }
    log.Println("Daemon terminated")
}
