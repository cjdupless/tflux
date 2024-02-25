package tflux

type linkingTask struct {
	dag  *dag
	task *Task
}

func (lt *linkingTask) GoTo(thisTask *Task) *linkingTask {
	err := lt.dag.tryAddLink(lt.task, thisTask)
	if err != nil {
		panic(err)
	}
	return &linkingTask{lt.dag, thisTask}
}
