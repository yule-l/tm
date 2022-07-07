package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/yule-l/tm"
)

var (
	maxRetries    = flag.Uint("retries", 3, "Max task retries")
	tasksFilePath = flag.String("input", "tasks.txt", "Input file path")
	force         = flag.Bool("force", false, "Force start from first task")
)

func main() {
	flag.Parse()
	taskManager, err := tm.NewTasksManager(tm.Config{
		Force:      *force,
		FilePath:   *tasksFilePath,
		MaxRetries: uint8(*maxRetries),
	})
	if err != nil {
		log.Fatalln(err)
	}

	doer := tm.NewDefaultDoer(taskManager)
	doer.Do(func(task string) error {
		if task == "error" {
			err = taskManager.Error(task, errors.New("some error"))
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Println(task)
		return nil
	})
}
