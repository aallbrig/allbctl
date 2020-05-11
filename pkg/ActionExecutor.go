package pkg

import "os/exec"

type Action struct {
	Name string
	Cmd  *exec.Cmd
}

var ActionsToExecute []Action

func AddActionToQueue(a Action) {
	ActionsToExecute = append(ActionsToExecute, a)
}
