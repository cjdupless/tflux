package tflux

import "time"

type Executable interface {
	Run()
}

type LinkedExecutable interface {
	Executable
	Upstream() map[*Executable]bool
	Downstream() []*Executable
}

type SchedulePlanner interface {
	NextEvent() time.Time
}

type ExecutionPlanner interface {
	NextExec() *Executable
}
