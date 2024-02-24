package tflux

type LinkingTask struct {
	dag  *dag
	task *Task
}

func (lt *LinkingTask) GoTo(thisTask *Task) *LinkingTask {
	err := lt.dag.tryAddLink(lt.task, thisTask)
	if err != nil {
		panic(err)
	}
	return &LinkingTask{lt.dag, thisTask}
}
