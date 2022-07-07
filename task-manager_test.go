package tm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"testing"
)

func Test_taskManager_Next(t1 *testing.T) {
	type fields struct {
		tasksInfo    *tasks
		delayedTasks *delayedTasks
		tasksQueue   *bufio.Scanner
		tasksResult  io.ReadWriter
		maxAttempts  uint8
	}
	tests := []struct {
		name        string
		fields      fields
		wantTask    string
		wantNotasks bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &taskManager{
				tasksInfo:    tt.fields.tasksInfo,
				delayedTasks: tt.fields.delayedTasks,
				tasksQueue:   tt.fields.tasksQueue,
				tasksResult:  tt.fields.tasksResult,
				maxAttempts:  tt.fields.maxAttempts,
			}
			gotTask, gotNotasks := t.Next()
			if gotTask != tt.wantTask {
				t1.Errorf("Next() gotTask = %v, want %v", gotTask, tt.wantTask)
			}
			if gotNotasks != tt.wantNotasks {
				t1.Errorf("Next() gotNotasks = %v, want %v", gotNotasks, tt.wantNotasks)
			}
		})
	}
}

func Test_taskManager_fillResultsWithErrorIfNotNil(t *testing.T) {
	tm := &taskManager{
		maxAttempts: 3,
	}
	t.Run("error is nil", func(t *testing.T) {
		r := &result{
			Payload: "sometask1",
			Attempt: 0,
			Errors:  nil,
			Status:  done,
		}
		tm.fillResultsWithErrorIfNotNil(r, nil)
		if r.Payload != "sometask1" {
			t.Errorf("unexpected payload, expected = sometask1, got = %s", r.Payload)
		}
		if r.Attempt != 0 {
			t.Errorf("unexpected attemnt, expected = 0, got = %d", r.Attempt)
		}
		if r.Errors != nil {
			t.Errorf("unexpected errors, expected = nil, got = %v", r.Errors)
		}
		if r.Status != done {
			t.Errorf("unexpected status, expected = 0, got = %d", r.Status)
		}
	})

	t.Run("error is not nil", func(t *testing.T) {
		r := &result{
			Payload: "sometask1",
			Attempt: 0,
			Errors:  nil,
			Status:  done,
		}
		tm.fillResultsWithErrorIfNotNil(r, errors.New("someerror"))
		if r.Payload != "sometask1" {
			t.Errorf("unexpected payload, expected = sometask1, got = %s", r.Payload)
		}
		if r.Attempt != 0 {
			t.Errorf("unexpected attemnt, expected = 0, got = %d", r.Attempt)
		}
		if r.Errors[0].Attempt != 0 {
			t.Errorf("unexpected errors[0].attempt, expected = 0, got = %v", r.Errors[0].Attempt)
		}
		if r.Errors[0].Message != "someerror" {
			t.Errorf("unexpected errors[0].message, expected = someerror, got = %s", r.Errors[0].Message)
		}
		if r.Status != delay {
			t.Errorf("unexpected status, expected = 0, got = %d", r.Status)
		}
	})

	t.Run("error is not nil and attempts reached", func(t *testing.T) {
		r := &result{
			Payload: "sometask1",
			Attempt: 1,
			Errors:  nil,
			Status:  done,
		}
		tm.maxAttempts = 0
		tm.fillResultsWithErrorIfNotNil(r, errors.New("someerror"))
		if r.Payload != "sometask1" {
			t.Errorf("unexpected payload, expected = sometask1, got = %s", r.Payload)
		}
		if r.Attempt != 1 {
			t.Errorf("unexpected attemnt, expected = 1, got = %d", r.Attempt)
		}
		if r.Errors[0].Attempt != 1 {
			t.Errorf("unexpected errors[0].attempt, expected = 1, got = %v", r.Errors[0].Attempt)
		}
		if r.Errors[0].Message != "someerror" {
			t.Errorf("unexpected errors[0].message, expected = someerror, got = %s", r.Errors[0].Message)
		}
		if r.Status != done {
			t.Errorf("unexpected status, expected = 1, got = %d", r.Status)
		}
	})
}

