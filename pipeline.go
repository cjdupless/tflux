package tflux

import (
	"github.com/google/uuid"
	"slices"
)

type Pipeline struct {
	id      string
	executionStages [][]*Task
	taskDAG *dag
}

func NewPipeline() *Pipeline {
	pl := Pipeline{}
	pl.taskDAG = newDAG()
	pl.id = uuid.New().String()
	return &pl
}

func (pl *Pipeline) buildExecutionStages() ([][]*Task, int) {
	stages := make([][]*Task, 0)
	var getDownstream func([]*Task) []*Task

	getDownstream = func(tasks []*Task) (nextTasks []*Task) {
		nextTasks = make([]*Task, 0)
		for _, task := range tasks {
			if len(task.downstream) == 0 {
				continue
			}
			nextTasks = append(nextTasks, task.downstream...)
		}
		return
	}

	var numTasks int
	upstream := []*Task{pl.taskDAG.root}
	for {
		numTasks += len(upstream)
		downstream := getDownstream(upstream)
		if len(downstream) == 0 {
			break
		}
		stages = append(stages, downstream)
		upstream = downstream
	}
	return stages, numTasks
}

func (pl *Pipeline) AddTask(upstream, task *Task) error {
	err := pl.taskDAG.addTask(upstream, task)
	if err != nil {
		return err
	}
	return nil
}
