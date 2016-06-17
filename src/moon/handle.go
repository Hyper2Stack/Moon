package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"

    "github.com/hyper2stack/mooncommon/protocol"
)

var version = "0.1"

func getAgentInfo() ([]byte, error) {
    info := new(protocol.Agent)
    info.Version = version
    payload, _ := json.Marshal(info)
    return payload, nil
}

func getNodeInfo() ([]byte, error) {
    info := new(protocol.Node)
    info.Hostname, _ = os.Hostname()
    info.Nics = make([]*protocol.Nic, 0)

    intfs, _ := net.Interfaces()
    for _, i := range intfs {
        // remove docker0 and those virtual interface belongs to containers
        if i.Name == "docker0" || strings.HasPrefix(i.Name, "veth") {
            continue
        }

        addrs, _ := i.Addrs()
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
                nic := new(protocol.Nic)
                nic.Ip4Addr = ipnet.IP.String()
                nic.Name = i.Name
                nic.IsPrimary = false
                info.Nics = append(info.Nics, nic)
                break
            }
        }
    }

    // find the interface attached to default gateway
    args := []string{"-c", "ip route get 8.8.8.8 | head -1 | awk '{print $5}'"}
    out, err := exec.Command("bash", args...).Output()
    if err == nil {
        for _, nic := range info.Nics {
            if nic.Name == strings.Trim(string(out), "\n") {
                nic.IsPrimary = true
                break
            }
        }
    } else {
        log.Printf("Error: exec commond, %s %v, %v\n", "bash", args, err)
    }

    payload, _ := json.Marshal(info)
    return payload, nil
}

func execScript(input []byte) ([]byte, error) {
    script := new(protocol.Script)
    if err := json.Unmarshal(input, script); err != nil {
        return nil, fmt.Errorf("invalid payload, %v", err)
    }

    result := new(protocol.ScriptResult)
    result.Ok = true
    result.CommandResults = make([]*protocol.CommandResult, 0)

    for _, c := range script.Commands {
        log.Printf("Executing command %v\n", c)

        exitCode := 0
        out, err := exec.Command(c.Command, c.Args...).CombinedOutput()
        if err != nil {
            exitCode = 99
            if exiterr, ok := err.(*exec.ExitError); ok {
                if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                    exitCode = status.ExitStatus()
                }
            }

            if c.Restrict {
                result.Ok = false
            }
        }

        r := new(protocol.CommandResult)
        r.Command = c
        r.ExitCode = exitCode
        r.Output = string(out)
        result.CommandResults = append(result.CommandResults, r)
    }

    payload, _ := json.Marshal(result)
    return payload, nil
}

func createFile(input []byte) ([]byte, error) {
    file := new(protocol.File)
    if err := json.Unmarshal(input, file); err != nil {
        return nil, fmt.Errorf("invalid payload, %v", err)
    }

    mode, err := strconv.ParseUint(file.Mode, 8, 32)
    if err != nil {
        return nil, err
    }

    if err := ioutil.WriteFile(file.Path, []byte(file.Content), os.FileMode(mode)); err != nil {
        return nil, err
    }

    return nil, nil
}

func handle(method string, payload []byte) ([]byte, error) {
    switch method {
    case protocol.MethodNodeInfo:
        return getNodeInfo()
    case protocol.MethodAgentInfo:
        return getAgentInfo()
    case protocol.MethodExecScript:
        return execScript(payload)
    case protocol.MethodCreateFile:
        return createFile(payload)
    default:
        return nil, fmt.Errorf("method %s not supported", method)
    }
}
