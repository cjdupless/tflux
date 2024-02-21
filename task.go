package tflux

import (
	"fmt"
	"sync"
	"github.com/google/uuid"
)

type linkDirection bool

const (
	down = false
	up   = true
)

type TaskResult struct {
	Info  string
	Error error
}

type Task struct {
	id         string
	status     Status
	function   func() TaskResult
	upstream   []*Task
	downstream []*Task
}

func NewTask() *Task {
	t := Task{
		upstream:   make([]*Task, 0),
		downstream: make([]*Task, 0),
	}
	t.status = NoneStatus
	t.id = uuid.New().String()
	return &t
}

func (t *Task) SetStatus(status Status) {
	t.status = status
	if t.status == FailedStatus || t.status == UpFailedStatus {
		for _, task := range t.downstream {
			task.SetStatus(UpFailedStatus)
		}
	}
}

func (t *Task) canRun() bool {
	if t.status == UpFailedStatus {
		return false
	}
	successCount := 0
	for _, task := range t.upstream {
		if task.status == SuccessStatus {
			successCount++
		}
	}
	return successCount == len(t.upstream)
}

func (t *Task) done() bool {
	return (
		t.status == FailedStatus || 
		t.status == UpFailedStatus || 
		t.status == SuccessStatus)
}

func (t *Task) Run() {
	result := t.function()
	if result.Error != nil {
		// Do some logging.
		t.SetStatus(FailedStatus)
		return
	}
	t.status = SuccessStatus
}

func (t *Task) Clone() *Task {
	clone := NewTask()
	clone.function = t.function
	return clone
}

func (t *Task) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *Task) addLink(lt *Task, direction linkDirection) {
	if direction == up {
		t.upstream = append(t.upstream, lt)
	} else {
		t.downstream = append(t.downstream, lt)
	}
}