func Test_taskManager_nextDelayedTask(t *testing.T) {
	tm := &taskManager{
		tasksInfo:    newTasks(),
		delayedTasks: newDelayedTasks(),
	}

	t.Run("check delayed tasks on init", func(t *testing.T) {
		if task, notasks := tm.nextDelayedTask(); task != "" || notasks != true {
			t.Errorf("expected task == \"\" and notasks = true, got %s and %v", task, notasks)
		}
	})

	t.Logf("fill delayed tasks with 3 elements")

	tm.delayedTasks.append("sometask1")
	tm.delayedTasks.append("sometask2")
	tm.delayedTasks.append("sometask3")

	t.Run("check tasks order", func(t *testing.T) {
		if task, notasks := tm.nextDelayedTask(); task != "sometask1" || notasks != false {
			t.Errorf("expected task == sometask1 and notasks = false, got %s and %v", task, notasks)
		}
		if task, notasks := tm.nextDelayedTask(); task != "sometask2" || notasks != false {
			t.Errorf("expected task == sometask2 and notasks = false, got %s and %v", task, notasks)
		}
		if task, notasks := tm.nextDelayedTask(); task != "sometask3" || notasks != false {
			t.Errorf("expected task == sometask3 and notasks = false, got %s and %v", task, notasks)
		}
	})

	t.Logf("no tasks, check next")

	t.Run("no available tasks", func(t *testing.T) {
		if task, notasks := tm.nextDelayedTask(); task != "" || notasks != true {
			t.Errorf("expected task == \"\" and notasks = true, got %s and %v", task, notasks)
		}
	})
}

func Test_taskManager_nextQueueTask(t *testing.T) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	scanner := bufio.NewScanner(buffer)
	scanner.Split(bufio.ScanLines)
	tm := &taskManager{
		tasksInfo:  newTasks(),
		tasksQueue: scanner,
	}

	t.Run("check queued tasks on init", func(t *testing.T) {
		if task, notasks := tm.nextQueueTask(); task != "" || notasks != true {
			t.Errorf("expected task == \"\" and notasks = true, got %s and %v", task, notasks)
		}
	})

	t.Logf("fill queued tasks with 3 elements")

	buffer.WriteString("sometask1\n")
	buffer.WriteString("sometask2\n")
	buffer.WriteString("sometask3\n")
	scanner = bufio.NewScanner(buffer)
	scanner.Split(bufio.ScanLines)
	tm.tasksQueue = scanner

	t.Run("check tasks order", func(t *testing.T) {
		if task, notasks := tm.nextQueueTask(); task != "sometask1" || notasks != false {
			t.Errorf("expected task == sometask1 and notasks = false, got %s and %v", task, notasks)
		}
		if task, notasks := tm.nextQueueTask(); task != "sometask2" || notasks != false {
			t.Errorf("expected task == sometask2 and notasks = false, got %s and %v", task, notasks)
		}
		if task, notasks := tm.nextQueueTask(); task != "sometask3" || notasks != false {
			t.Errorf("expected task == sometask3 and notasks = false, got %s and %v", task, notasks)
		}
	})

	t.Logf("no tasks, check next")

	t.Run("no available tasks", func(t *testing.T) {
		if task, notasks := tm.nextQueueTask(); task != "" || notasks != true {
			t.Errorf("expected task == \"\" and notasks = true, got %s and %v", task, notasks)
		}
	})
}

