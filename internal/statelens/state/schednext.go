package state

import (
	"bytes"
	"fmt"

	"schednext/internal/runtime"
)

func AllJobs() ([]byte, error) {

	var b bytes.Buffer

	runtime.State.RLock()
	defer runtime.State.RUnlock()

	for _, job := range runtime.State.Jobs {

		exit := -1

		if job.LastExitCode != nil {
			exit = *job.LastExitCode
		}

		fmt.Fprintf(
			&b,
			"%s | %s | enabled=%v | exit=%d\n",
			job.ID,
			job.Status,
			job.Enabled,
			exit,
		)
	}

	return b.Bytes(), nil
}

func JobDetails(jobName string) ([]byte, error) {

	var b bytes.Buffer

	runtime.State.RLock()
	defer runtime.State.RUnlock()

	job, ok := runtime.State.Jobs[jobName]

	if !ok {
		return []byte("job not found\n"), nil
	}

	fmt.Fprintf(&b, "id=%s\n", job.ID)
	fmt.Fprintf(&b, "status=%s\n", job.Status)
	fmt.Fprintf(&b, "enabled=%v\n", job.Enabled)

	if job.LastExitCode != nil {
		fmt.Fprintf(&b, "lastExitCode=%d\n", *job.LastExitCode)
	}

	if job.LastRunAt != nil {
		fmt.Fprintf(&b, "lastRunAt=%s\n", job.LastRunAt.String())
	}

	return b.Bytes(), nil
}