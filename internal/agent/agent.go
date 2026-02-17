package agent

import (
	"log"
	"os"
	"path/filepath"
	showtime "time"

	"schednext/internal/util"
)

const configName = "schednext.config"

func RunAgent(homeRoot string) {
	StartIPCServer(homeRoot)

	ticker := showtime.NewTicker(30 * showtime.Second)

	for {
		scanUsers(homeRoot)
		<-ticker.C
	}
}

func scanUsers(homeRoot string) {
	entries, err := os.ReadDir(homeRoot)
	if err != nil {
		log.Println("failed to read home dir:", err)
		return
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		userHome := filepath.Join(homeRoot, e.Name())
		cfgPath := filepath.Join(userHome, configName)

		processConfig(userHome, cfgPath)
	}
}

func processConfig(userHome, cfgPath string) {
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0644)
	if err != nil {
		return 
	}
	defer f.Close()

	if err := lockFile(f); err != nil {
		return
	}
	defer unlockFile(f)

	var cfg Config
	if err := util.ReadConfig(cfgPath, &cfg); err != nil {
		log.Println("invalid config:", err)
		return
	}

	now := showtime.Now()

	for i := range cfg.Jobs {
		job := &cfg.Jobs[i]

		if shouldRun(*job, now) {
			lockUntil := now.Add(showtime.Duration(job.MaxRuntimeSeconds) * showtime.Second)
			job.LockUntil = &lockUntil

			util.WriteConfigAtomic(cfgPath, &cfg)
			go executeJob(userHome, configName, job)
		}
	}
}
