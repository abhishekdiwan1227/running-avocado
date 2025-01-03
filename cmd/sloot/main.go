package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	avocado "github.com/abhishekdiwan1227/running-avocado/lib"
)

var helpText string = `invalid arguments
usage:
sloot import {filepath}
sloot import -f|--filename {filepath}
`

func main() {
	avocado.Start()
	avocado.GetConfig().Wagon.MigrateDatabase()

	args := os.Args[1:]

	if len(args) == 2 && args[0] == "import" {
		filePath := args[1]
		importFromFile(filePath)
	} else if len(args) == 3 && args[0] == "import" && (args[1] == "-f" || args[1] == "--filename") {
		filePath := args[1]
		importFromFile(filePath)
	} else {
		fmt.Print(helpText)
	}
}

func importFromFile(filePath string) {
	fmt.Printf("Imporing from file %s\n", filePath)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		panic(err.Error())
	}

	tdStream, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var task *avocado.Task = &avocado.Task{
		Active: true,
	}
	err = json.Unmarshal(tdStream, task)
	if err != nil {
		panic(err.Error())
	}

	db := avocado.GetConfig().Wagon

	db.AddTask(task)

	definition := &avocado.ScriptTaskDefinition{
		TaskID:     task.ID,
		Entrypoint: "/home/mango/source/repos/running-avocado/.dev/test.sh",
		Command:    "/bin/sh",
	}

	db.AddScriptDefinition(definition, task)

	taskType := avocado.Local

	task.TaskDefinitionID = &definition.ID
	task.TaskDefinitionType = &taskType
}
