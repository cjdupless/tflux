package tflux

import (
	"fmt"
	"github.com/google/uuid"
)

type Pipeline struct {
	id      string
	taskDAG *dag
}

func NewPipeline() *Pipeline {
	pl := Pipeline{}
	pl.taskDAG = newDAG()
	pl.id = uuid.New().String()
	return &pl
}

func (pl *Pipeline) String() string {
	return fmt.Sprintf("%v", pl.taskDAG.root)
}

func (pl *Pipeline) Queue() *PRQueue {
	root:= pl.cloneTasks()
	queue := NewPRQ(root)
	return queue
}

func (pl *Pipeline) cloneTasks() *Task {
	cloneMap := make(map[*Task]*Task)
	for task := range pl.taskDAG.taskRefList {
		cloneMap[task] = task.Clone()
	}
	for task := range pl.taskDAG.taskRefList {
		for i, utask := range task.upstream {
			cloneMap[task].upstream[i] = cloneMap[utask]
		}
		for i, dtask := range task.downstream {
			cloneMap[task].downstream[i] = cloneMap[dtask]
		}
	}
	return cloneMap[pl.taskDAG.root]
}

func (pl *Pipeline) AddLink(upstream, task *Task) error {
	err := pl.taskDAG.tryAddLink(upstream, task)
	if err != nil {
		return err
	}
	return nil
}

func (pl *Pipeline) AddStart(task *Task) error {
	err := pl.taskDAG.addStart(task)
	if err != nil {
		return err
	}
	return nil
}

func (pl *Pipeline) From(task *Task) *linkingTask {
	if task == nil {
		panic("task cannot be nil")
	}
	if pl.taskDAG.root == nil {
		pl.taskDAG.addStart(task)
	}
	return &linkingTask{pl.taskDAG, task}
}
