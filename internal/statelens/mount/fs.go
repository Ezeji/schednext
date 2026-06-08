package mount

import (
	"schednext/internal/statelens/lens"

	"bazil.org/fuse/fs"
)

type StateFS struct{}

func (StateFS) Root() (fs.Node, error) {
	return &RootDir{
		cpu: lens.NewCPUDir(),
		mem: lens.NewMemDir(),
		net: lens.NewNetDir(),
		schednext: lens.NewSchedNextDir(),
	}, nil
}
