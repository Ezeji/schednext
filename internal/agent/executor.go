package agent

import (
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"schednext/internal/util"
	"schednext/internal/runtime"
	"schednext/internal/model"
)

func executeJob(runtimeOptPath string, configName string, job *model.Job) int {
	cmdPath := filepath.Join(runtimeOptPath, job.Binary)
	cfgPath := filepath.Join(runtimeOptPath, configName)

	cmd := exec.Command(cmdPath)
	start := time.Now()

	runtime.State.Lock()

	if js, ok := runtime.State.Jobs[job.ID]; ok {
		js.Status = "running"
	}

	runtime.State.Unlock()

	err := cmd.Run()
	exitCode := 0

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 1
		}
	}

	now := time.Now()

	var cfg model.Config
    if err := util.ReadConfig(cfgPath, &cfg); err != nil {
        log.Println("failed to reload config:", err)
        return exitCode
    }

    for i := range cfg.Jobs {
        if cfg.Jobs[i].ID == job.ID {
            cfg.Jobs[i].LastRunAt = &now
            cfg.Jobs[i].LastExitCode = &exitCode
            cfg.Jobs[i].LockUntil = nil
        }
    }

    util.WriteConfigAtomic(cfgPath, &cfg)

	runtime.State.Lock()

	if js, ok := runtime.State.Jobs[job.ID]; ok {
		js.Status = "idle"
		js.LastRunAt = &now
		js.LastExitCode = &exitCode
		js.LockUntil = nil
	}

	runtime.State.Unlock()

	log.Printf("job %s finished in %s with exit %d",
		job.ID, time.Since(start), exitCode)

	return exitCode
}
