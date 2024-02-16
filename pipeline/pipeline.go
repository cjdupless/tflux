package pipeline

import (
	"github.com/google/uuid"
)

type taskSequence [][]*Task

type Pipeline struct {
	Schedule
	ExecScheme *ExecutionScheme
	id string
	dag *taskGraph
}

func NewPipeline() *Pipeline {
	pl := Pipeline{}
	pl.dag = newGraph()
	pl.ExecScheme = nil
	pl.id = uuid.New().String()
	return &pl
}

func (pl *Pipeline) buildTaskSequence() (result taskSequence){
	result = make(taskSequence, 0)
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
	
	upstream := []*Task{pl.dag.root}
	for { 
		downstream := getDownstream(upstream)
		if len(downstream) == 0 {
			break
		}
		result = append(result, downstream)
		upstream = downstream
	}
	return
}

// SetDAG - Not implemented yet
func (pl *Pipeline) SetDAG(DAG interface{}) {
	// DAG is a struct with all tasks in
	// relation to each other from start node to
	// leaf nodes.
	// DAG will be used to build pl.dag completely 
}

func (pl *Pipeline) AddTask(upstream, task *Task) error {
	err := pl.dag.addTask(upstream, task)
	if err != nil {
		return err
	}
	return nil
}