package main

import (
    "flag"
    "log"
    "os"
    "syscall"

    "github.com/sevlyar/go-daemon"
)

var (
    signal = flag.String("s", "", `send signal to the daemon
        quit — graceful shutdown
        stop — fast shutdown
        reload — reloading the configuration file`)
)

func main() {
    flag.Parse()
    daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
    daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
    daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

    cntxt := &daemon.Context{
        PidFileName: "/var/run/moon.pid",
        PidFilePerm: 0644,
        LogFileName: "/var/log/moon/moon.log",
        LogFilePerm: 0640,
        WorkDir:     "/",
        Umask:       027,
        Args:        []string{"[/usr/sbin/moon]"},
    }

    if len(daemon.ActiveFlags()) > 0 {
        d, err := cntxt.Search()
        if err != nil {
            log.Fatalln("Error: unable to send signal to the daemon,", err)
        }
        daemon.SendCommands(d)
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

var (
    interrupt = make(chan struct{})
    done = make(chan struct{})
)

func termHandler(sig os.Signal) error {
    log.Println("Terminating...")
    interrupt <- struct{}{}
    <-done
    return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
    log.Println("Reloading...")
    // TBD
    return nil
}
