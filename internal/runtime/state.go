package runtime

import (
	"sync"
	"time"
)

type JobState struct {
	ID                string
	Enabled           bool
	Status            string
	LastRunAt         *time.Time
	LastExitCode      *int
	LockUntil         *time.Time
	MaxRuntimeSeconds int
}

type RuntimeState struct {
	sync.RWMutex
	Jobs map[string]*JobState
}

var State = RuntimeState{
	Jobs: map[string]*JobState{},
}