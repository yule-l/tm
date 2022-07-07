package tm

import (
	"sync"
	"time"
)

type delayedTasks struct {
	mu sync.Mutex

	data   []string
	cursor int
}

func newDelayedTasks() *delayedTasks {
	d := new(delayedTasks)
	d.reset()
	return d
}

func (d *delayedTasks) append(task string) {
	d.mu.Lock()
	d.data = append(d.data, task)
	d.mu.Unlock()
}

func (d *delayedTasks) next() (task string, id int, ok bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	l := len(d.data)
	for {
		if d.cursor >= l {
			return "", -1, false
		}
		cur := d.cursor
		val := d.data[cur]
		d.cursor++
		if val == unsetValue {
			continue
		}
		return val, d.cursor, true
	}
}

const unsetValue = "$#%unset$#%"

func (d *delayedTasks) unset(id int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if id > len(d.data)-1 {
		return
	}
	d.data[id] = unsetValue
}

func (d *delayedTasks) reset() {
	d.mu.Lock()
	d.data = make([]string, 0)
	d.cursor = 0
	d.mu.Unlock()
}

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
