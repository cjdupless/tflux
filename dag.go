package tflux

import (
	"fmt"
	"slices"
)

type dag struct {
	taskRefList map[*Task]bool
	root        *Task
}

func newDAG() *dag {
	d := dag{}
	d.taskRefList = make(map[*Task]bool)
	return &d
}

func (d *dag) causesCycle(newNode *Task) bool {
	var isCyclic func(*Task) bool
	isCyclic = func(node *Task) bool {
		result := false
		for _, upsNode := range node.upstream {
			if upsNode == newNode {
				return true
			}
			if upsNode == d.root {
				return false
			}
			result = result || isCyclic(upsNode)
		}
		return result
	}
	return isCyclic(newNode)
}

func (d *dag) removeLink(usTask, dsTask *Task) {
	removeTask := func(list []*Task, task *Task) []*Task {
		delIndex := slices.Index(list, task)
		return append(
			list[0:delIndex],
			list[delIndex+1:]...,
		)
	}
	usTask.upstream = removeTask(usTask.upstream, dsTask)
	dsTask.downstream = removeTask(dsTask.downstream, usTask)
}

func (d *dag) addLink(usTask, dsTask *Task) {
	usTask.addLink(dsTask, down)
	dsTask.addLink(usTask, up)
}

func (d *dag) addStart(task *Task) error {
	if task == nil {
		return fmt.Errorf("start task cannot be nil")
	}
	d.taskRefList[task] = true
	d.root = task
	return nil
}

func (d *dag) tryAddLink(usTask, dsTask *Task) error {
	if usTask == nil {
		return fmt.Errorf("upstream task cannot be nil")
	}
	if dsTask == nil {
		return fmt.Errorf("downstream task cannot be nil")
	}

	if !d.taskRefList[usTask] {
		return fmt.Errorf("upstream task %v does not exist on DAG %v", *usTask, *d)
	}

	if d.taskRefList[dsTask] {
		return nil
	}

	d.addLink(usTask, dsTask)
	if d.causesCycle(dsTask) {
		d.removeLink(usTask, dsTask)
		return fmt.Errorf("link %v -> %v causes a cyclic graph", usTask, dsTask)
	}

	d.taskRefList[usTask] = true
	d.taskRefList[dsTask] = true
	return nil
}
