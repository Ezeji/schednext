package statelens

import (
	"schednext/internal/statelens/mount"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func Mount(path string) (*fuse.Conn, error) {

	c, err := fuse.Mount(
		path,
		fuse.ReadOnly(),
		fuse.FSName("statelens"),
	)

	if err != nil {
		return nil, err
	}

	go func() {
		fs.Serve(c, &mount.StateFS{})
	}()

	return c, nil
}