package main

import (
    "fmt"
    "log"
    "os/exec"
)

type Job interface {
    Run() interface{}
}

////////////////////////////////////////////////////////////////////////////////

type AgentInfoJob struct {}

func (j *AgentInfoJob) Run() interface{} {
    return getAgentInfo()
}

////////////////////////////////////////////////////////////////////////////////

type NodeInfoJob struct {}

func (j *NodeInfoJob) Run() interface{} {
    return getNodeInfo()
}

////////////////////////////////////////////////////////////////////////////////

type ShellCommand struct {
    Command  string
    Args     []string
    Restrict bool
}

type ScriptJob struct {
    Commands []*ShellCommand
}

type ScriptJobResult struct {
    Err        bool          `json:"err"`
    ErrCommand *ShellCommand `json:"err_command"`
}

func (j *ScriptJob) Run() interface{} {
    result := new(ScriptJobResult)
    for _, sc := range j.Commands {
        log.Printf("Executing command %s %v\n", sc.Command, sc.Args)
        if out, err := exec.Command(sc.Command, sc.Args...).Output(); err != nil {
            log.Printf("Command error, %v\noutput:\n%s\n", err, out)
            if sc.Restrict {
                result.Err = true
                result.ErrCommand = sc
                return result
            }
        }
    }

    return result
}

////////////////////////////////////////////////////////////////////////////////

func createJob(m *Message) (Job, error) {
    switch m.Type {
    case "get-agent-info":
        return new(AgentInfoJob), nil
    case "get-node-info":
        return new(NodeInfoJob), nil
    case "exec-shell-script":
        job, ok := m.Content.(*ScriptJob)
        if !ok {
            return nil, fmt.Errorf("invalid payload for type %s", m.Type)
        }
        return job, nil
    default:
        return nil, fmt.Errorf("not supported job type %s", m.Type)
    }
}
