# tm – task manager
[![GoDoc](https://pkg.go.dev/badge/github.com/yule-l/tm)](https://pkg.go.dev/github.com/yule-l/tm)
[![Go](https://github.com/yule-l/tm/actions/workflows/go.yml/badge.svg)](https://github.com/yule-l/tm/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/yule-l/tm/branch/master/graph/badge.svg?token=ZRL9IO6JNC)](https://codecov.io/gh/yule-l/tm)
[![Go Report Card](https://goreportcard.com/badge/github.com/yule-l/tm)](https://goreportcard.com/report/github.com/yule-l/tm)

The simplest task manager library in Go.

## Overview

`tm` is a library providing a simple mechanism for tasks control. 

### Features
* Read tasks line by line from file
* Mark tasks as completed
* Mark tasks as not completed and try complete this N times in the future

## Installation

If you want to install the latest version of the library, try this line:

```shell
go get -u github.com/yule-l/tm@latest
```

Then import `tm` in your code

```go
import "github.com/yule-l/tm"
```

## Quick Start

### Doers
Doer is a simple way to run your tasks in ordered, parallel or something else way, you can control them by `context.Context`.

#### Parallel Doer
```go
ctx := context.Background()
taskManager, _ := tm.NewTasksManager(tm.NewDefaultConfig("tasks.txt"))
doer := tm.NewDefaultParallelDoer(taskManager)
doer.Do(ctx, func(ctx context.Context, task string) error {
	// do tasks in goroutines
	return nil
})
```

#### Ordered Doer
```go
ctx := context.Background()
taskManager, _ := tm.NewTasksManager(tm.NewDefaultConfig("tasks.txt"))
doer := tm.NewDefaultOrderedDoer(taskManager)

doer.Do(ctx, func(ctx context.Context, task string) error {
	// do task one by one
	return nil
})
```

### Custom

or you can use task manager more specifically

```go
taskManager, _ := tm.NewTasksManager(tm.NewDefaultConfig("tasks.txt"))

for {
    task, nodata := taskManager.Next()
    if nodata {
        return
    }
    err := doTask(task)
	if err != nil {
		taskManager.Error(task, err)
		continue
    }
	
    taskManager.Finish(task)
}
```

## Examples

See [examples](examples)

## License

This library is licensed under either of

* [Apache License, Version 2.0](LICENSE-APACHE)
* [MIT license](LICENSE-MIT)

at your option.
