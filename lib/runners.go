package avocado

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Runner interface {
	Run(task *Task) *WorkloadResult
}

func CreateRunner(task *Task) Runner {
	switch *task.TaskDefinitionType {
	case Local:
		return &LocalScriptRunner{Task: *task}
	}
	return nil
}

type LocalScriptRunner struct {
	Task Task
}

func (runner *LocalScriptRunner) Run(task *Task) *WorkloadResult {
	wagon := GetConfig().Wagon
	definition := wagon.GetScriptDefinition(runner.Task.ID)

	if _, err := os.Stat(definition.Entrypoint); errors.Is(err, os.ErrNotExist) {
		panic(fmt.Errorf("file not found at %s", definition.Entrypoint))
	}

	var arguments []string
	if definition.Arguments != nil && len(*definition.Arguments) > 0 {
		if err := json.Unmarshal([]byte(*definition.Arguments), &arguments); err != nil {
			panic(err)
		}
	}

	cmd := exec.Command(definition.Command, append([]string{definition.Entrypoint}, arguments...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Default().Fatal(err.Error())
	} else {
		fmt.Printf("%d", cmd.ProcessState.ExitCode())
	}

	return &WorkloadResult{
		ReturnCode: cmd.ProcessState.ExitCode(),
	}
}
