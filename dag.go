package tflux

import (
	"fmt"
	"slices"
)

type dag struct {
	taskRefList map[*Task]bool
	root        *Task
}

func (d *dag) String() string {
	return fmt.Sprintf(
		"%v\n %v\n", d.taskRefList, d.root,
	)
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

func (d *dag) addTask(usTask, dsTask *Task) error {
	if !d.taskRefList[usTask] {
		return fmt.Errorf("task %v does not exist", *usTask)
	}
	if dsTask == nil {
		return fmt.Errorf("task %v does not exist", *usTask)
	}
	d.addLink(usTask, dsTask)
	if d.causesCycle(dsTask) {
		d.removeLink(usTask, dsTask)
		return fmt.Errorf("link %v -> %v causes a cyclic graph", usTask, dsTask)
	}
	d.taskRefList[dsTask] = true
	if d.root == nil {
		d.root = usTask
	}
	return nil
}