func Test_taskManager_writeResult(t *testing.T) {
	t.Run("just add one task", func(t *testing.T) {
		buf := bytes.NewBuffer(make([]byte, 0, 10240))
		tm := &taskManager{
			tasksInfo:    newTasks(),
			delayedTasks: newDelayedTasks(),
			tasksResult:  buf,
			maxAttempts:  3,
		}
		err := tm.writeResult("sometask1", nil)
		if err != nil {
			t.Errorf("unexpected error = %v", err)
		}
		r := &result{}
		_ = json.Unmarshal(buf.Bytes(), r)

		if r.Payload != "sometask1" {
			t.Errorf("unexpected r.Payload, expected = sometask1, got = %s", r.Payload)
		}
		if r.Attempt != 1 {
			t.Errorf("unexpected r.Attemnt, expected = 0, got = %d", r.Attempt)
		}
		if r.Errors != nil {
			t.Errorf("unexpected r.Errors, expected = nil, got = %v", r.Errors)
		}
		if r.Status != done {
			t.Errorf("unexpected r.Status, expected = 0, got = %d", r.Status)
		}

		rs, ok := tm.tasksInfo.load("sometask1")
		if !ok {
			t.Errorf("unexpected tasks info value, no data about task = sometask1")
		}

		if rs.Payload != "sometask1" {
			t.Errorf("unexpected rs.Payload, expected = sometask1, got = %s", rs.Payload)
		}
		if rs.Attempt != 1 {
			t.Errorf("unexpected rs.Attemnt, expected = 1, got = %d", rs.Attempt)
		}
		if rs.Errors != nil {
			t.Errorf("unexpected rs.Rrrors, expected = nil, got = %v", rs.Errors)
		}
		if rs.Status != done {
			t.Errorf("unexpected rs.Status, expected = 0, got = %d", rs.Status)
		}
	})

	t.Run("just add one task with an error", func(t *testing.T) {
		buf := bytes.NewBuffer(make([]byte, 0, 10240))
		tm := &taskManager{
			tasksInfo:    newTasks(),
			delayedTasks: newDelayedTasks(),
			tasksResult:  buf,
			maxAttempts:  3,
		}
		err := tm.writeResult("sometask1", errors.New("someerror"))
		if err != nil {
			t.Errorf("unexpected error = %v", err)
		}
		r := &result{}
		_ = json.Unmarshal(buf.Bytes(), r)

		if r.Payload != "sometask1" {
			t.Errorf("unexpected r.Payload, expected = sometask1, got = %s", r.Payload)
		}
		if r.Attempt != 1 {
			t.Errorf("unexpected r.Attemnt, expected = 0, got = %d", r.Attempt)
		}
		if r.Errors == nil {
			t.Errorf("unexpected r.Errors, expected = slice, got = nil")
		}
		if len(r.Errors) != 1 {
			t.Errorf("unexpected len(r.Errors), expected = 1, got = %d", len(r.Errors))
		}
		if r.Errors[0].Message != "someerror" {
			t.Errorf("unexpected r.Errors[0].Message, expected = someerror, got = %s", r.Errors[0].Message)
		}
		if r.Errors[0].Attempt != 1 {
			t.Errorf("unexpected r.Errors[0].Message, expected = 0, got = %d", r.Errors[0].Attempt)
		}
		if r.Status != delay {
			t.Errorf("unexpected r.Status, expected = 1, got = %d", r.Status)
		}

		rs, ok := tm.tasksInfo.load("sometask1")
		if !ok {
			t.Errorf("unexpected tasks info value, no data about task = sometask1")
		}

		if rs.Payload != "sometask1" {
			t.Errorf("unexpected rs.Payload, expected = sometask1, got = %s", rs.Payload)
		}
		if rs.Attempt != 1 {
			t.Errorf("unexpected rs.Attemnt, expected = 1, got = %d", rs.Attempt)
		}
		if rs.Errors == nil {
			t.Errorf("unexpected r.Errors, expected = slice, got = nil")
		}
		if len(rs.Errors) != 1 {
			t.Errorf("unexpected len(r.Errors), expected = 1, got = %d", len(rs.Errors))
		}
		if rs.Errors[0].Message != "someerror" {
			t.Errorf("unexpected r.Errors[0].Message, expected = someerror, got = %s", rs.Errors[0].Message)
		}
		if rs.Errors[0].Attempt != 1 {
			t.Errorf("unexpected r.Errors[0].Message, expected = 0, got = %d", rs.Errors[0].Attempt)
		}
		if rs.Status != delay {
			t.Errorf("unexpected rs.Status, expected = 1, got = %d", rs.Status)
		}
	})

}
