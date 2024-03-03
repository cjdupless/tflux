package tflux

import (
	"fmt"
)

// PRQueue is the Pipeline Run Queue specific to each pipeline
type PRQueue struct {
	tasksDone []*Task
	taskQueue []*Task
}

func NewPRQ(dagRoot *Task) *PRQueue {
	prq := PRQueue{}
	prq.buildQueue(dagRoot)
	return &prq
}

func (prq *PRQueue) buildQueue(dagRoot *Task) {
	prq.taskQueue = make([]*Task, 0)

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
	prq.taskQueue = append(prq.taskQueue, upstream...)
	for {
		downstream := getDownstream(upstream)
		if len(downstream) == 0 {
			break
		}
		prq.taskQueue = append(prq.taskQueue, downstream...)
		upstream = downstream
	}
}

func (prq *PRQueue) cleanup() {
	delIndices := make([]int, 0)
	for index, task := range prq.taskQueue {
		if task.done() {
			delIndices = append(delIndices, index)
		}
	}
	for i, di := range delIndices {
		// Everytime you remove an item the next index to 
		shiftedDI := di - i // delete shifts i positions to the left  
		prq.tasksDone = append(prq.tasksDone, prq.taskQueue[shiftedDI])
		prq.taskQueue = append(
			prq.taskQueue[0:shiftedDI],
			prq.taskQueue[shiftedDI+1:]...
		)
	}
}

func (prq *PRQueue) stillAlive() *Task {
	return NewTask(
		"STILL_ALIVE",
		func() OpResult {
			fmt.Println("...")
			return OpResult{
				Error: nil,
			}
		},
	)
}

func (prq *PRQueue) Next() *Task {
	prq.cleanup()

	if len(prq.taskQueue) == 0 {
		return nil
	}

	for _, task := range prq.taskQueue {
		if task.canRun() {
			return task
		}
	}
	return prq.stillAlive()
}
