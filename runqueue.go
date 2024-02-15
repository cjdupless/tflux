package tflux

// import "slices"

type TaskQueue struct {
	queue []*Task
}

func (tq *TaskQueue) push(task *Task) {
	tq.queue = append(tq.queue, task)
}

func (tq *TaskQueue) pop() *Task {
	if len(tq.queue) == 0 {
		return nil
	} 
	task := tq.queue[0]
	tq.queue = tq.queue[1:]
	return task
}

type RunQueue struct {
	queue *TaskQueue
}

func NewRunQueue() *RunQueue {
	rq := RunQueue{}
	rq.queue = &TaskQueue{queue: make([]*Task, 0)}
	return &rq
}

func (rq *RunQueue) Enlist(pl *RunScheme) {

}
