# Library
## Installation

If you want to install the latest version of the library, use this line:

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

Also see the [documentation](https://pkg.go.dev/github.com/yule-l/tm)

## Examples

See [examples](examples)
