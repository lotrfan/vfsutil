package vfsutil

import (
	"os"

	"github.com/blang/vfs"
)

type Chrootfs struct {
	fs   vfs.Filesystem
	root string
}

func Chroot(fs vfs.Filesystem, root string) *Chrootfs {
	return &Chrootfs{
		fs:   fs,
		root: root,
	}
}

func (fs *Chrootfs) concatPath(path string) string {
	if len(path) == 0 {
		return fs.root
	}
	if len(fs.root) == 0 {
		return path
	}
	if path[0] == fs.PathSeparator() {
		path = path[1:]
	}
	if fs.root[len(fs.root)-1] == fs.PathSeparator() {
		return fs.root + path
	} else {
		return fs.root + string(fs.PathSeparator()) + path
	}
}

func (fs *Chrootfs) PathSeparator() uint8 {
	return fs.fs.PathSeparator()
}
func (fs *Chrootfs) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	return fs.fs.OpenFile(fs.concatPath(name), flag, perm)
}
func (fs *Chrootfs) Remove(name string) error {
	return fs.fs.Remove(fs.concatPath(name))
}
func (fs *Chrootfs) Rename(oldpath, newpath string) error {
	return fs.fs.Rename(fs.concatPath(oldpath), fs.concatPath(newpath))
}
func (fs *Chrootfs) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(fs.concatPath(name), perm)
}
func (fs *Chrootfs) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(fs.concatPath(name))
}
func (fs *Chrootfs) Lstat(name string) (os.FileInfo, error) {
	return fs.fs.Lstat(fs.concatPath(name))
}
func (fs *Chrootfs) ReadDir(path string) ([]os.FileInfo, error) {
	return fs.fs.ReadDir(fs.concatPath(path))
}
