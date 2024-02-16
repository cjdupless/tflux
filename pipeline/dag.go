package pipeline

import (
	"fmt"
	"slices"
	flx "github.com/cjdupless/tflux"
)

type dag struct {
	taskRefList map[*flx.Executable]bool
	root *flx.Executable
}

func (d *dag) String() string {
	return fmt.Sprintf(
		"%v\n %v\n", d.taskRefList, d.root,
	)
}

func newGraph() *dag {
	tg := dag{}
	tg.taskRefList = make(map[*flx.Executable]bool)
	return &tg
}

func (tg *dag) causesCycle(newNode *flx.Executable) bool {
	var isCyclic func(*flx.Executable) bool
	isCyclic = func(node *flx.Executable) bool {
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

func (tg *dag) removeLink(usTask, dsTask *flx.Executable) {
	delIndex := slices.Index(usTask.downstream, dsTask)
	usTask.downstream = append(
		usTask.downstream[0: delIndex],
		usTask.downstream[delIndex + 1:]...
	)
	delete(dsTask.upstream, usTask)
}

func (tg *dag) addLink(usTask, dsTask *flx.Executable) {
	usTask.addDownstream(dsTask)
	dsTask.addUpstream(usTask)
}

func (tg *dag) addTask(usTask, dsTask *flx.Executable) error {
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

