package lens

import (
	"context"
	"os"

	"schednext/internal/runtime"
	"schednext/internal/statelens/state"
	"schednext/internal/statelens/vfs"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type SchedNextDir struct{}

func NewSchedNextDir() *SchedNextDir {
	return &SchedNextDir{}
}

func (d *SchedNextDir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0555
	return nil
}

func (d *SchedNextDir) Lookup(ctx context.Context, name string) (fs.Node, error) {

	if name == "all" {
		return vfs.NewDynamicFile(state.AllJobs), nil
	}

	runtime.State.RLock()
	defer runtime.State.RUnlock()

	if _, ok := runtime.State.Jobs[name]; ok {

		return vfs.NewDynamicFile(
			func() ([]byte, error) {
				return state.JobDetails(name)
			},
		), nil
	}

	return nil, fuse.ENOENT
}

func (d *SchedNextDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {

	entries := []fuse.Dirent{
		{
			Name: "all",
			Type: fuse.DT_File,
		},
	}

	runtime.State.RLock()
	defer runtime.State.RUnlock()

	for jobName := range runtime.State.Jobs {
		entries = append(entries, fuse.Dirent{
			Name: jobName,
			Type: fuse.DT_File,
		})
	}

	return entries, nil
}