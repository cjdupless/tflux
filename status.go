package tflux

import "fmt"

type Status uint

const (
	NoneStatus = iota
	QueuedStatus
	ExecutingStatus
	SuccessStatus
	FailedStatus
)

func (s Status) String() string{
	switch s {
	case NoneStatus:
		return "Status: None"
	case QueuedStatus:
		return "Status: Queued"
	case ExecutingStatus:
		return "Status: Executing"
	case SuccessStatus:
		return "Status: Success"
	case FailedStatus:
		return "Status: Failed"
	default:
		return "Status: Unknown"
	}
}

func (s Status) Check() error {
	if s > 4 {
		return fmt.Errorf("invalid status, %v", s)
	}
	return nil
}