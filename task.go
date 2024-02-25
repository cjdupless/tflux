package tflux

import (
	"fmt"
	"github.com/google/uuid"
)

type linkDirection uint8

const (
	down = iota
	up
)

type OpResult struct {
	Info  string
	Error error
}

type Task struct {
	id         string
	status     Status
	name       string
	function   func() OpResult
	upstream   []*Task
	downstream []*Task
}

func NewTask(name string, function func() OpResult) *Task {
	t := Task{
		upstream:   make([]*Task, 0),
		downstream: make([]*Task, 0),
	}
	t.name = name
	t.function = function
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
	return (t.status == FailedStatus ||
		t.status == UpFailedStatus ||
		t.status == SuccessStatus)
}

func (t *Task) Run() {
	defer func() {
		err := recover()
		if err != nil {
			t.SetStatus(FailedStatus)
		}
	}()

	result := t.function()
	if result.Error != nil {
		t.SetStatus(FailedStatus)
	} else {
		t.SetStatus(SuccessStatus)
	}
}

func (t *Task) Clone() *Task {
	clone := Task{
		name: t.name,
		function: t.function,
		upstream: make([]*Task, len(t.upstream)),
		downstream: make([]*Task, len(t.downstream)),
	}
	return &clone
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
