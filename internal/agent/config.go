package agent

import "time"

type Config struct {
	Version int   `json:"version"`
	Jobs    []Job `json:"jobs"`
}

type Job struct {
	ID                 string     `json:"id"`
	Binary             string     `json:"binary"`
	Cron               string     `json:"cron"`
	Enabled            bool       `json:"enabled"`
	LastRunAt          *time.Time `json:"lastRunAt,omitempty"`
	LastExitCode       *int       `json:"lastExitCode,omitempty"`
	LockUntil          *time.Time `json:"lockUntil,omitempty"`
	MaxRuntimeSeconds  int        `json:"maxRuntimeSeconds"`
}
