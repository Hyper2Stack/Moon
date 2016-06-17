package protocol

import (
)


const (
    Req = "req"
    Res = "res"

    MethodNodeInfo   = "node.info"
    MethodAgentInfo  = "agent.info"
    MethodExecScript = "script.exec"
    MethodCreateFile = "file.create"

    StatusOK    = "ok"
    StatusError = "error"
)

type Node struct {
    Hostname string `json:"hostname"`
    Nics     []*Nic `json:"nics"`
}

type Nic struct {
    Name      string `json:"name"`
    Ip4Addr   string `json:"ip4addr"`
    IsPrimary bool   `json:"is_primary"`
}

type Agent struct {
    Version string `json:"version"`
}

type Command struct {
    Command  string   `json:"command"`
    Args     []string `json:"args"`
    Restrict bool     `json:"restrict"`
}

type CommandResult struct {
    Command  *Command `json:"command"`
    Output   string   `json:"output"`
    ExitCode int      `json:"exit_code"`
}

type Script struct {
    Commands []*Command `json:"commands"`
}

type ScriptResult struct {
    CommandResults []*CommandResult `json:"command_results"`
    Ok             bool             `json:"ok"`
}

type File struct {
    Path    string `json:"path"`
    Mode    string `json:"mode"`
    Content string `json:"content"`
}
