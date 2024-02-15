package tflux

import (
	"fmt"
	"slices"
)

type taskGraph struct {
	taskRefList map[*Task]bool
	root *Task
}

func (tg *taskGraph) String() string {
	return fmt.Sprintf(
		"%v\n %v\n", tg.taskRefList, tg.root,
	)
}

func newGraph() *taskGraph {
	tg := taskGraph{}
	tg.taskRefList = make(map[*Task]bool)
	return &tg
}

func (tg *taskGraph) causesCycle(newNode *Task) bool {
	var isCyclic func(*Task) bool
	isCyclic = func(node *Task) bool {
		result := false
		for upsNode := range node.upstream {
			if upsNode == newNode {
				return true
			}
			if upsNode == tg.root {
				return false
			}
			result = result || isCyclic(upsNode)
		}
		return result
	}
	return isCyclic(newNode)
}

func (tg *taskGraph) removeLink(usTask, dsTask *Task) {
	delIndex := slices.Index(usTask.downstream, dsTask)
	usTask.downstream = append(
		usTask.downstream[0: delIndex],
		usTask.downstream[delIndex + 1:]...
	)
	delete(dsTask.upstream, usTask)
}

func (tg *taskGraph) addLink(usTask, dsTask *Task) {
	usTask.addDownstream(dsTask)
	dsTask.addUpstream(usTask)
}

func (tg *taskGraph) addTask(usTask, dsTask *Task) error {
	if !tg.taskRefList[usTask] {
		return fmt.Errorf("task %v does not exist", *usTask)
	}
	if dsTask == nil {
		return fmt.Errorf("task %v does not exist", *usTask)
	}
	tg.addLink(usTask, dsTask)
	if tg.causesCycle(dsTask) {
		tg.removeLink(usTask, dsTask)
		return fmt.Errorf("link %v -> %v causes a cyclic graph", usTask, dsTask)
	}
	tg.taskRefList[dsTask] = true
	if tg.root == nil {
		tg.root = usTask
	}
	return nil
}

