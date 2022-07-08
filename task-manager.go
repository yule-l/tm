package tm

import (
	"bufio"
	"encoding/json"
	"io"
	"time"
)

// TaskManager is a task manager interface
type TaskManager interface {
	Next() (task string, finished bool)
	Finish(task string) error
	Error(task string, e error) error
}

var _ TaskManager = (*taskManager)(nil)

type taskManager struct {
	tasksInfo    *tasks
	delayedTasks *delayedTasks
	tasksQueue   *bufio.Scanner
	tasksResult  io.ReadWriter
	maxAttempts  uint8
}

// NewTasksManager will create task manager, open tasks file, open done file, preload data from done file
func NewTasksManager(cfg Config) (*taskManager, error) {
	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	doneFile, err := openDoneFile(cfg.FilePath, cfg.Force)
	if err != nil {
		return nil, err
	}

	queueFile, err := openQueueFile(cfg.FilePath)
	if err != nil {
		return nil, err
	}

	return &taskManager{
		tasksInfo:    loadTasksInfo(doneFile),
		delayedTasks: newDelayedTasks(),
		tasksQueue:   queueFile,
		tasksResult:  doneFile,
		maxAttempts:  cfg.MaxRetries,
	}, nil
}

// Next returns task from queue and notasks variable, if notasks is true, queue is empty
func (t *taskManager) Next() (task string, notasks bool) {
	task, notasks = t.nextQueueTask()
	if notasks {
		return t.nextDelayedTask()
	}
	return task, false
}

func (t *taskManager) nextQueueTask() (task string, notasks bool) {
	for t.tasksQueue.Scan() {
		task := t.tasksQueue.Text()
		r, ok := t.tasksInfo.load(task)
		if ok {
			if r.Status == done {
				continue
			}
		} else {
			t.tasksInfo.store(task, result{
				Payload: task,
				Attempt: 0,
				Errors:  nil,
				Status:  inProgress,
			})
		}
		return task, false
	}
	return "", true
}

func (t *taskManager) nextDelayedTask() (task string, notasks bool) {
	for {
		task, cursor, ok := t.delayedTasks.next()
		if !ok {
			return "", true
		}
		r, ok := t.tasksInfo.load(task)
		if ok {
			if r.Status == done {
				t.delayedTasks.unset(cursor)
				continue
			}
		} else {
			t.tasksInfo.store(task, result{
				Payload: task,
				Attempt: 0,
				Errors:  nil,
				Status:  inProgress,
			})
		}
		return task, false
	}
}

// Finish marks task as done
func (t *taskManager) Finish(task string) error {
	return t.writeResult(task, nil)
}

// Error marks task as not completed
// If Config.MaxRetries not reached task will be marked with delayed status
func (t *taskManager) Error(task string, e error) error {
	return t.writeResult(task, e)
}

func (t *taskManager) writeResult(task string, e error) error {
	r, _ := t.tasksInfo.load(task)
	r.Status = done
	r.Payload = task
	r.Attempt++

	t.fillResultsWithErrorIfNotNil(&r, e)

	if r.Status == delay {
		t.delayedTasks.append(task)
	}

	t.tasksInfo.store(task, r)
	rb, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = t.tasksResult.Write(append(rb, []byte("\n")...))
	return err
}

func (t *taskManager) fillResultsWithErrorIfNotNil(r *result, e error) {
	if e == nil {
		return
	}
	r.Errors = append(r.Errors, rerr{
		Attempt:   r.Attempt,
		Message:   e.Error(),
		Timestamp: time.Now(),
	})
	if r.Attempt < t.maxAttempts {
		r.Status = delay
	}
}
