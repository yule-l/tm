package tm

import (
	"sync"
	"time"
)

type tasks struct {
	data map[string]result
	mu   sync.Mutex
}

func newTasks() *tasks {
	return &tasks{
		data: make(map[string]result),
	}
}

func (t *tasks) store(task string, r result) {
	t.mu.Lock()
	t.data[task] = r
	t.mu.Unlock()
}

func (t *tasks) load(task string) (r result, ok bool) {
	t.mu.Lock()
	r, ok = t.data[task]
	t.mu.Unlock()
	return r, ok
}

type taskStatus uint8

const (
	done       taskStatus = iota
	delay      taskStatus = 1
	inProgress taskStatus = 2
)

type rerr struct {
	Attempt   uint8     `json:"attempt"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type result struct {
	Payload string     `json:"payload"`
	Attempt uint8      `json:"attempt"`
	Errors  []rerr     `json:"errors"`
	Status  taskStatus `json:"status"`
}
