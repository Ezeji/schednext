package agent

import (
	"log"
	"os"
	"path/filepath"
	showtime "time"

	"schednext/internal/util"
	"schednext/internal/runtime"
	"schednext/internal/model"
)

const configName = "schednext.config"
const runtimePathName = "schednext-runtime"

var stopChan = make(chan struct{})

func RunAgent(optPath string) {
	StartIPCServer(optPath)

	ticker := showtime.NewTicker(30 * showtime.Second)

	for {
		select {
		case <-ticker.C:
			locateRuntime(optPath)

		case <-stopChan:
			log.Println("agent shutting down")
			return
		}
	}
}

func Shutdown() {
	close(stopChan)
}

func locateRuntime(optPath string) {
	runtimeOptPath := filepath.Join(optPath, runtimePathName)
	cfgPath := filepath.Join(runtimeOptPath, configName)

	processConfig(runtimeOptPath, cfgPath)
}

func processConfig(runtimeOptPath, cfgPath string) {
	f, err := os.OpenFile(cfgPath, os.O_RDWR, 0644)
	if err != nil {
		return 
	}
	defer f.Close()

	if err := lockFile(f); err != nil {
		return
	}
	defer unlockFile(f)

	var cfg model.Config
	if err := util.ReadConfig(cfgPath, &cfg); err != nil {
		log.Println("invalid config:", err)
		return
	}

	runtime.SyncFromConfig(cfg)

	now := showtime.Now()

	for i := range cfg.Jobs {
		job := &cfg.Jobs[i]

		if shouldRun(*job, now) {
			lockUntil := now.Add(showtime.Duration(job.MaxRuntimeSeconds) * showtime.Second)
			job.LockUntil = &lockUntil

			util.WriteConfigAtomic(cfgPath, &cfg)
			go executeJob(runtimeOptPath, configName, job)
		}
	}
}
