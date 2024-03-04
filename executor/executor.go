/*
	This package serves as the singleton executor that will execute pipelines' PRQueues assigned to it.
*/
package executor

import (
	"sync"
	"github.com/cjdupless/tflux"
)

var (
	once  sync.Once
	queue chan *tflux.Task
)

func init() {
	once.Do(
		func() {
			queue = make(chan *tflux.Task)
		},
	)
	go processQueue()
}

func processQueue() {
	for task := range queue {
		task.Run()
	}
}

func Assign(prq *tflux.PRQueue) {
	go func() {
		for {
			task := prq.Next()
			if task == nil {
				break
			}
			task.SetStatus(tflux.QueuedStatus)
			queue <- task
		}
	}()
}
