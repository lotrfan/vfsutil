package vfsutil

import (
	"io/ioutil"

	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
)

var (
	// fs1 has:
	// /
	//   file -> contains "test"
	fs1 vfs.Filesystem

	// fs1 has:
	// /
	//   file -> contains "test"
	//	 directory/
	//     file -> contains "test in a directory"
	fs2 vfs.Filesystem
)

func init() {
	fs1 = memfs.Create()
	writeFile(fs1, "/file", []byte("test"))
	fs1 = vfs.ReadOnly(fs1)

	fs2 = memfs.Create()
	writeFile(fs2, "/file", []byte("test"))
	fs2.Mkdir("/directory", 0700)
	writeFile(fs2, "/directory/file", []byte("test in a directory"))
	fs2 = vfs.ReadOnly(fs2)
}

func writeFile(fs vfs.Filesystem, path string, b []byte) error {
	f, err := vfs.Create(fs, path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}

func readFile(fs vfs.Filesystem, path string) ([]byte, error) {
	f, err := vfs.Open(fs, path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}
