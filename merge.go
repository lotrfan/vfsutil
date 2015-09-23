package vfsutil

import (
	"io"
	"os"

	"github.com/blang/vfs"
)

func Merge(dest, src vfs.Filesystem) error {
	return Walk(src, "/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return vfs.MkdirAll(dest, path, info.Mode()|0700)
		}

		fs, err := vfs.Open(src, path)
		if err != nil {
			return err
		}
		defer fs.Close()

		fd, err := dest.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, info.Mode())
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = io.Copy(fd, fs)
		return err
	})
}
