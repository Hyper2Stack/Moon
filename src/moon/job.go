package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os/exec"
)

type Job interface {
    Run() *Response
}

////////////////////////////////////////////////////////////////////////////////

type AgentInfoJob struct {}

func (j *AgentInfoJob) Run() *Response {
    res := new(Response)
    res.Status = StatusOK
    res.Content, _ = json.Marshal(getAgentInfo())
    return res
}

////////////////////////////////////////////////////////////////////////////////

type NodeInfoJob struct {}

func (j *NodeInfoJob) Run() *Response {
    res := new(Response)
    res.Status = StatusOK
    res.Content, _ = json.Marshal(getNodeInfo())
    return res
}

////////////////////////////////////////////////////////////////////////////////

type ShellCommand struct {
    Command  string   `json:"command"`
    Args     []string `json:"args"`
    Restrict bool     `json:"restrict"`
}

type ScriptJob struct {
    Commands []*ShellCommand `json:"commands"`
}

type ScriptJobResult struct {
    ErrCommand *ShellCommand `json:"err_command"`
}

func (j *ScriptJob) Run() *Response {
    res := new(Response)
    for _, sc := range j.Commands {
        log.Printf("Executing command %s %v\n", sc.Command, sc.Args)
        if out, err := exec.Command(sc.Command, sc.Args...).Output(); err != nil {
            log.Printf("Command error, %v\noutput:\n%s\n", err, out)
            if sc.Restrict {
                res.Status = StatusError
                jobResult := new(ScriptJobResult)
                jobResult.ErrCommand = sc
                res.Content, _ = json.Marshal(jobResult)
                return res
            }
        }
    }

    res.Status = StatusOK
    return res
}

////////////////////////////////////////////////////////////////////////////////

func createJob(r *Request) (Job, error) {
    switch r.Action {
    case ActionAgentInfo:
        return new(AgentInfoJob), nil
    case ActionNodeInfo:
        return new(NodeInfoJob), nil
    case ActionExecShell:
        job := new(ScriptJob)
        if err := json.Unmarshal(r.Content, job); err != nil {
            return nil, fmt.Errorf("invalid request payload of action %s", r.Action)
        }
        return job, nil
    default:
        return nil, fmt.Errorf("not supported request action %s", r.Action)
    }
}
