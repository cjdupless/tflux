package tflux

type ExecutionScheme interface {
	GetNextTask() *Task
}

type DefaultExecScheme struct {
	taskSeq taskSequence
}

func NewDefaultExecScheme(seq taskSequence) *DefaultExecScheme {
	des := DefaultExecScheme{taskSeq: seq}
	return &des
}

func (des *DefaultExecScheme) GetNextTask() *Task {
	return nil
}

