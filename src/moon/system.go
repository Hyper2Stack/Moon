package main

import (
    "log"
    "net"
    "os"
    "os/exec"
    "strings"
)

var version string = "0.1"

type AgentInfo struct {
    Version string `json:"version"`
}

func getAgentInfo() *AgentInfo {
    info := new(AgentInfo)
    info.Version = version
    return info
}

////////////////////////////////////////////////////////////////////////////////

type NodeInfo struct {
    Hostname string `json:"hostname"`
    Nics     []*Nic `json:"nics"`
}

type Nic struct {
    Name    string   `json:"name"`
    Ip4Addr string   `json:"ip4addr"`
    Tags    []string `json:"tags"`
}

func getNodeInfo() *NodeInfo {
    info := new(NodeInfo)
    info.Hostname, _ = os.Hostname()
    info.Nics = make([]*Nic, 0)

    intfs, _ := net.Interfaces()
    for _, i := range intfs {
        // remove docker0 and those virtual interface belongs to containers
        if i.Name == "docker0" || strings.HasPrefix(i.Name, "veth") {
            continue
        }

        addrs, _ := i.Addrs()
        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
                nic := new(Nic)
                nic.Ip4Addr = ipnet.IP.String()
                nic.Name = i.Name
                info.Nics = append(info.Nics, nic)
                break
            }
        }
    }

    // find the interface attached to default gateway
    args := []string{"-c", "ip route get 8.8.8.8 | head -1 | awk '{print $5}'"}
    out, err := exec.Command("bash", args...).Output()
    if err != nil {
        log.Printf("Error: exec commond, %s %v, %v\n", "bash", args, err)
        return info
    }

    for _, nic := range info.Nics {
        if nic.Name == strings.Trim(string(out), "\n") {
            nic.Tags = append(nic.Tags, "default")
        }
    }

    return info
}
