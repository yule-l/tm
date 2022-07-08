package tm

import (
	"context"
	"strings"
	"sync"
)

const DefaultMaxParallelTasks = 5
const DefaultTrimCutSet = " "

type Doer interface {
	Do(ctx context.Context, f func(ctx context.Context, task string) error)
}

// ParallelDoer is a simple and powerful way to do your tasks fast.
type ParallelDoer struct {
	tm               TaskManager
	trimCutSet       string
	maxParallelTasks int
}

// NewDefaultParallelDoer returns ParallelDoer with predefined trim cut sets and max parallel tasks
func NewDefaultParallelDoer(tm TaskManager) *ParallelDoer {
	return &ParallelDoer{
		tm:               tm,
		trimCutSet:       DefaultTrimCutSet,
		maxParallelTasks: DefaultMaxParallelTasks,
	}
}

// NewParallelDoer returns a new instance of ParallelDoer
func NewParallelDoer(
	tm TaskManager,
	trimCutSet string,
	maxParallelTasks int,
) *ParallelDoer {
	return &ParallelDoer{
		tm:               tm,
		trimCutSet:       trimCutSet,
		maxParallelTasks: maxParallelTasks,
	}
}

// Do the tasks parallel without strict order
//
// If f returns an error, it will be delayed
func (d *ParallelDoer) Do(ctx context.Context, f func(ctx context.Context, task string) error) {
	ch := make(chan string, d.maxParallelTasks)
	go d.producer(ctx, ch)

	d.consumer(ctx, f, ch)
}

func (d *ParallelDoer) producer(ctx context.Context, ch chan<- string) {
	defer close(ch)
	for {
		task, nodata := d.tm.Next()
		if nodata {
			return
		}
		task = strings.Trim(task, d.trimCutSet)
		select {
		case ch <- task:
		case <-ctx.Done():
			return
		}
	}
}

func (d *ParallelDoer) consumer(
	ctx context.Context,
	f func(ctx context.Context, task string) error,
	ch <-chan string,
) {
	wg := &sync.WaitGroup{}
	for i := 0; i < d.maxParallelTasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-ch:
					if !ok {
						return
					}
					err := f(ctx, task)
					if err != nil {
						_ = d.tm.Error(task, err)
					}
					_ = d.tm.Finish(task)
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
}

// OrderedDoer is a simple and powerful way to do your tasks in strict order.
type OrderedDoer struct {
	tm         TaskManager
	trimCutSet string
}

// NewDefaultOrderedDoer returns OrderedDoer with predefined trim cut sets
func NewDefaultOrderedDoer(tm TaskManager) *OrderedDoer {
	return &OrderedDoer{
		tm:         tm,
		trimCutSet: DefaultTrimCutSet,
	}
}

// NewOrderedDoer returns a new instance of OrderedDoer
func NewOrderedDoer(
	tm TaskManager,
	trimCutSet string,
) *OrderedDoer {
	return &OrderedDoer{
		tm:         tm,
		trimCutSet: trimCutSet,
	}
}

// Do the tasks one by one, with strict order
//
// If f returns an error, it will be delayed
func (d *OrderedDoer) Do(ctx context.Context, f func(ctx context.Context, task string) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		task, nodata := d.tm.Next()
		if nodata {
			return
		}
		task = strings.Trim(task, d.trimCutSet)
		err := f(ctx, task)
		if err != nil {
			_ = d.tm.Error(task, err)
			continue
		}
		_ = d.tm.Finish(task)
	}
}

