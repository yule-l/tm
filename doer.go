package tm

import (
	"strings"
)

type doer struct {
	tm         TaskManager
	trimCutSet string
}

// NewDefaultDoer returns doer with predefined trim cut sets
func NewDefaultDoer(tm TaskManager) *doer {
	return &doer{
		tm:         tm,
		trimCutSet: " ",
	}
}

// NewDoer returns a new instance of doer
func NewDoer(tm TaskManager, trimCutSet string) *doer {
	return &doer{
		tm:         tm,
		trimCutSet: trimCutSet,
	}
}

// Do the tasks on default way
//
// If f returns an error, it will be delayed
func (d *doer) Do(f func(task string) error) {
	for {
		task, nodata := d.tm.Next()
		if nodata {
			return
		}
		if d.trimCutSet != "" {
			task = strings.Trim(task, d.trimCutSet)
		}
		err := f(task)
		if err != nil {
			_ = d.tm.Error(task, err)
		}
		_ = d.tm.Finish(task)
	}
}
