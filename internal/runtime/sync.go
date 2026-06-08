package runtime

import "schednext/internal/model"

func SyncFromConfig(cfg model.Config) {
	State.Lock()
	defer State.Unlock()

	for _, job := range cfg.Jobs {
		js, ok := State.Jobs[job.ID]

		if !ok {
			js = &JobState{
				ID: job.ID,
			}
			State.Jobs[job.ID] = js
		}

		js.Enabled = job.Enabled
		js.LastRunAt = job.LastRunAt
		js.LastExitCode = job.LastExitCode
		js.LockUntil = job.LockUntil
		js.MaxRuntimeSeconds = job.MaxRuntimeSeconds

		if job.LockUntil != nil {
			js.Status = "running"
		} else {
			js.Status = "idle"
		}
	}
}