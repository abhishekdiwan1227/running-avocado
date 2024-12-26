package avo

import (
	"encoding/json"
	"os"
	"os/exec"
)

type Runner interface {
	Run()
}

func CreateRunner(task Task) Runner {
	switch *task.TaskDefinitionType {
	case Local:
		return &LocalScriptRunner{Task: task}
	}
	return nil
}

type LocalScriptRunner struct {
	Task Task
}

func (runner *LocalScriptRunner) Run() {
	wagon := GetConfig().Wagon
	definition := wagon.GetScriptDefinition(runner.Task.ID)

	var arguments []string
	if definition.Arguments != nil && len(*definition.Arguments) > 0 {
		if err := json.Unmarshal([]byte(*definition.Arguments), &arguments); err != nil {
			panic(err)
		}
	}

	cmd := exec.Command(definition.Path, arguments...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
