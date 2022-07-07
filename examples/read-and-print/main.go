package main

import (
	"errors"
	"flag"
	"log"
	"strings"

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

	for {
		task, nodata := taskManager.Next()
		if nodata {
			return
		}
		task = strings.Trim(task, " ")
		if task == "error" {
			err = taskManager.Error(task, errors.New("some error"))
			if err != nil {
				log.Println(err)
			}
			continue
		}
		err = taskManager.Finish(task)
		if err != nil {
			log.Println(err)
		}
	}
}
