package tflux

import (
	"github.com/google/uuid"
)

// PRQueue is the Pipeline Run Queue specific to each pipeline
type PRQueue struct {
	// mx sync.Mutex
	runID string
	taskList []*Task
	executionStages [][]*Task
}

func NewPRQ(dagRoot *Task, taskList []*Task) *PRQueue {
	prq := PRQueue{}
	prq.runID = uuid.New().String()
	prq.taskList = make([]*Task, len(taskList))
	copy(prq.taskList, taskList)
	prq.setExecutionStages(dagRoot)
	return &prq
}

func (prq *PRQueue) setExecutionStages(dagRoot *Task) {
	prq.executionStages = make([][]*Task, 0)

	getDownstream := func(tasks []*Task) (nextTasks []*Task) {
		nextTasks = make([]*Task, 0)
		for _, task := range tasks {
			if len(task.downstream) == 0 {
				continue
			}
			nextTasks = append(nextTasks, task.downstream...)
		}
		return
	}

	upstream := []*Task{dagRoot}
	prq.executionStages = append(prq.executionStages, upstream)
	for {
		downstream := getDownstream(upstream)
		if len(downstream) == 0 {
			break
		}
		prq.executionStages = append(prq.executionStages, downstream)
		upstream = downstream
	}
}

func (prq *PRQueue) Done() bool {
	count := 0
	for _, task := range prq.taskList {
		if task.done() {
			count++
		}
	}
	return count == len(prq.taskList)
}

func (prq *PRQueue) Next() *Task {
	for _, stage := range prq.executionStages {
		for _, task := range stage {
			if task.done() {
				continue
			}
			if task.canRun() {
				return task
			}
		}
	}
	return nil
}