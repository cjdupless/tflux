package tflux

import (
	"fmt"
	"github.com/google/uuid"
)

type TaskResult struct {
	ExtraInfo string
	ExecResult interface{}
	Status Status
	Error error
}

type ExecFunc func() TaskResult

type Task struct {
	id string
	status Status 
	executable ExecFunc
	upstream map[*Task]bool	
	downstream []*Task
}

func NewTask() *Task {
	t := Task{
		upstream: make(map[*Task]bool),
		downstream: make([]*Task, 0),
	}
	t.status = NoneStatus
	t.id = uuid.New().String()  
	return &t
}

func (t *Task) SetExec(exec ExecFunc) error {
	if exec == nil {
		return fmt.Errorf("the executable cannot be nil")
	}
	t.executable = exec
}

func (t *Task) Clone() *Task {
	clone := NewTask()
	clone.executable = t.executable
	return clone
}

func (t *Task) String() string {
	return fmt.Sprintf("%v", *t)
}

func (t *Task) Run() {
	// Pre-execution instructions
	err := t.executable()
	if err != nil {
		// Do some logging.
		t.status = FailedStatus
		return
	}
	t.status = SuccessStatus
	// Post-execution instructions
}

func (t *Task) addDownstream(dt *Task) {
	t.downstream = append(t.downstream, dt)
}

func (t *Task) addUpstream(ut *Task) {
	t.upstream[ut] = true
}