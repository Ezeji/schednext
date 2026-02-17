package agent

import (
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"schednext/internal/util"
)

func executeJob(userHome string, configName string, job *Job) int {
	cmdPath := filepath.Join(userHome, job.Binary)
	cfgPath := filepath.Join(userHome, configName)

	cmd := exec.Command(cmdPath)
	start := time.Now()

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

	var cfg Config
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

	log.Printf("job %s finished in %s with exit %d",
		job.ID, time.Since(start), exitCode)

	return exitCode
}
